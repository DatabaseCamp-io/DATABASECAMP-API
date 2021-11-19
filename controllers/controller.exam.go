package controllers

import (
	loader "DatabaseCamp/controllers/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
)

type examController struct {
	examRepo repositories.IExamRepository
	userRepo repositories.IUserRepository
}

type IExamController interface {
	GetExam(examID int, userID int) (*models.ExamResponse, error)
	GetOverview(userID int) (*models.ExamOverviewResponse, error)
	CheckExam(userID int, request models.ExamAnswerRequest) (*models.ExamResultOverviewResponse, error)
	GetExamResult(userID int, examResultID int) (*models.ExamResultOverviewResponse, error)
}

func NewExamController(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) examController {
	return examController{examRepo: examRepo, userRepo: userRepo}
}

func (c examController) GetExam(examID int, userID int) (*models.ExamResponse, error) {
	examActivities, err := c.examRepo.GetExamActivity(examID)
	if err != nil || len(examActivities) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	exam := models.NewExam()
	exam.Prepare(examActivities)

	response := exam.ToResponse()
	return response, nil
}

func (c examController) GetOverview(userID int) (*models.ExamOverviewResponse, error) {
	loader := loader.NewExamOverviewLoader(c.examRepo, c.userRepo)
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	examOverview := models.NewExamOverview()
	examOverview.PrepareExamOverview(loader.ExamResultsDB, loader.CorrectedBadgeDB, loader.ExamDB)

	response := examOverview.ToResponse()
	return response, nil
}

func (c examController) CheckExam(userID int, request models.ExamAnswerRequest) (*models.ExamResultOverviewResponse, error) {
	examActivities, err := c.examRepo.GetExamActivity(*request.ExamID)
	if err != nil || len(examActivities) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	exam := models.NewExam()
	exam.Prepare(examActivities)
	if len(request.Activities) != len(exam.Activities) {
		return nil, errs.ErrActivitiesNumberIncorrect
	}

	_, err = exam.CheckAnswer(request.Activities)
	if err != nil {
		return nil, err
	}

	userBadgeDB := models.UserBadgeDB{
		UserID:  userID,
		BadgeID: exam.Info.BadgeID,
	}
	examResultDB := exam.ToExamResultDB(userID)
	examResultActivities := exam.ToExamResultActivitiesDB()

	err = c.saveExamResult(exam.Info.Type, userBadgeDB, examResultDB, examResultActivities)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrInsertError
	}

	exam.Result.ExamResultID = examResultDB.ID
	response := exam.ToExamResultOverviewResponse()
	return response, nil
}

func (c examController) addExamResultID(examResultID int, examResultActivity []models.ExamResultActivityDB) []models.ExamResultActivityDB {
	newExamResultActivity := make([]models.ExamResultActivityDB, 0)
	for _, v := range examResultActivity {
		newExamResultActivity = append(newExamResultActivity, models.ExamResultActivityDB{
			ExamResultID: examResultID,
			ActivityID:   v.ActivityID,
			Score:        v.Score,
		})
	}
	return newExamResultActivity
}

func (c examController) saveExamResult(examType string, userBadgeDB models.UserBadgeDB, examResultDB *models.ExamResultDB, resultActivitiesDB []models.ExamResultActivityDB) error {
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

	if examType == string(models.Exam.MiniExam) && examResultDB.IsPassed {
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

func (c examController) GetExamResult(userID int, examResultID int) (*models.ExamResultOverviewResponse, error) {

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

	exam := models.NewExam()
	exam.Prepare(examActivities)
	exam.PrepareResult(examResults[0])
	response := exam.ToExamResultOverviewResponse()

	return response, nil
}
