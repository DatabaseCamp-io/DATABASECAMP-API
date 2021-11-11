package controller

import (
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

type examOverviewInfo struct {
	correctedBadge []models.CorrectedBadgeDB
	exam           []models.ExamDB
	examResults    []models.ExamResultDB
}

type examController struct {
	examRepo repository.IExamRepository
	userRepo repository.IUserRepository
}

type IExamController interface {
	GetExam(examID int, userID int) (interface{}, error)
	GetOverview(userID int) (*models.ExamOverviewResponse, error)
	CheckExam(userID int, request models.ExamAnswerRequest) (*models.ExamResultOverview, error)
}

func NewExamController(examRepo repository.IExamRepository, userRepo repository.IUserRepository) examController {
	return examController{examRepo: examRepo, userRepo: userRepo}
}

func (c examController) loadOverviewInfo(userID int) (*examOverviewInfo, error) {
	var wg sync.WaitGroup
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	correctedBadge := make([]models.CorrectedBadgeDB, 0)
	exam := make([]models.ExamDB, 0)
	examResults := make([]models.ExamResultDB, 0)
	wg.Add(3)
	go c.loadCorrectedBadgeAsync(&concurrent, userID, &correctedBadge)
	go c.loadExamAsync(&concurrent, &exam)
	go c.loadExamResultAsync(&concurrent, userID, &examResults)
	wg.Wait()
	if err != nil {
		return nil, err
	}

	info := examOverviewInfo{
		correctedBadge: correctedBadge,
		exam:           exam,
		examResults:    examResults,
	}

	return &info, nil
}

func (c examController) loadExamResultAsync(concurrent *models.Concurrent, userID int, examResults *[]models.ExamResultDB) {
	defer concurrent.Wg.Done()
	result, err := c.userRepo.GetExamResult(userID)
	if err != nil {
		*concurrent.Err = err
	}
	*examResults = append(*examResults, result...)
}

func (c examController) loadExamAsync(concurrent *models.Concurrent, exam *[]models.ExamDB) {
	defer concurrent.Wg.Done()
	result, err := c.examRepo.GetExamOverview()
	if err != nil {
		*concurrent.Err = err
	}
	*exam = append(*exam, result...)
}

func (c examController) loadCorrectedBadgeAsync(concurrent *models.Concurrent, userID int, correctedBadge *[]models.CorrectedBadgeDB) {
	defer concurrent.Wg.Done()
	result, err := c.userRepo.GetCollectedBadge(userID)
	if err != nil {
		*concurrent.Err = err
	}
	*correctedBadge = append(*correctedBadge, result...)
}

func (c examController) countExamScore(examResults []models.ExamResultDB) map[int]int {
	examCountScore := map[int]int{}
	for _, v := range examResults {
		examCountScore[v.ID] += v.Score
	}
	return examCountScore
}

func (c examController) canDoFianlExam(info examOverviewInfo) bool {
	for _, v := range info.correctedBadge {
		if v.UserID == nil {
			return false
		}
	}

	return true
}

func (c examController) prepareExamResultMap(info examOverviewInfo) map[int]*[]models.ExamResultOverview {
	examResultMap := map[int]*[]models.ExamResultOverview{}
	examCountScore := c.countExamScore(info.examResults)
	for _, v := range info.examResults {
		if examResultMap[v.ExamID] == nil {
			temp := make([]models.ExamResultOverview, 0)
			examResultMap[v.ExamID] = &temp
		}
		*examResultMap[v.ExamID] = append(*examResultMap[v.ExamID], models.ExamResultOverview{
			CreatedTimestamp: v.CreatedTimestamp,
			Score:            examCountScore[v.ID],
			IsPassed:         v.IsPassed,
		})
	}
	return examResultMap
}

func (c examController) prepareExamOverview(info examOverviewInfo) *models.ExamOverviewResponse {
	res := models.ExamOverviewResponse{}
	examResultMap := c.prepareExamResultMap(info)

	for _, v := range info.exam {
		if v.Type == string(models.Exam.Pretest) {
			res.PreExam = &models.ExamOverview{
				ExamID:   v.ID,
				ExamType: v.Type,
				Results:  examResultMap[v.ID],
			}
		} else if v.Type == string(models.Exam.MiniExam) {
			if res.MiniExam == nil {
				temp := make([]models.ExamOverview, 0)
				res.MiniExam = &temp
			}

			_v := v

			*res.MiniExam = append(*res.MiniExam, models.ExamOverview{
				ExamID:           v.ID,
				ExamType:         v.Type,
				ContentGroupID:   &_v.ContentGroupID,
				ContentGroupName: &_v.ContentGroupName,
				Results:          examResultMap[v.ID],
			})
		} else if v.Type == string(models.Exam.Posttest) {
			cando := c.canDoFianlExam(info)
			res.FinalExam = &models.ExamOverview{
				ExamID:   v.ID,
				ExamType: v.Type,
				CanDo:    &cando,
				Results:  examResultMap[v.ID],
			}
		}
	}

	return &res
}

func (c examController) GetOverview(userID int) (*models.ExamOverviewResponse, error) {
	info, err := c.loadOverviewInfo(userID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}

	return c.prepareExamOverview(*info), nil
}

func (c examController) GetExam(examID int, userID int) (interface{}, error) {
	examActivity, err := c.examRepo.GetExamActivity(examID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}

	if len(examActivity) == 0 {
		return nil, errs.NewNotFoundError("ไม่พบข้อสอบ", "Exam Not Found")
	}

	preparedExam, _ := c.prepareExam(examActivity)

	return preparedExam, nil
}

func (c examController) initialActivityChoiceMap(activityID int, typeID int, activityChoiceMap map[int]interface{}) {
	if typeID == 1 {
		temp := make([]models.MatchingChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 2 {
		temp := make([]models.MultipleChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 3 {
		temp := make([]models.CompletionChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else {
		temp := make([]interface{}, 0)
		activityChoiceMap[activityID] = temp
	}
}

func (c examController) appendActivityChoice(examActivity models.ExamActivity, activityChoiceMap map[int]interface{}) {
	choices := activityChoiceMap[examActivity.ActivityID]
	if examActivity.ActivityTypeID == 1 {
		temp := choices.([]models.MatchingChoiceDB)
		temp = append(temp, c.getChoice(examActivity).(models.MatchingChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 2 {
		temp := choices.([]models.MultipleChoiceDB)
		temp = append(temp, c.getChoice(examActivity).(models.MultipleChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 3 {
		temp := choices.([]models.CompletionChoiceDB)
		temp = append(temp, c.getChoice(examActivity).(models.CompletionChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	}
}

func (c examController) prepareChoices(activityID int, typeID int, activityChoiceMap map[int]interface{}) interface{} {
	choices := activityChoiceMap[activityID]
	activityManager := services.NewActivityManager()
	if typeID == 1 {
		return activityManager.PrepareMatchingChoice(choices.([]models.MatchingChoiceDB))
	} else if typeID == 2 {
		return activityManager.PrepareMultipleChoice(choices.([]models.MultipleChoiceDB))
	} else if typeID == 3 {
		return activityManager.PrepareCompletionChoice(choices.([]models.CompletionChoiceDB))
	}
	return nil
}

func (c examController) getChoice(examActivity models.ExamActivity) interface{} {
	if examActivity.ActivityTypeID == 1 {
		choice := models.MatchingChoiceDB{}
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 2 {
		choice := models.MultipleChoiceDB{}
		examActivity.Content = examActivity.MultipleChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 3 {
		choice := models.CompletionChoiceDB{}
		examActivity.Content = examActivity.CompletionChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else {
		return nil
	}
}

func (c examController) prepareExam(examActivity []models.ExamActivity) (models.ExamResponse, map[int]interface{}) {
	exam := models.ExamDB{}
	activityChoiceMap := map[int]interface{}{}
	activityMap := map[int]models.ActivityDB{}
	activityResponse := make([]models.ExamActivityResponse, 0)
	for _, v := range examActivity {
		activity := models.ActivityDB{}
		utils.NewType().StructToStruct(v, &exam)
		utils.NewType().StructToStruct(v, &activity)
		activityMap[v.ActivityID] = activity
		if activityChoiceMap[v.ActivityID] == nil {
			c.initialActivityChoiceMap(v.ActivityID, v.ActivityTypeID, activityChoiceMap)
		}
		c.appendActivityChoice(v, activityChoiceMap)
	}

	for k := range activityChoiceMap {
		activityResponse = append(activityResponse, models.ExamActivityResponse{
			Info:    activityMap[k],
			Choices: c.prepareChoices(k, activityMap[k].TypeID, activityChoiceMap),
		})
	}

	return models.ExamResponse{
		Exam:       exam,
		Activities: activityResponse,
	}, activityChoiceMap
}

func (c examController) checkExamActivityAsync(concurrent *models.Concurrent, preparedActivity models.ExamActivityResponse, answer models.ExamActivityAnswer, examResultActivity *[]models.ExamResultActivityDB) {
	defer concurrent.Wg.Done()
	score := 0
	activityManager := services.NewActivityManager()
	isCorrect, err := activityManager.IsAnswerCorrect(preparedActivity.Info.TypeID, preparedActivity.Choices, answer.Answer)
	if err != nil {
		*concurrent.Err = err
		return
	}
	if isCorrect {
		score += preparedActivity.Info.Point
	}
	concurrent.Mutex.Lock()
	*examResultActivity = append(*examResultActivity, models.ExamResultActivityDB{
		ActivityID: preparedActivity.Info.ID,
		Score:      score,
	})
	concurrent.Mutex.Unlock()
}

func (c examController) checkExamActivities(preparedExam models.ExamResponse, activityChoiceMap map[int]interface{}, request models.ExamAnswerRequest) ([]models.ExamResultActivityDB, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err, Mutex: &mutex}
	answerMap := map[int]models.ExamActivityAnswer{}
	examResultActivity := make([]models.ExamResultActivityDB, 0)

	for _, v := range request.Activities {
		answerMap[v.ActivityID] = v
	}

	wg.Add(len(preparedExam.Activities))
	for _, v := range preparedExam.Activities {
		_v := v
		_v.Choices = activityChoiceMap[v.Info.ID]
		go c.checkExamActivityAsync(&concurrent, _v, answerMap[v.Info.ID], &examResultActivity)
	}
	wg.Wait()

	return examResultActivity, err
}

func (c examController) sumScore(examResultActivity []models.ExamResultActivityDB) int {
	sum := 0
	for _, v := range examResultActivity {
		sum += v.Score
	}
	return sum
}

func (c examController) isPassed(sumScore int, totalScore int) bool {
	passedRate := 0.5
	if totalScore == 0 {
		return true
	} else {
		return (float64)(sumScore/totalScore) > passedRate
	}

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

func (c examController) saveExamResult(userID int, preparedExam models.ExamResponse, examResultActivity []models.ExamResultActivityDB) (*models.ExamResultOverview, error) {
	tx := database.NewTransaction()
	tx.Begin()

	sumScore := c.sumScore(examResultActivity)
	totalScore := c.calculateTotalScore(preparedExam.Activities)
	isPassed := c.isPassed(sumScore, totalScore)
	examResult := models.ExamResultDB{
		ExamID:           preparedExam.Exam.ID,
		UserID:           userID,
		IsPassed:         isPassed,
		CreatedTimestamp: time.Now().Local(),
	}

	insertedExamResult, err := c.examRepo.InsertExamResultTransaction(tx, examResult)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = c.examRepo.InsertExamResultActivityTransaction(tx, c.addExamResultID(insertedExamResult.ID, examResultActivity))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if preparedExam.Exam.Type == string(models.Exam.MiniExam) && isPassed {
		badge := models.UserBadgeDB{
			UserID:  userID,
			BadgeID: preparedExam.Exam.BadgeID,
		}
		_, err = c.userRepo.InsertUserBadgeTransaction(tx, badge)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	tx.Close()

	result := models.ExamResultOverview{
		ExamID:           preparedExam.Exam.ID,
		ExamResultID:     insertedExamResult.ID,
		CreatedTimestamp: insertedExamResult.CreatedTimestamp,
		Score:            sumScore,
		IsPassed:         isPassed,
	}

	return &result, nil
}

func (c examController) calculateTotalScore(activities []models.ExamActivityResponse) int {
	sum := 0
	for _, v := range activities {
		sum += v.Info.Point
	}
	return sum
}

func (c examController) CheckExam(userID int, request models.ExamAnswerRequest) (*models.ExamResultOverview, error) {
	examActivity, err := c.examRepo.GetExamActivity(*request.ExamID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}

	preparedExam, activityChoiceMap := c.prepareExam(examActivity)

	if len(request.Activities) != len(preparedExam.Activities) {
		return nil, errs.NewBadRequestError("จำนวนของกิจกรรมไม่ถูกต้อง", "Number of Activity Incorrect")
	}

	examResultActivity, err := c.checkExamActivities(preparedExam, activityChoiceMap, request)
	if err != nil {
		return nil, err
	}

	result, err := c.saveExamResult(userID, preparedExam, examResultActivity)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาดในการบันทึกข้อมูล", "Internal Server Error")
	}

	return result, nil
}
