package response

import "DatabaseCamp/models/entities"

type ActivityResponse struct {
	Activity entities.ActivityDetail `json:"activity"`
	Choices  interface{}             `json:"choice"`
	Hint     entities.ActivityHint   `json:"hint"`
}

func NewActivityResponse(activity entities.Activity) *ActivityResponse {
	return &ActivityResponse{
		Activity: activity.Info,
		Choices:  activity.PropositionChoices,
		Hint:     *activity.Hint,
	}
}
