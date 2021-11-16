package models

type examDetailOverview struct {
	ExamID           int                   `json:"exam_id"`
	ExamType         string                `json:"exam_type"`
	ContentGroupID   *int                  `json:"content_group_id,omitempty"`
	ContentGroupName *string               `json:"content_group_name,omitempty"`
	CanDo            *bool                 `json:"can_do,omitempty"`
	Results          *[]examResultOverview `json:"results,omitempty"`
}

type examOverview struct {
	PreExam   *examDetailOverview   `json:"pre_exam"`
	MiniExam  *[]examDetailOverview `json:"mini_exam"`
	FinalExam *examDetailOverview   `json:"final_exam"`
}

func NewExamOverview() *examOverview {
	return &examOverview{}
}

func (o *examOverview) ToResponse() *ExamOverviewResponse {
	return &ExamOverviewResponse{
		PreExam:   o.PreExam,
		MiniExam:  o.MiniExam,
		FinalExam: o.FinalExam,
	}
}

func (o *examOverview) PrepareExamOverview(examResultsDB []ExamResultDB, correctedBadgesDB []CorrectedBadgeDB, examsDB []ExamDB) {
	examResultMap := o.createExamResultMap(examResultsDB)
	for _, examDB := range examsDB {
		if examDB.Type == string(Exam.Pretest) {
			o.PreExam = &examDetailOverview{
				ExamID:   examDB.ID,
				ExamType: examDB.Type,
				Results:  examResultMap[examDB.ID],
			}
		} else if examDB.Type == string(Exam.MiniExam) {
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
		} else if examDB.Type == string(Exam.Posttest) {
			cando := o.canDoFianlExam(correctedBadgesDB)
			o.FinalExam = &examDetailOverview{
				ExamID:   examDB.ID,
				ExamType: examDB.Type,
				CanDo:    &cando,
				Results:  examResultMap[examDB.ID],
			}
		}
	}
}

func (o *examOverview) createExamResultMap(examResultsDB []ExamResultDB) map[int]*[]examResultOverview {
	examResultMap := map[int]*[]examResultOverview{}
	examScoreCount := o.countExamScore(examResultsDB)
	for _, examResult := range examResultsDB {
		if examResultMap[examResult.ExamID] == nil {
			temp := make([]examResultOverview, 0)
			examResultMap[examResult.ExamID] = &temp
		}
		*examResultMap[examResult.ExamID] = append(*examResultMap[examResult.ExamID], examResultOverview{
			ExamResultID:     examResult.ID,
			TotalScore:       examScoreCount[examResult.ID],
			IsPassed:         examResult.IsPassed,
			CreatedTimestamp: examResult.CreatedTimestamp,
		})
	}
	return examResultMap
}

func (o *examOverview) countExamScore(examResultsDB []ExamResultDB) map[int]int {
	examCountScore := map[int]int{}
	for _, examResultDB := range examResultsDB {
		examCountScore[examResultDB.ID] += examResultDB.Score
	}
	return examCountScore
}

func (o *examOverview) canDoFianlExam(correctedBadgesDB []CorrectedBadgeDB) bool {
	for _, correctedBadgeDB := range correctedBadgesDB {
		if correctedBadgeDB.UserID == nil {
			return false
		}
	}
	return true
}
