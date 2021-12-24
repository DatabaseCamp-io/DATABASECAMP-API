package response

// response.exam.go
/**
 * 	This file is a part of models, used to collect response of exam
 */

import (
	"DatabaseCamp/models/entities"
)

/**
 * This class represent exam response
 */
type ExamResponse struct {
	Exam       entities.ExamInfo  `json:"exam"`
	Activities []ActivityResponse `json:"activities"`
}

/**
 * Constructor creates a new ExamResponse instance
 *
 * @param exam		Entities exam for create exam response
 *
 * @return 	instance of ExamOverviewResponse
 */
func NewExamResponse(exam entities.Exam) *ExamResponse {
	response := ExamResponse{}
	response.prepare(exam)
	return &response
}

/**
* Prepare exam response
*
* @param exam		Entities exam for create exam response
 */
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
