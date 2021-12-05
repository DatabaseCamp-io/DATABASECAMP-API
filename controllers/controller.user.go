package controllers

// controller.user.go
/**
 * 	This file is a part of controllers, used to do business logic of user
 */

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/response"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
)

/**
 * This class do business logic of user
 */
type userController struct {
	Repo repositories.IUserRepository
}

/**
 * Constructor creates a new userController instance
 *
 * @param   UserRepo        	User Repository for load user data
 *
 * @return 	instance of userController
 */
func NewUserController(repo repositories.IUserRepository) userController {
	return userController{Repo: repo}
}

/**
 * Register for learning in platform
 *
 * @param 	request 		User request for create account
 *
 * @return response of user
 * @return error of register
 */
func (c userController) Register(request request.UserRequest) (*response.UserResponse, error) {

	// Create User
	user := entities.NewUserByRequest(request)

	// Insert User and Check error
	userDB, err := c.Repo.InsertUser(user.ToDB())
	if err != nil {
		logs.New().Error(err)
		if utils.NewHelper().IsSqlDuplicateError(err) {
			return nil, errs.ErrEmailAlreadyExists
		} else {
			return nil, errs.ErrInsertError
		}
	}

	// Set user ID
	user.SetID(userDB.ID)

	// Create user response
	response := response.NewUserReponse(user)
	return &response, nil
}

/**
 * Login
 *
 * @param 	request 		User request for login to the platform
 *
 * @return response of user
 * @return error of login
 */
func (c userController) Login(request request.UserRequest) (*response.UserResponse, error) {

	// Get user from the repository
	userDB, err := c.Repo.GetUserByEmail(request.Email)

	// Create User
	user := entities.NewUserByUserDB(*userDB)

	// Check email and password
	if err != nil || !user.IsPasswordCorrect(request.Password) {
		logs.New().Error(err)
		return nil, errs.ErrEmailOrPasswordNotCorrect
	}

	// Create user response
	response := response.NewUserReponse(user)

	return &response, nil
}

/**
 * Get user profile
 *
 * @param 	userID 		   	User ID for getting user profile
 *
 * @return response of get profile
 * @return error of getting profile
 */
func (c userController) GetProfile(id int) (*response.GetProfileResponse, error) {

	// Get user profile from repository
	profileDB, err := c.Repo.GetProfile(id)
	if err != nil || profileDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	// Get all Badge from repository
	allBadge, err := c.Repo.GetAllBadge()
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	// Get user badges from repository
	userBadgeGain, err := c.Repo.GetUserBadge(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	// Create user and set data
	user := entities.User{}
	user.SetCorrectedBadges(allBadge, userBadgeGain)

	// Create get profile response
	response := response.NewGetProfileResponse(*profileDB, user.GetBadges())

	return &response, nil
}

/**
 * Edit profile data
 *
 * @param 	userID 		   	User ID for getting user profile
 * @param 	request 		User request for editing the user profile
 *
 * @return response of edit profile
 * @return error of edit profile
 */
func (c userController) EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error) {

	// Update user data and check error
	err := c.Repo.UpdatesByID(userID, map[string]interface{}{"name": request.Name})
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrUpdateError
	}

	// Create edit profile response
	response := response.EditProfileResponse{UpdatedName: request.Name}

	return &response, nil
}

/**
 * Get others and own ranking
 *
 * @param 	id		User ID for getting user ranking
 *
 * @return response of ranking
 * @return error of getting ranking
 */
func (c userController) GetRanking(userID int) (*response.RankingResponse, error) {

	// Get user ranking from repository
	userRankingDB, err := c.Repo.GetPointRanking(userID)
	if err != nil || userRankingDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	// Get others ranking from repository
	leaderBoardDB, err := c.Repo.GetRankingLeaderBoard()
	if err != nil || leaderBoardDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrLeaderBoardNotFound
	}

	// Create ranking response
	response := response.RankingResponse{
		UserRanking: *userRankingDB,
		LeaderBoard: leaderBoardDB,
	}

	return &response, nil
}
