package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/utils"
	"time"
)

type userController struct {
	repo repository.IUserRepository
}

type IUserController interface {
	Register(request models.UserRequest) (*models.UserResponse, error)
	Login(request models.UserRequest) (*models.UserResponse, error)
	GetProfile(userID int) (*models.ProfileResponse, error)
}

func NewUserController(repo repository.IUserRepository) userController {
	return userController{repo: repo}
}

func (c userController) setupUserModel(request models.UserRequest) models.User {
	hashedPassword := utils.NewHelper().HashAndSalt(request.Password)
	return models.User{
		Name:                  request.Name,
		Email:                 request.Email,
		Password:              hashedPassword,
		ExpiredTokenTimestamp: time.Now().Local(),
		CreatedTimestamp:      time.Now().Local(),
		UpdatedTimestamp:      time.Now().Local(),
	}
}

func (c userController) Register(request models.UserRequest) (*models.UserResponse, error) {
	response := models.UserResponse{}
	user := c.setupUserModel(request)

	user, err := c.repo.Insert(user)
	if err != nil {
		logs.New().Error(err)

		if utils.NewHelper().IsSqlDuplicateError(err) {
			return nil, errs.NewBadRequestError("อีเมลมีการใช้งานแล้ว", "Email is already exists")
		} else {
			return nil, errs.NewInternalServerError("ลงทะเบียนไม่สำเร็จ", "Register Failed")
		}
	}

	utils.NewType().StructToStruct(user, &response)
	return &response, nil
}

func (c userController) Login(request models.UserRequest) (*models.UserResponse, error) {
	response := models.UserResponse{}
	user, err := c.repo.GetUserByEmail(request.Email)

	if err != nil || !utils.NewHelper().ComparePasswords(user.Password, request.Password) {
		logs.New().Error(err)
		return nil, errs.NewBadRequestError("อีเมลหรือรหัสผ่านไม่ถูกต้อง", "Email or Password Not Correct")
	}

	utils.NewType().StructToStruct(user, &response)
	return &response, nil
}

func (c userController) GetProfile(id int) (*models.ProfileResponse, error) {
	response := models.ProfileResponse{}
	profileDB, err := c.repo.GetProfile(id)
	if err != nil || profileDB == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบผู้ใช้", "Profile Not Found")
	}
	err = utils.NewType().StructToStruct(profileDB, &response)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
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
	response.Badges = c.calculateUserBadge(allBadge, userBadgeGain)
	return &response, nil
}

func (c userController) calculateUserBadge(allBadge []models.Badge, userBadgeGain []models.UserBadgeIDPair) []models.Badge {
	for i, v := range allBadge {
		allBadge[i].IsCollect = c.isCollectBadge(v.ID, userBadgeGain)
	}
	return allBadge
}

func (c userController) isCollectBadge(badgeID int, userBadgeGain []models.UserBadgeIDPair) bool {
	for _, v := range userBadgeGain {
		if v.BadgeID == badgeID {
			return true
		}
	}
	return false
}

func (c userController) GetUserRanking(id int) (*models.PointRanking, error) {
	response := models.PointRanking{}
	user, err := c.repo.UserPointranking(id)
	if err != nil || user == nil{
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบผู้ใช้", "Profile Not Found")
	}
	return &response, nil
}

func(c userController) LeaderBoard() ([]models.PointRanking, error){
	response := make([]models.PointRanking, 0)
	ranking, err := c.repo.GetAllPointranking(){
		if err != nil || ranking == nil{
			logs.New().Error(err)
			return nil, errs.NewNotFoundError("ไม่มีตารางคะแนน", "LeaderBoard Not Found")
		}
	return response, nil
	}
}