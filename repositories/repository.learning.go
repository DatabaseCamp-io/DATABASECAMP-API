package repositories

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/general"
)

type learningRepository struct {
	database database.IDatabase
}

type ILearningReader interface {
	GetContent(id int) (*general.ContentDB, error)
	GetOverview() ([]general.OverviewDB, error)
	GetContentExam(examType string) ([]general.ContentExamDB, error)
	GetActivity(id int) (*general.ActivityDB, error)
	GetMatchingChoice(activityID int) ([]general.MatchingChoiceDB, error)
	GetMultipleChoice(activityID int) ([]general.MultipleChoiceDB, error)
	GetCompletionChoice(activityID int) ([]general.CompletionChoiceDB, error)
	GetActivityHints(activityID int) ([]general.HintDB, error)
	GetContentActivity(contentID int) ([]general.ActivityDB, error)
}

type ILearningRepository interface {
	ILearningReader
}

func NewLearningRepository(db database.IDatabase) learningRepository {
	return learningRepository{database: db}
}

func (r learningRepository) GetContent(id int) (*general.ContentDB, error) {
	content := general.ContentDB{}
	err := r.database.GetDB().
		Table(general.TableName.Content).
		Where(general.IDName.Content+" = ?", id).
		Find(&content).
		Error
	return &content, err
}

func (r learningRepository) GetOverview() ([]general.OverviewDB, error) {
	overview := make([]general.OverviewDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.ContentGroup).
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

func (r learningRepository) GetContentExam(examType string) ([]general.ContentExamDB, error) {
	contentExam := make([]general.ContentExamDB, 0)
	db := r.database.GetDB()
	examSubquery := db.Table(general.TableName.Exam).
		Select("exam_id").
		Where("type = ?", string(examType)).
		Order("created_timestamp desc").
		Limit(1)
	err := r.database.GetDB().
		Table(general.TableName.ContentExam).
		Where("exam_id = (?)", examSubquery).
		Find(&contentExam).
		Error
	return contentExam, err
}

func (r learningRepository) GetContentActivity(contentID int) ([]general.ActivityDB, error) {
	activity := make([]general.ActivityDB, 0)

	err := r.database.GetDB().
		Table(general.TableName.Activity).
		Where(general.IDName.Content+" = ?", contentID).
		Find(&activity).
		Error

	return activity, err
}

func (r learningRepository) GetActivity(id int) (*general.ActivityDB, error) {
	activity := general.ActivityDB{}

	err := r.database.GetDB().
		Table(general.TableName.Activity).
		Where(general.IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	return &activity, err
}

func (r learningRepository) GetMatchingChoice(activityID int) ([]general.MatchingChoiceDB, error) {
	matchingChoice := make([]general.MatchingChoiceDB, 0)

	err := r.database.GetDB().
		Table(general.TableName.MatchingChoice).
		Where(general.IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	return matchingChoice, err
}

func (r learningRepository) GetMultipleChoice(activityID int) ([]general.MultipleChoiceDB, error) {
	multipleChoice := make([]general.MultipleChoiceDB, 0)

	err := r.database.GetDB().
		Table(general.TableName.MultipleChoice).
		Where(general.IDName.Activity+" = ?", activityID).
		Find(&multipleChoice).
		Error

	return multipleChoice, err
}

func (r learningRepository) GetCompletionChoice(activityID int) ([]general.CompletionChoiceDB, error) {
	completionChoice := make([]general.CompletionChoiceDB, 0)

	err := r.database.GetDB().
		Table(general.TableName.CompletionChoice).
		Where(general.IDName.Activity+" = ?", activityID).
		Find(&completionChoice).
		Error

	return completionChoice, err
}

func (r learningRepository) GetActivityHints(activityID int) ([]general.HintDB, error) {
	hints := make([]general.HintDB, 0)

	err := r.database.GetDB().
		Table(general.TableName.Hint).
		Where(general.IDName.Activity+" = ?", activityID).
		Order("level ASC").
		Find(&hints).
		Error

	return hints, err
}
