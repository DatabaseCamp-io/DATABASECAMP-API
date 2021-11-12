package models

import (
	"DatabaseCamp/errs"
	"time"
)

type ExamType string

var Exam = struct {
	Pretest  ExamType
	MiniExam ExamType
	Posttest ExamType
}{
	"PRE",
	"MINI",
	"POST",
}

type ExamResultOverview struct {
	ExamID           int       `json:"exam_id"`
	ExamResultID     int       `json:"exam_result_id"`
	ExamType         string    `json:"exam_type"`
	ContentGroupName string    `json:"content_group_name"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	Score            int       `json:"score"`
	IsPassed         bool      `json:"is_passed"`
}

type ExamOverview struct {
	ExamID           int                   `json:"exam_id"`
	ExamType         string                `json:"exam_type"`
	ContentGroupID   *int                  `json:"content_group_id,omitempty"`
	ContentGroupName *string               `json:"content_group_name,omitempty"`
	CanDo            *bool                 `json:"can_do,omitempty"`
	Results          *[]ExamResultOverview `json:"results"`
}

type ExamActivityInfo struct {
	ActivityID     int    `json:"activity_id"`
	ActivityTypeID int    `json:"activity_type_id"`
	Point          int    `json:"point"`
	Story          string `json:"story"`
	Question       string `json:"question"`
}

type ExamOverviewResponse struct {
	PreExam   *ExamOverview   `json:"pre_exam"`
	MiniExam  *[]ExamOverview `json:"mini_exam"`
	FinalExam *ExamOverview   `json:"final_exam"`
}

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

type ExamActivityResponse struct {
	Info    ActivityDB  `json:"info"`
	Choices interface{} `json:"choices"`
}

type ExamResponse struct {
	Exam       ExamDB                 `json:"exam"`
	Activities []ExamActivityResponse `json:"activities"`
}
