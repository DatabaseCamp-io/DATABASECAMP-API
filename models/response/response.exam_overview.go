package response

import (
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
)

type examDetailOverview struct {
	ExamID           int                            `json:"exam_id"`
	ExamType         string                         `json:"exam_type"`
	ContentGroupID   *int                           `json:"content_group_id,omitempty"`
	ContentGroupName *string                        `json:"content_group_name,omitempty"`
	CanDo            *bool                          `json:"can_do,omitempty"`
	Results          *[]entities.ExamResultOverview `json:"results,omitempty"`
}

type ExamOverviewResponse struct {
	PreExam   *examDetailOverview   `json:"pre_exam"`
	MiniExam  *[]examDetailOverview `json:"mini_exam"`
	FinalExam *examDetailOverview   `json:"final_exam"`
}

func NewExamOverviewResponse(examResultsDB []storages.ExamResultDB, examsDB []storages.ExamDB, canDoFinalExam bool) *ExamOverviewResponse {
	response := ExamOverviewResponse{}
	response.prepare(examResultsDB, examsDB, canDoFinalExam)
	return &ExamOverviewResponse{}
}

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

func (o *ExamOverviewResponse) countExamScore(examResultsDB []storages.ExamResultDB) map[int]int {
	examCountScore := map[int]int{}
	for _, examResultDB := range examResultsDB {
		examCountScore[examResultDB.ID] += examResultDB.Score
	}
	return examCountScore
}
