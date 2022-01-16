package activity

import "time"

type HintRoadMap struct {
	Level       int `json:"level"`
	ReducePoint int `json:"reduce_point"`
}

type ActivityHint struct {
	TotalHint   int           `json:"total_hint"`
	UsedHints   []Hint        `json:"used_hints"`
	HintRoadMap []HintRoadMap `json:"hint_roadmap"`
}

type UserHint struct {
	UserID           int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	HintID           int       `gorm:"primaryKey;column:hint_id" json:"hint_id"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type UserHints []UserHint

func (hints UserHints) IsUsed(hintID int) bool {
	for _, hint := range hints {
		if hint.HintID == hintID {
			return true
		}
	}
	return false
}

type Hint struct {
	ID          int    `gorm:"primaryKey;column:hint_id" json:"hint_id"`
	ActivityID  int    `gorm:"column:activity_id" json:"activity_id"`
	Content     string `gorm:"column:content" json:"content"`
	PointReduce int    `gorm:"column:point_reduce" json:"point_reduce"`
	Level       int    `gorm:"column:level" json:"level"`
}

type Hints []Hint

func (hints Hints) GetUsedHints(userHints UserHints) (usedHints Hints) {
	for _, hint := range hints {
		if userHints.IsUsed(hint.ID) {
			usedHints = append(usedHints, hint)
		}
	}
	return
}

func (hints Hints) CreateRoadmap() (roadmap []HintRoadMap) {
	for _, hint := range hints {
		roadmap = append(roadmap, HintRoadMap{
			Level:       hint.Level,
			ReducePoint: hint.PointReduce,
		})
	}
	return
}

func (hints Hints) GetNextLevelHint(userHints UserHints) *Hint {
	for _, hint := range hints {
		if !userHints.IsUsed(hint.ID) {
			return &hint
		}
	}
	return nil
}
