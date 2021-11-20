package response

import (
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
	"time"
)

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

func NewUserReponse(user entities.User) UserResponse {
	response := UserResponse{
		ID:               user.GetID(),
		Name:             user.GetName(),
		Email:            user.GetEmail(),
		Point:            0,
		AccessToken:      "",
		CreatedTimestamp: user.GetCreatedTimestamp(),
		UpdatedTimestamp: user.GetUpdatedTimstamp(),
	}
	return response
}

// Model for response get profile
type GetProfileResponse struct {
	ID               int              `json:"user_id"`
	Name             string           `json:"name"`
	Point            int              `json:"point"`
	ActivityCount    int              `json:"activity_count"`
	Badges           []entities.Badge `json:"badges"`
	CreatedTimestamp time.Time        ` json:"created_timestamp"`
}

func NewGetProfileResponse(profileDB storages.ProfileDB, badges []entities.Badge) GetProfileResponse {
	response := GetProfileResponse{
		ID:               profileDB.ID,
		Name:             profileDB.Name,
		Point:            profileDB.Point,
		ActivityCount:    profileDB.ActivityCount,
		Badges:           badges,
		CreatedTimestamp: profileDB.CreatedTimestamp,
	}
	return response
}

// Model for response edit profile
type EditProfileResponse struct {
	UpdatedName string `json:"updated_name"`
}

// Model for response ranking
type RankingResponse struct {
	UserRanking storages.RankingDB   `json:"user_ranking"`
	LeaderBoard []storages.RankingDB `json:"leader_board"`
}
