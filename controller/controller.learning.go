package controller

import (
	loader "DatabaseCamp/controller/loaders"
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

type learningController struct {
	learningRepo repository.ILearningRepository
	userRepo     repository.IUserRepository
	service      services.IAwsService
}

type ILearningController interface {
	GetVideoLecture(id int) (*models.VideoLectureResponse, error)
	GetOverview(userID int) (*models.OverviewResponse, error)
	GetActivity(userID int, activityID int) (*models.ActivityResponse, error)
	UseHint(userID int, activityID int) (*models.HintDB, error)
	GetContentRoadmap(userID int, contentID int) (*models.ContentRoadmapResponse, error)
	CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*models.AnswerResponse, error)
}

func NewLearningController(
	learningRepo repository.ILearningRepository,
	userRepo repository.IUserRepository,
	service services.IAwsService,
) learningController {
	return learningController{
		learningRepo: learningRepo,
		userRepo:     userRepo,
		service:      service,
	}
}

func (c learningController) GetVideoLecture(id int) (*models.VideoLectureResponse, error) {
	contentDB, err := c.learningRepo.GetContent(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบเนื้อหา", "Content Not Found")
	}

	videoLink, err := c.service.GetFileLink(contentDB.VideoPath)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewServiceUnavailableError("Service ไม่พร้อมใช้งาน", "Service Unavailable")
	}

	res := models.VideoLectureResponse{
		ContentID:   contentDB.ID,
		ContentName: contentDB.Name,
		VideoLink:   videoLink,
	}

	return &res, nil
}

func (c learningController) GetOverview(userID int) (*models.OverviewResponse, error) {
	loader := loader.NewLearningOverviewLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
	overview := models.NewOverview()
	overview.Prepare(loader.OverviewDB, loader.LearningProgressionDB)
	response := overview.ToResponse()
	return response, nil
}

func (c learningController) GetActivity(userID int, activityID int) (*models.ActivityResponse, error) {
	loader := loader.NewActivityLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID, activityID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}
	choiceDB, err := c.getChoices(loader.ActivityDB.ID, loader.ActivityDB.TypeID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}
	activity := models.NewActivity()
	activity.PrepareActivity(*loader.ActivityDB)
	activity.PrepareChoicesByChoiceDB(choiceDB)
	activity.PrepareHint(loader.ActivityHintsDB, loader.UserHintsDB)
	response := activity.ToPropositionResponse()
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
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
}

func (c learningController) finishActivityTrasaction(userID int, activityID int, addPoint int) error {
	tx := database.NewTransaction()
	tx.Begin()

	progression := models.LearningProgressionDB{
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
		err = c.userRepo.ChangePointTransaction(tx, userID, addPoint, models.Mode.Add)
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	tx.Close()

	return nil
}

func (c learningController) CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*models.AnswerResponse, error) {
	loader := loader.NewCheckAnswerLoader(c.learningRepo)
	err := loader.Load(activityID, typeID, c.getChoices)
	if err != nil || loader.ActivityDB == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}

	if loader.ActivityDB.TypeID != typeID {
		return nil, errs.NewNotFoundError("ประเภทของกิจกรรมไม่ถูกต้อง", "Invalid Activity Type")
	}

	activity := models.NewActivity()
	activity.PrepareActivity(*loader.ActivityDB)
	activity.PrepareChoicesByChoiceDB(loader.ChoicesDB)
	isCorrect, err := activity.IsAnswerCorrect(answer)
	if err != nil {
		return nil, err
	}

	if isCorrect {
		err = c.finishActivityTrasaction(userID, activityID, activity.Info.Point)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			logs.New().Error(err)
			return nil, errs.NewInternalServerError("เกิดข้อผิดพลาดในการบันทึกกิจกรรม", "Saved Activity Failed")
		}
	}

	userDB, err := c.userRepo.GetUserByID(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}

	response := activity.ToAnswerResponse(userDB.Point, isCorrect)

	return response, nil
}

func (c learningController) isUsedHint(userHints []models.UserHintDB, hintID int) bool {
	for _, userHint := range userHints {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

func (c learningController) getNextLevelHint(ActivityHintsDB []models.HintDB, userHintsDB []models.UserHintDB) *models.HintDB {
	for _, activityHint := range ActivityHintsDB {
		if !c.isUsedHint(userHintsDB, activityHint.ID) {
			return &activityHint
		}
	}
	return nil
}

func (c learningController) UseHint(userID int, activityID int) (*models.HintDB, error) {
	loader := loader.NewHintLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID, activityID)
	if err != nil || len(loader.ActivityHintsDB) == 0 {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบคำใบ้ของกิจกรรม", "Activity Hints Not Found")
	}
	nextLevelHint := c.getNextLevelHint(loader.ActivityHintsDB, loader.UserHintsDB)

	if nextLevelHint == nil {
		return nil, errs.NewBadRequestError("ได้ใช้คำใบ้ทั้งหมดของกิจกรรมแล้ว", "Activity Hints Has Been Used")
	}

	if loader.UserDB.Point < nextLevelHint.PointReduce {
		return nil, errs.NewBadRequestError("แต้มไม่เพียงพอในการขอคำใบ้", "Not Enough Points")
	}

	err = c.useHintTransaction(userID, nextLevelHint.PointReduce, nextLevelHint.ID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
	return nextLevelHint, nil
}

func (c learningController) useHintTransaction(userID int, reducePoint int, hintID int) error {
	var wg sync.WaitGroup
	var err error
	tx := database.NewTransaction()

	ct := models.ConcurrentTransaction{
		Concurrent: &models.Concurrent{
			Wg:  &wg,
			Err: &err,
		},
		Transaction: tx,
	}

	wg.Add(2)
	go c.updateUserPointAsyncTrasaction(&ct, userID, reducePoint, models.Mode.Reduce)
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

func (c learningController) updateUserPointAsyncTrasaction(ct *models.ConcurrentTransaction, userID int, updatePoint int, mode models.ChangePointMode) {
	defer ct.Concurrent.Wg.Done()
	err := c.userRepo.ChangePointTransaction(ct.Transaction, userID, updatePoint, mode)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

func (c learningController) insertUserHintAsyncTransaction(ct *models.ConcurrentTransaction, userID int, hintID int) {
	defer ct.Concurrent.Wg.Done()
	hint := models.UserHintDB{
		UserID:           userID,
		HintID:           hintID,
		CreatedTimestamp: time.Now().Local(),
	}
	_, err := c.userRepo.InsertUserHintTransaction(ct.Transaction, hint)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

func (c learningController) GetContentRoadmap(userID int, contentID int) (*models.ContentRoadmapResponse, error) {
	loader := loader.NewContentRoadmapLoader(c.learningRepo, c.userRepo)
	err := loader.Load(userID, contentID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบเนื้อหา", "Content Not Found")
	}
	roadmap := models.NewContentRoadmap()
	roadmap.Prepare(*loader.ContentDB, loader.ContentActivityDB, loader.LearningProgressionDB)
	response := roadmap.ToResponse()
	return response, nil
}
