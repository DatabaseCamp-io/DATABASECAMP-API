package request

import (
	"database-camp/internal/errs"
	"database-camp/internal/models/entities/exam"
)

type ExamAnswerRequest struct {
	ExamID     *int                      `json:"exam_id"`
	Activities []exam.ExamActivityAnswer `json:"activities"`
}

func (r ExamAnswerRequest) Validate() error {
	if r.ExamID == nil {
		return errs.NewBadRequestError("ไม่พบรหัสของข้อสอบในคำร้องขอ", "Exam ID Not Found")
	} else if len(r.Activities) == 0 {
		return errs.NewBadRequestError("ไม่พบกิจกรรมของข้อสอบในคำร้องขอ", "Activities Exam Not Found")
	}
	return nil
}
