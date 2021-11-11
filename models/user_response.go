package models

import "time"

// Model for response user
type UserResponse struct {
	ID               int       `json:"user_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Point            int       `json:"point"`
	AccessToken      string    `json:"access_token"`
	CreatedTimestamp time.Time ` json:"created_timestamp"`
	UpdatedTimestamp time.Time ` json:"updated_timestamp"`
}

// Model for response profile
type ProfileResponse struct {
	ID               int       `json:"user_id"`
	Name             string    `json:"name"`
	Point            int       `json:"point"`
	ActivityCount    int       `json:"activity_count"`
	Badges           []BadgeDB `json:"badges"`
	CreatedTimestamp time.Time ` json:"created_timestamp"`
}
