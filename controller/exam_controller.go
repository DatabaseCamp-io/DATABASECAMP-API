package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"
	"DatabaseCamp/utils"
)

type examController struct {
	examRepo repository.IExamRepository
}

type IExamController interface {
	GetExam(examID int, userID int) (interface{}, error)
}

func NewExamController(examRepo repository.IExamRepository) examController {
	return examController{examRepo: examRepo}
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
