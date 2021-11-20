package request

import "DatabaseCamp/errs"

type ExamActivityAnswer struct {
	ActivityID int         `json:"activity_id"`
	Answer     interface{} `json:"answer"`
}

type ExamAnswerRequest struct {
	ExamID     *int                 `json:"exam_id"`
	Activities []ExamActivityAnswer `json:"activities"`
}

func (r ExamAnswerRequest) Validate() error {
	if r.ExamID == nil {
		return errs.NewBadRequestError("ไม่พบรหัสของข้อสอบในคำร้องขอ", "Exam ID Not Found")
	} else if len(r.Activities) == 0 {
		return errs.NewBadRequestError("ไม่พบกิจกรรมของข้อสอบในคำร้องขอ", "Activities Exam Not Found")
	}
	return nil
}
