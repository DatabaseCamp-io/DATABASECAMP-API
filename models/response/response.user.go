package response

import (
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/general"
	"DatabaseCamp/utils"
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
	response := UserResponse{}
	utils.NewType().StructToStruct(user, &response)
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

func NewGetProfileResponse(user entities.User) GetProfileResponse {
	response := GetProfileResponse{}
	utils.NewType().StructToStruct(user, &response)
	return response
}

// Model for response edit profile
type EditProfileResponse struct {
	UpdatedName string `json:"updated_name"`
}

// Model for response ranking
type RankingResponse struct {
	UserRanking general.RankingDB   `json:"user_ranking"`
	LeaderBoard []general.RankingDB `json:"leader_board"`
}
