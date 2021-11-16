package models

type contentRoadmapItem struct {
	ActivityID int  `json:"activity_id"`
	IsLearned  bool `json:"is_learned"`
	Order      int  `json:"order"`
}

type contentRoadmap struct {
	ContentID   int                  `json:"content_id"`
	ContentName string               `json:"content_name"`
	Items       []contentRoadmapItem `json:"items"`
}

func NewContentRoadmap() *contentRoadmap {
	return &contentRoadmap{}
}

func (c *contentRoadmap) ToResponse() *ContentRoadmapResponse {
	response := ContentRoadmapResponse{
		ContentID:   c.ContentID,
		ContentName: c.ContentName,
		Items:       c.Items,
	}
	return &response
}

func (c *contentRoadmap) Prepare(contentDB ContentDB, contentActivitiesDB []ActivityDB, learningProgressionDB []LearningProgressionDB) {
	c.ContentID = contentDB.ID
	c.ContentName = contentDB.Name
	for _, activity := range contentActivitiesDB {
		isLearned := c.isLearnedActivity(learningProgressionDB, activity.ID)
		c.Items = append(c.Items, contentRoadmapItem{
			ActivityID: activity.ID,
			IsLearned:  isLearned,
			Order:      activity.Order,
		})
	}
}

func (c *contentRoadmap) isLearnedActivity(progression []LearningProgressionDB, activityID int) bool {
	for _, v := range progression {
		if v.ActivityID == activityID {
			return true
		}
	}
	return false
}
