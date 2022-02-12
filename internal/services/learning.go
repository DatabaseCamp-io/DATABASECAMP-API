package services

import (
	"database-camp/internal/errs"
	"database-camp/internal/logs"
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/models/request"
	"database-camp/internal/models/response"
	"database-camp/internal/repositories"
	"database-camp/internal/services/loaders"
	"database-camp/internal/utils"
)

type LearningService interface {
	GetVideoLecture(id int) (*response.VideoLectureResponse, error)
	GetOverview(userID int) (*response.ContentOverviewResponse, error)
	GetActivity(userID int, activityID int) (*response.ActivityResponse, error)
	GetRecommend(userID int) (*response.RecommendResponse, error)
	UseHint(userID int, activityID int) (*response.UsedHintResponse, error)
	GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error)
	CheckAnswer(userID int, request request.CheckAnswerRequest) (*response.AnswerResponse, error)
	CheckPeerReview(userID int, request request.PeerReviewRequest) (*response.AnswerResponse, error)
}

type learningService struct {
	learningRepo repositories.LearningRepository
	userRepo     repositories.UserRepository
}

func NewLearningService(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *learningService {
	return &learningService{learningRepo: learningRepo, userRepo: userRepo}
}

func (s learningService) GetVideoLecture(id int) (*response.VideoLectureResponse, error) {
	contentDB, err := s.learningRepo.GetContent(id)
	if err != nil || contentDB == nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrContentNotFound
	}

	videoLink, err := s.learningRepo.GetVideoFileLink(contentDB.VideoPath)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrServiceUnavailableError
	}

	response := response.VideoLectureResponse{
		ContentID:   contentDB.ID,
		ContentName: contentDB.Name,
		VideoLink:   videoLink,
	}

	return &response, nil
}

func (s learningService) GetOverview(userID int) (*response.ContentOverviewResponse, error) {
	loader := loaders.NewLearningOverviewLoader(s.learningRepo, s.userRepo)

	err := loader.Load(userID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	overview := loader.GetOverview()
	preExamID := loader.GetPreExamID()
	learningProgression := loader.GetLearningProgression()
	lastedGroup, contentGroup := overview.GetLearningOverview(learningProgression)

	response := response.ContentOverviewResponse{
		PreExam:              preExamID,
		LastedGroup:          lastedGroup,
		ContentGroupOverview: contentGroup,
	}

	return &response, nil
}

func (s learningService) GetActivity(userID int, activityID int) (*response.ActivityResponse, error) {
	loader := loaders.NewActivityLoader(s.learningRepo, s.userRepo)

	err := loader.Load(userID, activityID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	_activity := loader.GetActivity()
	activityHints := loader.GetActivityHints()
	userHints := loader.GetUserHints()

	if *_activity == (activity.Activity{}) {
		return nil, errs.ErrActivitiesNotFound
	}

	choices, err := s.learningRepo.GetActivityChoices(_activity.ID, _activity.TypeID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrActivitiesNotFound
	}

	response := response.ActivityResponse{
		Activity: *_activity,
		Choices:  choices.CreatePropositionChoices(),
		Hint: &activity.ActivityHint{
			TotalHint:   len(activityHints),
			UsedHints:   activityHints.GetUsedHints(userHints),
			HintRoadMap: activityHints.CreateRoadmap(),
		},
	}

	return &response, nil
}

func (s learningService) UseHint(userID int, activityID int) (*response.UsedHintResponse, error) {
	loader := loaders.NewHintLoader(s.learningRepo, s.userRepo)

	err := loader.Load(userID, activityID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	activityHints := loader.GetActivityHintsDB()
	userHints := loader.GetUserHintsDB()
	user := loader.GetUser()

	if len(activityHints) == 0 {
		return nil, errs.ErrLoadError
	}

	nextLevelHint := activityHints.GetNextLevelHint(userHints)
	if nextLevelHint == nil {
		return nil, errs.ErrHintAlreadyUsed
	}

	if user.Point < nextLevelHint.PointReduce {
		return nil, errs.ErrHintPointsNotEnough
	}

	err = s.learningRepo.UseHint(userID, nextLevelHint.PointReduce, nextLevelHint.ID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrInsertError
	}

	response := response.UsedHintResponse{
		HintDB: *nextLevelHint,
	}

	return &response, nil
}

func (s learningService) GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error) {
	loader := loaders.NewContentRoadmapLoader(s.learningRepo, s.userRepo)

	err := loader.Load(userID, contentID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrContentNotFound
	}

	content := loader.GetContent()
	contentActivity := loader.GetContentActivity()
	learningProgression := loader.GetLearningProgression()

	if content == nil {
		return nil, errs.ErrContentNotFound
	}

	roadmapItems := contentActivity.GetContentRoadmap(learningProgression)

	response := response.ContentRoadmapResponse{
		ContentID:   content.ID,
		ContentName: content.Name,
		Items:       roadmapItems,
	}

	return &response, nil
}

func (s learningService) finishActivityAnswer(
	progression *content.LearningProgression,
	activityID int,
	activityTypeID int,
	activityPoint int,
	isCorrect bool,
	userID int,
) (int, error) {
	hasProgression := *progression != (content.LearningProgression{})

	err := s.userRepo.InsertLearningProgression(userID, activityID, activityPoint, isCorrect, hasProgression)
	if err != nil {
		logs.GetInstance().Error(err)
		return 0, errs.ErrInsertError
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil || user == nil {
		logs.GetInstance().Error(err)
		return 0, errs.ErrUserNotFound
	}

	return user.Point, nil
}

func (s learningService) CheckAnswer(userID int, request request.CheckAnswerRequest) (*response.AnswerResponse, error) {

	loader := loaders.NewCheckAnswerLoader(s.learningRepo)

	err := loader.Load(*request.ActivityID, *request.ActivityTypeID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	_activity := loader.GetActivity()
	choices := loader.GetChoices()
	progression := loader.GetProgression()

	if _activity == nil {
		return nil, errs.ErrLoadError
	}

	if _activity.TypeID != *request.ActivityTypeID {
		return nil, errs.ErrActivityTypeInvalid
	}

	var isCorrect bool
	var errMessage *string
	var erChoiceAnswer activity.ERChoiceAnswer
	var choice activity.ERChoice

	if *request.ActivityTypeID == 6 {
		var ok bool
		choice, ok = choices.(activity.ERChoice)
		if !ok {
			return nil, errs.ErrAnswerInvalid
		}

		err := utils.StructToStruct(request.Answer, &erChoiceAnswer)

		if err != nil {
			logs.GetInstance().Error(err)
			return nil, errs.ErrInternalServerError
		}
		var message string
		isCorrect, message = erChoiceAnswer.IsCorrect(choice)

		errMessage = &message
	} else {
		formatedAnswer, err := activity.FormatAnswer(request.Answer, *request.ActivityTypeID)
		if err != nil {
			logs.GetInstance().Error(err)
			return nil, err
		}

		isCorrect, err = formatedAnswer.IsCorrect(choices)
		if err != nil {
			return nil, err
		}
	}

	updatedPoint, err := s.finishActivityAnswer(
		progression,
		*request.ActivityID,
		*request.ActivityTypeID,
		_activity.Point,
		isCorrect,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if *request.ActivityTypeID == activity.PEER_ACTIVITY_TYPE_ID && *request.ActivityID == activity.PEER_ACTIVITY_ID && choice.Type == activity.ER_CHOICE_DRAW {

		erAnswer := activity.ERAnswer{
			UserID:        userID,
			Tables:        erChoiceAnswer.Tables,
			Relationships: erChoiceAnswer.Relationships,
		}

		err = s.learningRepo.InsertERAnswer(erAnswer, userID)
		if err != nil && !utils.IsSqlDuplicateError(err) {
			logs.GetInstance().Error(err)
			return nil, errs.ErrInsertError
		}
	}

	response := response.AnswerResponse{
		ActivityID:   _activity.ID,
		IsCorrect:    isCorrect,
		UpdatedPoint: updatedPoint,
		ErrMessage:   errMessage,
	}

	return &response, nil
}

func (s learningService) GetRecommend(userID int) (*response.RecommendResponse, error) {
	loader := loaders.NewRecommendLoader(s.learningRepo, s.userRepo)

	err := loader.Load(userID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	contentGroups := loader.GetContentGroups()
	preTestResults := loader.GetPreTestResults()

	recommend := preTestResults.GetRecommend(contentGroups)

	response := response.RecommendResponse{
		Recommend: recommend,
	}

	return &response, err
}

func (s learningService) CheckPeerReview(userID int, request request.PeerReviewRequest) (*response.AnswerResponse, error) {
	loader := loaders.NewCheckPeerReviewLoader(s.learningRepo)

	err := loader.Load(*request.ERAnswerID, activity.PEER_ACTIVITY_ID)
	if err != nil {
		logs.GetInstance().Error(err)

		if err == errs.ErrNotFoundError {
			return nil, errs.ErrNotFoundError
		} else {
			return nil, errs.ErrLoadError
		}
	}

	_erAnswer := loader.GetERAnswer()
	_erChoice := loader.GetERChoice()
	_progression := loader.GetProgression()

	suggestionsList := _erChoice.GetSuggestionsList(activity.ERChoiceAnswer{
		Tables:        _erAnswer.Tables,
		Relationships: _erAnswer.Relationships,
	})

	correct := suggestionsList.Compare(request.Reviews)

	updatedPoint, err := s.finishActivityAnswer(
		_progression,
		activity.PEER_ACTIVITY_ID,
		activity.PEER_ACTIVITY_TYPE_ID,
		activity.PEER_ACTIVITY_POINT,
		correct,
		userID,
	)

	if err != nil {
		return nil, err
	}

	res := response.AnswerResponse{
		ActivityID:   activity.PEER_ACTIVITY_ID,
		IsCorrect:    correct,
		UpdatedPoint: updatedPoint,
	}

	return &res, nil
}
