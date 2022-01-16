package activity

import "database-camp/internal/models/entities/content"

type ContentRoadmapItem struct {
	ActivityID int  `json:"activity_id"`
	IsLearned  bool `json:"is_learned"`
	Order      int  `json:"order"`
}

type Activity struct {
	ID        int    `gorm:"primaryKey;column:activity_id" json:"activity_id"`
	TypeID    int    `gorm:"column:activity_type_id" json:"activity_type_id"`
	ContentID *int   `gorm:"column:content_id" json:"content_id"`
	Order     int    `gorm:"column:activity_order" json:"activity_order"`
	Story     string `gorm:"column:story" json:"story"`
	Point     int    `gorm:"column:point" json:"point"`
	Question  string `gorm:"column:question" json:"question"`
}

func (activity Activity) IsAnswerCorrect(answer interface{}, choices Choices) bool {
	return true
}

type Activities []Activity

func (activities Activities) GetContentRoadmap(progression content.LearningProgressionList) (items []ContentRoadmapItem) {
	for _, activity := range activities {
		isLearned := progression.IsLearned(activity.ID)
		items = append(items, ContentRoadmapItem{
			ActivityID: activity.ID,
			IsLearned:  isLearned,
			Order:      activity.Order,
		})
	}
	return
}
