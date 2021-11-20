package entities

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

var ExamType = struct {
	Pretest  string
	MiniExam string
	Posttest string
}{
	"PRE",
	"MINI",
	"POST",
}

type ExamActivityResult struct {
	ActivityID int
	Score      int
}

type ExamResultOverview struct {
	ExamResultID     int                  `json:"exam_result_id"`
	TotalScore       int                  `json:"score"`
	IsPassed         bool                 `json:"is_passed"`
	ActivitiesResult []ExamActivityResult `json:"activities_result,omitempty"`
	CreatedTimestamp time.Time            `json:"created_timestamp"`
}

type ExamInfo struct {
	ID               int       `json:"exam_id"`
	Type             string    `json:"exam_type"`
	Instruction      string    `json:"instruction"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	ContentGroupID   int       `json:"content_group_id"`
	ContentGroupName string    `json:"content_group_name"`
	BadgeID          int       `json:"badge_id"`
}

type Exam struct {
	Info       ExamInfo            `json:"exam"`
	Activities []Activity          `json:"activities"`
	Result     *ExamResultOverview `json:"result"`
}

func (e *Exam) GetInfo() ExamInfo {
	return e.Info
}

func (e *Exam) ToExamResultActivitiesDB() []general.ExamResultActivityDB {
	resultActivities := make([]general.ExamResultActivityDB, 0)
	for _, resultActivity := range e.Result.ActivitiesResult {
		resultActivities = append(resultActivities, general.ExamResultActivityDB{
			ActivityID: resultActivity.ActivityID,
			Score:      resultActivity.Score,
		})
	}
	return resultActivities
}

func (e *Exam) ToExamResultDB(userID int) *general.ExamResultDB {
	return &general.ExamResultDB{
		ExamID:           e.Info.ID,
		UserID:           userID,
		Score:            e.Result.TotalScore,
		IsPassed:         e.Result.IsPassed,
		CreatedTimestamp: time.Now().Local(),
	}
}

func (e *Exam) PrepareResult(examResultDB general.ExamResultDB) {
	e.Result = &ExamResultOverview{
		ExamResultID:     examResultDB.ID,
		TotalScore:       examResultDB.Score,
		IsPassed:         examResultDB.IsPassed,
		CreatedTimestamp: examResultDB.CreatedTimestamp,
	}
}

func (e *Exam) Prepare(examActivitiesDB []general.ExamActivityDB) {
	activityChoiceDBMap := map[int]interface{}{}
	examActivityDBMap := map[int]general.ActivityDB{}
	for _, examActivityDB := range examActivitiesDB {
		activity := general.ActivityDB{}
		utils.NewType().StructToStruct(examActivityDB, &e.Info)
		if examActivityDB.ExamType == string(ExamType.Posttest) {
			e.Info.BadgeID = 3
		}
		utils.NewType().StructToStruct(examActivityDB, &activity)
		examActivityDBMap[examActivityDB.ActivityID] = activity
		if activityChoiceDBMap[examActivityDB.ActivityID] == nil {
			e.initialActivityChoiceMap(examActivityDB.ActivityID, examActivityDB.ActivityTypeID, activityChoiceDBMap)
		}
		e.appendExamActivityChoice(examActivityDB, activityChoiceDBMap)
	}

	for _, examActivityDB := range examActivityDBMap {
		activity := Activity{}
		activity.PrepareActivity(examActivityDB)
		activity.PrepareChoicesByChoiceDB(activityChoiceDBMap[examActivityDB.ID])
		e.Activities = append(e.Activities, activity)
	}
}

func (e *Exam) initialActivityChoiceMap(activityID int, typeID int, activityChoiceMap map[int]interface{}) {
	if typeID == 1 {
		temp := make([]general.MatchingChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 2 {
		temp := make([]general.MultipleChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 3 {
		temp := make([]general.CompletionChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else {
		temp := make([]interface{}, 0)
		activityChoiceMap[activityID] = temp
	}
}

func (e *Exam) appendExamActivityChoice(examActivity general.ExamActivityDB, activityChoiceMap map[int]interface{}) {
	choices := activityChoiceMap[examActivity.ActivityID]
	if examActivity.ActivityTypeID == 1 {
		temp := choices.([]general.MatchingChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(general.MatchingChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 2 {
		temp := choices.([]general.MultipleChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(general.MultipleChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 3 {
		temp := choices.([]general.CompletionChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(general.CompletionChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	}
}

func (e *Exam) examActivityToChoice(examActivity general.ExamActivityDB) interface{} {
	if examActivity.ActivityTypeID == 1 {
		choice := general.MatchingChoiceDB{}
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 2 {
		choice := general.MultipleChoiceDB{}
		examActivity.Content = examActivity.MultipleChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 3 {
		choice := general.CompletionChoiceDB{}
		examActivity.Content = examActivity.CompletionChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else {
		return nil
	}
}

func (e *Exam) CheckAnswer(answers []request.ExamActivityAnswer) (*ExamResultOverview, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err, Mutex: &mutex}

	activityMap := map[int]*Activity{}
	for _, activity := range e.Activities {
		activityMap[activity.Info.ID] = &activity
	}

	e.Result = &ExamResultOverview{
		TotalScore:       0,
		IsPassed:         false,
		CreatedTimestamp: time.Now().Local(),
	}

	wg.Add(len(answers))
	for _, answer := range answers {
		go e.checkActivityAsync(&concurrent, *activityMap[answer.ActivityID], answer.Answer)
	}
	wg.Wait()

	e.summaryResult()
	return e.Result, err
}

func (e *Exam) summaryResult() {
	answerTotalScore := e.GetAnswerTotalScore()
	activitiesTotalScore := e.GetActivitiesTotalScore()
	e.Result.IsPassed = e.isPassed(answerTotalScore, activitiesTotalScore)
}

func (e *Exam) GetAnswerTotalScore() int {
	sum := 0
	for _, activityResult := range e.Result.ActivitiesResult {
		sum += activityResult.Score
	}
	e.Result.TotalScore = sum
	return sum
}

func (e *Exam) GetActivitiesTotalScore() int {
	sum := 0
	for _, activity := range e.Activities {
		sum += activity.Info.Point
	}
	return sum
}

func (e *Exam) isPassed(answerTotalScore int, activitiesTotalScore int) bool {
	passedRate := 0.5
	if activitiesTotalScore == 0 {
		return true
	} else {
		return (float64)(answerTotalScore/activitiesTotalScore) > passedRate
	}

}

func (e *Exam) checkActivityAsync(concurrent *general.Concurrent, activity Activity, answer interface{}) {
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
	e.Result.ActivitiesResult = append(e.Result.ActivitiesResult, ExamActivityResult{
		ActivityID: activity.Info.ID,
		Score:      score,
	})
	concurrent.Mutex.Unlock()
}
