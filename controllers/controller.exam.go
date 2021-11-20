package controllers

import (
	"DatabaseCamp/controllers/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/response"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
)

type examController struct {
	examRepo repositories.IExamRepository
	userRepo repositories.IUserRepository
}

type IExamController interface {
	GetExam(examID int, userID int) (*response.ExamResponse, error)
	GetOverview(userID int) (*response.ExamOverviewResponse, error)
	CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error)
	GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error)
}

func NewExamController(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) examController {
	return examController{examRepo: examRepo, userRepo: userRepo}
}

func (c examController) GetExam(examID int, userID int) (*response.ExamResponse, error) {
	loader := loaders.NewExamLoader(c.examRepo, c.userRepo)
	err := loader.Load(userID, examID)
	if err != nil || len(loader.ExamActivitiesDB) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	exam := entities.Exam{}
	exam.Prepare(loader.ExamActivitiesDB)

	if exam.GetInfo().Type == string(entities.ExamType.Posttest) && !c.canDoFianlExam(loader.CorrectedBadgeDB) {
		return nil, errs.ErrFinalExamBadgesNotEnough
	}

	response := response.NewExamResponse(exam)
	return response, nil
}

func (c examController) GetOverview(userID int) (*response.ExamOverviewResponse, error) {
	loader := loaders.NewExamOverviewLoader(c.examRepo, c.userRepo)
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	response := response.NewExamOverviewResponse(loader.ExamResultsDB, loader.ExamDB, c.canDoFianlExam(loader.CorrectedBadgeDB))
	return response, nil
}

func (c examController) canDoFianlExam(correctedBadgesDB []general.CorrectedBadgeDB) bool {
	for _, correctedBadgeDB := range correctedBadgesDB {
		if correctedBadgeDB.UserID == nil && correctedBadgeDB.BadgeID != 3 {
			return false
		}
	}
	return true
}

func (c examController) CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error) {
	examActivities, err := c.examRepo.GetExamActivity(*request.ExamID)
	if err != nil || len(examActivities) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	exam := entities.Exam{}
	exam.Prepare(examActivities)
	if len(request.Activities) != len(exam.GetActivities()) {
		return nil, errs.ErrActivitiesNumberIncorrect
	}

	_, err = exam.CheckAnswer(request.Activities)
	if err != nil {
		return nil, err
	}

	userBadgeDB := general.UserBadgeDB{
		UserID:  userID,
		BadgeID: exam.GetInfo().BadgeID,
	}
	examResultDB := exam.ToExamResultDB(userID)
	examResultActivities := exam.ToExamResultActivitiesDB()

	err = c.saveExamResult(exam.GetInfo().Type, userBadgeDB, examResultDB, examResultActivities)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrInsertError
	}

	exam.GetResult().ExamResultID = examResultDB.ID
	response := response.NewExamResultOverviewResponse(exam)
	return response, nil
}

func (c examController) addExamResultID(examResultID int, examResultActivity []general.ExamResultActivityDB) []general.ExamResultActivityDB {
	newExamResultActivity := make([]general.ExamResultActivityDB, 0)
	for _, v := range examResultActivity {
		newExamResultActivity = append(newExamResultActivity, general.ExamResultActivityDB{
			ExamResultID: examResultID,
			ActivityID:   v.ActivityID,
			Score:        v.Score,
		})
	}
	return newExamResultActivity
}

func (c examController) saveExamResult(examType string, userBadgeDB general.UserBadgeDB, examResultDB *general.ExamResultDB, resultActivitiesDB []general.ExamResultActivityDB) error {
	var err error
	tx := database.NewTransaction()
	tx.Begin()

	*examResultDB, err = c.examRepo.InsertExamResultTransaction(tx, *examResultDB)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = c.examRepo.InsertExamResultActivityTransaction(tx, c.addExamResultID(examResultDB.ID, resultActivitiesDB))
	if err != nil {
		tx.Rollback()
		return err
	}

	if examType == entities.ExamType.MiniExam || examType == entities.ExamType.Posttest && examResultDB.IsPassed {
		_, err = c.userRepo.InsertUserBadgeTransaction(tx, userBadgeDB)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	tx.Close()
	return nil
}

func (c examController) GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error) {

	examResults, err := c.userRepo.GetExamResultByID(userID, examResultID)
	if err != nil || len(examResults) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	examActivities, err := c.examRepo.GetExamActivity(examResults[0].ExamID)
	if err != nil || len(examActivities) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	exam := entities.Exam{}
	exam.Prepare(examActivities)
	exam.PrepareResult(examResults[0])
	response := response.NewExamResultOverviewResponse(exam)
	return response, nil
}
