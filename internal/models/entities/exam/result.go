package exam

import "time"

type ExamResult struct {
	ID                   int       `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	ExamID               int       `gorm:"column:exam_id" json:"exam_id"`
	UserID               int       `gorm:"column:user_id" json:"user_id"`
	Score                int       `gorm:"->;column:score" json:"score"`
	ExamType             string    `gorm:"->;column:type" json:"exam_type"`
	ExamContentGroupName string    `gorm:"->;column:content_group_name" json:"content_group_name"`
	IsPassed             bool      `gorm:"column:is_passed" json:"is_passed"`
	CreatedTimestamp     time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type ExamResults []ExamResult

func (results ExamResults) FinishedExam(examID int) bool {
	for _, result := range results {
		if result.ExamID == examID {
			return true
		}
	}
	return false
}

func (results ExamResults) GetResultOverview(examID int) []ResultOverview {
	overview := make([]ResultOverview, 0)
	for _, result := range results {
		if result.ExamID == examID {
			overview = append(overview, ResultOverview{
				ExamResultID:     result.ID,
				TotalScore:       result.Score,
				IsPassed:         result.IsPassed,
				CreatedTimestamp: result.CreatedTimestamp,
			})
		}
	}
	return overview
}

type Result struct {
	ActivitiesResult ResultActivities
	ExamResult       ExamResult
}

type ResultActivity struct {
	ExamResultID int `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	ActivityID   int `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	Score        int `gorm:"column:score" json:"score"`
}

type ResultActivities []ResultActivity

func (activities *ResultActivities) SetExamResultID(id int) {
	newActivities := make(ResultActivities, 0)

	for _, activity := range *activities {
		newActivities = append(newActivities, ResultActivity{
			ExamResultID: id,
			ActivityID:   activity.ActivityID,
			Score:        activity.Score,
		})
	}

	*activities = newActivities
}
