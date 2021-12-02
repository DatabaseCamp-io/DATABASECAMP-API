package repositories

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/services"
)

type learningRepository struct {
	Database database.IDatabase
	Service  services.IAwsService
}

//Interface that show how others function call and use function in module learning respository
type ILearningRepository interface {
	GetContent(id int) (*storages.ContentDB, error)
	GetOverview() ([]storages.OverviewDB, error)
	GetContentExam(examType string) ([]storages.ContentExamDB, error)
	GetActivity(id int) (*storages.ActivityDB, error)
	GetMatchingChoice(activityID int) ([]storages.MatchingChoiceDB, error)
	GetMultipleChoice(activityID int) ([]storages.MultipleChoiceDB, error)
	GetCompletionChoice(activityID int) ([]storages.CompletionChoiceDB, error)
	GetActivityHints(activityID int) ([]storages.HintDB, error)
	GetContentActivity(contentID int) ([]storages.ActivityDB, error)
	GetVideoFileLink(imagekey string) (string, error)
}

// Create learning repository in database
func NewLearningRepository(db database.IDatabase, service services.IAwsService) learningRepository {
	return learningRepository{Database: db, Service: service}
}

// Get content from database
func (r learningRepository) GetContent(id int) (*storages.ContentDB, error) {
	content := storages.ContentDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.Content).
		Where(storages.IDName.Content+" = ?", id).
		Find(&content).
		Error
	return &content, err
}

// Get overview from database
func (r learningRepository) GetOverview() ([]storages.OverviewDB, error) {
	overview := make([]storages.OverviewDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.ContentGroup).
		Select("ContentGroup.content_group_id AS content_group_id",
			"Content.content_id AS content_id",
			"Activity.activity_id AS activity_id",
			"ContentGroup.name AS group_name",
			"Content.name AS content_name",
		).
		Joins("LEFT JOIN Content ON ContentGroup.content_group_id = Content.content_group_id").
		Joins("LEFT JOIN Activity ON Content.content_id = Activity.content_id").
		Order("content_group_id ASC").
		Find(&overview).
		Error
	return overview, err
}

// Get exam content from database
func (r learningRepository) GetContentExam(examType string) ([]storages.ContentExamDB, error) {
	contentExam := make([]storages.ContentExamDB, 0)
	db := r.Database.GetDB()
	examSubquery := db.Table(storages.TableName.Exam).
		Select("exam_id").
		Where("type = ?", string(examType)).
		Order("created_timestamp desc").
		Limit(1)
	err := r.Database.GetDB().
		Table(storages.TableName.ContentExam).
		Where("exam_id = (?)", examSubquery).
		Find(&contentExam).
		Error
	return contentExam, err
}

// Get activity content from database
func (r learningRepository) GetContentActivity(contentID int) ([]storages.ActivityDB, error) {
	activity := make([]storages.ActivityDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.Activity).
		Where(storages.IDName.Content+" = ?", contentID).
		Find(&activity).
		Error

	return activity, err
}

// Get activity from database
func (r learningRepository) GetActivity(id int) (*storages.ActivityDB, error) {
	activity := storages.ActivityDB{}

	err := r.Database.GetDB().
		Table(storages.TableName.Activity).
		Where(storages.IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	return &activity, err
}

// Get matching choice activity from database
func (r learningRepository) GetMatchingChoice(activityID int) ([]storages.MatchingChoiceDB, error) {
	matchingChoice := make([]storages.MatchingChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.MatchingChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	return matchingChoice, err
}

// Get multiple choice activity from database
func (r learningRepository) GetMultipleChoice(activityID int) ([]storages.MultipleChoiceDB, error) {
	multipleChoice := make([]storages.MultipleChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.MultipleChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&multipleChoice).
		Error

	return multipleChoice, err
}

// Get completion choice activity from database
func (r learningRepository) GetCompletionChoice(activityID int) ([]storages.CompletionChoiceDB, error) {
	completionChoice := make([]storages.CompletionChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.CompletionChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&completionChoice).
		Error

	return completionChoice, err
}

// Get activity hint from database
func (r learningRepository) GetActivityHints(activityID int) ([]storages.HintDB, error) {
	hints := make([]storages.HintDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.Hint).
		Where(storages.IDName.Activity+" = ?", activityID).
		Order("level ASC").
		Find(&hints).
		Error

	return hints, err
}

// Get video link from AWS service
func (r learningRepository) GetVideoFileLink(imagekey string) (string, error) {
	return r.Service.GetFileLink(imagekey)
}
