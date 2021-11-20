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

func NewLearningRepository(db database.IDatabase, service services.IAwsService) learningRepository {
	return learningRepository{Database: db, Service: service}
}

func (r learningRepository) GetContent(id int) (*storages.ContentDB, error) {
	content := storages.ContentDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.Content).
		Where(storages.IDName.Content+" = ?", id).
		Find(&content).
		Error
	return &content, err
}

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

func (r learningRepository) GetContentActivity(contentID int) ([]storages.ActivityDB, error) {
	activity := make([]storages.ActivityDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.Activity).
		Where(storages.IDName.Content+" = ?", contentID).
		Find(&activity).
		Error

	return activity, err
}

func (r learningRepository) GetActivity(id int) (*storages.ActivityDB, error) {
	activity := storages.ActivityDB{}

	err := r.Database.GetDB().
		Table(storages.TableName.Activity).
		Where(storages.IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	return &activity, err
}

func (r learningRepository) GetMatchingChoice(activityID int) ([]storages.MatchingChoiceDB, error) {
	matchingChoice := make([]storages.MatchingChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.MatchingChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	return matchingChoice, err
}

func (r learningRepository) GetMultipleChoice(activityID int) ([]storages.MultipleChoiceDB, error) {
	multipleChoice := make([]storages.MultipleChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.MultipleChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&multipleChoice).
		Error

	return multipleChoice, err
}

func (r learningRepository) GetCompletionChoice(activityID int) ([]storages.CompletionChoiceDB, error) {
	completionChoice := make([]storages.CompletionChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.CompletionChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&completionChoice).
		Error

	return completionChoice, err
}

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

func (r learningRepository) GetVideoFileLink(imagekey string) (string, error) {
	return r.Service.GetFileLink(imagekey)
}
