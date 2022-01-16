package services

import (
	"database-camp/internal/errs"
	"database-camp/internal/logs"
	"database-camp/internal/models/entities/badge"
	"database-camp/internal/models/entities/exam"
	"database-camp/internal/models/request"
	"database-camp/internal/models/response"
	"database-camp/internal/repositories"
	"database-camp/internal/services/loaders"
)

type ExamService interface {
	GetExam(examID int, userID int) (*response.ExamResponse, error)
	GetOverview(userID int) (*response.ExamOverviewResponse, error)
	CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error)
	GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error)
}

type examService struct {
	examRepo     repositories.ExamRepository
	userRepo     repositories.UserRepository
	learningRepo repositories.LearningRepository
}

func NewExamService(
	examRepo repositories.ExamRepository,
	userRepo repositories.UserRepository,
	learningRepo repositories.LearningRepository,
) *examService {
	return &examService{
		examRepo:     examRepo,
		userRepo:     userRepo,
		learningRepo: learningRepo,
	}
}

func (s examService) GetExam(examID int, userID int) (*response.ExamResponse, error) {
	examLoader := loaders.NewExamLoader(s.examRepo, s.userRepo)

	err := examLoader.Load(userID, examID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrExamNotFound
	}

	examActivities := examLoader.GetExamActivities()
	correctedBadges := examLoader.GetCorrectedBadge()
	_exam := examLoader.GetExam()

	if len(examActivities) == 0 || *_exam == (exam.Exam{}) {
		return nil, errs.ErrExamNotFound
	}

	if _exam.Type == string(exam.POST) && !correctedBadges.CanDoFianlExam() {
		return nil, errs.ErrFinalExamBadgesNotEnough
	}

	activitiesExamLoader := loaders.NewActivityExamLoader(s.learningRepo)

	err = activitiesExamLoader.Load(examActivities)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrExamNotFound
	}

	activities := activitiesExamLoader.GetActivities()

	activitiesResponse := make([]response.ActivityResponse, 0)

	for _, activity := range activities {
		choices := activity.Choices.CreatePropositionChoices()
		activitiesResponse = append(activitiesResponse, response.ActivityResponse{
			Activity: activity.Activity,
			Choices:  choices,
			Hint:     nil,
		})
	}

	response := response.ExamResponse{
		Exam:       *_exam,
		Activities: activitiesResponse,
	}

	return &response, nil
}

func (s examService) GetOverview(userID int) (*response.ExamOverviewResponse, error) {
	loader := loaders.NewExamOverviewLoader(s.examRepo, s.userRepo)

	err := loader.Load(userID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	examResults := loader.GetExamResults()
	exams := loader.GetExams()
	correctedBadges := loader.GetCorrectedBadges()

	canDo := correctedBadges.CanDoFianlExam()

	response := response.ExamOverviewResponse{
		PreExam:   exams.GetPreExam(examResults),
		MiniExam:  exams.GetMiniExam(examResults),
		FinalExam: exams.GetFinalExam(examResults, canDo),
	}

	return &response, nil
}

func (s examService) CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error) {
	checkExamLoader := loaders.NewCheckExamLoader(s.examRepo)

	err := checkExamLoader.Load(*request.ExamID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrExamNotFound
	}

	examActivities := checkExamLoader.GetExamActivities()
	exam := checkExamLoader.GetExam()

	if len(examActivities) == 0 {
		return nil, errs.ErrExamNotFound
	}

	activitiesExamLoader := loaders.NewActivityExamLoader(s.learningRepo)

	err = activitiesExamLoader.Load(examActivities)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrExamNotFound
	}

	activities := activitiesExamLoader.GetActivities()

	if len(request.Activities) != len(activities) {
		return nil, errs.ErrActivitiesNumberIncorrect
	}

	result, err := activities.CheckAnswers(*request.ExamID, userID, request.Activities)
	if err != nil {
		return nil, err
	}

	result, err = s.examRepo.SaveResult(*result)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.InsertBadge(badge.UserBadge{
		UserID:  userID,
		BadgeID: exam.BadgeID,
	})
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrInsertError
	}

	response := response.ExamResultOverviewResponse{
		ExamID:           exam.ID,
		ExamResultID:     result.ExamResult.ID,
		ExamType:         exam.Type,
		ContentGroupName: exam.ContentGroupName,
		CreatedTimestamp: result.ExamResult.CreatedTimestamp,
		Score:            result.ExamResult.Score,
		IsPassed:         result.ExamResult.IsPassed,
		ActivitiesResult: result.ActivitiesResult,
	}

	return &response, nil
}

func (s examService) GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error) {
	loader := loaders.NewExamResultLoader(s.examRepo)

	err := loader.Load(userID, examResultID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrExamNotFound
	}

	examResult := loader.GetExamResult()
	resultActivities := loader.GetResultActivities()

	response := response.ExamResultOverviewResponse{
		ExamResultID:     examResult.ID,
		ExamID:           examResult.ExamID,
		ExamType:         examResult.ExamType,
		ContentGroupName: examResult.ExamType,
		CreatedTimestamp: examResult.CreatedTimestamp,
		Score:            examResult.Score,
		IsPassed:         examResult.IsPassed,
		ActivitiesResult: resultActivities,
	}

	return &response, nil
}
