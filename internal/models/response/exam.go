package response

import (
	"database-camp/internal/models/entities/exam"
	"time"
)

type ExamOverviewResponse struct {
	PreExam   *exam.DetailOverview   `json:"pre_exam"`
	MiniExam  *[]exam.DetailOverview `json:"mini_exam"`
	FinalExam *exam.DetailOverview   `json:"final_exam"`
}

type ExamResultOverviewResponse struct {
	ExamID           int                   `json:"exam_id"`
	ExamResultID     int                   `json:"exam_result_id"`
	ExamType         string                `json:"exam_type"`
	ContentGroupName string                `json:"content_group_name"`
	CreatedTimestamp time.Time             `json:"created_timestamp"`
	Score            int                   `json:"score"`
	IsPassed         bool                  `json:"is_passed"`
	ActivitiesResult exam.ResultActivities `json:"activities_result"`
}

type ExamResponse struct {
	Exam       exam.Exam          `json:"exam"`
	Activities []ActivityResponse `json:"activities"`
}
