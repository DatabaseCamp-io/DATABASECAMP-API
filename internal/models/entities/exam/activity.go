package exam

import (
	"database-camp/internal/models/entities/activity"
	"time"
)

type ExamActivityAnswer struct {
	ActivityID int         `json:"activity_id"`
	Answer     interface{} `json:"answer"`
}

type ExamActivity struct {
	ActivityID     int `gorm:"column:activity_id"`
	ActivityTypeID int `gorm:"column:activity_type_id"`
}

type ExamActivities []ExamActivity

type Activity struct {
	activity.Activity
	activity.Choices
}

type Activities []Activity

func (activities Activities) CheckAnswers(examID int, userID int, answers []ExamActivityAnswer) (*Result, error) {
	answerScore := 0
	totalScore := 0

	activitiesResult := make([]ResultActivity, 0)

	for _, examActivity := range activities {
		for _, answer := range answers {
			if examActivity.Activity.ID == answer.ActivityID {
				formatedAnswer, err := activity.FormatAnswer(answer.Answer, examActivity.Activity.TypeID)
				if err != nil {
					return nil, err
				}

				isCorrect, err := formatedAnswer.IsCorrect(examActivity.Choices)
				if err != nil {
					return nil, err
				}

				if isCorrect {
					answerScore += examActivity.Activity.Point
				}

				totalScore += examActivity.Activity.Point

				activitiesResult = append(activitiesResult, ResultActivity{
					ActivityID: examActivity.Activity.ID,
					Score:      answerScore,
				})
			}
		}
	}

	return &Result{
		ActivitiesResult: activitiesResult,
		ExamResult: ExamResult{
			ExamID:           examID,
			UserID:           userID,
			Score:            answerScore,
			IsPassed:         isPassed(answerScore, totalScore),
			CreatedTimestamp: time.Now().Local(),
		},
	}, nil
}

func isPassed(answerTotalScore int, activitiesTotalScore int) bool {
	passedRate := 0.5
	if activitiesTotalScore == 0 {
		return true
	} else {
		return (float64)(answerTotalScore/activitiesTotalScore) > passedRate
	}

}
