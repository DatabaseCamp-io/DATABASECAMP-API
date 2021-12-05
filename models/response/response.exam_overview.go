package response

// response.exam_overview.go
/**
 * 	This file is a part of models, used to collect response of exam overview
 */

import (
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
)

// Model of exam overview item to prepare exam overview response
type examDetailOverview struct {
	ExamID           int                            `json:"exam_id"`
	ExamType         string                         `json:"exam_type"`
	ContentGroupID   *int                           `json:"content_group_id,omitempty"`
	ContentGroupName *string                        `json:"content_group_name,omitempty"`
	CanDo            *bool                          `json:"can_do,omitempty"`
	Results          *[]entities.ExamResultOverview `json:"results,omitempty"`
}

/**
 * This class represent exam overview response
 */
type ExamOverviewResponse struct {
	PreExam   *examDetailOverview   `json:"pre_exam"`
	MiniExam  *[]examDetailOverview `json:"mini_exam"`
	FinalExam *examDetailOverview   `json:"final_exam"`
}

/**
 * Constructor creates a new ExamOverviewResponse instance
 *
 * @param examResultsDB		Exam result from database to prepare exam overview response
 * @param examsDB			Exam from database to prepare exam overview response
 * @param canDoFinalExam	Boolean that indicate user can do the final exam
 *
 * @return 	instance of ExamOverviewResponse
 */
func NewExamOverviewResponse(examResultsDB []storages.ExamResultDB, examsDB []storages.ExamDB, canDoFinalExam bool) *ExamOverviewResponse {
	response := ExamOverviewResponse{}
	response.prepare(examResultsDB, examsDB, canDoFinalExam)
	return &response
}

/**
* Prepare exam overview response
*
* @param examResultsDB		Exam result from database to prepare exam overview response
* @param examsDB			Exam from database to prepare exam overview response
* @param canDoFinalExam		Boolean that indicate user can do the final exam
 */
func (o *ExamOverviewResponse) prepare(examResultsDB []storages.ExamResultDB, examsDB []storages.ExamDB, canDoFinalExam bool) {
	examResultMap := o.createExamResultMap(examResultsDB)
	for _, examDB := range examsDB {
		if examDB.Type == entities.ExamType.Pretest {
			o.PreExam = &examDetailOverview{
				ExamID:   examDB.ID,
				ExamType: examDB.Type,
				Results:  examResultMap[examDB.ID],
			}
		} else if examDB.Type == entities.ExamType.MiniExam {
			if o.MiniExam == nil {
				temp := make([]examDetailOverview, 0)
				o.MiniExam = &temp
			}
			contentGroupID := examDB.ContentGroupID
			contentGroupName := examDB.ContentGroupName
			*o.MiniExam = append(*o.MiniExam, examDetailOverview{
				ExamID:           examDB.ID,
				ExamType:         examDB.Type,
				ContentGroupID:   &contentGroupID,
				ContentGroupName: &contentGroupName,
				Results:          examResultMap[examDB.ID],
			})
		} else if examDB.Type == entities.ExamType.Posttest {
			o.FinalExam = &examDetailOverview{
				ExamID:   examDB.ID,
				ExamType: examDB.Type,
				CanDo:    &canDoFinalExam,
				Results:  examResultMap[examDB.ID],
			}
		}
	}
}

/**
 * Create map of exam result
 * [key]	exam id
 * [value]	exam result overview
 *
 * @param examResultsDB	 exam result from database to create exam result map
 *
 * @return map of exam result
 */
func (o *ExamOverviewResponse) createExamResultMap(examResultsDB []storages.ExamResultDB) map[int]*[]entities.ExamResultOverview {
	examResultMap := map[int]*[]entities.ExamResultOverview{}
	examScoreCount := o.countExamScore(examResultsDB)
	for _, examResult := range examResultsDB {
		if examResultMap[examResult.ExamID] == nil {
			temp := make([]entities.ExamResultOverview, 0)
			examResultMap[examResult.ExamID] = &temp
		}
		*examResultMap[examResult.ExamID] = append(*examResultMap[examResult.ExamID], entities.ExamResultOverview{
			ExamResultID:     examResult.ID,
			TotalScore:       examScoreCount[examResult.ID],
			IsPassed:         examResult.IsPassed,
			CreatedTimestamp: examResult.CreatedTimestamp,
		})
	}
	return examResultMap
}

/**
 * Create map of count exam score
 * [key]  	exam result id
 * [value]	score of the exam result
 *
 * @param examResultsDB	 exam result from database to calculate exam score
 *
 * @return exam score
 */
func (o *ExamOverviewResponse) countExamScore(examResultsDB []storages.ExamResultDB) map[int]int {
	examCountScore := map[int]int{}
	for _, examResultDB := range examResultsDB {
		examCountScore[examResultDB.ID] += examResultDB.Score
	}
	return examCountScore
}
