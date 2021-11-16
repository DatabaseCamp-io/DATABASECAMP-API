package models

type ActivityResponse struct {
	Activity activityDetail `json:"activity"`
	Choices  interface{}    `json:"choice"`
	Hint     activityHint   `json:"hint"`
}

type OverviewResponse struct {
	LastedGroup          *lastedGroupOverview   `json:"lasted_group"`
	ContentGroupOverview []contentGroupOverview `json:"content_group_overview"`
}

type ContentRoadmapResponse struct {
	ContentID   int                  `json:"content_id"`
	ContentName string               `json:"content_name"`
	Items       []contentRoadmapItem `json:"items"`
}

type AnswerResponse struct {
	ActivityID   int  `json:"activity_id"`
	IsCorrect    bool `json:"is_correct"`
	UpdatedPoint int  `json:"updated_point"`
}

type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}
