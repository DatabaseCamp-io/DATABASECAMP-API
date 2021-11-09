package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"
	"DatabaseCamp/utils"
	"sync"
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

			*res.MiniExam = append(*res.MiniExam, models.ExamOverview{
				ExamID:           v.ID,
				ExamType:         v.Type,
				ContentGroupID:   &v.ContentGroupID,
				ContentGroupName: &v.ContentGroupName,
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

	return c.prepareExam(examActivity), nil
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

func (c examController) prepareExam(examActivity []models.ExamActivity) interface{} {
	exam := models.ExamDB{}
	activityChoiceMap := map[int]interface{}{}
	activityMap := map[int]models.ActivityDB{}
	activityResponse := make([]interface{}, 0)
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
		activityResponse = append(activityResponse, map[string]interface{}{
			"info":    activityMap[k],
			"choices": c.prepareChoices(k, activityMap[k].TypeID, activityChoiceMap),
		})
	}

	return map[string]interface{}{
		"exam":     exam,
		"activity": activityResponse,
	}
}
