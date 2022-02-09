package repositories

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
	ExamResultActivity  string
	VocabGroupChoice    string
	VocabGroup          string
	DependencyChoice    string
	Dependency          string
	Determinant         string
	ERChoice            string
	ERChoiceTables      string
	Tables              string
	Attributes          string
	Relationship        string
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
	"ExamResultActivity",
	"VocabGroupChoice",
	"VocabGroup",
	"DependencyChoice",
	"Dependency",
	"Determinant",
	"ERChoice",
	"ERChoiceTables",
	"Tables",
	"Attributes",
	"Relationship",
}

var IDName = struct {
	User             string
	Activity         string
	Hint             string
	Content          string
	ContentGroup     string
	Exam             string
	MiniExam         string
	Badge            string
	ExamResult       string
	VocabGroup       string
	Dependency       string
	Determinant      string
	DependencyChoice string
	ERChoice         string
	Table            string
	Attribute        string
	Relationship     string
}{
	"user_id",
	"activity_id",
	"hint_id",
	"content_id",
	"content_group_id",
	"exam_id",
	"mini_exam_id",
	"badge_id",
	"exam_result_id",
	"vocab_group_id",
	"dependency_id",
	"determinant_id",
	"dependency_choice_id",
	"er_choice_id",
	"table_id",
	"attribute_id",
	"relationship_id",
}

var ViewName = struct {
	Profile           string
	Ranking           string
	ExamInfo          string
	ExamResultSummary string
	UserPreTest       string
	UserPreTestResult string
	SpiderData        string
	RandomERAnswer    string
}{
	"Profile",
	"Ranking",
	"ExamInfo",
	"ExamResultSummary",
	"UserPreTest",
	"UserPreTestResult",
	"SpiderData",
	"RandomERAnswer",
}
