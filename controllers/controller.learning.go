package controllers

import (
	"DatabaseCamp/controllers/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/response"
	"DatabaseCamp/repositories"
	"DatabaseCamp/services"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

type learningController struct {
	learningRepo repositories.ILearningRepository
	userRepo     repositories.IUserRepository
	service      services.IAwsService
}

type ILearningController interface {
	GetVideoLecture(id int) (*response.VideoLectureResponse, error)
	GetOverview(userID int) (*response.ContentOverviewResponse, error)
	GetActivity(userID int, activityID int) (*response.ActivityResponse, error)
	UseHint(userID int, activityID int) (*general.HintDB, error)
	GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error)
	CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*response.AnswerResponse, error)
}

func NewLearningController(
	learningRepo repositories.ILearningRepository,
	userRepo repositories.IUserRepository,
	service services.IAwsService,
) learningController {
	return learningController{
		learningRepo: learningRepo,
		userRepo:     userRepo,
		service:      service,
	}
}

func (c learningController) GetVideoLecture(id int) (*response.VideoLectureResponse, error) {
	contentDB, err := c.learningRepo.GetContent(id)
	if err != nil || contentDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrContentNotFound
	}

	videoLink, err := c.service.GetFileLink(contentDB.VideoPath)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrServiceUnavailableError
	}

	res := response.VideoLectureResponse{
		ContentID:   contentDB.ID,
		ContentName: contentDB.Name,
		VideoLink:   videoLink,
	}

	return &res, nil
}

func (c learningController) GetOverview(userID int) (*response.ContentOverviewResponse, error) {
	loader := loaders.NewLearningOverviewLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}
	response := response.NewContentOverviewResponse(loader.OverviewDB, loader.LearningProgressionDB)
	return response, nil
}

func (c learningController) GetActivity(userID int, activityID int) (*response.ActivityResponse, error) {
	loader := loaders.NewActivityLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID, activityID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	choiceDB, err := c.getChoices(loader.ActivityDB.ID, loader.ActivityDB.TypeID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrActivitiesNotFound
	}

	activity := entities.Activity{}
	activity.SetActivity(*loader.ActivityDB)
	activity.SetChoicesByChoiceDB(choiceDB)
	activity.SetHint(loader.ActivityHintsDB, loader.UserHintsDB)

	response := response.NewActivityResponse(activity)
	return response, nil
}

func (c learningController) getChoices(activityID int, typeID int) (interface{}, error) {
	if typeID == 1 {
		return c.learningRepo.GetMatchingChoice(activityID)
	} else if typeID == 2 {
		return c.learningRepo.GetMultipleChoice(activityID)
	} else if typeID == 3 {
		return c.learningRepo.GetCompletionChoice(activityID)
	} else {
		return nil, errs.ErrActivityTypeInvalid
	}
}

func (c learningController) finishActivityTrasaction(userID int, activityID int, addPoint int) error {
	tx := database.NewTransaction()
	tx.Begin()

	progression := general.LearningProgressionDB{
		UserID:           userID,
		ActivityID:       activityID,
		CreatedTimestamp: time.Now().Local(),
	}

	_, err := c.userRepo.InsertLearningProgressionTransaction(tx, progression)
	if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
		tx.Rollback()
		return err
	}

	if !utils.NewHelper().IsSqlDuplicateError(err) {
		err = c.userRepo.ChangePointTransaction(tx, userID, addPoint, entities.Mode.Add)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	tx.Close()

	return nil
}

func (c learningController) UseHint(userID int, activityID int) (*general.HintDB, error) {
	loader := loaders.NewHintLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID, activityID)
	if err != nil || len(loader.ActivityHintsDB) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	nextLevelHint := c.getNextLevelHint(loader.ActivityHintsDB, loader.UserHintsDB)
	if nextLevelHint == nil {
		return nil, errs.ErrHintAlreadyUsed
	}

	if loader.UserDB.Point < nextLevelHint.PointReduce {
		return nil, errs.ErrHintPointsNotEnough
	}

	err = c.useHintTransaction(userID, nextLevelHint.PointReduce, nextLevelHint.ID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrInsertError
	}

	return nextLevelHint, nil
}

func (c learningController) isUsedHint(userHints []general.UserHintDB, hintID int) bool {
	for _, userHint := range userHints {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

func (c learningController) getNextLevelHint(ActivityHintsDB []general.HintDB, userHintsDB []general.UserHintDB) *general.HintDB {
	for _, activityHint := range ActivityHintsDB {
		if !c.isUsedHint(userHintsDB, activityHint.ID) {
			return &activityHint
		}
	}
	return nil
}

func (c learningController) useHintTransaction(userID int, reducePoint int, hintID int) error {
	var wg sync.WaitGroup
	var err error
	tx := database.NewTransaction()

	ct := general.ConcurrentTransaction{
		Concurrent: &general.Concurrent{
			Wg:  &wg,
			Err: &err,
		},
		Transaction: tx,
	}

	wg.Add(2)
	go c.updateUserPointAsyncTrasaction(&ct, userID, reducePoint, entities.Mode.Reduce)
	go c.insertUserHintAsyncTransaction(&ct, userID, hintID)
	wg.Wait()

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	tx.Close()

	return nil
}

func (c learningController) updateUserPointAsyncTrasaction(ct *general.ConcurrentTransaction, userID int, updatePoint int, mode entities.ChangePointMode) {
	defer ct.Concurrent.Wg.Done()
	err := c.userRepo.ChangePointTransaction(ct.Transaction, userID, updatePoint, mode)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

func (c learningController) insertUserHintAsyncTransaction(ct *general.ConcurrentTransaction, userID int, hintID int) {
	defer ct.Concurrent.Wg.Done()
	hint := general.UserHintDB{
		UserID:           userID,
		HintID:           hintID,
		CreatedTimestamp: time.Now().Local(),
	}
	_, err := c.userRepo.InsertUserHintTransaction(ct.Transaction, hint)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

func (c learningController) GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error) {
	loader := loaders.NewContentRoadmapLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID, contentID)
	if err != nil || loader.ContentDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrContentNotFound
	}
	response := response.NewContentRoadmapResponse(*loader.ContentDB, loader.ContentActivityDB, loader.LearningProgressionDB)
	return response, nil
}

func (c learningController) CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*response.AnswerResponse, error) {
	loader := loaders.NewCheckAnswerLoader(c.learningRepo)
	err := loader.Load(activityID, typeID, c.getChoices)
	if err != nil || loader.ActivityDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	if loader.ActivityDB.TypeID != typeID {
		return nil, errs.ErrActivityTypeInvalid
	}

	activity := entities.Activity{}
	activity.SetActivity(*loader.ActivityDB)
	activity.SetChoicesByChoiceDB(loader.ChoicesDB)

	isCorrect, err := activity.IsAnswerCorrect(answer)
	if err != nil {
		return nil, err
	}

	if isCorrect {
		err = c.finishActivityTrasaction(userID, activityID, activity.GetInfo().Point)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			logs.New().Error(err)
			return nil, errs.ErrInsertError
		}
	}

	userDB, err := c.userRepo.GetUserByID(userID)
	if err != nil || userDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	response := response.NewActivityAnswerResponse(activity, userDB.Point, isCorrect)
	return response, nil
}
