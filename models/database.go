package models

var TableName = struct {
	User                string
	Content             string
	ContentGroup        string
	LearningProgression string
	Exam                string
	ExamResult          string
	ContentExam         string
	UserBadge           string
	Badge               string
}{
	"User",
	"Content",
	"ContentGroup",
	"LearningProgression",
	"Exam",
	"ExamResult",
	"ContentExam",
	"UserBadge",
	"Badge",
}

var IDName = struct {
	User string
}{
	"user_id",
}

var ViewName = struct {
	Profile string
	Ranking string
}{
	"Profile",
	"Ranking",
}
