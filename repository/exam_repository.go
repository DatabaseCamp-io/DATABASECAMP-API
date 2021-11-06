package repository

import "DatabaseCamp/database"

type examRepository struct {
	database database.IDatabase
}

type IExamRepository interface {
}

func NewExamRepository(db database.IDatabase) examRepository {
	return examRepository{database: db}
}
