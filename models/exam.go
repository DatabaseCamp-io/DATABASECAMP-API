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

type ExamResultDB struct {
	ID               int       `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	ExamID           int       `gorm:"column:exam_id" json:"exam_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	Score            int       `gorm:"->;column:score" json:"score"`
	IsPassed         bool      `gorm:"column:is_passed" json:"is_passed"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type ExamActivity struct {
	ExamID                  int       `gorm:"column:exam_id" json:"exam_id"`
	ExamType                string    `gorm:"column:exam_type" json:"exam_type"`
	Instruction             string    `gorm:"column:instruction" json:"instruction"`
	CreatedTimestamp        time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	ActivityID              int       `gorm:"column:activity_id" json:"activity_id"`
	Point                   int       `gorm:"column:point" json:"point"`
	ActivityTypeID          int       `gorm:"column:activity_type_id" json:"activity_type_id"`
	Question                string    `gorm:"column:question" json:"question"`
	Story                   string    `gorm:"column:story" json:"story"`
	PairItem1               string    `gorm:"column:pair_item1" json:"pair_item1"`
	PairItem2               string    `gorm:"column:pair_item2" json:"pair_item2"`
	CompletionChoiceID      int       `gorm:"column:completion_choice_id" json:"completion_choice_id"`
	CompletionChoiceContent string    `gorm:"column:completion_choice_content" json:"completion_choice_content"`
	QuestionFirst           string    `gorm:"column:question_first" json:"question_first"`
	QuestionLast            string    `gorm:"column:question_last" json:"question_last"`
	MultipleChoiceID        int       `gorm:"column:multiple_choice_id" json:"multiple_choice_id"`
	MultipleChoiceContent   string    `gorm:"column:multiple_choice_content" json:"multiple_choice_content"`
	IsCorrect               bool      `gorm:"column:is_correct" json:"is_correct"`
	Content                 string    `json:"content"`
	ContentGroupID          int       `gorm:"column:content_group_id" json:"content_group_id"`
	ContentGroupName        string    `gorm:"column:content_group_name" json:"content_group_name"`
	BadgeID                 int       `gorm:"column:badge_id" json:"badge_id"`
}

type ExamDB struct {
	ID               int       `gorm:"primaryKey;column:exam_id" json:"exam_id"`
	Type             string    `gorm:"column:type" json:"exam_type"`
	Instruction      string    `gorm:"column:instruction" json:"instruction"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	ContentGroupID   int       `gorm:"column:content_group_id" json:"content_group_id"`
	ContentGroupName string    `gorm:"column:content_group_name" json:"content_group_name"`
	BadgeID          int       `gorm:"column:badge_id" json:"badge_id"`
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

type ExamResultActivityDB struct {
	ExamResultID int `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	ActivityID   int `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	Score        int `gorm:"column:score" json:"score"`
}
