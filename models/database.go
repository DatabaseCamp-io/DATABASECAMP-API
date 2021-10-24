package models

var TableName = struct {
	User                string
	Content             string
	ContentGroup        string
	LearningProgression string
	Exam                string
}{
	"User",
	"Content",
	"ContentGroup",
	"LearningProgression",
	"Exam",
}

var IDName = struct {
	User string
}{
	"user_id",
}
var ViewName = struct {
	Profile string
}{
	"Profile",
}
