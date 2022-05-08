package exam

import "time"

type ResultOverview struct {
	ExamResultID     int       `json:"exam_result_id"`
	TotalScore       int       `json:"score"`
	IsPassed         bool      `json:"is_passed"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
}

type DetailOverview struct {
	ExamID           int               `json:"exam_id"`
	ExamType         string            `json:"exam_type"`
	ContentGroupID   *int              `json:"content_group_id,omitempty"`
	ContentGroupName *string           `json:"content_group_name,omitempty"`
	CanDo            *bool             `json:"can_do,omitempty"`
	Results          *[]ResultOverview `json:"results,omitempty"`
}
