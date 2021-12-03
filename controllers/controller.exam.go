package controllers

import (
	"DatabaseCamp/controllers/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/response"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
)

type examController struct {
	ExamRepo repositories.IExamRepository
	UserRepo repositories.IUserRepository
}

// Interface that show how others function call and use function in module exam controller
type IExamController interface {
	GetExam(examID int, userID int) (*response.ExamResponse, error)
	GetOverview(userID int) (*response.ExamOverviewResponse, error)
	CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error)
	GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error)
}

// Create exam controller instance
func NewExamController(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) examController {
	return examController{ExamRepo: examRepo, UserRepo: userRepo}
}

func (c examController) GetExam(examID int, userID int) (*response.ExamResponse, error) {
	loader := loaders.NewExamLoader(c.ExamRepo, c.UserRepo)
	err := loader.Load(userID, examID)
	if err != nil || len(loader.GetExamActivitiesDB()) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	exam := entities.Exam{}
	exam.Prepare(loader.GetExamActivitiesDB())

	if exam.GetInfo().Type == string(entities.ExamType.Posttest) && !c.canDoFianlExam(loader.GetCorrectedBadgeDB()) {
		return nil, errs.ErrFinalExamBadgesNotEnough
	}

	response := response.NewExamResponse(exam)
	return response, nil
}

func (c examController) GetOverview(userID int) (*response.ExamOverviewResponse, error) {
	loader := loaders.NewExamOverviewLoader(c.ExamRepo, c.UserRepo)
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	response := response.NewExamOverviewResponse(loader.GetExamResultsDB(), loader.GetExamDB(), c.canDoFianlExam(loader.GetCorrectedBadgeDB()))
	return response, nil
}

func (c examController) canDoFianlExam(correctedBadgesDB []storages.CorrectedBadgeDB) bool {
	for _, correctedBadgeDB := range correctedBadgesDB {
		if correctedBadgeDB.UserID == nil && correctedBadgeDB.BadgeID != 3 {
			return false
		}
	}
	return true
}

func (c examController) CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error) {
	examActivities, err := c.ExamRepo.GetExamActivity(*request.ExamID)
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

	userBadgeDB := storages.UserBadgeDB{
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

	exam.SetResultID(examResultDB.ID)
	response := response.NewExamResultOverviewResponse(exam)
	return response, nil
}

func (c examController) addExamResultID(examResultID int, examResultActivity []storages.ExamResultActivityDB) []storages.ExamResultActivityDB {
	newExamResultActivity := make([]storages.ExamResultActivityDB, 0)
	for _, v := range examResultActivity {
		newExamResultActivity = append(newExamResultActivity, storages.ExamResultActivityDB{
			ExamResultID: examResultID,
			ActivityID:   v.ActivityID,
			Score:        v.Score,
		})
	}
	return newExamResultActivity
}

func (c examController) saveExamResult(examType string, userBadgeDB storages.UserBadgeDB, examResultDB *storages.ExamResultDB, resultActivitiesDB []storages.ExamResultActivityDB) error {
	var err error
	tx := database.NewTransaction()
	tx.Begin()

	*examResultDB, err = c.ExamRepo.InsertExamResultTransaction(tx, *examResultDB)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = c.ExamRepo.InsertExamResultActivityTransaction(tx, c.addExamResultID(examResultDB.ID, resultActivitiesDB))
	if err != nil {
		tx.Rollback()
		return err
	}

	if (examType == entities.ExamType.MiniExam || examType == entities.ExamType.Posttest) && examResultDB.IsPassed {
		_, err = c.UserRepo.InsertUserBadgeTransaction(tx, userBadgeDB)
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

	examResults, err := c.UserRepo.GetExamResultByID(userID, examResultID)
	if err != nil || len(examResults) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	examActivities, err := c.ExamRepo.GetExamActivity(examResults[0].ExamID)
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
