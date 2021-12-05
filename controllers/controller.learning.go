package controllers

import (
	"DatabaseCamp/controllers/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/response"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

type learningController struct {
	LearningRepo repositories.ILearningRepository
	UserRepo     repositories.IUserRepository
}

// Create learning controller instance
func NewLearningController(
	learningRepo repositories.ILearningRepository, userRepo repositories.IUserRepository) learningController {
	return learningController{LearningRepo: learningRepo, UserRepo: userRepo}
}

// Get video lecture
func (c learningController) GetVideoLecture(id int) (*response.VideoLectureResponse, error) {
	contentDB, err := c.LearningRepo.GetContent(id)
	if err != nil || contentDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrContentNotFound
	}

	videoLink, err := c.LearningRepo.GetVideoFileLink(contentDB.VideoPath)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrServiceUnavailableError
	}

	response := response.NewVideoLectureResponse(contentDB.ID, contentDB.Name, videoLink)
	return response, nil
}

func (c learningController) GetOverview(userID int) (*response.ContentOverviewResponse, error) {
	loader := loaders.NewLearningOverviewLoader(c.LearningRepo, c.UserRepo)
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}
	response := response.NewContentOverviewResponse(loader.GetOverviewDB(), loader.GetLearningProgressionDB())
	return response, nil
}

func (c learningController) GetActivity(userID int, activityID int) (*response.ActivityResponse, error) {
	loader := loaders.NewActivityLoader(c.LearningRepo, c.UserRepo)
	err := loader.Load(userID, activityID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	choiceDB, err := c.getChoices(loader.GetActivityDB().ID, loader.GetActivityDB().TypeID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrActivitiesNotFound
	}

	activity := entities.Activity{}
	activity.SetActivity(*loader.GetActivityDB())
	activity.SetChoicesByChoiceDB(choiceDB)
	activity.SetHint(loader.GetActivityHintsDB(), loader.GetUserHintsDB())

	response := response.NewActivityResponse(activity)
	return response, nil
}

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

func (c learningController) finishActivityTrasaction(userID int, activityID int, addPoint int) error {
	tx := database.NewTransaction()
	tx.Begin()

	progression := storages.LearningProgressionDB{
		UserID:           userID,
		ActivityID:       activityID,
		CreatedTimestamp: time.Now().Local(),
	}

	_, err := c.UserRepo.InsertLearningProgressionTransaction(tx, progression)
	if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
		tx.Rollback()
		return err
	}

	if !utils.NewHelper().IsSqlDuplicateError(err) {
		err = c.UserRepo.ChangePointTransaction(tx, userID, addPoint, entities.Mode.Add)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	tx.Close()

	return nil
}

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

func (c learningController) isUsedHint(userHints []storages.UserHintDB, hintID int) bool {
	for _, userHint := range userHints {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

func (c learningController) getNextLevelHint(ActivityHintsDB []storages.HintDB, userHintsDB []storages.UserHintDB) *storages.HintDB {
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
	err := c.UserRepo.ChangePointTransaction(ct.Transaction, userID, updatePoint, mode)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

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

func (c learningController) GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error) {
	loader := loaders.NewContentRoadmapLoader(c.LearningRepo, c.UserRepo)
	err := loader.Load(userID, contentID)
	if err != nil || loader.GetContentDB() == nil {
		logs.New().Error(err)
		return nil, errs.ErrContentNotFound
	}
	response := response.NewContentRoadmapResponse(*loader.GetContentDB(), loader.GetContentActivityDB(), loader.GetLearningProgressionDB())
	return response, nil
}

func (c learningController) CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*response.AnswerResponse, error) {
	loader := loaders.NewCheckAnswerLoader(c.LearningRepo)
	err := loader.Load(activityID, typeID, c.getChoices)
	if err != nil || loader.GetActivityDB() == nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	if loader.GetActivityDB().TypeID != typeID {
		return nil, errs.ErrActivityTypeInvalid
	}

	activity := entities.Activity{}
	activity.SetActivity(*loader.GetActivityDB())
	activity.SetChoicesByChoiceDB(loader.GetChoicesDB())

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

	userDB, err := c.UserRepo.GetUserByID(userID)
	if err != nil || userDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	response := response.NewActivityAnswerResponse(activity, userDB.Point, isCorrect)
	return response, nil
}
