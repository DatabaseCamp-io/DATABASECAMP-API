package models

import (
	"DatabaseCamp/utils"
	"time"
)

type ChangePointMode string

var Mode = struct {
	Add    ChangePointMode
	Reduce ChangePointMode
}{
	"+",
	"-",
}

type CorrectedBadge struct {
	ID        int    `json:"badge_id"`
	ImagePath string `json:"icon_path"`
	Name      string `json:"name"`
	IsCollect bool   `json:"is_collect"`
}

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

func NewUser(request interface{}) User {
	user := User{}
	utils.NewType().StructToStruct(request, &user)
	return user
}

func NewUserWithHashPassword(request UserRequest) User {
	user := User{}
	utils.NewType().StructToStruct(request, &user)
	hashedPassword := utils.NewHelper().HashAndSalt(request.Password)
	user.Password = hashedPassword
	return user
}

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

func (u *User) isCorrectedBadge(allBadgesDB []BadgeDB, badgeID int) bool {
	for _, badgeDB := range allBadgesDB {
		if badgeID == badgeDB.ID {
			return true
		}
	}
	return false
}

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

func (u *User) ToUserResponse() UserResponse {
	response := UserResponse{}
	utils.NewType().StructToStruct(u, &response)
	return response
}

func (u *User) ToProfileResponse() GetProfileResponse {
	response := GetProfileResponse{}
	utils.NewType().StructToStruct(u, &response)
	return response
}

func (u *User) IsPasswordCorrect(password string) bool {
	return utils.NewHelper().ComparePasswords(u.Password, password)
}
