package general

import "time"

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
}

var IDName = struct {
	User       string
	Activity   string
	Hint       string
	Content    string
	Exam       string
	MiniExam   string
	Badge      string
	ExamResult string
}{
	"user_id",
	"activity_id",
	"hint_id",
	"content_id",
	"exam_id",
	"mini_exam_id",
	"badge_id",
	"exam_result_id",
}

var ViewName = struct {
	Profile string
	Ranking string
}{
	"Profile",
	"Ranking",
}

// Model mapped User table in the database
type UserDB struct {
	ID                    int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name                  string    `gorm:"column:name" json:"name"`
	Email                 string    `gorm:"column:email" json:"email"`
	Password              string    `gorm:"column:password" json:"password"`
	AccessToken           string    `gorm:"column:access_token" json:"access_token"`
	Point                 int       `gorm:"column:point" json:"point"`
	ExpiredTokenTimestamp time.Time `gorm:"column:expired_token_timestamp" json:"expired_token_timestamp"`
	CreatedTimestamp      time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	UpdatedTimestamp      time.Time `gorm:"column:updated_timestamp" json:"updated_timestamp"`
}

// Model mapped UserBage table in the database
type UserBadgeDB struct {
	UserID  int `gorm:"primaryKey;column:user_id" json:"user_id"`
	BadgeID int `gorm:"primaryKey;column:badge_id" json:"badge_id"`
}

// Model mapped Bage table in the database
type BadgeDB struct {
	ID        int    `gorm:"primaryKey;column:badge_id" json:"badge_id"`
	ImagePath string `gorm:"column:icon_path" json:"icon_path"`
	Name      string `gorm:"column:name" json:"name"`
}

// Model mapped Probile view in the database
type ProfileDB struct {
	ID               int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name             string    `gorm:"column:name" json:"name"`
	Point            int       `gorm:"column:point" json:"point"`
	ActivityCount    int       `gorm:"column:activity_count" json:"activity_count"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

// Model mapped Ranking view in the database
type RankingDB struct {
	ID      int    `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name    string `gorm:"column:name" json:"name"`
	Point   int    `gorm:"column:point" json:"point"`
	Ranking int    `gorm:"column:ranking" json:"ranking"`
}

// Model mapped LearningProgression table in the database
type LearningProgressionDB struct {
	ID               int       `gorm:"primaryKey;column:learning_progression_id" json:"learning_progression_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	ActivityID       int       `gorm:"column:activity_id" json:"activity_id"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

// Model mapped joined table in the database
// 		Table - UserBadge
// 		Table - Badge
type CorrectedBadgeDB struct {
	BadgeID int    `gorm:"column:badge_id" json:"badge_id"`
	Name    string `gorm:"column:badge_name" json:"badge_name"`
	UserID  *int   `gorm:"column:user_id" json:"user_id"`
}

// Model mapped Hint table in the database
type HintDB struct {
	ID          int    `gorm:"primaryKey;column:hint_id" json:"hint_id"`
	ActivityID  int    `gorm:"column:activity_id" json:"activity_id"`
	Content     string `gorm:"column:content" json:"content"`
	PointReduce int    `gorm:"column:point_reduce" json:"point_reduce"`
	Level       int    `gorm:"column:level" json:"level"`
}

// Model mapped UserHint table in the database
type UserHintDB struct {
	UserID           int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	HintID           int       `gorm:"primaryKey;column:hint_id" json:"hint_id"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

// Model mapped Content table in the database
type ContentDB struct {
	ID        int    `gorm:"primaryKey;column:content_id" json:"content_id"`
	GroupID   int    `gorm:"column:content_group_id" json:"content_group_id"`
	Name      string `gorm:"column:name" json:"name"`
	VideoPath string `gorm:"column:video_path" json:"video_path"`
	SlidePath string `gorm:"column:slide_path" json:"slide"`
}

// Model mapped Activity table in the database
type ActivityDB struct {
	ID        int    `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	TypeID    int    `gorm:"column:activity_type_id" json:"activity_type_id"`
	ContentID *int   `gorm:"column:content_id" json:"content_id"`
	Order     int    `gorm:"column:activity_order" json:"activity_order"`
	Story     string `gorm:"column:story" json:"story"`
	Point     int    `gorm:"column:point" json:"point"`
	Question  string `gorm:"column:question" json:"question"`
}

// Model mapped MultipleChoice table in the database
type MultipleChoiceDB struct {
	ID        int    `gorm:"primaryKey;column:multiple_choice_id" json:"multiple_choice_id"`
	Content   string `gorm:"column:content" json:"content"`
	IsCorrect bool   `gorm:"column:is_correct" json:"is_correct"`
}

// Model mapped CompletionChoice table in the database
type CompletionChoiceDB struct {
	ID            int    `gorm:"primaryKey;column:completion_choice_id" json:"completion_choice_id"`
	Content       string `gorm:"column:content" json:"content"`
	QuestionFirst string `gorm:"column:question_first" json:"question_first"`
	QuestionLast  string `gorm:"column:question_last" json:"question_last"`
}

// Model mapped MatchingChoice table in the database
type MatchingChoiceDB struct {
	ID        int    `gorm:"primaryKey;column:matching_choice_id" json:"matching_choice_id"`
	PairItem1 string `gorm:"column:pair_item1" json:"pair_item1"`
	PairItem2 string `gorm:"column:pair_item2" json:"pair_item2"`
}

// Model mapped joined table in the database
// 		Table - ContentGroup
// 		Table - Content
//		Table - Activity
type OverviewDB struct {
	GroupID     int    `gorm:"column:content_group_id" json:"group_id"`
	ContentID   int    `gorm:"column:content_id" json:"content_id"`
	ActivityID  *int   `gorm:"column:activity_id" json:"activity_id"`
	GroupName   string `gorm:"column:group_name" json:"group_name"`
	ContentName string `gorm:"column:content_name" json:"content_name"`
}

// Model mapped ContentExam table in the database
type ContentExamDB struct {
	ExamID     int `gorm:"primaryKey;column:exam_id" json:"exam_id"`
	GroupID    int `gorm:"primaryKey;column:content_group_id" json:"group_id"`
	ActivityID int `gorm:"primaryKey;column:activity_id" json:"activity_id"`
}

// Model mapped ExamResult table in the database
type ExamResultDB struct {
	ID               int       `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	ExamID           int       `gorm:"column:exam_id" json:"exam_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	Score            int       `gorm:"->;column:score" json:"score"`
	IsPassed         bool      `gorm:"column:is_passed" json:"is_passed"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

// Model mapped Exam table in the database
type ExamDB struct {
	ID               int       `gorm:"primaryKey;column:exam_id" json:"exam_id"`
	Type             string    `gorm:"column:type" json:"exam_type"`
	Instruction      string    `gorm:"column:instruction" json:"instruction"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	ContentGroupID   int       `gorm:"column:content_group_id" json:"content_group_id"`
	ContentGroupName string    `gorm:"column:content_group_name" json:"content_group_name"`
	BadgeID          int       `gorm:"column:badge_id" json:"badge_id"`
}

// Model mapped ExamResultActivity table in the database
type ExamResultActivityDB struct {
	ExamResultID int `gorm:"primaryKey;column:exam_result_id" json:"exam_result_id"`
	ActivityID   int `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	Score        int `gorm:"column:score" json:"score"`
}

// Model mapped joined table in the database
// 		Table - Exam
// 		Table - ContentExam
//		Table - ContentGroup
//		Table - Activity
//		Table - MatchingChoice
//		Table - MultipleChoice
//		Table - CompletionChoice
type ExamActivityDB struct {
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
