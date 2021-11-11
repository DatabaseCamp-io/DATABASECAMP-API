package models

import (
	"DatabaseCamp/utils"
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

// CorrectedBadge data class
type CorrectedBadge struct {
	ID        int    `json:"badge_id"`
	ImagePath string `json:"icon_path"`
	Name      string `json:"name"`
	IsCollect bool   `json:"is_collect"`
}

// User Class
type User struct {
	ID                    int              `json:"user_id"`
	Name                  string           `json:"name"`
	Email                 string           `json:"email"`
	Password              string           `json:"password"`
	AccessToken           string           `json:"access_token"`
	Point                 int              `json:"point"`
	ActivityCount         int              `json:"activity_count"`
	Ranking               int              `json:"ranking"`
	Badges                []CorrectedBadge `json:"badges"`
	ExpiredTokenTimestamp time.Time        `json:"expired_token_timestamp"`
	CreatedTimestamp      time.Time        `json:"created_timestamp"`
	UpdatedTimestamp      time.Time        `json:"updated_timestamp"`
}

// New user by any request type
func NewUser(request interface{}) User {
	user := User{}
	utils.NewType().StructToStruct(request, &user)
	return user
}

// New user by request
func NewUserByRequest(request UserRequest) User {
	user := User{}
	utils.NewType().StructToStruct(request, &user)
	hashedPassword := utils.NewHelper().HashAndSalt(request.Password)
	user.Password = hashedPassword
	return user
}

// Setter for set corrected badge
func (u *User) SetCorrectedBadges(allBadgesDB []BadgeDB, correctedBadgesDB []UserBadgeDB) {
	for _, badgeDB := range allBadgesDB {
		u.Badges = append(u.Badges, CorrectedBadge{
			ID:        badgeDB.ID,
			ImagePath: badgeDB.ImagePath,
			Name:      badgeDB.Name,
			IsCollect: u.isCorrectedBadge(allBadgesDB, badgeDB.ID),
		})
	}
}

// Private method for checking which badge is corrected
func (u *User) isCorrectedBadge(allBadgesDB []BadgeDB, badgeID int) bool {
	for _, badgeDB := range allBadgesDB {
		if badgeID == badgeDB.ID {
			return true
		}
	}
	return false
}

// To UserDB model
func (u *User) ToDB() UserDB {
	return UserDB{
		Name:                  u.Name,
		Email:                 u.Email,
		Password:              u.Password,
		ExpiredTokenTimestamp: time.Now().Local(),
		CreatedTimestamp:      time.Now().Local(),
		UpdatedTimestamp:      time.Now().Local(),
	}
}

// To UserResponse model
func (u *User) ToUserResponse() UserResponse {
	response := UserResponse{}
	utils.NewType().StructToStruct(u, &response)
	return response
}

// To ProfileResponse model
func (u *User) ToProfileResponse() GetProfileResponse {
	response := GetProfileResponse{}
	utils.NewType().StructToStruct(u, &response)
	return response
}

// Check password with hashpassword
func (u *User) IsPasswordCorrect(password string) bool {
	return utils.NewHelper().ComparePasswords(u.Password, password)
}
