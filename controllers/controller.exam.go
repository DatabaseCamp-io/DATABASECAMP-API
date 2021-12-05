package controllers

// controller.exam.go
/**
 * 	This file is a part of controllers, used to do business logic of exam
 */

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

/**
 * This class do business logic of exam
 */
type examController struct {
	ExamRepo repositories.IExamRepository // repository for load exam data
	UserRepo repositories.IUserRepository // repository for load user data
}

/**
 * Constructor creates a new examController instance
 *
 * @param   examRepo    	Exam Repository for load exam data
 * @param   userRepo        User Repository for load user data
 *
 * @return 	instance of examController
 */
func NewExamController(examRepo repositories.IExamRepository, userRepo repositories.IUserRepository) examController {
	return examController{ExamRepo: examRepo, UserRepo: userRepo}
}

/**
 * Get the exam to use for the test
 *
 * @param 	examID 		   Exam ID for getting activities of the exam
 * @param 	userID 		   User ID for getting detail of user about the exam
 *
 * @return response of the exam
 * @return error of getting the exam
 */
func (c examController) GetExam(examID int, userID int) (*response.ExamResponse, error) {

	// Create exam loader instance
	loader := loaders.NewExamLoader(c.ExamRepo, c.UserRepo)

	// Load exam data and check the error
	err := loader.Load(userID, examID)
	if err != nil || len(loader.GetExamActivitiesDB()) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	// Create and prepare exam
	exam := entities.Exam{}
	exam.Prepare(loader.GetExamActivitiesDB())

	// User need to have all badges before they can take the final exam
	if exam.GetInfo().Type == string(entities.ExamType.Posttest) && !c.canDoFianlExam(loader.GetCorrectedBadgeDB()) {
		return nil, errs.ErrFinalExamBadgesNotEnough
	}

	// Create exam response
	response := response.NewExamResponse(exam)

	return response, nil
}

/**
 * Get overview of the exam
 *
 * @param 	userID 		   User ID for getting detail of user about the exam overview
 *
 * @return response of the exam overview
 * @return error of getting exam overview
 */
func (c examController) GetOverview(userID int) (*response.ExamOverviewResponse, error) {

	// Create exam overview loader instance
	loader := loaders.NewExamOverviewLoader(c.ExamRepo, c.UserRepo)

	// Load exam overview and check the error
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	// Get data from loader
	examResultDB := loader.GetExamResultsDB()
	examDB := loader.GetExamDB()
	correctedBadgeDB := loader.GetCorrectedBadgeDB()

	// Check that user can do final exam
	canDo := c.canDoFianlExam(correctedBadgeDB)

	// Create exam overview response
	response := response.NewExamOverviewResponse(examResultDB, examDB, canDo)

	return response, nil
}

/**
 * Check user can take the final exam
 *
 * @param 	correctedBadgesDB 		all badges that user corrected from the database
 *
 * @return true if user can take the final exam, false otherwise
 */
func (c examController) canDoFianlExam(correctedBadgesDB []storages.CorrectedBadgeDB) bool {
	for _, correctedBadgeDB := range correctedBadgesDB {
		if correctedBadgeDB.UserID == nil && correctedBadgeDB.BadgeID != 3 {
			return false
		}
	}
	return true
}

/**
 * Check answer of the exam
 *
 * @param 	userID 		   	User ID for record user exam
 * @param 	request 		Exam answer request for check answer of the exam
 *
 * @return response of the exam result overview
 * @return error of checking exam
 */
func (c examController) CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error) {

	// Get activities of the exam from the repository
	examActivities, err := c.ExamRepo.GetExamActivity(*request.ExamID)

	// Check error and empty exam
	if err != nil || len(examActivities) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	// Create and prepare exam solution
	exam := entities.Exam{}
	exam.Prepare(examActivities)

	// The number of activities of the solution and activities of the user's answer should be equal
	if len(request.Activities) != len(exam.GetActivities()) {
		return nil, errs.ErrActivitiesNumberIncorrect
	}

	// Check exam answer
	_, err = exam.CheckAnswer(request.Activities)
	if err != nil {
		return nil, err
	}

	// Create user badge
	userBadgeDB := storages.UserBadgeDB{
		UserID:  userID,
		BadgeID: exam.GetInfo().BadgeID,
	}

	// Get exam result from exam
	examResultDB := exam.ToExamResultDB(userID)

	// Get exam result activities from exam
	examResultActivities := exam.ToExamResultActivitiesDB()

	// Save exam result and check the error
	err = c.saveExamResult(exam.GetInfo().Type, userBadgeDB, examResultDB, examResultActivities)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrInsertError
	}

	// Set exam result ID
	exam.SetResultID(examResultDB.ID)

	// Create exam result response
	response := response.NewExamResultOverviewResponse(exam)

	return response, nil
}

/**
 * Add exam result ID to exam result activity
 *
 * @param 	examResultID 		Exam result ID to add
 * @param 	examResultActivity 	Exam result activities for add exam result ID
 *
 * @return  exam result activity that added exam result id
 */
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

/**
 * Save exam result by database transaction
 *
 * @param 	examType 				Type of the exam
 * @param 	userBadgeDB 			Badge of the exam to be given to the user
 * @param 	examResultDB 			Exam result to save into the database
 * @param 	resultActivitiesDB 		Exam result activities to save into the database
 *
 * @return  the error of saving data
 */
func (c examController) saveExamResult(examType string, userBadgeDB storages.UserBadgeDB, examResultDB *storages.ExamResultDB, resultActivitiesDB []storages.ExamResultActivityDB) error {

	var err error

	// Create database transaction
	tx := database.NewTransaction()

	// Begin transaction
	tx.Begin()

	// Insert Exam result
	*examResultDB, err = c.ExamRepo.InsertExamResultTransaction(tx, *examResultDB)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert Exam result Activity
	_, err = c.ExamRepo.InsertExamResultActivityTransaction(tx, c.addExamResultID(examResultDB.ID, resultActivitiesDB))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Give exam badge if user passed the exam
	if (examType == entities.ExamType.MiniExam || examType == entities.ExamType.Posttest) && examResultDB.IsPassed {
		_, err = c.UserRepo.InsertUserBadgeTransaction(tx, userBadgeDB)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	tx.Commit()

	// Close transaction
	tx.Close()

	return nil
}

/**
 * Get exam result of the user
 *
 * @param 	userID 		   	User ID for getting user data
 * @param 	examResultID 	Exam result ID for getting exam results of the user
 *
 * @return response of the exam result overview
 * @return error of getting exam result
 */
func (c examController) GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error) {

	// Get exam result of the user by exam result id
	examResults, err := c.UserRepo.GetExamResultByID(userID, examResultID)
	if err != nil || len(examResults) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	// Get exam result activities of the exam
	examActivities, err := c.ExamRepo.GetExamActivity(examResults[0].ExamID)
	if err != nil || len(examActivities) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrExamNotFound
	}

	// Create and prepare exam
	exam := entities.Exam{}
	exam.Prepare(examActivities)
	exam.PrepareResult(examResults[0])

	// Create exam overview response
	response := response.NewExamResultOverviewResponse(exam)

	return response, nil
}
