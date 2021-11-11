package models

import "time"

// Model mapped User table in the database
type UserDB struct {
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

// Model mapped UserBage table in the database
type UserBadgeDB struct {
	UserID  int `gorm:"primaryKey;column:user_id" json:"user_id"`
	BadgeID int `gorm:"primaryKey;column:badge_id" json:"badge_id"`
}

// Model mapped Bage table in the database
type BadgeDB struct {
	ID        int    `gorm:"primaryKey;column:badge_id" json:"badge_id"`
	ImagePath string `gorm:"column:icon_path" json:"icon_path"`
	Name      string `gorm:"column:name" json:"name"`
}

// Model mapped Probile view in the database
type ProfileDB struct {
	ID               int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name             string    `gorm:"column:name" json:"name"`
	Point            int       `gorm:"column:point" json:"point"`
	ActivityCount    int       `gorm:"column:activity_count" json:"activity_count"`
	CreatedTimestamp time.Time `gorm:"column:created_timestamp" json:"created_timestamp"`
}

// Model mapped Ranking view in the database
type RankingDB struct {
	ID      int    `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name    string `gorm:"column:name" json:"name"`
	Point   int    `gorm:"column:point" json:"point"`
	Ranking int    `gorm:"column:ranking" json:"ranking"`
}

// Model mapped joined table in the database
// 		Table - UserBadge
// 		Table - Badge
type CorrectedBadgeDB struct {
	BadgeID int  `gorm:"column:badge_id" json:"badge_id"`
	Name    int  `gorm:"column:badge_name" json:"badge_name"`
	UserID  *int `gorm:"column:user_id" json:"user_id"`
}
