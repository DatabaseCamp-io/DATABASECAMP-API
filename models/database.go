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
	Activity            string
	ActivityType        string
	MatchingChoice      string
	CompletionChoice    string
	MultipleChoice      string
	Hint                string
	UserHint            string
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
	"Activity",
	"ActivityType",
	"MatchingChoice",
	"CompletionChoice",
	"MultipleChoice",
	"Hint",
	"UserHint",
}

var IDName = struct {
	User     string
	Activity string
	Hint     string
	Content  string
	Exam     string
	MiniExam string
}{
	"user_id",
	"activity_id",
	"hint_id",
	"content_id",
	"exam_id",
	"mini_exam_id",
}

var ViewName = struct {
	Profile string
	Ranking string
}{
	"Profile",
	"Ranking",
}
