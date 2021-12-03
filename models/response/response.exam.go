package response

import (
	"DatabaseCamp/models/entities"
)

type ExamResponse struct {
	Exam       entities.ExamInfo  `json:"exam"`
	Activities []ActivityResponse `json:"activities"`
}

// Create exam response instance
func NewExamResponse(exam entities.Exam) *ExamResponse {
	response := ExamResponse{}
	response.prepare(exam)
	return &response
}

// Prepare exam entities
func (e *ExamResponse) prepare(exam entities.Exam) {
	activitiesResponse := make([]ActivityResponse, 0)
	for _, activity := range exam.GetActivities() {
		activitiesResponse = append(activitiesResponse, ActivityResponse{
			Activity: activity.GetInfo(),
			Choices:  activity.GetPropositionChoices(),
		})
	}
	e.Activities = activitiesResponse
	e.Exam = exam.GetInfo()
}
