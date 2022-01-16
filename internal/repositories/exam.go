package repositories

import (
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/models/entities/exam"
	"fmt"
)

type ExamRepository interface {
	GetExam(id int) (*exam.Exam, error)
	GetExams() ([]exam.Exam, error)
	GetExamActivities(id int) ([]exam.ExamActivity, error)
	GetExamResult(userID int, examResultID int) (*exam.ExamResult, error)
	GetExamResults(userID int) ([]exam.ExamResult, error)
	GetActivitiesResult(examResultID int) ([]exam.ResultActivity, error)
	SaveResult(result exam.Result) (*exam.Result, error)
}

type examRepository struct {
	db database.MysqlDB
}

func NewExamRepository(db database.MysqlDB) *examRepository {
	return &examRepository{db: db}
}

func (r examRepository) GetExam(id int) (*exam.Exam, error) {
	exam := exam.Exam{}
	err := r.db.GetDB().
		Table(ViewName.ExamInfo).
		Where(IDName.Exam+" = ?", id).
		Find(&exam).
		Error
	return &exam, err
}

func (r examRepository) GetExams() ([]exam.Exam, error) {
	exam := make([]exam.Exam, 0)
	err := r.db.GetDB().
		Table(TableName.Exam).
		Select(
			TableName.Exam+".exam_id AS exam_id",
			TableName.Exam+".type AS type",
			TableName.Exam+".instruction AS instruction",
			TableName.Exam+".created_timestamp AS created_timestamp",
			TableName.ContentGroup+".content_group_id AS content_group_id",
			TableName.ContentGroup+".name AS content_group_name",
			TableName.ContentGroup+".badge_id AS badge_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			TableName.ContentGroup,
			TableName.ContentGroup,
			IDName.MiniExam,
			TableName.Exam,
			IDName.Exam,
		)).
		Order(TableName.ContentGroup + ".content_group_id").
		Order(TableName.Exam + ".created_timestamp DESC").
		Find(&exam).
		Error
	return exam, err
}

func (r examRepository) GetExamActivities(id int) ([]exam.ExamActivity, error) {
	var activities []exam.ExamActivity
	err := r.db.GetDB().
		Table(TableName.ContentExam).
		Select(
			TableName.ContentExam+".activity_id AS activity_id",
			TableName.Activity+".activity_type_id AS activity_type_id",
		).
		Joins(fmt.Sprintf(
			"INNER JOIN %s ON %s.%s = %s.%s",
			TableName.Activity,
			TableName.ContentExam,
			IDName.Activity,
			TableName.Activity,
			IDName.Activity,
		)).
		Where(IDName.Exam+" = ?", id).
		Find(&activities).
		Error
	return activities, err
}

func (r examRepository) GetExamResult(userID int, examResultID int) (*exam.ExamResult, error) {
	result := exam.ExamResult{}
	err := r.db.GetDB().
		Table(ViewName.ExamResultSummary).
		Where(IDName.ExamResult+" = ?", examResultID).
		Where(IDName.User+" = ?", userID).
		Find(&result).
		Error
	return &result, err
}

func (r examRepository) GetExamResults(userID int) ([]exam.ExamResult, error) {
	examResults := make([]exam.ExamResult, 0)
	err := r.db.GetDB().
		Table(TableName.ExamResult).
		Select(
			TableName.ExamResult+".exam_result_id AS exam_result_id",
			TableName.ExamResult+".exam_id AS exam_id",
			TableName.ExamResult+".user_id AS user_id",
			TableName.ExamResult+".is_passed AS is_passed",
			TableName.ExamResult+".created_timestamp AS created_timestamp",
			fmt.Sprintf("COUNT(%s.score) AS score", TableName.ExamResultActivity),
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			TableName.ExamResultActivity,
			TableName.ExamResultActivity,
			IDName.ExamResult,
			TableName.ExamResult,
			IDName.ExamResult,
		)).
		Where(IDName.User+" = ?", userID).
		Group(TableName.ExamResult + ".exam_result_id").
		Find(&examResults).
		Error
	return examResults, err
}

func (r examRepository) GetActivitiesResult(examResultID int) ([]exam.ResultActivity, error) {
	activities := make([]exam.ResultActivity, 0)
	err := r.db.GetDB().
		Table(TableName.ExamResultActivity).
		Where(IDName.ExamResult+" = ?", examResultID).
		Find(&activities).
		Error
	return activities, err
}

func (r examRepository) SaveResult(result exam.Result) (*exam.Result, error) {
	tx := r.db.GetDB().Begin()

	err := tx.Table(TableName.ExamResult).Create(&result.ExamResult).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	examResultID := result.ExamResult.ID
	result.ActivitiesResult.SetExamResultID(examResultID)

	err = tx.Table(TableName.ExamResultActivity).Create(&result.ActivitiesResult).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &result, err
}
