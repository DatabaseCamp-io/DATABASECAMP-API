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
	GetExam(id int) ([]models.ExamDB, error)
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

func (r userRepository) GetExam(id int) ([]models.ExamDB, error) {
	exam := make([]models.ExamDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.Exam).
		Where(models.IDName.User+" = ?", id).
		Find(&exam).
		Error
	return exam, err
}
