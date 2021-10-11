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
