package request

// request.exam.go
/**
 * 	This file is a part of models, used to collect request of exam
 */

import "DatabaseCamp/errs"

// Model for correct exam activity answer request
type ExamActivityAnswer struct {
	ActivityID int         `json:"activity_id"`
	Answer     interface{} `json:"answer"`
}

/**
 * 	This class represent exam answer request
 */
type ExamAnswerRequest struct {
	ExamID     *int                 `json:"exam_id"`
	Activities []ExamActivityAnswer `json:"activities"`
}

/**
 * Validate exam answer request
 *
 * @return the error of validating request
 */
func (r ExamAnswerRequest) Validate() error {
	if r.ExamID == nil {
		return errs.NewBadRequestError("ไม่พบรหัสของข้อสอบในคำร้องขอ", "Exam ID Not Found")
	} else if len(r.Activities) == 0 {
		return errs.NewBadRequestError("ไม่พบกิจกรรมของข้อสอบในคำร้องขอ", "Activities Exam Not Found")
	}
	return nil
}
