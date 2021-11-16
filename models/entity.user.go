package models

import (
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

// badge data class
type badge struct {
	ID          int    `json:"badge_id"`
	ImagePath   string `json:"icon_path"`
	Name        string `json:"name"`
	IsCollected bool   `json:"is_collect"`
}

// user Class
type user struct {
	ID                    int       `json:"user_id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	Password              string    `json:"password"`
	AccessToken           string    `json:"access_token"`
	Point                 int       `json:"point"`
	ActivityCount         int       `json:"activity_count"`
	Ranking               int       `json:"ranking"`
	Badges                []badge   `json:"badges"`
	ExpiredTokenTimestamp time.Time `json:"expired_token_timestamp"`
	CreatedTimestamp      time.Time `json:"created_timestamp"`
	UpdatedTimestamp      time.Time `json:"updated_timestamp"`
}

// New user by any request type
func NewUser(request interface{}) *user {
	user := user{
		CreatedTimestamp: time.Now().Local(),
		UpdatedTimestamp: time.Now().Local(),
	}
	utils.NewType().StructToStruct(request, &user)
	return &user
}

// New user by request
func NewUserByRequest(request UserRequest) user {
	user := user{
		CreatedTimestamp: time.Now().Local(),
		UpdatedTimestamp: time.Now().Local(),
	}
	utils.NewType().StructToStruct(request, &user)
	hashedPassword := utils.NewHelper().HashAndSalt(request.Password)
	user.Password = hashedPassword
	return user
}

// Setter for set corrected badge
func (u *user) SetCorrectedBadges(allBadgesDB []BadgeDB, correctedBadgesDB []UserBadgeDB) {
	for _, badgeDB := range allBadgesDB {
		u.Badges = append(u.Badges, badge{
			ID:          badgeDB.ID,
			ImagePath:   badgeDB.ImagePath,
			Name:        badgeDB.Name,
			IsCollected: u.isCorrectedBadge(correctedBadgesDB, badgeDB.ID),
		})
	}
}

// Private method for checking which badge is corrected
func (u *user) isCorrectedBadge(allBadgesDB []UserBadgeDB, badgeID int) bool {
	for _, correctedBadgeDB := range allBadgesDB {
		if badgeID == correctedBadgeDB.BadgeID {
			return true
		}
	}
	return false
}

// To UserDB model
func (u *user) ToDB() UserDB {
	tokenExpireHour := time.Hour * utils.NewType().ParseDuration(os.Getenv("TOKEN_EXPIRE_HOUR"))
	expiredTokenTimestamp := time.Now().Local().Add(tokenExpireHour)
	return UserDB{
		Name:                  u.Name,
		Email:                 u.Email,
		Password:              u.Password,
		CreatedTimestamp:      u.CreatedTimestamp,
		UpdatedTimestamp:      u.CreatedTimestamp,
		ExpiredTokenTimestamp: expiredTokenTimestamp,
	}
}

// To UserResponse model
func (u *user) ToUserResponse() UserResponse {
	response := UserResponse{}
	utils.NewType().StructToStruct(u, &response)
	return response
}

// To ProfileResponse model
func (u *user) ToProfileResponse() GetProfileResponse {
	response := GetProfileResponse{}
	utils.NewType().StructToStruct(u, &response)
	return response
}

// Check password with hashpassword
func (u *user) IsPasswordCorrect(password string) bool {
	return utils.NewHelper().ComparePasswords(u.Password, password)
}
