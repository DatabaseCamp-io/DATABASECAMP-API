package models

type ActivityResponse struct {
	Activity ActivityDB   `json:"activity"`
	Choice   interface{}  `json:"choice"`
	Hint     ActivityHint `json:"hint"`
}

type AnswerResponse struct {
	ActivityID int  `json:"activity_id"`
	IsCorrect  bool `json:"is_correct"`
}

type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}

type OverviewResponse struct {
	LastedGroup          *LastedGroup           `json:"lasted_group"`
	ContentGroupOverview []ContentGroupOverview `json:"content_group_overview"`
}

type RoadmapResponse struct {
	ContentID   int           `json:"content_id"`
	ContentName string        `json:"content_name"`
	Items       []RoadmapItem `json:"items"`
}
