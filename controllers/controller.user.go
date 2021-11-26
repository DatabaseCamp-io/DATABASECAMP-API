package controllers

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/response"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
)

type userController struct {
	Repo repositories.IUserRepository
}

type IUserController interface {
	Register(request request.UserRequest) (*response.UserResponse, error)
	Login(request request.UserRequest) (*response.UserResponse, error)
	GetProfile(userID int) (*response.GetProfileResponse, error)
	EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error)
	GetRanking(id int) (*response.RankingResponse, error)
}

func NewUserController(repo repositories.IUserRepository) userController {
	return userController{Repo: repo}
}

func (c userController) Register(request request.UserRequest) (*response.UserResponse, error) {
	user := entities.NewUserByRequest(request)
	userDB, err := c.Repo.InsertUser(user.ToDB())
	if err != nil {
		logs.New().Error(err)
		if utils.NewHelper().IsSqlDuplicateError(err) {
			return nil, errs.ErrEmailAlreadyExists
		} else {
			return nil, errs.ErrInsertError
		}
	}
	user.SetID(userDB.ID)
	response := response.NewUserReponse(user)
	return &response, nil
}

func (c userController) Login(request request.UserRequest) (*response.UserResponse, error) {
	userDB, err := c.Repo.GetUserByEmail(request.Email)
	user := entities.NewUserByUserDB(*userDB)
	if err != nil || !user.IsPasswordCorrect(request.Password) {
		logs.New().Error(err)
		return nil, errs.ErrEmailOrPasswordNotCorrect
	}
	response := response.NewUserReponse(user)
	return &response, nil
}

func (c userController) GetProfile(id int) (*response.GetProfileResponse, error) {
	profileDB, err := c.Repo.GetProfile(id)
	if err != nil || profileDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	allBadge, err := c.Repo.GetAllBadge()
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	userBadgeGain, err := c.Repo.GetUserBadge(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrLoadError
	}

	user := entities.User{}
	user.SetCorrectedBadges(allBadge, userBadgeGain)
	response := response.NewGetProfileResponse(*profileDB, user.GetBadges())
	return &response, nil
}

func (c userController) EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error) {
	err := c.Repo.UpdatesByID(userID, map[string]interface{}{"name": request.Name})
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrUpdateError
	}
	response := response.EditProfileResponse{UpdatedName: request.Name}
	return &response, nil
}

func (c userController) GetRanking(userID int) (*response.RankingResponse, error) {
	userRankingDB, err := c.Repo.GetPointRanking(userID)
	if err != nil || userRankingDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrUserNotFound
	}

	leaderBoardDB, err := c.Repo.GetRankingLeaderBoard()
	if err != nil || leaderBoardDB == nil {
		logs.New().Error(err)
		return nil, errs.ErrLeaderBoardNotFound
	}

	response := response.RankingResponse{
		UserRanking: *userRankingDB,
		LeaderBoard: leaderBoardDB,
	}
	return &response, nil
}
