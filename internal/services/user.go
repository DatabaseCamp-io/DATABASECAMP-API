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
	"database-camp/internal/services/loaders"
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
	userRepo     repositories.UserRepository
	learningRepo repositories.LearningRepository
}

func NewUserService(userRepo repositories.UserRepository, learningRepo repositories.LearningRepository) *userService {
	return &userService{userRepo: userRepo, learningRepo: learningRepo}
}

func (s userService) Register(request request.UserRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.InsertUser(user.User{
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

	jwt := jwt.New(s.userRepo)

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
	user, err := s.userRepo.GetUserByEmail(request.Email)

	correctPassword := utils.ComparePasswords(user.Password, request.Password)

	if err != nil || !correctPassword {
		logs.GetInstance().Error(err)
		return nil, errs.ErrEmailOrPasswordNotCorrect
	}

	jwt := jwt.New(s.userRepo)

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

	loader := loaders.NewProfileLoader(s.learningRepo, s.userRepo)

	err := loader.Load(id)
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrContentNotFound
	}

	spiderDataset := loader.GetSpiderDataset()
	contentGroups := loader.GetContentGroups()
	badges := loader.GetBadges()
	profile := loader.GetProfile()
	userBadges := loader.GetUserBadges()

	spiderDataset.FillContentGroups(contentGroups)

	_badges := badge.NewBadges(badges, userBadges)

	response := response.GetProfileResponse{
		ID:               profile.ID,
		Name:             profile.Name,
		Point:            profile.Point,
		ActivityCount:    profile.ActivityCount,
		Badges:           _badges,
		SpiderDataset:    spiderDataset,
		CreatedTimestamp: profile.CreatedTimestamp,
	}

	return &response, nil
}

func (s userService) EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error) {
	err := s.userRepo.UpdatesByID(userID, map[string]interface{}{"name": request.Name})
	if err != nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrUpdateError
	}

	response := response.EditProfileResponse{UpdatedName: request.Name}

	return &response, nil
}

func (s userService) GetRanking(userID int) (*response.RankingResponse, error) {
	userRankingDB, err := s.userRepo.GetPointRanking(userID)
	if err != nil || userRankingDB == nil {
		logs.GetInstance().Error(err)
		return nil, errs.ErrUserNotFound
	}

	leaderBoardDB, err := s.userRepo.GetRankingLeaderBoard()
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
