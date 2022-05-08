package content

import "time"

type LearningProgression struct {
	ID               int       `gorm:"primaryKey;column:learning_progression_id" json:"learning_progression_id"`
	UserID           int       `gorm:"column:user_id" json:"user_id"`
	ActivityID       int       `gorm:"column:activity_id" json:"activity_id"`
	IsCorrect        bool      `gorm:"column:is_correct" json:"is_correct"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type LearningProgressionList []LearningProgression

func (l LearningProgressionList) getLastedActivityID() *int {
	if len(l) == 0 {
		return nil
	} else {
		return &l[0].ActivityID
	}
}

func (l LearningProgressionList) createUserActivityCountByContentID(activityContentIDMap ActivityContentIDMap) map[int]int {
	userActivityCount := map[int]int{}
	for _, learningProgression := range l {
		userActivityCount[activityContentIDMap[learningProgression.ActivityID]]++
	}
	return userActivityCount
}

func (l LearningProgressionList) IsLearned(activityID int) bool {
	for _, v := range l {
		if v.ActivityID == activityID {
			return true
		}
	}
	return false
}
