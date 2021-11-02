package models

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/utils"
	"time"
)

type User struct {
	ID                    int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name                  string    `gorm:"column:name" json:"name"`
	Email                 string    `gorm:"column:email" json:"email"`
	Password              string    `gorm:"column:password" json:"password"`
	AccessToken           string    `gorm:"column:access_token" json:"access_token"`
	Point                 int       `gorm:"column:point" json:"point"`
	ExpiredTokenTimestamp time.Time `gorm:"column:expired_token_timestamp" json:"expired_token_timestamp"`
	CreatedTimestamp      time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
	UpdatedTimestamp      time.Time `gorm:"column:updated_timestamp" json:"updated_timestamp"`
}

type UserResponse struct {
	ID               int       `json:"user_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Point            int       `json:"point"`
	AccessToken      string    `gorm:"column:access_token" json:"access_token"`
	CreatedTimestamp time.Time ` json:"created_timestamp"`
	UpdatedTimestamp time.Time ` json:"updated_timestamp"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r UserRequest) ValidateRegister() error {
	var err error
	if r.Name == "" {
		err = errs.NewBadRequestError("ไม่พบชื่อในคำร้องขอ", "Name Not Found")
	} else {
		err = r.ValidateLogin()
	}
	return err
}

func (r UserRequest) ValidateLogin() error {
	var err error
	if r.Email == "" {
		err = errs.NewBadRequestError("ไม่พบอีเมลในคำร้องขอ", "Email Not Found")
	} else if !utils.NewHelper().IsEmailValid(r.Email) {
		err = errs.NewBadRequestError("รูปแบบ email ไม่ถูกต้อง", "Email Invalid")
	} else if r.Password == "" {
		err = errs.NewBadRequestError("ไม่พบรหัสผ่านในคำร้องขอ", "Password Not Found")
	} else if len(r.Password) < 8 {
		err = errs.NewBadRequestError("ความยาวของรหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร", "Password length must be at least 8 characters")
	}
	return err
}

type ProfileResponse struct {
	ID               int       `json:"user_id"`
	Name             string    `json:"name"`
	Point            int       `json:"point"`
	ActivityCount    int       `json:"activity_count"`
	Badges           []Badge   `json:"badges"`
	CreatedTimestamp time.Time ` json:"created_timestamp"`
}

type Badge struct {
	ID        int    `gorm:"primaryKey;column:badge_id" json:"badge_id"`
	ImagePath string `gorm:"column:icon_path" json:"icon_path"`
	Name      string `gorm:"column:name" json:"name"`
	IsCollect bool   `json:"is_collect"`
}

type ProfileDB struct {
	ID               int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name             string    `gorm:"column:name" json:"name"`
	Point            int       `gorm:"column:point" json:"point"`
	ActivityCount    int       `gorm:"column:activity_count" json:"activity_count"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

type UserBadgeIDPair struct {
	UserID  int `gorm:"primaryKey;column:user_id" json:"user_id"`
	BadgeID int `gorm:"primaryKey;column:badge_id" json:"badge_id"`
}

type PointRanking struct {
	ID      int    `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name    string `gorm:"column:name" json:"name"`
	Point   int    `gorm:"column:point" json:"point"`
	Ranking int    `gorm:"column:ranking" json:"ranking"`
}
