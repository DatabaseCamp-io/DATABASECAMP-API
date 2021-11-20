package response

import "DatabaseCamp/models/general"

type contentRoadmapItem struct {
	ActivityID int  `json:"activity_id"`
	IsLearned  bool `json:"is_learned"`
	Order      int  `json:"order"`
}

type ContentRoadmapResponse struct {
	ContentID   int                  `json:"content_id"`
	ContentName string               `json:"content_name"`
	Items       []contentRoadmapItem `json:"items"`
}

func NewContentRoadmapResponse(contentDB general.ContentDB, contentActivitiesDB []general.ActivityDB, learningProgressionDB []general.LearningProgressionDB) *ContentRoadmapResponse {
	response := ContentRoadmapResponse{}
	response.prepare(contentDB, contentActivitiesDB, learningProgressionDB)
	return &response
}

func (c *ContentRoadmapResponse) prepare(contentDB general.ContentDB, contentActivitiesDB []general.ActivityDB, learningProgressionDB []general.LearningProgressionDB) {
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

func (c *ContentRoadmapResponse) isLearnedActivity(progression []general.LearningProgressionDB, activityID int) bool {
	for _, v := range progression {
		if v.ActivityID == activityID {
			return true
		}
	}
	return false
}