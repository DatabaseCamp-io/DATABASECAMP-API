package response

// response.content_roadmap.go
/**
 * 	This file is a part of models, used to collect response of content roadmap
 */

import "DatabaseCamp/models/storages"

// Model of content roadmap item to prepare content roadmap response
type contentRoadmapItem struct {
	ActivityID int  `json:"activity_id"`
	IsLearned  bool `json:"is_learned"`
	Order      int  `json:"order"`
}

/**
 * This class represent content roadmap response
 */
type ContentRoadmapResponse struct {
	ContentID   int                  `json:"content_id"`
	ContentName string               `json:"content_name"`
	Items       []contentRoadmapItem `json:"items"`
}

/**
 * Constructor creates a new ContentRoadmapResponse instance
 *
 * @param contentDB					Content model from database to prepare overview response
 * @param contentActivitiesDB		Content Activity progression from database to prepare overview response
 * @param learningProgressionDB		Learning progression from database to prepare overview response
 *
 * @return 	instance of ContentRoadmapResponse
 */
func NewContentRoadmapResponse(contentDB storages.ContentDB, contentActivitiesDB []storages.ActivityDB, learningProgressionDB []storages.LearningProgressionDB) *ContentRoadmapResponse {
	response := ContentRoadmapResponse{}
	response.prepare(contentDB, contentActivitiesDB, learningProgressionDB)
	return &response
}

/**
 * Prepare content roadmap response
 *
 * @param contentDB					Content model from database to prepare overview response
 * @param contentActivitiesDB		Content Activity progression from database to prepare overview response
 * @param learningProgressionDB		Learning progression from database to prepare overview response
 */
func (c *ContentRoadmapResponse) prepare(contentDB storages.ContentDB, contentActivitiesDB []storages.ActivityDB, learningProgressionDB []storages.LearningProgressionDB) {
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

/**
 * Check learned activity
 *
 * @param activityID				Activity ID to check
 * @param learningProgressionDB		Learning progression from database to prepare overview response
 *
 * @return true if the activity is learned, false otherwise
 */
func (c *ContentRoadmapResponse) isLearnedActivity(progression []storages.LearningProgressionDB, activityID int) bool {
	for _, v := range progression {
		if v.ActivityID == activityID {
			return true
		}
	}
	return false
}
