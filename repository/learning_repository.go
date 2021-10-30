package repository

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models"
)

type learningRepository struct {
	database database.IDatabase
}

type ILearningRepository interface {
	GetContent(id int) (*models.ContentDB, error)
	GetOverview() ([]models.OverviewDB, error)
	GetContentExam(examType models.ExamType) ([]models.ContentExamDB, error)
	GetActivity(id int) (*models.ActivityDB, error)
	GetMatchingChoice(activityID int) ([]models.MatchingChoiceDB, error)
}

func NewLearningRepository(db database.IDatabase) learningRepository {
	return learningRepository{database: db}
}

func (r learningRepository) GetContent(id int) (*models.ContentDB, error) {
	content := models.ContentDB{}
	err := r.database.GetDB().
		Table(models.TableName.Content).
		Create(&content).
		Error
	return &content, err
}

func (r learningRepository) GetOverview() ([]models.OverviewDB, error) {
	overview := make([]models.OverviewDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.ContentGroup).
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

func (r learningRepository) GetContentExam(examType models.ExamType) ([]models.ContentExamDB, error) {
	contentExam := make([]models.ContentExamDB, 0)
	db := r.database.GetDB()
	examSubquery := db.Table(models.TableName.Exam).
		Select("exam_id").
		Where("type = ?", string(examType)).
		Order("created_timestamp desc").
		Limit(1)
	err := r.database.GetDB().
		Table(models.TableName.ContentExam).
		Where("exam_id = (?)", examSubquery).
		Find(&contentExam).
		Error
	return contentExam, err
}

func (r learningRepository) GetActivity(id int) (*models.ActivityDB, error) {
	activity := models.ActivityDB{}

	err := r.database.GetDB().
		Table(models.TableName.Activity).
		Where(models.IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	return &activity, err
}

func (r learningRepository) GetMatchingChoice(activityID int) ([]models.MatchingChoiceDB, error) {
	matchingChoice := make([]models.MatchingChoiceDB, 0)

	err := r.database.GetDB().
		Table(models.TableName.MatchingChoice).
		Where(models.IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	return matchingChoice, err
}
