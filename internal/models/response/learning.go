package response

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
)

type ContentOverviewResponse struct {
	PreExam              *int                           `json:"pre_exam_id"`
	LastedGroup          *content.LastedGroupOverview   `json:"lasted_group"`
	ContentGroupOverview []content.ContentGroupOverview `json:"content_group_overview"`
}

type ContentRoadmapResponse struct {
	ContentID   int                           `json:"content_id"`
	ContentName string                        `json:"content_name"`
	Items       []activity.ContentRoadmapItem `json:"items"`
}

type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}

type ActivityResponse struct {
	Activity activity.Activity      `json:"activity"`
	Choices  interface{}            `json:"choice"`
	Hint     *activity.ActivityHint `json:"hint"`
}

type AnswerResponse struct {
	ActivityID   int     `json:"activity_id"`
	IsCorrect    bool    `json:"is_correct"`
	UpdatedPoint int     `json:"updated_point"`
	ErrMessage   *string `json:"err_message"`
}

type UsedHintResponse struct {
	HintDB activity.Hint `json:"hint"`
}
