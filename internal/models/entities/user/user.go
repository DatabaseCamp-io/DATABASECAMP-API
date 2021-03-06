package user

import (
	"time"
)

type User struct {
	ID                    int       `gorm:"primaryKey;column:user_id"`
	Name                  string    `gorm:"column:name"`
	Email                 string    `gorm:"column:email"`
	Password              string    `gorm:"column:password"`
	AccessToken           string    `gorm:"column:access_token"`
	Point                 int       `gorm:"column:point"`
	ExpiredTokenTimestamp time.Time `gorm:"column:expired_token_timestamp"`
	CreatedTimestamp      time.Time `gorm:"column:created_timestamp"`
	UpdatedTimestamp      time.Time `gorm:"column:updated_timestamp"`
}

type Profile struct {
	ID               int       `gorm:"primaryKey;column:user_id"`
	Name             string    `gorm:"column:name"`
	Point            int       `gorm:"column:point"`
	ActivityCount    int       `gorm:"column:activity_count"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp"`
}

type Ranking struct {
	ID      int    `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name    string `gorm:"column:name" json:"name"`
	Point   int    `gorm:"column:point" json:"point"`
	Ranking int    `gorm:"column:ranking" json:"ranking"`
}
