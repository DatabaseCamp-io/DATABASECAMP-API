package models

import (
	"DatabaseCamp/errs"
	"fmt"
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

type ContentDB struct {
	ID        int    `gorm:"primaryKey;column:content_id" json:"content_id"`
	GroupID   int    `gorm:"column:content_group_id" json:"content_group_id"`
	Name      string `gorm:"column:name" json:"name"`
	VideoPath string `gorm:"column:view_path" json:"view_path"`
	SlidePath string `gorm:"column:slide_path" json:"slide"`
}

type OverviewDB struct {
	GroupID     int    `gorm:"column:content_group_id" json:"group_id"`
	ContentID   int    `gorm:"column:content_id" json:"content_id"`
	ActivityID  *int   `gorm:"column:activity_id" json:"activity_id"`
	GroupName   string `gorm:"column:group_name" json:"group_name"`
	ContentName string `gorm:"column:content_name" json:"content_name"`
}

type MultipleChoiceDB struct {
	ID        int    `gorm:"primaryKey;column:multiple_choice_id" json:"multiple_choice_id"`
	Content   string `gorm:"column:content" json:"content"`
	IsCorrect bool   `gorm:"column:is_correct" json:"is_correct"`
}

type CompletionChoiceDB struct {
	ID            int    `gorm:"primaryKey;column:completion_choice_id" json:"completion_choice_id"`
	Content       string `gorm:"column:content" json:"content"`
	QuestionFirst string `gorm:"column:question_first" json:"question_first"`
	QuestionLast  string `gorm:"column:question_last" json:"question_last"`
}

type MatchingChoiceDB struct {
	ID        int    `gorm:"primaryKey;column:matching_choice_id" json:"matching_choice_id"`
	PairItem1 string `gorm:"column:pair_item1" json:"pair_item1"`
	PairItem2 string `gorm:"column:pair_item2" json:"pair_item2"`
}

type ActivityDB struct {
	ID       int    `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	TypeID   int    `gorm:"column:activity_type_id" json:"activity_type_id"`
	Order    int    `gorm:"column:activity_order" json:"activity_order"`
	Story    string `gorm:"column:story" json:"story"`
	Question string `gorm:"column:question" json:"question"`
}

type LearningProgressionDB struct {
	ID               int       `gorm:"primaryKey;column:learning_progression_id" json:"learning_progression_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	ActivityID       int       `gorm:"column:activity_id" json:"activity_id"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type ContentExamDB struct {
	ExamID     int `gorm:"primaryKey;column:exam_id" json:"exam_id"`
	GroupID    int `gorm:"primaryKey;column:content_group_id" json:"group_id"`
	ActivityID int `gorm:"primaryKey;column:activity_id" json:"activity_id"`
}

type OverviewInfo struct {
	Overview            []OverviewDB
	LearningProgression []LearningProgressionDB
	ExamResult          []ExamResultDB
	ContentExam         []ContentExamDB
}

type ExamResultDB struct {
	ID               int       `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	ActivityID       int       `gorm:"column:activity_id" json:"activity_id"`
	Score            int       `gorm:"column:score" json:"score"`
	IsPassed         bool      `gorm:"column:is_passed" json:"is_passed"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type ActivityResponse struct {
	Activity ActivityDB  `json:"activity"`
	Choice   interface{} `json:"choice"`
}

type PairItem struct {
	Item1 *string `json:"item1"`
	Item2 *string `json:"item2"`
}

type PairContent struct {
	ID      *int    `json:"completion_choice_id"`
	Content *string `json:"content"`
}

type MatchingChoiceAnswerRequest struct {
	ActivityID *int       `json:"activity_id"`
	Answer     []PairItem `json:"answer"`
}

func (r MatchingChoiceAnswerRequest) validatePairItem(pairItem PairItem) error {
	if pairItem.Item1 == nil || pairItem.Item2 == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

func (r MatchingChoiceAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของกิจกรรมในคำร้องขอ", "Activity ID Not Found")
	} else if len(r.Answer) == 0 {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	} else {
		for _, v := range r.Answer {
			e := r.validatePairItem(v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

type CompletionAnswerRequest struct {
	ActivityID *int          `json:"activity_id"`
	Answer     []PairContent `json:"answer"`
}

func (r CompletionAnswerRequest) validatePairItem(pairContent PairContent) error {
	if pairContent.Content == nil || pairContent.ID == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

func (r CompletionAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของกิจกรรมในคำร้องขอ", "Activity ID Not Found")
	} else if len(r.Answer) == 0 {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	} else {
		for _, v := range r.Answer {
			e := r.validatePairItem(v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

type AnswerResponse struct {
	ActivityID int  `json:"activity_id"`
	IsCorrect  bool `json:"is_correct"`
}

type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}

type LastedGroup struct {
	GroupID     int    `json:"group_id"`
	ContentID   int    `json:"content_id"`
	GroupName   string `json:"group_name"`
	ContentName string `json:"content_name"`
	Progress    int    `json:"progress"`
}

type ContentOverview struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	IsLasted    bool   `json:"is_lasted"`
	Progress    int    `json:"progress"`
}

type ContentGroupOverview struct {
	GroupID     int               `json:"group_id"`
	IsRecommend bool              `json:"is_recommend"`
	IsLasted    bool              `json:"is_lasted"`
	GroupName   string            `json:"group_name"`
	Progress    int               `json:"progress"`
	Contents    []ContentOverview `json:"contents"`
}

type OverviewResponse struct {
	LastedGroup          *LastedGroup           `json:"lasted_group"`
	ContentGroupOverview []ContentGroupOverview `json:"content_group_overview"`
}
