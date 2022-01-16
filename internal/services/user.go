package services

import (
	"database-camp/internal/errs"
	"database-camp/internal/logs"
	"database-camp/internal/middleware/jwt"
	"database-camp/internal/models/entities/badge"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/models/request"
	"database-camp/internal/models/response"
	"database-camp/internal/repositories"
	"database-camp/internal/utils"
	"time"
)

type UserService interface {
	Register(request request.UserRequest) (*response.UserResponse, error)
	Login(request request.UserRequest) (*response.UserResponse, error)
	GetProfile(userID int) (*response.GetProfileResponse, error)
	EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error)
	GetRanking(id int) (*response.RankingResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s userService) Register(request request.UserRequest) (*response.UserResponse, error) {
	user, err := s.repo.InsertUser(user.User{
		Name:             request.Name,
		Email:            request.Email,
		Password:         utils.HashAndSalt(request.Password),
		Point:            0,
		CreatedTimestamp: time.Now().Local(),
		UpdatedTimestamp: time.Now().Local(),
	})

	if err != nil {
		logs.GetInstance().Error(err)

		if utils.IsSqlDuplicateError(err) {
			return nil, errs.ErrEmailAlreadyExists
		} else {
			return nil, errs.ErrInsertError
		}
	}

	jwt := jwt.New(s.repo)

	token, err := jwt.Sign(user.ID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrInternalServerError
	}

	response := response.UserResponse{
		ID:               user.ID,
		Name:             user.Name,
		Email:            user.Email,
		Point:            user.Point,
		AccessToken:      token,
		CreatedTimestamp: user.CreatedTimestamp,
		UpdatedTimestamp: user.UpdatedTimestamp,
	}

	return &response, nil
}

func (s userService) Login(request request.UserRequest) (*response.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(request.Email)

	correctPassword := utils.ComparePasswords(user.Password, request.Password)

	if err != nil || !correctPassword {
		logs.GetInstance().Error(err)
		return nil, errs.ErrEmailOrPasswordNotCorrect
	}

	jwt := jwt.New(s.repo)

	token, err := jwt.Sign(user.ID)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrInternalServerError
	}

	response := response.UserResponse{
		ID:               user.ID,
		Name:             user.Name,
		Email:            user.Email,
		Point:            user.Point,
		AccessToken:      token,
		CreatedTimestamp: user.CreatedTimestamp,
		UpdatedTimestamp: user.UpdatedTimestamp,
	}

	return &response, nil
}

func (s userService) GetProfile(id int) (*response.GetProfileResponse, error) {
	profile, err := s.repo.GetProfile(id)
	if err != nil || profile == nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrUserNotFound
	}

	allBadgeDB, err := s.repo.GetAllBadge()
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	userBadgeDB, err := s.repo.GetUserBadge(id)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLoadError
	}

	badges := badge.NewBadges(allBadgeDB, userBadgeDB)

	response := response.GetProfileResponse{
		ID:               profile.ID,
		Name:             profile.Name,
		Point:            profile.Point,
		ActivityCount:    profile.ActivityCount,
		Badges:           badges,
		CreatedTimestamp: profile.CreatedTimestamp,
	}

	return &response, nil
}

func (s userService) EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error) {
	err := s.repo.UpdatesByID(userID, map[string]interface{}{"name": request.Name})
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrUpdateError
	}

	response := response.EditProfileResponse{UpdatedName: request.Name}

	return &response, nil
}

func (s userService) GetRanking(userID int) (*response.RankingResponse, error) {
	userRankingDB, err := s.repo.GetPointRanking(userID)
	if err != nil || userRankingDB == nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrUserNotFound
	}

	leaderBoardDB, err := s.repo.GetRankingLeaderBoard()
	if err != nil || leaderBoardDB == nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrLeaderBoardNotFound
	}

	response := response.RankingResponse{
		UserRanking: *userRankingDB,
		LeaderBoard: leaderBoardDB,
	}

	return &response, nil
}
