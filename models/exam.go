package models

import "time"

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

type ExamActivity struct {
	ExamID                  int       `gorm:"column:exam_id" json:"exam_id"`
	ExamType                string    `gorm:"column:exam_type" json:"exam_type"`
	Instruction             string    `gorm:"column:instruction" json:"instruction"`
	CreatedTimestamp        time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	ActivityID              int       `gorm:"column:activity_id" json:"activity_id"`
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
