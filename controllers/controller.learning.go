package controllers

// controller.learning.go
/**
 * 	This file is a part of controllers, used to do business logic of learning
 */

import (
	"DatabaseCamp/controllers/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/entities/activity"
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/response"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

/**
 * This class do business logic of learning
 */
type learningController struct {
	LearningRepo repositories.ILearningRepository // repository for load learning data
	UserRepo     repositories.IUserRepository     // repository for load user data
}

/**
 * Constructor creates a new learningController instance
 *
 * @param   learningRepo    	Learning Repository for load learning data
 * @param   userRepo        	User Repository for load user data
 *
 * @return 	instance of learningController
 */
func NewLearningController(
	learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) learningController {
	return learningController{LearningRepo: learningRepo, UserRepo: userRepo}
}

/**
 * Get video lecture of the content
 *
 * @param 	id 		   	Content ID for getting video lecture of the content
 *
 * @return response of the video lecture
 * @return error of getting video lecture
 */
func (c learningController) GetVideoLecture(id int) (*response.VideoLectureResponse, error) {

	// Get content from the repository
	contentDB, err := c.LearningRepo.GetContent(id)
	if err != nil || contentDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrContentNotFound
	}

	// Get video link of the content from the repository
	videoLink, err := c.LearningRepo.GetVideoFileLink(contentDB.VideoPath)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrServiceUnavailableError
	}

	// Create video lecture response
	response := response.NewVideoLectureResponse(contentDB.ID, contentDB.Name, videoLink)

	return response, nil
}

/**
 * Get content overview of thew user
 *
 * @param 	userID 		   	User ID for getting content overview of the user
 *
 * @return response of the content overview
 * @return error of getting content overview
 */
func (c learningController) GetOverview(userID int) (*response.ContentOverviewResponse, error) {

	// Create learning overview loader
	loader := loaders.NewLearningOverviewLoader(c.LearningRepo, c.UserRepo)

	// Load content overview and check error
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	// Create content overview response
	response := response.NewContentOverviewResponse(loader.GetOverviewDB(), loader.GetLearningProgressionDB())

	return response, nil
}

/**
 * Get activity for user to do
 *
 * @param 	userID 		   		User ID for getting user hints
 * @param 	activityID 			Activity ID for getting activity data
 *
 * @return response of the activity
 * @return error of getting activity
 */
func (c learningController) GetActivity(userID int, activityID int) (*activity.Response, error) {

	// Create activity loader
	loader := loaders.NewActivityLoader(c.LearningRepo, c.UserRepo)

	// Load activity data and check error
	err := loader.Load(userID, activityID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	// Get activity data from loader
	activityDB := loader.GetActivityDB()

	userHints := loader.GetUserHintsDB()

	activityHintsDB := loader.GetActivityHintsDB()

	// Get choices of the activity from the database
	choiceDB, err := c.getChoices(activityDB.ID, activityDB.TypeID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrActivitiesNotFound
	}

	activity := activity.New(
		activityDB,
		userHints,
		activityHintsDB,
		choiceDB,
	)

	return activity.NewResponse()
}

/**
 * Get Choices of activity by activity type id
 *
 * @param 	activityID 			Activity ID for indicate activity
 * @param 	typeID 				Activity type ID for indicate type of the activity
 *
 * @return choices of the activity
 * @return error of getting choices
 */
func (c learningController) getChoices(activityID int, typeID int) (interface{}, error) {
	if typeID == 1 {
		return c.LearningRepo.GetMatchingChoice(activityID)
	} else if typeID == 2 {
		return c.LearningRepo.GetMultipleChoice(activityID)
	} else if typeID == 3 {
		return c.LearningRepo.GetCompletionChoice(activityID)
	} else {
		return nil, errs.ErrActivityTypeInvalid
	}
}

/**
 * Finish activity by database transaction
 *
 * @param 	userID 				User ID for adding activity point
 * @param 	activityID 			Activity ID for indicate activity
 * @param 	addPoint 			Point for add to the user
 *
 * @return error of finishing choices
 */
func (c learningController) finishActivityTrasaction(userID int, activityID int, addPoint int) error {

	// Create database transaction
	tx := database.NewTransaction()

	// Begin transaction
	tx.Begin()

	// Create progression
	progression := storages.LearningProgressionDB{
		UserID:           userID,
		ActivityID:       activityID,
		CreatedTimestamp: time.Now().Local(),
	}

	// Insert progression
	_, err := c.UserRepo.InsertLearningProgressionTransaction(tx, progression)
	if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
		tx.Rollback()
		return err
	}

	// Give a point to the user if do this activity for the first time
	if !utils.NewHelper().IsSqlDuplicateError(err) {
		err = c.UserRepo.ChangePointTransaction(tx, userID, addPoint, entities.Mode.Add)
	}

	// Check error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()

	// Close transaction
	tx.Close()

	return nil
}

/**
 * Use hint of the activity
 *
 * @param 	userID 		   		User ID for getting user hints
 * @param 	activityID 			Activity ID for getting all hints of the activity
 *
 * @return hint of the activity that user can use
 * @return error of using hint
 */
func (c learningController) UseHint(userID int, activityID int) (*storages.HintDB, error) {
	loader := loaders.NewHintLoader(c.LearningRepo, c.UserRepo)
	err := loader.Load(userID, activityID)
	if err != nil || len(loader.GetActivityHintsDB()) == 0 {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	nextLevelHint := c.getNextLevelHint(loader.GetActivityHintsDB(), loader.GetUserHintsDB())
	if nextLevelHint == nil {
		return nil, errs.ErrHintAlreadyUsed
	}

	if loader.GetUserDB().Point < nextLevelHint.PointReduce {
		return nil, errs.ErrHintPointsNotEnough
	}

	err = c.useHintTransaction(userID, nextLevelHint.PointReduce, nextLevelHint.ID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrInsertError
	}

	return nextLevelHint, nil
}

/**
 * Check for used hint of the user
 *
 * @param 	userHints 		User hints
 * @param 	hintID 			Hint ID for check to user hints
 *
 * @return true if that hint is used, false otherwise
 */
func (c learningController) isUsedHint(userHints []storages.UserHintDB, hintID int) bool {
	for _, userHint := range userHints {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

/**
 * Get the next level of hint that user can use
 *
 * @param 	userHints 			User hints
 * @param 	ActivityHintsDB 	Hints of the activity
 *
 * @return next level of hint that user can use
 */
func (c learningController) getNextLevelHint(ActivityHintsDB []storages.HintDB, userHintsDB []storages.UserHintDB) *storages.HintDB {
	for _, activityHint := range ActivityHintsDB {
		if !c.isUsedHint(userHintsDB, activityHint.ID) {
			return &activityHint
		}
	}
	return nil
}

/**
 * Use hint by database transaction
 *
 * @param 	userID 				User ID that used hint
 * @param 	reducePoint 		Point reduce when used hint
 * @param 	hintID 				Hint ID that user used
 *
 * @return error of using hint
 */
func (c learningController) useHintTransaction(userID int, reducePoint int, hintID int) error {

	var wg sync.WaitGroup
	var err error

	// Create transaction
	tx := database.NewTransaction()

	// Create concurrent
	ct := general.ConcurrentTransaction{
		Concurrent: &general.Concurrent{
			Wg:  &wg,
			Err: &err,
		},
		Transaction: tx,
	}

	// Add thread for concurrency
	wg.Add(2)

	// Doing 2 thread concurrency

	// Update user point
	go c.updateUserPointAsyncTrasaction(&ct, userID, reducePoint, entities.Mode.Reduce)

	// Insert user hint
	go c.insertUserHintAsyncTransaction(&ct, userID, hintID)

	// Waiting for all thread
	wg.Wait()

	// Check error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()

	// Close transaction
	tx.Close()

	return nil
}

/**
 * Update concurrency user point by database transaction
 *
 * @param 	ct 					Concurrent transaction model for do concurrent with database transaction
 * @param 	userID 				User ID for update user point
 * @param 	updatePoint 		Point to be update
 * @param 	mode 				mode for update user point
 */
func (c learningController) updateUserPointAsyncTrasaction(ct *general.ConcurrentTransaction, userID int, updatePoint int, mode entities.ChangePointMode) {
	defer ct.Concurrent.Wg.Done()
	err := c.UserRepo.ChangePointTransaction(ct.Transaction, userID, updatePoint, mode)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

/**
 * Insert concurrency user hint by database transaction
 *
 * @param 	ct 					Concurrent transaction model for do concurrent with database transaction
 * @param 	userID 				User ID for insert user hint
 * @param 	hintID 				Hint ID to insert
 */
func (c learningController) insertUserHintAsyncTransaction(ct *general.ConcurrentTransaction, userID int, hintID int) {
	defer ct.Concurrent.Wg.Done()
	hint := storages.UserHintDB{
		UserID:           userID,
		HintID:           hintID,
		CreatedTimestamp: time.Now().Local(),
	}
	_, err := c.UserRepo.InsertUserHintTransaction(ct.Transaction, hint)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

/**
 * Get roadmap of the content
 *
 * @param 	userID 		   		User ID for getting learning progression of the user
 * @param 	contentID 			Content ID for getting roadmap of the content
 *
 * @return response of the content roadmap
 * @return error of getting content roadmap
 */
func (c learningController) GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error) {

	// Create content roadmap loader
	loader := loaders.NewContentRoadmapLoader(c.LearningRepo, c.UserRepo)

	// Load content roadmap and check error
	err := loader.Load(userID, contentID)
	if err != nil || loader.GetContentDB() == nil {
		logs.New().Error(err)
		return nil, errs.ErrContentNotFound
	}

	// Get content roadmap from loader
	contentDB := loader.GetContentDB()
	contentActivity := loader.GetContentActivityDB()
	learningProgressionDB := loader.GetLearningProgressionDB()

	// Create content roadmap response
	response := response.NewContentRoadmapResponse(*contentDB, contentActivity, learningProgressionDB)

	return response, nil
}

/**
 * Check activity answer
 *
 * @param 	userID 		   		User ID for record user activity
 * @param 	activityID 			Activity ID for getting activity solution
 * @param 	typeID 				Activity ID for indicate type of the activity
 * @param 	answer 				Answer of the user
 *
 * @return response of the answer
 * @return error of checking activity answer
 */
func (c learningController) CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*response.AnswerResponse, error) {

	// Create check answer loader
	loader := loaders.NewCheckAnswerLoader(c.LearningRepo)

	// Load activity solution data and check error
	err := loader.Load(activityID, typeID, c.getChoices)
	if err != nil || loader.GetActivityDB() == nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	// type of the activity solution and type of the activity answer should be equal
	if loader.GetActivityDB().TypeID != typeID {
		return nil, errs.ErrActivityTypeInvalid
	}

	// Create activity and set data
	activity := entities.Activity{}
	activity.SetActivity(*loader.GetActivityDB())
	activity.SetChoicesByChoiceDB(loader.GetChoicesDB())

	// Check answer
	isCorrect, err := activity.IsAnswerCorrect(answer)
	if err != nil {
		return nil, err
	}

	// Finish Activity if answer is correct
	if isCorrect {
		err = c.finishActivityTrasaction(userID, activityID, activity.GetInfo().Point)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			logs.New().Error(err)
			return nil, errs.ErrInsertError
		}
	}

	// Get user from the repository
	userDB, err := c.UserRepo.GetUserByID(userID)
	if err != nil || userDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	// Create activity answer response
	response := response.NewActivityAnswerResponse(activity, userDB.Point, isCorrect)

	return response, nil
}
