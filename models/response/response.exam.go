package response

import (
	"DatabaseCamp/models/entities"
)

type ExamResponse struct {
	Exam       entities.ExamInfo  `json:"exam"`
	Activities []ActivityResponse `json:"activities"`
}

func NewExamResponse(exam entities.Exam) *ExamResponse {
	response := ExamResponse{}
	response.prepare(exam)
	return &response
}

func (e *ExamResponse) prepare(exam entities.Exam) {
	activitiesResponse := make([]ActivityResponse, 0)
	for _, activity := range exam.Activities {
		activitiesResponse = append(activitiesResponse, ActivityResponse{
			Activity: activity.Info,
			Choices:  activity.PropositionChoices,
		})
	}
	e.Activities = activitiesResponse
	e.Exam = exam.Info
}
