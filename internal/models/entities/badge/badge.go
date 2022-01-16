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

func (b *Badge) setIsCollected(correctedBadgesDB []UserBadge) {
	for _, correctedBadgeDB := range correctedBadgesDB {
		if b.ID == correctedBadgeDB.BadgeID {
			b.IsCollected = true
		}
	}

	b.IsCollected = false
}

func NewBadges(allBadgesDB []Badge, correctedBadgesDB []UserBadge) []Badge {
	badges := make([]Badge, 0, len(allBadgesDB))

	for _, badgeDB := range allBadgesDB {
		badge := Badge{
			ID:        badgeDB.ID,
			ImagePath: badgeDB.ImagePath,
			Name:      badgeDB.Name,
		}

		badge.setIsCollected(correctedBadgesDB)

		badges = append(badges)
	}

	return badges
}
