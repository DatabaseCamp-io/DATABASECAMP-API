package user

import "database-camp/internal/models/entities/content"

type RecommendGroup struct {
	ContentGroupID int  `json:"content_group_id"`
	IsRecommend    bool `json:"is_recommend"`
}

type PreTestResult struct {
	ActivityID     int `gorm:"column:activity_id"`
	Score          int `gorm:"column:score"`
	Point          int `gorm:"column:point"`
	ContentGroupID int `gorm:"column:content_group_id"`
}

func (result PreTestResult) IsRecommend() bool {
	if result.Point == 0 {
		return false
	}

	return float64(result.Score)/float64(result.Point) < 0.3
}

type PreTestResultsMap map[int]*PreTestResult

func NewPreTestResultsMap(results PreTestResults) (resultMap PreTestResultsMap) {
	resultMap = make(PreTestResultsMap)
	for _, result := range results {
		resultMap[result.ContentGroupID] = &result
	}
	return
}

func (resultsMap PreTestResultsMap) IsRecommend(groupID int) bool {
	if resultsMap[groupID] == nil {
		return true
	} else {
		return resultsMap[groupID].IsRecommend()
	}
}

type PreTestResults []PreTestResult

func (results PreTestResults) GetRecommend(contentGroups content.ContentGroups) []RecommendGroup {

	resultsMap := NewPreTestResultsMap(results)

	recommend := make([]RecommendGroup, 0, len(contentGroups))
	for _, group := range contentGroups {
		recommend = append(recommend, RecommendGroup{
			ContentGroupID: group.ID,
			IsRecommend:    resultsMap.IsRecommend(group.ID),
		})
	}

	return recommend
}
