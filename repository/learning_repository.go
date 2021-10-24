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
		Find(&overview).
		Error
	return overview, err
}
