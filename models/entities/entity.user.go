package entities

import (
	"DatabaseCamp/models/general"
	"DatabaseCamp/utils"
	"os"
	"time"
)

// Type for changes user points mode
type ChangePointMode string

// Change mode
var Mode = struct {
	Add    ChangePointMode
	Reduce ChangePointMode
}{
	"+",
	"-",
}

// Badge data class
type Badge struct {
	ID          int    `json:"badge_id"`
	ImagePath   string `json:"icon_path"`
	Name        string `json:"name"`
	IsCollected bool   `json:"is_collect"`
}

// User Class
type User struct {
	id               int
	name             string
	email            string
	password         string
	badges           []Badge
	createdTimestamp time.Time
	updatedTimestamp time.Time
}

func (u *User) SetID(id int) {
	u.id = id
}

func (u *User) HashPassword() {
	u.password = utils.NewHelper().HashAndSalt(u.password)
}

func (u *User) SetTimestamp() {
	u.createdTimestamp = time.Now().Local()
	u.updatedTimestamp = time.Now().Local()
}

// Setter for set corrected badge
func (u *User) SetCorrectedBadges(allBadgesDB []general.BadgeDB, correctedBadgesDB []general.UserBadgeDB) {
	for _, badgeDB := range allBadgesDB {
		u.badges = append(u.badges, Badge{
			ID:          badgeDB.ID,
			ImagePath:   badgeDB.ImagePath,
			Name:        badgeDB.Name,
			IsCollected: u.isCorrectedBadge(correctedBadgesDB, badgeDB.ID),
		})
	}
}

// Private method for checking which badge is corrected
func (u *User) isCorrectedBadge(allBadgesDB []general.UserBadgeDB, badgeID int) bool {
	for _, correctedBadgeDB := range allBadgesDB {
		if badgeID == correctedBadgeDB.BadgeID {
			return true
		}
	}
	return false
}

// To UserDB model
func (u *User) ToDB() general.UserDB {
	tokenExpireHour := time.Hour * utils.NewType().ParseDuration(os.Getenv("TOKEN_EXPIRE_HOUR"))
	expiredTokenTimestamp := time.Now().Local().Add(tokenExpireHour)
	return general.UserDB{
		Name:                  u.name,
		Email:                 u.email,
		Password:              u.password,
		CreatedTimestamp:      u.createdTimestamp,
		UpdatedTimestamp:      u.updatedTimestamp,
		ExpiredTokenTimestamp: expiredTokenTimestamp,
	}
}

// Check password with hashpassword
func (u *User) IsPasswordCorrect(password string) bool {
	return utils.NewHelper().ComparePasswords(u.password, password)
}
