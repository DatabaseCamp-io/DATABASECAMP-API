package models

import (
	"DatabaseCamp/utils"
	"sync"
	"time"
)

type ExamType string

var Exam = struct {
	Pretest  ExamType
	MiniExam ExamType
	Posttest ExamType
}{
	"PRE",
	"MINI",
	"POST",
}

type examActivityResult struct {
	ActivityID int
	Score      int
}

type examResultOverview struct {
	ExamResultID     int                  `json:"exam_result_id"`
	TotalScore       int                  `json:"score"`
	IsPassed         bool                 `json:"is_passed"`
	ActivitiesResult []examActivityResult `json:"activities_result,omitempty"`
	CreatedTimestamp time.Time            `json:"created_timestamp"`
}

type examInfo struct {
	ID               int       `json:"exam_id"`
	Type             string    `json:"exam_type"`
	Instruction      string    `json:"instruction"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	ContentGroupID   int       `json:"content_group_id"`
	ContentGroupName string    `json:"content_group_name"`
	BadgeID          int       `json:"badge_id"`
}

type exam struct {
	Info       examInfo            `json:"exam"`
	Activities []activity          `json:"activities"`
	Result     *examResultOverview `json:"result"`
}

func NewExam() *exam {
	return &exam{}
}

func (e *exam) GetInfo() examInfo {
	return e.Info
}

func (e *exam) ToExamResultOverviewResponse() *ExamResultOverviewResponse {
	return &ExamResultOverviewResponse{
		ExamID:           e.Info.ID,
		ExamResultID:     e.Result.ExamResultID,
		ExamType:         e.Info.Type,
		ContentGroupName: e.Info.ContentGroupName,
		CreatedTimestamp: e.Result.CreatedTimestamp,
		Score:            e.Result.TotalScore,
		IsPassed:         e.Result.IsPassed,
		ActivitiesResult: e.Result.ActivitiesResult,
	}
}

func (e *exam) ToExamResultActivitiesDB() []ExamResultActivityDB {
	resultActivities := make([]ExamResultActivityDB, 0)
	for _, resultActivity := range e.Result.ActivitiesResult {
		resultActivities = append(resultActivities, ExamResultActivityDB{
			ActivityID: resultActivity.ActivityID,
			Score:      resultActivity.Score,
		})
	}
	return resultActivities
}

func (e *exam) ToExamResultDB(userID int) *ExamResultDB {
	return &ExamResultDB{
		ExamID:           e.Info.ID,
		UserID:           userID,
		Score:            e.Result.TotalScore,
		IsPassed:         e.Result.IsPassed,
		CreatedTimestamp: time.Now().Local(),
	}
}

func (e *exam) ToResponse() *ExamResponse {
	activitiesResponse := make([]ActivityResponse, 0)
	for _, activity := range e.Activities {
		activitiesResponse = append(activitiesResponse, ActivityResponse{
			Activity: activity.Info,
			Choices:  activity.PropositionChoices,
		})
	}
	response := ExamResponse{
		Exam:       e.Info,
		Activities: activitiesResponse,
	}
	return &response
}

func (e *exam) PrepareResult(examResultDB ExamResultDB) {
	e.Result = &examResultOverview{
		ExamResultID:     examResultDB.ID,
		TotalScore:       examResultDB.Score,
		IsPassed:         examResultDB.IsPassed,
		CreatedTimestamp: examResultDB.CreatedTimestamp,
	}
}

func (e *exam) Prepare(examActivitiesDB []ExamActivityDB) {
	activityChoiceDBMap := map[int]interface{}{}
	for _, examActivityDB := range examActivitiesDB {
		utils.NewType().StructToStruct(examActivityDB, &e.Info)
		if activityChoiceDBMap[examActivityDB.ActivityID] == nil {
			e.initialActivityChoiceMap(examActivityDB.ActivityID, examActivityDB.ActivityTypeID, activityChoiceDBMap)
		}
		e.appendExamActivityChoice(examActivityDB, activityChoiceDBMap)
	}

	for _, examActivityDB := range examActivitiesDB {
		activityDB := ActivityDB{}
		utils.NewType().StructToStruct(examActivityDB, &activityDB)
		activity := NewActivity()
		activity.PrepareActivity(activityDB)
		activity.PrepareChoicesByChoiceDB(activityChoiceDBMap[activityDB.ID])
		e.Activities = append(e.Activities, *activity)
	}
}

func (e *exam) initialActivityChoiceMap(activityID int, typeID int, activityChoiceMap map[int]interface{}) {
	if typeID == 1 {
		temp := make([]MatchingChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 2 {
		temp := make([]MultipleChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 3 {
		temp := make([]CompletionChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else {
		temp := make([]interface{}, 0)
		activityChoiceMap[activityID] = temp
	}
}

func (e *exam) appendExamActivityChoice(examActivity ExamActivityDB, activityChoiceMap map[int]interface{}) {
	choices := activityChoiceMap[examActivity.ActivityID]
	if examActivity.ActivityTypeID == 1 {
		temp := choices.([]MatchingChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(MatchingChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 2 {
		temp := choices.([]MultipleChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(MultipleChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 3 {
		temp := choices.([]CompletionChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(CompletionChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	}
}

func (e *exam) examActivityToChoice(examActivity ExamActivityDB) interface{} {
	if examActivity.ActivityTypeID == 1 {
		choice := MatchingChoiceDB{}
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 2 {
		choice := MultipleChoiceDB{}
		examActivity.Content = examActivity.MultipleChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 3 {
		choice := CompletionChoiceDB{}
		examActivity.Content = examActivity.CompletionChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else {
		return nil
	}
}

func (e *exam) CheckAnswer(answers []ExamActivityAnswer) (*examResultOverview, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err, Mutex: &mutex}

	activityMap := map[int]*activity{}
	for _, activity := range e.Activities {
		activityMap[activity.Info.ID] = &activity
	}

	e.Result = &examResultOverview{
		TotalScore:       0,
		IsPassed:         false,
		CreatedTimestamp: time.Now().Local(),
	}

	wg.Add(len(answers))
	for _, answer := range answers {
		go e.checkActivityAsync(&concurrent, *activityMap[answer.ActivityID], answer)
	}
	wg.Wait()

	e.summaryResult()
	return e.Result, err
}

func (e *exam) summaryResult() {
	answerTotalScore := e.GetAnswerTotalScore()
	activitiesTotalScore := e.GetActivitiesTotalScore()
	e.Result.IsPassed = e.isPassed(answerTotalScore, activitiesTotalScore)
}

func (e *exam) GetAnswerTotalScore() int {
	sum := 0
	for _, activityResult := range e.Result.ActivitiesResult {
		sum += activityResult.Score
	}
	e.Result.TotalScore = sum
	return sum
}

func (e *exam) GetActivitiesTotalScore() int {
	sum := 0
	for _, activity := range e.Activities {
		sum += activity.Info.Point
	}
	return sum
}

func (e *exam) isPassed(answerTotalScore int, activitiesTotalScore int) bool {
	passedRate := 0.5
	if activitiesTotalScore == 0 {
		return true
	} else {
		return (float64)(answerTotalScore/activitiesTotalScore) > passedRate
	}

}

func (e *exam) checkActivityAsync(concurrent *Concurrent, activity activity, answer interface{}) {
	defer concurrent.Wg.Done()
	score := 0
	isCorrect := true
	isCorrect, err := activity.IsAnswerCorrect(answer)
	if err != nil {
		*concurrent.Err = err
		return
	}

	if isCorrect {
		score += activity.Info.Point
	}

	concurrent.Mutex.Lock()
	e.Result.ActivitiesResult = append(e.Result.ActivitiesResult, examActivityResult{
		ActivityID: activity.Info.ID,
		Score:      score,
	})
	concurrent.Mutex.Unlock()
}
