package response

// response.exam_result_overview.go
/**
 * 	This file is a part of models, used to collect response of exam result overview
 */

import (
	"DatabaseCamp/models/entities"
	"time"
)

/**
 * This class represent exam result overview response
 */
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

/**
 * Constructor creates a new ExamResultOverviewResponse instance
 *
 * @param exam		Entities exam for create exam result overview response
 *
 * @return 	instance of ExamOverviewResponse
 */
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
