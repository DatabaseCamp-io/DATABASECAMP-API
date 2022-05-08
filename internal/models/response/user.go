package response

import (
	"database-camp/internal/models/entities/badge"
	"database-camp/internal/models/entities/user"
	"time"
)

type UserResponse struct {
	ID               int       `json:"user_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Point            int       `json:"point"`
	AccessToken      string    `json:"access_token"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	UpdatedTimestamp time.Time `json:"updated_timestamp"`
}

type GetProfileResponse struct {
	ID               int                `json:"user_id"`
	Name             string             `json:"name"`
	Point            int                `json:"point"`
	ActivityCount    int                `json:"activity_count"`
	Badges           []badge.Badge      `json:"badges"`
	SpiderDataset    user.SpiderDataset `json:"spider"`
	CreatedTimestamp time.Time          ` json:"created_timestamp"`
}

type EditProfileResponse struct {
	UpdatedName string `json:"updated_name"`
}

type RankingResponse struct {
	UserRanking user.Ranking   `json:"user_ranking"`
	LeaderBoard []user.Ranking `json:"leader_board"`
}

type RecommendResponse struct {
	Recommend []user.RecommendGroup `json:"recommend_group"`
}
