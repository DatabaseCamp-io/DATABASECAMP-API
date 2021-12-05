package response

// response.learning.go
/**
 * 	This file is a part of models, used to collect response of learning
 */

import "DatabaseCamp/models/entities"

// Model of video lecture item to prepare video lecture response
type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}

/**
 * Constructor creates a new VideoLectureResponse instance
 *
 * @param contentID				Content id from database to create video lecture response
 * @param contentName			Content name from database to create video lecture response
 * @param videoLink				Video link from database to create video lecture response
 *
 * @return 	instance of VideoLectureResponse
 */
func NewVideoLectureResponse(contentID int, contentName string, videoLink string) *VideoLectureResponse {
	return &VideoLectureResponse{
		ContentID:   contentID,
		ContentName: contentName,
		VideoLink:   videoLink,
	}
}

/**
 * This class represent activity response
 */
type ActivityResponse struct {
	Activity entities.ActivityDetail `json:"activity"`
	Choices  interface{}             `json:"choice"`
	Hint     entities.ActivityHint   `json:"hint"`
}

/**
 * Constructor creates a new ActivityResponse instance
 *
 * @param activity			Entities activity from database to create activity response
 *
 * @return 	instance of ActivityResponse
 */
func NewActivityResponse(activity entities.Activity) *ActivityResponse {
	return &ActivityResponse{
		Activity: activity.GetInfo(),
		Choices:  activity.GetPropositionChoices(),
		Hint:     *activity.GetHint(),
	}
}

/**
 * This class represent answer response
 */
type AnswerResponse struct {
	ActivityID   int  `json:"activity_id"`
	IsCorrect    bool `json:"is_correct"`
	UpdatedPoint int  `json:"updated_point"`
}

/**
 * Constructor creates a new ActivityAnswerResponse instance
 *
 * @param activity			Entities activity from database to create activity answer response
 * @param updatedPoint		Update point from database to create activity answer response
 * @param isCorrect			Is correct from database to create activity answer response
 *
 * @return 	instance of ActivityAnswerResponse
 */
func NewActivityAnswerResponse(activity entities.Activity, updatedPoint int, isCorrect bool) *AnswerResponse {
	return &AnswerResponse{
		ActivityID:   activity.GetInfo().ID,
		IsCorrect:    isCorrect,
		UpdatedPoint: updatedPoint,
	}
}
