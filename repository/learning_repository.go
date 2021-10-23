package repository

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models"
)

type learningRepository struct {
	database database.IDatabase
}

type ILearningRepository interface {
	GetContent(id int) (*models.Content, error)
}

func NewLearningRepository(db database.IDatabase) learningRepository {
	return learningRepository{database: db}
}

func (r learningRepository) GetContent(id int) (*models.Content, error) {
	content := models.Content{}
	err := r.database.GetDB().
		Table(models.TableName.Content).
		Create(&content).
		Error
	return &content, err
}
