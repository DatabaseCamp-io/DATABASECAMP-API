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

func NewExamResultOverviewResponse(exam entities.Exam) *ExamResultOverviewResponse {
	return &ExamResultOverviewResponse{
		ExamID:           exam.Info.ID,
		ExamResultID:     exam.Result.ExamResultID,
		ExamType:         exam.Info.Type,
		ContentGroupName: exam.Info.ContentGroupName,
		CreatedTimestamp: exam.Result.CreatedTimestamp,
		Score:            exam.Result.TotalScore,
		IsPassed:         exam.Result.IsPassed,
		ActivitiesResult: exam.Result.ActivitiesResult,
	}
}
