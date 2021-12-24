package entities

// entity.exam.go
/**
 * 	This file is a part of models, used to collect model for entities of exam
 */

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

// Exam type
var ExamType = struct {
	Pretest  string
	MiniExam string
	Posttest string
}{
	"PRE",
	"MINI",
	"POST",
}

// Model of exam activity result for Exam entity
type ExamActivityResult struct {
	ActivityID int
	Score      int
}

// Model of exam result overview for Exam entity
type ExamResultOverview struct {
	ExamResultID     int                  `json:"exam_result_id"`
	TotalScore       int                  `json:"score"`
	IsPassed         bool                 `json:"is_passed"`
	ActivitiesResult []ExamActivityResult `json:"activities_result,omitempty"`
	CreatedTimestamp time.Time            `json:"created_timestamp"`
}

// Model of exam info for Exam entity
type ExamInfo struct {
	ID               int       `json:"exam_id"`
	Type             string    `json:"exam_type"`
	Instruction      string    `json:"instruction"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	ContentGroupID   int       `json:"content_group_id"`
	ContentGroupName string    `json:"content_group_name"`
	BadgeID          int       `json:"badge_id"`
}

/**
 * This class manage exam model
 */
type Exam struct {
	info       ExamInfo
	activities []Activity
	result     *ExamResultOverview
}

/**
 * Getter for getting exam result
 *
 * @return exam result
 */
func (e *Exam) GetResult() *ExamResultOverview {
	return e.result
}

/**
 * Getter for getting activities of the exam
 *
 * @return activities of the exam
 */
func (e *Exam) GetActivities() []Activity {
	return e.activities
}

/**
 * Getter for getting information of the exam
 *
 * @return information of the exam
 */
func (e *Exam) GetInfo() ExamInfo {
	return e.info
}

/**
 * Setter for set exam result ID
 *
 * @param  id 	Exam result id to set
 */
func (e *Exam) SetResultID(id int) {
	e.result.ExamResultID = id
}

/**
 * To exam result activities database model
 *
 * @return exam result activities database model
 */
func (e *Exam) ToExamResultActivitiesDB() []storages.ExamResultActivityDB {
	resultActivities := make([]storages.ExamResultActivityDB, 0)
	for _, resultActivity := range e.result.ActivitiesResult {
		resultActivities = append(resultActivities, storages.ExamResultActivityDB{
			ActivityID: resultActivity.ActivityID,
			Score:      resultActivity.Score,
		})
	}
	return resultActivities
}

/**
 * To exam result database model
 *
 * @param 	userID 		User id to set exam result model
 *
 * @return exam result database model
 */
func (e *Exam) ToExamResultDB(userID int) *storages.ExamResultDB {
	return &storages.ExamResultDB{
		ExamID:           e.info.ID,
		UserID:           userID,
		Score:            e.result.TotalScore,
		IsPassed:         e.result.IsPassed,
		CreatedTimestamp: time.Now().Local(),
	}
}

/**
 * Prepare exam result
 *
 * @param 	examResultDB 		Exam result to prepare
 */
func (e *Exam) PrepareResult(examResultDB storages.ExamResultDB) {
	e.result = &ExamResultOverview{
		ExamResultID:     examResultDB.ID,
		TotalScore:       examResultDB.Score,
		IsPassed:         examResultDB.IsPassed,
		CreatedTimestamp: examResultDB.CreatedTimestamp,
	}
}

/**
 * Prepare exam
 *
 * @param 	examActivitiesDB 	Exam activities to prepare
 */
func (e *Exam) Prepare(examActivitiesDB []storages.ExamActivityDB) {
	activityChoiceDBMap := map[int]interface{}{}
	examActivityDBMap := map[int]storages.ActivityDB{}
	for _, examActivityDB := range examActivitiesDB {
		activity := storages.ActivityDB{}
		utils.NewType().StructToStruct(examActivityDB, &e.info)
		if examActivityDB.ExamType == string(ExamType.Posttest) {
			e.info.BadgeID = 3
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
		activity.SetActivity(examActivityDB)
		activity.SetChoicesByChoiceDB(activityChoiceDBMap[examActivityDB.ID])
		e.activities = append(e.activities, activity)
	}
}

/**
 * Initail Map of activity choice
 *
 * @param 	activityID 			Activity ID to initialize
 * @param 	typeID 				Activity type ID to initialize
 * @param 	activityChoiceMap 	Map of activity choice
 */
func (e *Exam) initialActivityChoiceMap(activityID int, typeID int, activityChoiceMap map[int]interface{}) {
	if typeID == 1 {
		temp := make([]storages.MatchingChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 2 {
		temp := make([]storages.MultipleChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else if typeID == 3 {
		temp := make([]storages.CompletionChoiceDB, 0)
		activityChoiceMap[activityID] = temp
	} else {
		temp := make([]interface{}, 0)
		activityChoiceMap[activityID] = temp
	}
}

/**
 * Append exam activity choice to map of activity choice
 *
 * @param 	examActivity 			exam activity to append
 * @param 	activityChoiceMap 		Map of activity choice
 */
func (e *Exam) appendExamActivityChoice(examActivity storages.ExamActivityDB, activityChoiceMap map[int]interface{}) {
	choices := activityChoiceMap[examActivity.ActivityID]
	if examActivity.ActivityTypeID == 1 {
		temp := choices.([]storages.MatchingChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(storages.MatchingChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 2 {
		temp := choices.([]storages.MultipleChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(storages.MultipleChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	} else if examActivity.ActivityTypeID == 3 {
		temp := choices.([]storages.CompletionChoiceDB)
		temp = append(temp, e.examActivityToChoice(examActivity).(storages.CompletionChoiceDB))
		activityChoiceMap[examActivity.ActivityID] = temp
	}
}

/**
 * Convert exam activity model from the database to choice
 *
 * @param 	examActivity 			exam activity model from the database
 *
 * @return  choice that converted from exam activity model
 */
func (e *Exam) examActivityToChoice(examActivity storages.ExamActivityDB) interface{} {
	if examActivity.ActivityTypeID == 1 {
		choice := storages.MatchingChoiceDB{}
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 2 {
		choice := storages.MultipleChoiceDB{}
		examActivity.Content = examActivity.MultipleChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else if examActivity.ActivityTypeID == 3 {
		choice := storages.CompletionChoiceDB{}
		examActivity.Content = examActivity.CompletionChoiceContent
		utils.NewType().StructToStruct(examActivity, &choice)
		return choice
	} else {
		return nil
	}
}

/**
 * Check exam answers
 *
 * @param 	answers 	Answers of the exam
 *
 * @return  exam result overview
 * @return  the error of the checking exam
 */
func (e *Exam) CheckAnswer(answers []request.ExamActivityAnswer) (*ExamResultOverview, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error
	concurrent := general.Concurrent{Wg: &wg, Err: &err, Mutex: &mutex}

	activityMap := map[int]Activity{}
	for _, activity := range e.activities {
		activityMap[activity.GetInfo().ID] = activity
	}

	e.result = &ExamResultOverview{
		TotalScore:       0,
		IsPassed:         false,
		CreatedTimestamp: time.Now().Local(),
	}

	wg.Add(len(answers))
	for _, answer := range answers {
		go e.checkActivityAsync(&concurrent, activityMap[answer.ActivityID], answer.Answer)
	}
	wg.Wait()

	e.summaryResult()
	return e.result, err
}

/**
 * Summary exam result
 */
func (e *Exam) summaryResult() {
	answerTotalScore := e.GetAnswerTotalScore()
	activitiesTotalScore := e.GetActivitiesTotalScore()
	e.result.IsPassed = e.isPassed(answerTotalScore, activitiesTotalScore)
}

/**
 * Calculate total score for the answer
 *
 * @return score of the answers
 */
func (e *Exam) GetAnswerTotalScore() int {
	sum := 0
	for _, activityResult := range e.result.ActivitiesResult {
		sum += activityResult.Score
	}
	e.result.TotalScore = sum
	return sum
}

/**
 * Get total score of the activities
 *
 * @return total score of the activities
 */
func (e *Exam) GetActivitiesTotalScore() int {
	sum := 0
	for _, activity := range e.activities {
		sum += activity.GetInfo().Point
	}
	return sum
}

/**
 * Check pass exam
 *
 * @param answerTotalScore			Total score of the Answers
 * @param activitiesTotalScore		Total score of the activities
 *
 * @return true of passed the exam, false otherwise
 */
func (e *Exam) isPassed(answerTotalScore int, activitiesTotalScore int) bool {
	passedRate := 0.5
	if activitiesTotalScore == 0 {
		return true
	} else {
		return (float64)(answerTotalScore/activitiesTotalScore) > passedRate
	}

}

/**
 * Check concurrency activity answer
 *
 * @param concurrent	Concurrency model for doing concurrent
 * @param activity		Activity to check
 * @param answer		Answer of the activity
 */
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
		score += activity.GetInfo().Point
	}

	concurrent.Mutex.Lock()
	e.result.ActivitiesResult = append(e.result.ActivitiesResult, ExamActivityResult{
		ActivityID: activity.GetInfo().ID,
		Score:      score,
	})
	concurrent.Mutex.Unlock()
}
