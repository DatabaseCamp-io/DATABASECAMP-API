package content

type content struct {
	id         int
	name       string
	activities []int
}

type group struct {
	id       int
	name     string
	contents map[int]*content
}

type contentOverview struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	IsLasted    bool   `json:"is_lasted"`
	Progress    int    `json:"progress"`
}

type LastedGroupOverview struct {
	GroupID     int    `json:"group_id"`
	ContentID   int    `json:"content_id"`
	ActivityID  int    `json:"activity_id"`
	GroupName   string `json:"group_name"`
	ContentName string `json:"content_name"`
	Progress    int    `json:"progress"`
}

type ContentGroupOverview struct {
	GroupID   int               `json:"group_id"`
	IsLasted  bool              `json:"is_lasted"`
	GroupName string            `json:"group_name"`
	Progress  int               `json:"progress"`
	Contents  []contentOverview `json:"contents"`
}

type Overview struct {
	GroupID     int    `gorm:"column:content_group_id" json:"group_id"`
	ContentID   int    `gorm:"column:content_id" json:"content_id"`
	ActivityID  *int   `gorm:"column:activity_id" json:"activity_id"`
	GroupName   string `gorm:"column:group_name" json:"group_name"`
	ContentName string `gorm:"column:content_name" json:"content_name"`
}

type ActivityContentIDMap map[int]int

func (m ActivityContentIDMap) getLastedContentID(lastedActivityID *int) *int {
	if lastedActivityID == nil {
		return nil
	} else {
		contentID := m[*lastedActivityID]
		return &contentID
	}
}

type OverviewList []Overview

func (l OverviewList) GetLearningOverview(progressionList LearningProgressionList) (*LastedGroupOverview, []ContentGroupOverview) {

	var lastedGroupOverview *LastedGroupOverview
	var contentGroupOverview []ContentGroupOverview

	groupMap := l.createGroupMap()
	activityContentIDMap := l.createActivityContentIDMap()
	userActivityCountByContentID := progressionList.createUserActivityCountByContentID(activityContentIDMap)
	lastedActivityID := progressionList.getLastedActivityID()
	lastedContentID := activityContentIDMap.getLastedContentID(lastedActivityID)

	for _, group := range groupMap {
		groupActivityCount := 0
		groupUserActivityCount := 0
		isGroupLasted := false
		contents := make([]contentOverview, 0)
		for _, content := range group.contents {
			groupActivityCount += len(content.activities)
			groupUserActivityCount += userActivityCountByContentID[content.id]
			isContentLasted := lastedContentID != nil && *lastedContentID == content.id
			contentProgress := calculateProgress(userActivityCountByContentID[content.id], len(content.activities))
			contents = append(contents, contentOverview{
				ContentID:   content.id,
				ContentName: content.name,
				IsLasted:    isContentLasted,
				Progress:    contentProgress,
			})
			if isContentLasted {
				lastedGroupOverview = &LastedGroupOverview{
					GroupID:     group.id,
					ContentID:   content.id,
					GroupName:   group.name,
					ActivityID:  *lastedActivityID,
					ContentName: content.name,
					Progress:    contentProgress,
				}
			}
		}
		groupProgress := calculateProgress(groupUserActivityCount, groupActivityCount)
		contentGroupOverview = append(contentGroupOverview, ContentGroupOverview{
			GroupID:   group.id,
			IsLasted:  isGroupLasted,
			GroupName: group.name,
			Progress:  groupProgress,
			Contents:  contents,
		})
	}

	return lastedGroupOverview, contentGroupOverview
}

func (l OverviewList) createGroupMap() map[int]*group {
	groupMap := map[int]*group{}
	for _, overview := range l {
		_group := groupMap[overview.GroupID]
		if _group == nil {
			_group = &group{
				id:       overview.GroupID,
				name:     overview.GroupName,
				contents: map[int]*content{},
			}
			groupMap[overview.GroupID] = _group
		}

		_content := _group.contents[overview.ContentID]
		if _content == nil {
			_content = &content{
				id:         overview.ContentID,
				name:       overview.ContentName,
				activities: []int{},
			}
			_group.contents[overview.ContentID] = _content
		}

		if overview.ActivityID != nil {
			_content.activities = append(_content.activities, *overview.ActivityID)
		}
	}

	return groupMap
}

func (l OverviewList) createActivityContentIDMap() ActivityContentIDMap {
	activityContentIDMap := map[int]int{}
	for _, overview := range l {
		if overview.ActivityID != nil {
			activityContentIDMap[*overview.ActivityID] = overview.ContentID
		}
	}
	return activityContentIDMap
}

func calculateProgress(progress int, total int) int {
	if total == 0 {
		return 0
	} else {
		ratio := float64(progress) / float64(total)
		return int(ratio * 100)
	}

}
