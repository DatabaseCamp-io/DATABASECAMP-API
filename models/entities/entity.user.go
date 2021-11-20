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
	ID                    int       `json:"user_id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	Password              string    `json:"password"`
	AccessToken           string    `json:"access_token"`
	Point                 int       `json:"point"`
	ActivityCount         int       `json:"activity_count"`
	Ranking               int       `json:"ranking"`
	Badges                []Badge   `json:"badges"`
	ExpiredTokenTimestamp time.Time `json:"expired_token_timestamp"`
	CreatedTimestamp      time.Time `json:"created_timestamp"`
	UpdatedTimestamp      time.Time `json:"updated_timestamp"`
}

func (u *User) HashPassword() {
	u.Password = utils.NewHelper().HashAndSalt(u.Password)
}

func (u *User) SetTimestamp() {
	u.CreatedTimestamp = time.Now().Local()
	u.UpdatedTimestamp = time.Now().Local()
}

// Setter for set corrected badge
func (u *User) SetCorrectedBadges(allBadgesDB []general.BadgeDB, correctedBadgesDB []general.UserBadgeDB) {
	for _, badgeDB := range allBadgesDB {
		u.Badges = append(u.Badges, Badge{
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
		Name:                  u.Name,
		Email:                 u.Email,
		Password:              u.Password,
		CreatedTimestamp:      u.CreatedTimestamp,
		UpdatedTimestamp:      u.CreatedTimestamp,
		ExpiredTokenTimestamp: expiredTokenTimestamp,
	}
}

// Check password with hashpassword
func (u *User) IsPasswordCorrect(password string) bool {
	return utils.NewHelper().ComparePasswords(u.Password, password)
}
