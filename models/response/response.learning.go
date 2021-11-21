package response

import "DatabaseCamp/models/entities"

type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}

func NewVideoLectureResponse(contentID int, contentName string, videoLink string) *VideoLectureResponse {
	return &VideoLectureResponse{
		ContentID:   contentID,
		ContentName: contentName,
		VideoLink:   videoLink,
	}
}

type ActivityResponse struct {
	Activity entities.ActivityDetail `json:"activity"`
	Choices  interface{}             `json:"choice"`
	Hint     entities.ActivityHint   `json:"hint"`
}

func NewActivityResponse(activity entities.Activity) *ActivityResponse {
	return &ActivityResponse{
		Activity: activity.GetInfo(),
		Choices:  activity.GetPropositionChoices(),
		Hint:     *activity.GetHint(),
	}
}

type AnswerResponse struct {
	ActivityID   int  `json:"activity_id"`
	IsCorrect    bool `json:"is_correct"`
	UpdatedPoint int  `json:"updated_point"`
}

func NewActivityAnswerResponse(activity entities.Activity, updatedPoint int, isCorrect bool) *AnswerResponse {
	return &AnswerResponse{
		ActivityID:   activity.GetInfo().ID,
		IsCorrect:    isCorrect,
		UpdatedPoint: updatedPoint,
	}
}
