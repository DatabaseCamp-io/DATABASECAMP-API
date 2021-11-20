package repositories

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/general"
	"fmt"
)

type examRepository struct {
	database database.IDatabase
}

type IExamReader interface {
	GetExamActivity(examID int) ([]general.ExamActivityDB, error)
	GetExamOverview() ([]general.ExamDB, error)
}

type IExamTransaction interface {
	InsertExamResultTransaction(tx database.ITransaction, examResult general.ExamResultDB) (general.ExamResultDB, error)
	InsertExamResultActivityTransaction(tx database.ITransaction, examResultActivity []general.ExamResultActivityDB) ([]general.ExamResultActivityDB, error)
}

type IExamRepository interface {
	IExamReader
	IExamTransaction
}

func NewExamRepository(db database.IDatabase) examRepository {
	return examRepository{database: db}
}

func (r examRepository) GetExamOverview() ([]general.ExamDB, error) {
	exam := make([]general.ExamDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.Exam).
		Select(
			general.TableName.Exam+".exam_id AS exam_id",
			general.TableName.Exam+".type AS type",
			general.TableName.Exam+".instruction AS instruction",
			general.TableName.Exam+".created_timestamp AS created_timestamp",
			general.TableName.ContentGroup+".content_group_id AS content_group_id",
			general.TableName.ContentGroup+".name AS content_group_name",
			general.TableName.ContentGroup+".badge_id AS badge_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.ContentGroup,
			general.TableName.ContentGroup,
			general.IDName.MiniExam,
			general.TableName.Exam,
			general.IDName.Exam,
		)).
		Find(&exam).
		Error
	return exam, err
}

func (r examRepository) GetExamActivity(examID int) ([]general.ExamActivityDB, error) {
	examActivity := make([]general.ExamActivityDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.Exam).
		Select(
			general.TableName.Exam+".exam_id AS exam_id",
			general.TableName.Exam+".type AS exam_type",
			general.TableName.Exam+".instruction AS instruction",
			general.TableName.Exam+".created_timestamp AS created_timestamp",
			general.TableName.Activity+".activity_id AS activity_id",
			general.TableName.Activity+".activity_type_id AS activity_type_id",
			general.TableName.Activity+".question AS question",
			general.TableName.Activity+".story AS story",
			general.TableName.Activity+".point AS point",
			general.TableName.MatchingChoice+".pair_item1 AS pair_item1",
			general.TableName.MatchingChoice+".pair_item2 AS pair_item2",
			general.TableName.CompletionChoice+".content AS completion_choice_content",
			general.TableName.CompletionChoice+".completion_choice_id AS completion_choice_id",
			general.TableName.CompletionChoice+".question_first AS question_first",
			general.TableName.CompletionChoice+".question_last AS question_last",
			general.TableName.MultipleChoice+".multiple_choice_id AS multiple_choice_id",
			general.TableName.MultipleChoice+".content AS multiple_choice_content",
			general.TableName.MultipleChoice+".is_correct AS is_correct",
			general.TableName.ContentGroup+".content_group_id AS content_group_id",
			general.TableName.ContentGroup+".name AS content_group_name",
			general.TableName.ContentGroup+".badge_id AS badge_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.ContentExam,
			general.TableName.ContentExam,
			general.IDName.Exam,
			general.TableName.Exam,
			general.IDName.Exam,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.Activity,
			general.TableName.ContentExam,
			general.IDName.Activity,
			general.TableName.Activity,
			general.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.MatchingChoice,
			general.TableName.MatchingChoice,
			general.IDName.Activity,
			general.TableName.Activity,
			general.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.MultipleChoice,
			general.TableName.MultipleChoice,
			general.IDName.Activity,
			general.TableName.Activity,
			general.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.CompletionChoice,
			general.TableName.CompletionChoice,
			general.IDName.Activity,
			general.TableName.Activity,
			general.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.ContentGroup,
			general.TableName.ContentGroup,
			general.IDName.MiniExam,
			general.TableName.Exam,
			general.IDName.Exam,
		)).
		Where(general.TableName.Exam+"."+general.IDName.Exam+"  = ?", examID).
		Find(&examActivity).
		Error
	return examActivity, err
}

func (r examRepository) InsertExamResultTransaction(tx database.ITransaction, examResult general.ExamResultDB) (general.ExamResultDB, error) {
	err := tx.GetDB().
		Table(general.TableName.ExamResult).
		Create(&examResult).
		Error
	return examResult, err
}

func (r examRepository) InsertExamResultActivityTransaction(tx database.ITransaction, examResultActivity []general.ExamResultActivityDB) ([]general.ExamResultActivityDB, error) {
	err := tx.GetDB().
		Table(general.TableName.ExamResultActivity).
		Create(&examResultActivity).
		Error

	return examResultActivity, err
}
