package badge

type UserBadge struct {
	UserID  int `gorm:"primaryKey;column:user_id" json:"user_id"`
	BadgeID int `gorm:"primaryKey;column:badge_id" json:"badge_id"`
}

type Badge struct {
	ID          int    `gorm:"primaryKey;column:badge_id"`
	ImagePath   string `gorm:"column:icon_path"`
	Name        string `gorm:"column:name"`
	IsCollected bool   `gorm:"-"`
}

func (b *Badge) isCollected(correctedBadgesDB []UserBadge) bool {
	for _, correctedBadgeDB := range correctedBadgesDB {
		if b.ID == correctedBadgeDB.BadgeID {
			return true
		}
	}

	return false
}

func NewBadges(allBadges []Badge, collectedBadges []UserBadge) []Badge {
	badges := make([]Badge, 0, len(allBadges))

	for _, badgeDB := range allBadges {
		badge := Badge{
			ID:        badgeDB.ID,
			ImagePath: badgeDB.ImagePath,
			Name:      badgeDB.Name,
		}

		badge.IsCollected = badge.isCollected(collectedBadges)

		badges = append(badges, badge)
	}

	return badges
}
