package repository

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models"
	"DatabaseCamp/utils"
)

type userRepository struct {
	database database.IDatabase
}

type IUserRepository interface {
	Insert(user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id int) (models.User, error)
	UpdatesByID(id int, updateData map[string]interface{}) error
	GetProfile(id int) (*models.ProfileDB, error)
	GetLearningProgression(id int) ([]models.LearningProgressionDB, error)
	GetFailedExam(id int) ([]models.ExamResultDB, error)
	GetAllBadge() ([]models.Badge, error)
	GetUserBadgeIDPair(id int) ([]models.UserBadgeIDPair, error)
	GetAllPointranking() ([]models.PointRanking, error)
	UserPointranking(id int) (*models.PointRanking, error)
}

func NewUserRepository(db database.IDatabase) userRepository {
	return userRepository{database: db}
}

func (r userRepository) Insert(user models.User) (models.User, error) {
	err := r.database.GetDB().
		Table(models.TableName.User).
		Create(&user).
		Error
	return user, err
}

func (r userRepository) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}
	err := r.database.GetDB().
		Table(models.TableName.User).
		Where("email = ?", email).
		Find(&user).
		Error
	return user, err
}

func (r userRepository) GetUserByID(id int) (models.User, error) {
	user := models.User{}
	err := r.database.GetDB().
		Table(models.TableName.User).
		Where(models.IDName.User+" = ?", id).
		Find(&user).
		Error
	return user, err
}

func (r userRepository) UpdatesByID(id int, updateData map[string]interface{}) error {
	err := r.database.GetDB().
		Table(models.TableName.User).
		Select("", utils.NewHelper().GetKeyList(updateData)).
		Where(models.IDName.User+" = ?", id).
		Updates(updateData).
		Error
	return err
}

func (r userRepository) GetProfile(id int) (*models.ProfileDB, error) {
	profile := models.ProfileDB{}
	err := r.database.GetDB().
		Table(models.ViewName.Profile).
		Where(models.IDName.User+" = ?", id).
		Find(&profile).
		Error
	if profile == (models.ProfileDB{}) {
		return nil, nil
	}
	return &profile, err
}

func (r userRepository) GetLearningProgression(id int) ([]models.LearningProgressionDB, error) {
	learningProgrogression := make([]models.LearningProgressionDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.LearningProgression).
		Where(models.IDName.User+" = ?", id).
		Order("created_timestamp desc").
		Find(&learningProgrogression).
		Error
	return learningProgrogression, err
}

func (r userRepository) GetFailedExam(id int) ([]models.ExamResultDB, error) {
	exam := make([]models.ExamResultDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.ExamResult).
		Where(models.IDName.User+" = ? AND is_passed = false", id).
		Find(&exam).
		Error

	return exam, err
}

func (r userRepository) GetAllBadge() ([]models.Badge, error) {
	badge := make([]models.Badge, 0)
	err := r.database.GetDB().
		Table(models.TableName.Badge).
		Find(&badge).
		Error
	return badge, err
}

func (r userRepository) GetUserBadgeIDPair(id int) ([]models.UserBadgeIDPair, error) {
	badgePair := make([]models.UserBadgeIDPair, 0)
	err := r.database.GetDB().
		Table(models.TableName.UserBadge).
		Where(models.IDName.User+" = ?", id).
		Find(&badgePair).
		Error
	return badgePair, err
}

func (r userRepository) GetAllPointranking() ([]models.PointRanking, error) {
	ranking := make([]models.PointRanking, 0)
	err := r.database.GetDB().
		Table(models.ViewName.Ranking).
		Limit(20).
		Find(&ranking).
		Error
	return ranking, err
}

func (r userRepository) UserPointranking(id int) (*models.PointRanking, error) {
	ranking := models.PointRanking{}
	err := r.database.GetDB().
		Table(models.ViewName.Ranking).
		Where(models.IDName.User+" = ?", id).
		Find(&ranking).
		Error
	return &ranking, err
}
