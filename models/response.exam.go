package models

import "time"

type ExamOverviewResponse struct {
	PreExam   *examDetailOverview   `json:"pre_exam"`
	MiniExam  *[]examDetailOverview `json:"mini_exam"`
	FinalExam *examDetailOverview   `json:"final_exam"`
}

type ExamResponse struct {
	Exam       examInfo           `json:"exam"`
	Activities []ActivityResponse `json:"activities"`
}

type ExamResultOverviewResponse struct {
	ExamID           int                  `json:"exam_id"`
	ExamResultID     int                  `json:"exam_result_id"`
	ExamType         string               `json:"exam_type"`
	ContentGroupName string               `json:"content_group_name"`
	CreatedTimestamp time.Time            `json:"created_timestamp"`
	Score            int                  `json:"score"`
	IsPassed         bool                 `json:"is_passed"`
	ActivitiesResult []examActivityResult `json:"activities_result"`
}
