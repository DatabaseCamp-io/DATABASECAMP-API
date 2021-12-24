package repositories

// repository.exam.go
/**
 * 	This file is a part of repositories, used to do data manipulation of exam
 */

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/storages"
	"fmt"
)

/**
 * 	This class manipulation exam data to other application
 */
type examRepository struct {
	Database database.IDatabase // Database to do database manipulation
}

/**
 * Constructor creates a new examRepository instance
 *
 * @param   db    Database to data manipulation
 *
 * @return 	instance of examRepository
 */
func NewExamRepository(db database.IDatabase) examRepository {
	return examRepository{Database: db}
}

/**
 * Get activities of the exam from the database
 *
 * @param 	examID  Exam ID for getting exam data
 *
 * @return activities of the exam
 * @return the error of getting exam
 */
func (r examRepository) GetExamActivity(examID int) ([]storages.ExamActivityDB, error) {
	examActivity := make([]storages.ExamActivityDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.Exam).
		Select(
			storages.TableName.Exam+".exam_id AS exam_id",
			storages.TableName.Exam+".type AS exam_type",
			storages.TableName.Exam+".instruction AS instruction",
			storages.TableName.Exam+".created_timestamp AS created_timestamp",
			storages.TableName.Activity+".activity_id AS activity_id",
			storages.TableName.Activity+".activity_type_id AS activity_type_id",
			storages.TableName.Activity+".question AS question",
			storages.TableName.Activity+".story AS story",
			storages.TableName.Activity+".point AS point",
			storages.TableName.MatchingChoice+".pair_item1 AS pair_item1",
			storages.TableName.MatchingChoice+".pair_item2 AS pair_item2",
			storages.TableName.CompletionChoice+".content AS completion_choice_content",
			storages.TableName.CompletionChoice+".completion_choice_id AS completion_choice_id",
			storages.TableName.CompletionChoice+".question_first AS question_first",
			storages.TableName.CompletionChoice+".question_last AS question_last",
			storages.TableName.MultipleChoice+".multiple_choice_id AS multiple_choice_id",
			storages.TableName.MultipleChoice+".content AS multiple_choice_content",
			storages.TableName.MultipleChoice+".is_correct AS is_correct",
			storages.TableName.ContentGroup+".content_group_id AS content_group_id",
			storages.TableName.ContentGroup+".name AS content_group_name",
			storages.TableName.ContentGroup+".badge_id AS badge_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.ContentExam,
			storages.TableName.ContentExam,
			storages.IDName.Exam,
			storages.TableName.Exam,
			storages.IDName.Exam,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.Activity,
			storages.TableName.ContentExam,
			storages.IDName.Activity,
			storages.TableName.Activity,
			storages.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.MatchingChoice,
			storages.TableName.MatchingChoice,
			storages.IDName.Activity,
			storages.TableName.Activity,
			storages.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.MultipleChoice,
			storages.TableName.MultipleChoice,
			storages.IDName.Activity,
			storages.TableName.Activity,
			storages.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.CompletionChoice,
			storages.TableName.CompletionChoice,
			storages.IDName.Activity,
			storages.TableName.Activity,
			storages.IDName.Activity,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.ContentGroup,
			storages.TableName.ContentGroup,
			storages.IDName.MiniExam,
			storages.TableName.Exam,
			storages.IDName.Exam,
		)).
		Where(storages.TableName.Exam+"."+storages.IDName.Exam+"  = ?", examID).
		Find(&examActivity).
		Error
	return examActivity, err
}

/**
 * Get exam overview from the database
 *
 * @return all exam in the application
 * @return the error of getting exam
 */
func (r examRepository) GetExamOverview() ([]storages.ExamDB, error) {
	exam := make([]storages.ExamDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.Exam).
		Select(
			storages.TableName.Exam+".exam_id AS exam_id",
			storages.TableName.Exam+".type AS type",
			storages.TableName.Exam+".instruction AS instruction",
			storages.TableName.Exam+".created_timestamp AS created_timestamp",
			storages.TableName.ContentGroup+".content_group_id AS content_group_id",
			storages.TableName.ContentGroup+".name AS content_group_name",
			storages.TableName.ContentGroup+".badge_id AS badge_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.ContentGroup,
			storages.TableName.ContentGroup,
			storages.IDName.MiniExam,
			storages.TableName.Exam,
			storages.IDName.Exam,
		)).
		Find(&exam).
		Error
	return exam, err
}

/**
 * Insert exam result into the database by database transaction
 *
 * @param 	tx  		Transaction model to do database transaction
 * @param 	examResult  Exam result for insert into the database
 *
 * @return inserted exam result
 * @return the error of inserting exam result
 */
func (r examRepository) InsertExamResultTransaction(tx database.ITransaction, examResult storages.ExamResultDB) (storages.ExamResultDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.ExamResult).
		Create(&examResult).
		Error
	return examResult, err
}

/**
 * Insert activities of the exam result into the database by database transaction
 *
 * @param 	tx  				Transaction model to do database transaction
 * @param 	examResultActivity  Activities of the exam result for insert into the database
 *
 * @return inserted activities of the exam result
 * @return the error of inserting activities of the exam result
 */
func (r examRepository) InsertExamResultActivityTransaction(tx database.ITransaction, examResultActivity []storages.ExamResultActivityDB) ([]storages.ExamResultActivityDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.ExamResultActivity).
		Create(&examResultActivity).
		Error

	return examResultActivity, err
}
