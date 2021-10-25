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

type LearningProgressionDB struct {
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
