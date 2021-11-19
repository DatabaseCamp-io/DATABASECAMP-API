package controllers

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
)

type userController struct {
	repo repositories.IUserRepository
}

type IUserController interface {
	Register(request models.UserRequest) (*models.UserResponse, error)
	Login(request models.UserRequest) (*models.UserResponse, error)
	GetProfile(userID int) (*models.GetProfileResponse, error)
	EditProfile(userID int, request models.UserRequest) (*models.EditProfileResponse, error)
	GetRanking(id int) (*models.RankingResponse, error)
}

func NewUserController(repo repositories.IUserRepository) userController {
	return userController{repo: repo}
}

func (c userController) Register(request models.UserRequest) (*models.UserResponse, error) {
	user := models.NewUserByRequest(request)

	userDB, err := c.repo.InsertUser(user.ToDB())
	if err != nil {
		logs.New().Error(err)
		if utils.NewHelper().IsSqlDuplicateError(err) {
			return nil, errs.ErrEmailAlreadyExists
		} else {
			return nil, errs.ErrInsertError
		}
	}

	user.ID = userDB.ID
	response := user.ToUserResponse()
	return &response, nil
}

func (c userController) Login(request models.UserRequest) (*models.UserResponse, error) {
	userDB, err := c.repo.GetUserByEmail(request.Email)
	user := models.NewUser(userDB)
	if err != nil || !user.IsPasswordCorrect(request.Password) {
		logs.New().Error(err)
		return nil, errs.ErrEmailOrPasswordNotCorrect
	}
	response := user.ToUserResponse()
	return &response, nil
}

func (c userController) GetProfile(id int) (*models.GetProfileResponse, error) {
	profileDB, err := c.repo.GetProfile(id)
	if err != nil || profileDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	allBadge, err := c.repo.GetAllBadge()
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	userBadgeGain, err := c.repo.GetUserBadge(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	user := models.NewUser(profileDB)
	user.SetCorrectedBadges(allBadge, userBadgeGain)

	response := user.ToProfileResponse()
	return &response, nil
}

func (c userController) EditProfile(userID int, request models.UserRequest) (*models.EditProfileResponse, error) {
	err := c.repo.UpdatesByID(userID, map[string]interface{}{"name": request.Name})
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrUpdateError
	}
	response := models.EditProfileResponse{UpdatedName: request.Name}
	return &response, nil
}

func (c userController) GetRanking(userID int) (*models.RankingResponse, error) {
	userRanking, err := c.repo.GetPointRanking(userID)
	if err != nil || userRanking == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	leaderBoard, err := c.repo.GetRankingLeaderBoard()
	if err != nil || leaderBoard == nil {
		logs.New().Error(err)
		return nil, errs.ErrLeaderBoardNotFound
	}

	response := models.RankingResponse{
		UserRanking: *userRanking,
		LeaderBoard: leaderBoard,
	}
	return &response, nil
}
