package response

import "DatabaseCamp/models/entities"

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
