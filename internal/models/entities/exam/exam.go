package exam

import (
	"time"
)

type ExamType string

const (
	POST ExamType = "POST"
	PRE  ExamType = "PRE"
	MINI ExamType = "MINI"
)

type Exam struct {
	ID               int       `gorm:"primaryKey;column:exam_id" json:"exam_id"`
	Type             string    `gorm:"column:type" json:"exam_type"`
	Instruction      string    `gorm:"column:instruction" json:"instruction"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	ContentGroupID   int       `gorm:"column:content_group_id" json:"content_group_id"`
	ContentGroupName string    `gorm:"column:content_group_name" json:"content_group_name"`
	BadgeID          int       `gorm:"column:badge_id" json:"badge_id"`
}

type Exams []Exam

func (exams Exams) GetPreExam(examResults ExamResults) *DetailOverview {
	var detail DetailOverview

	for _, exam := range exams {
		cando := !examResults.FinishedExam(exam.ID)
		results := examResults.GetResultOverview(exam.ID)

		if exam.Type == string(PRE) {
			detail = DetailOverview{
				ExamID:   exam.ID,
				ExamType: exam.Type,
				CanDo:    &cando,
				Results:  &results,
			}
			return &detail
		}

	}

	return nil
}

func (exams Exams) GetMiniExam(examResults ExamResults) *[]DetailOverview {
	found := false
	details := make([]DetailOverview, 0)

	for _, exam := range exams {
		if exam.Type == string(MINI) {
			found = true
			results := examResults.GetResultOverview(exam.ID)
			_exam := exam
			details = append(details, DetailOverview{
				ExamID:           _exam.ID,
				ExamType:         _exam.Type,
				ContentGroupID:   &_exam.ContentGroupID,
				ContentGroupName: &_exam.ContentGroupName,
				Results:          &results,
			})
		}
	}

	if !found {
		return nil
	} else {
		return &details
	}
}

func (exams Exams) GetFinalExam(examResults ExamResults, canDo bool) *DetailOverview {
	var detail DetailOverview

	for _, exam := range exams {
		if exam.Type == string(POST) {
			results := examResults.GetResultOverview(exam.ID)
			detail = DetailOverview{
				ExamID:   exam.ID,
				ExamType: exam.Type,
				CanDo:    &canDo,
				Results:  &results,
			}
			return &detail
		}
	}

	return nil
}
