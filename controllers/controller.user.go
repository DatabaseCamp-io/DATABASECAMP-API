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
	repo repositories.IUserRepository
}

type IUserController interface {
	Register(request request.UserRequest) (*response.UserResponse, error)
	Login(request request.UserRequest) (*response.UserResponse, error)
	GetProfile(userID int) (*response.GetProfileResponse, error)
	EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error)
	GetRanking(id int) (*response.RankingResponse, error)
}

func NewUserController(repo repositories.IUserRepository) userController {
	return userController{repo: repo}
}

func (c userController) Register(request request.UserRequest) (*response.UserResponse, error) {
	user := entities.User{}
	utils.NewType().StructToStruct(request, &user)
	user.SetTimestamp()
	user.HashPassword()
	userDB, err := c.repo.InsertUser(user.ToDB())
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
	userDB, err := c.repo.GetUserByEmail(request.Email)
	user := entities.User{}
	utils.NewType().StructToStruct(userDB, &user)
	if err != nil || !user.IsPasswordCorrect(request.Password) {
		logs.New().Error(err)
		return nil, errs.ErrEmailOrPasswordNotCorrect
	}
	response := response.NewUserReponse(user)
	return &response, nil
}

func (c userController) GetProfile(id int) (*response.GetProfileResponse, error) {
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

	user := entities.User{}
	utils.NewType().StructToStruct(profileDB, &user)
	user.SetCorrectedBadges(allBadge, userBadgeGain)

	response := response.NewGetProfileResponse(user)
	return &response, nil
}

func (c userController) EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error) {
	err := c.repo.UpdatesByID(userID, map[string]interface{}{"name": request.Name})
	if err != nil {
		logs.New().Error(err)
		return nil, errs.ErrUpdateError
	}
	response := response.EditProfileResponse{UpdatedName: request.Name}
	return &response, nil
}

func (c userController) GetRanking(userID int) (*response.RankingResponse, error) {
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

	response := response.RankingResponse{
		UserRanking: *userRanking,
		LeaderBoard: leaderBoard,
	}
	return &response, nil
}
