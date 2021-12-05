package response

// response.user.go
/**
 * 	This file is a part of models, used to collect response of user
 */
import (
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
	"time"
)

// Model of user item to prepare user response
type UserResponse struct {
	ID               int       `json:"user_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Point            int       `json:"point"`
	AccessToken      string    `json:"access_token"`
	CreatedTimestamp time.Time ` json:"created_timestamp"`
	UpdatedTimestamp time.Time ` json:"updated_timestamp"`
}

/**
 * Constructor creates a new VideoLectureResponse instance
 *
 * @param user			Entities user from database to create user response
 *
 *
 * @return 	instance of UserReponse
 */
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

// Model of get profile item to prepare get profile response
type GetProfileResponse struct {
	ID               int              `json:"user_id"`
	Name             string           `json:"name"`
	Point            int              `json:"point"`
	ActivityCount    int              `json:"activity_count"`
	Badges           []entities.Badge `json:"badges"`
	CreatedTimestamp time.Time        ` json:"created_timestamp"`
}
/**
 * Constructor creates a new VideoLectureResponse instance
 *
 * @param profileDB			Profile model from database to create get profile response
 * @param badges			Badge model from database to create get profile response
 *
 *
 * @return 	instance of GetProfileResponse
 */
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

/**
 * This class represent EditProfile response
 */
type EditProfileResponse struct {
	UpdatedName string `json:"updated_name"`
}

/**
 * This class represent Ranking response
 */
type RankingResponse struct {
	UserRanking storages.RankingDB   `json:"user_ranking"`
	LeaderBoard []storages.RankingDB `json:"leader_board"`
}
