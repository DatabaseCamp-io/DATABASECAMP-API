package repository

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models"
	"fmt"
)

type examRepository struct {
	database database.IDatabase
}

type IExamRepository interface {
	GetExamActivity(examID int) ([]models.ExamActivity, error)
	GetExamOverview() ([]models.ExamDB, error)
}

func NewExamRepository(db database.IDatabase) examRepository {
	return examRepository{database: db}
}

func (r examRepository) GetExamOverview() ([]models.ExamDB, error) {
	exam := make([]models.ExamDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.Exam).
		Select(
			models.TableName.Exam+".exam_id AS exam_id",
			models.TableName.Exam+".type AS type",
			models.TableName.Exam+".instruction AS instruction",
			models.TableName.Exam+".created_timestamp AS created_timestamp",
			models.TableName.ContentGroup+".content_group_id AS content_group_id",
			models.TableName.ContentGroup+".name AS content_group_name",
			models.TableName.ContentGroup+".badge_id AS badge_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.ContentGroup,
			models.TableName.ContentGroup,
			models.IDName.MiniExam,
			models.TableName.Exam,
			models.IDName.Exam,
		)).
		Find(&exam).
		Error
	return exam, err
}

func (r examRepository) GetExamActivity(examID int) ([]models.ExamActivity, error) {
	examActivity := make([]models.ExamActivity, 0)
	err := r.database.GetDB().
		Table(models.TableName.Exam).
		Select(
			models.TableName.Exam+".exam_id AS exam_id",
			models.TableName.Exam+".type AS exam_type",
			models.TableName.Exam+".instruction AS instruction",
			models.TableName.Exam+".created_timestamp AS created_timestamp",
			models.TableName.Activity+".activity_id AS activity_id",
			models.TableName.Activity+".activity_type_id AS activity_type_id",
			models.TableName.Activity+".question AS question",
			models.TableName.Activity+".story AS story",
			models.TableName.MatchingChoice+".pair_item1 AS pair_item1",
			models.TableName.MatchingChoice+".pair_item2 AS pair_item2",
			models.TableName.CompletionChoice+".content AS completion_choice_content",
			models.TableName.CompletionChoice+".completion_choice_id AS completion_choice_id",
			models.TableName.CompletionChoice+".question_first AS question_first",
			models.TableName.CompletionChoice+".question_last AS question_last",
			models.TableName.MultipleChoice+".multiple_choice_id AS multiple_choice_id",
			models.TableName.MultipleChoice+".content AS multiple_choice_content",
			models.TableName.MultipleChoice+".is_correct AS is_correct",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.ContentExam,
			models.TableName.ContentExam,
			models.IDName.Exam,
			models.TableName.Exam,
			models.IDName.Exam,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.Activity,
			models.TableName.ContentExam,
			models.IDName.Activity,
			models.TableName.Activity,
			models.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.MatchingChoice,
			models.TableName.MatchingChoice,
			models.IDName.Activity,
			models.TableName.Activity,
			models.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.MultipleChoice,
			models.TableName.MultipleChoice,
			models.IDName.Activity,
			models.TableName.Activity,
			models.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.CompletionChoice,
			models.TableName.CompletionChoice,
			models.IDName.Activity,
			models.TableName.Activity,
			models.IDName.Activity,
		)).
		Where(models.TableName.Exam+"."+models.IDName.Exam+"  = ?", examID).
		Find(&examActivity).
		Error
	return examActivity, err
}
