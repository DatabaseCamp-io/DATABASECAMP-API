package models

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

type lastedGroupOverview struct {
	GroupID     int    `json:"group_id"`
	ContentID   int    `json:"content_id"`
	GroupName   string `json:"group_name"`
	ContentName string `json:"content_name"`
	Progress    int    `json:"progress"`
}

type contentOverview struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	IsLasted    bool   `json:"is_lasted"`
	Progress    int    `json:"progress"`
}

type contentGroupOverview struct {
	GroupID     int               `json:"group_id"`
	IsRecommend bool              `json:"is_recommend"`
	IsLasted    bool              `json:"is_lasted"`
	GroupName   string            `json:"group_name"`
	Progress    int               `json:"progress"`
	Contents    []contentOverview `json:"contents"`
}

type overview struct {
	LastedGroup          *lastedGroupOverview   `json:"lasted_group"`
	ContentGroupOverview []contentGroupOverview `json:"content_group_overview"`
}

func NewOverview() *overview {
	return &overview{}
}

func (o *overview) ToResponse() *OverviewResponse {
	response := OverviewResponse{
		LastedGroup:          o.LastedGroup,
		ContentGroupOverview: o.ContentGroupOverview,
	}
	return &response
}

func (o *overview) Prepare(overviewDB []OverviewDB, learningProgressionDB []LearningProgressionDB) {
	groupMap := o.createGroupMap(overviewDB)
	activityContentIDMap := o.createContentActivityCount(overviewDB)
	userActivityCountByContentID := o.createUserActivityCountByContentID(learningProgressionDB, activityContentIDMap)
	lastedActivityID := o.getLastedActivityID(learningProgressionDB)
	lastedContentID := o.getLastedContentID(activityContentIDMap, lastedActivityID)
	for _, group := range groupMap {
		groupActivityCount := 0
		groupUserActivityCount := 0
		isGroupLasted := false
		contents := make([]contentOverview, 0)
		for _, content := range group.contents {
			groupActivityCount += len(content.activities)
			groupUserActivityCount += userActivityCountByContentID[content.id]
			isContentLasted := lastedContentID != nil && *lastedContentID == content.id
			contentProgress := o.calculateProgress(userActivityCountByContentID[content.id], len(content.activities))
			contents = append(contents, contentOverview{
				ContentID:   content.id,
				ContentName: content.name,
				IsLasted:    isContentLasted,
				Progress:    contentProgress,
			})
			if isContentLasted {
				o.LastedGroup = &lastedGroupOverview{
					GroupID:     group.id,
					ContentID:   content.id,
					GroupName:   group.name,
					ContentName: content.name,
					Progress:    contentProgress,
				}
			}
		}
		groupProgress := o.calculateProgress(groupUserActivityCount, groupActivityCount)
		o.ContentGroupOverview = append(o.ContentGroupOverview, contentGroupOverview{
			GroupID:     group.id,
			IsRecommend: false,
			IsLasted:    isGroupLasted,
			GroupName:   group.name,
			Progress:    groupProgress,
			Contents:    contents,
		})
	}
}

func (o *overview) createGroupMap(overviewDB []OverviewDB) map[int]*group {
	groupMap := map[int]*group{}
	for _, overview := range overviewDB {
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

func (o *overview) createContentActivityCount(overviewDB []OverviewDB) map[int]int {
	contentActivityCount := map[int]int{}
	for _, overview := range overviewDB {
		contentActivityCount[overview.ContentID] += 1
	}
	return contentActivityCount
}

func (o *overview) getLastedActivityID(learningProgressionDB []LearningProgressionDB) *int {
	if len(learningProgressionDB) == 0 {
		return nil
	} else {
		return &learningProgressionDB[0].ActivityID
	}
}

func (o *overview) createUserActivityCountByContentID(learningProgressionDB []LearningProgressionDB, activityContentIDMap map[int]int) map[int]int {
	userActivityCount := map[int]int{}
	for _, learningProgression := range learningProgressionDB {
		userActivityCount[activityContentIDMap[learningProgression.ActivityID]]++
	}
	return userActivityCount
}

func (o *overview) calculateProgress(progress int, total int) int {
	if total == 0 {
		return 0
	} else {
		ratio := float64(progress) / float64(total)
		return int(ratio * 100)
	}

}

func (o *overview) getLastedContentID(activityContentIDMap map[int]int, lastedActivityID *int) *int {
	if lastedActivityID == nil {
		return nil
	} else {
		contentID := activityContentIDMap[*lastedActivityID]
		return &contentID
	}
}