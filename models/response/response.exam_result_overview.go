package response

import (
	"DatabaseCamp/models/entities"
	"time"
)

type ExamResultOverviewResponse struct {
	ExamID           int                           `json:"exam_id"`
	ExamResultID     int                           `json:"exam_result_id"`
	ExamType         string                        `json:"exam_type"`
	ContentGroupName string                        `json:"content_group_name"`
	CreatedTimestamp time.Time                     `json:"created_timestamp"`
	Score            int                           `json:"score"`
	IsPassed         bool                          `json:"is_passed"`
	ActivitiesResult []entities.ExamActivityResult `json:"activities_result"`
}

// Create exan result overview response instance
func NewExamResultOverviewResponse(exam entities.Exam) *ExamResultOverviewResponse {
	return &ExamResultOverviewResponse{
		ExamID:           exam.GetInfo().ID,
		ExamResultID:     exam.GetResult().ExamResultID,
		ExamType:         exam.GetInfo().Type,
		ContentGroupName: exam.GetInfo().ContentGroupName,
		CreatedTimestamp: exam.GetResult().CreatedTimestamp,
		Score:            exam.GetResult().TotalScore,
		IsPassed:         exam.GetResult().IsPassed,
		ActivitiesResult: exam.GetResult().ActivitiesResult,
	}
}
