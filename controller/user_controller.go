package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/utils"
)

type userController struct {
	repo repository.IUserRepository
}

type IUserController interface {
	Register(request models.UserRequest) (*models.UserResponse, error)
	Login(request models.UserRequest) (*models.UserResponse, error)
	GetProfile(userID int) (*models.GetProfileResponse, error)
	EditProfile(userID int, request models.UserRequest) (*models.EditProfileResponse, error)
	GetRanking(id int) (*models.RankingResponse, error)
}

func NewUserController(repo repository.IUserRepository) userController {
	return userController{repo: repo}
}

func (c userController) Register(request models.UserRequest) (*models.UserResponse, error) {
	var err error
	user := models.NewUserByRequest(request)
	user.ID, err = c.repo.Insert(user.ToDB())
	if err != nil {
		logs.New().Error(err)
		if utils.NewHelper().IsSqlDuplicateError(err) {
			return nil, errs.NewBadRequestError("อีเมลมีการใช้งานแล้ว", "Email is already exists")
		} else {
			return nil, errs.NewInternalServerError("ลงทะเบียนไม่สำเร็จ", "Register Failed")
		}
	}
	response := user.ToUserResponse()
	return &response, nil
}

func (c userController) Login(request models.UserRequest) (*models.UserResponse, error) {
	userDB, err := c.repo.GetUserByEmail(request.Email)
	user := models.NewUser(userDB)
	if err != nil || !user.IsPasswordCorrect(request.Password) {
		logs.New().Error(err)
		return nil, errs.NewBadRequestError("อีเมลหรือรหัสผ่านไม่ถูกต้อง", "Email or Password Not Correct")
	}
	response := user.ToUserResponse()
	return &response, nil
}

func (c userController) GetProfile(id int) (*models.GetProfileResponse, error) {
	profileDB, err := c.repo.GetProfile(id)
	if err != nil || profileDB == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบผู้ใช้", "Profile Not Found")
	}
	allBadge, err := c.repo.GetAllBadge()
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
	userBadgeGain, err := c.repo.GetUserBadgeIDPair(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
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
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
	response := models.EditProfileResponse{UpdatedName: request.Name}
	return &response, nil
}

func (c userController) GetRanking(userID int) (*models.RankingResponse, error) {
	userRanking, err := c.getUserRanking(userID)
	if err != nil {
		return nil, err
	}
	leaderBoard, err := c.getLeaderBoard()
	if err != nil {
		return nil, err
	}
	response := models.RankingResponse{
		UserRanking: *userRanking,
		LeaderBoard: leaderBoard,
	}
	return &response, nil
}

func (c userController) getUserRanking(id int) (*models.RankingDB, error) {
	user, err := c.repo.UserPointranking(id)
	if err != nil || user == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบผู้ใช้", "Profile Not Found")
	}
	return user, nil
}

func (c userController) getLeaderBoard() ([]models.RankingDB, error) {
	ranking, err := c.repo.GetAllPointranking()
	if err != nil || ranking == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่มีตารางคะแนน", "LeaderBoard Not Found")
	}
	return ranking, nil
}
