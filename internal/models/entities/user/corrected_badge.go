package user

type CorrectedBadge struct {
	BadgeID int    `gorm:"column:badge_id" json:"badge_id"`
	Name    string `gorm:"column:badge_name" json:"badge_name"`
	UserID  *int   `gorm:"column:user_id" json:"user_id"`
}

type CorrectedBadges []CorrectedBadge

func (badges CorrectedBadges) CanDoFianlExam() bool {
	for _, badge := range badges {
		if badge.UserID == nil && badge.BadgeID != 3 {
			return false
		}
	}
	return true
}
