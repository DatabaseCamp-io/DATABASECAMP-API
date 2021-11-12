package repository

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models"
	"DatabaseCamp/utils"
	"fmt"
)

type userRepository struct {
	database database.IDatabase
}

type IUserRepository interface {
	Insert(user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id int) (models.User, error)
	GetProfile(id int) (*models.ProfileDB, error)
	GetLearningProgression(id int) ([]models.LearningProgressionDB, error)
	GetFailedExam(id int) ([]models.ExamResultDB, error)
	GetAllBadge() ([]models.BadgeDB, error)
	UpdatesByID(id int, updateData map[string]interface{}) error
	GetUserBadgeIDPair(id int) ([]models.UserBadgeIDPair, error)
	GetAllPointranking() ([]models.PointRanking, error)
	UserPointranking(id int) (*models.PointRanking, error)
	GetUserHint(userID int, activityID int) ([]models.UserHintDB, error)
	InsertUserHint(userHint models.UserHintDB) (*models.UserHintDB, error)
	InsertUserHintTransaction(tx database.ITransaction, userHint models.UserHintDB) (*models.UserHintDB, error)
	InsertLearningProgressionTransaction(tx database.ITransaction, progression models.LearningProgressionDB) (*models.LearningProgressionDB, error)
	ChangePointTransaction(tx database.ITransaction, userID int, point int, mode models.ChangePointMode) error
	GetCollectedBadge(userID int) ([]models.CorrectedBadgeDB, error)
	GetExamResult(userID int) ([]models.ExamResultDB, error)
	GetExamResultByID(userID int, examResultID int) ([]models.ExamResultDB, error)
	InsertUserBadgeTransaction(tx database.ITransaction, userBadge models.UserBadgeDB) (models.UserBadgeDB, error)
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

func (r userRepository) UpdatesByID(id int, updateData map[string]interface{}) error {
	err := r.database.GetDB().
		Table(models.TableName.User).
		Select("", utils.NewHelper().GetKeyList(updateData)).
		Where(models.IDName.User+" = ?", id).
		Updates(updateData).
		Error
	return err
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

func (r userRepository) GetAllBadge() ([]models.BadgeDB, error) {
	badge := make([]models.BadgeDB, 0)
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
		Order("ranking ASC").
		Order("name ASC").
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

func (r userRepository) GetUserHint(userID int, activityID int) ([]models.UserHintDB, error) {
	userhint := make([]models.UserHintDB, 0)

	hintSubquery := r.database.GetDB().
		Select("hint_id").
		Table(models.TableName.Hint).
		Where(models.IDName.Activity+" = ?", activityID)

	err := r.database.GetDB().
		Table(models.TableName.UserHint).
		Where(models.IDName.Hint+" IN (?)", hintSubquery).
		Where(models.IDName.User+" = ?", userID).
		Find(&userhint).
		Error

	return userhint, err
}

func (r userRepository) InsertUserHint(userHint models.UserHintDB) (*models.UserHintDB, error) {
	err := r.database.GetDB().
		Table(models.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

func (r userRepository) InsertUserHintTransaction(tx database.ITransaction, userHint models.UserHintDB) (*models.UserHintDB, error) {
	err := tx.GetDB().
		Table(models.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

func (r userRepository) InsertLearningProgressionTransaction(tx database.ITransaction, progression models.LearningProgressionDB) (*models.LearningProgressionDB, error) {
	err := tx.GetDB().
		Table(models.TableName.LearningProgression).
		Create(&progression).
		Error
	return &progression, err
}

func (r userRepository) ChangePointTransaction(tx database.ITransaction, userID int, point int, mode models.ChangePointMode) error {
	statement := fmt.Sprintf("UPDATE %s SET point = point %s %d WHERE %s = %d",
		models.TableName.User,
		mode,
		point,
		models.IDName.User,
		userID,
	)
	temp := map[string]interface{}{}
	err := tx.GetDB().Raw(statement).Find(&temp).Error
	return err
}

func (r userRepository) GetCollectedBadge(userID int) ([]models.CorrectedBadgeDB, error) {
	correctedBadge := make([]models.CorrectedBadgeDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.Badge).
		Select(
			models.TableName.Badge+".badge_id AS badge_id",
			models.TableName.Badge+".name AS badge_name",
			models.TableName.UserBadge+".user_id AS user_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s AND %s.%s = %d",
			models.TableName.UserBadge,
			models.TableName.UserBadge,
			models.IDName.Badge,
			models.TableName.Badge,
			models.IDName.Badge,
			models.TableName.UserBadge,
			models.IDName.User,
			userID,
		)).
		Find(&correctedBadge).
		Error
	return correctedBadge, err
}

func (r userRepository) GetExamResultByID(userID int, examResultID int) ([]models.ExamResultDB, error) {
	examResults := make([]models.ExamResultDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.ExamResult).
		Select(
			models.TableName.ExamResult+".exam_result_id AS exam_result_id",
			models.TableName.ExamResult+".exam_id AS exam_id",
			models.TableName.ExamResult+".user_id AS user_id",
			models.TableName.ExamResult+".is_passed AS is_passed",
			models.TableName.ExamResult+".created_timestamp AS created_timestamp",
			models.TableName.ExamResultActivity+".score AS score",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.ExamResultActivity,
			models.TableName.ExamResultActivity,
			models.IDName.ExamResult,
			models.TableName.ExamResult,
			models.IDName.ExamResult,
		)).
		Where(models.IDName.User+" = ?", userID).
		Where(models.TableName.ExamResult+"."+models.IDName.ExamResult+" = ?", examResultID).
		Find(&examResults).
		Error
	return examResults, err
}

func (r userRepository) GetExamResult(userID int) ([]models.ExamResultDB, error) {
	examResults := make([]models.ExamResultDB, 0)
	err := r.database.GetDB().
		Table(models.TableName.ExamResult).
		Select(
			models.TableName.ExamResult+".exam_result_id AS exam_result_id",
			models.TableName.ExamResult+".exam_id AS exam_id",
			models.TableName.ExamResult+".user_id AS user_id",
			models.TableName.ExamResult+".is_passed AS is_passed",
			models.TableName.ExamResult+".created_timestamp AS created_timestamp",
			models.TableName.ExamResultActivity+".score AS score",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			models.TableName.ExamResultActivity,
			models.TableName.ExamResultActivity,
			models.IDName.ExamResult,
			models.TableName.ExamResult,
			models.IDName.ExamResult,
		)).
		Where(models.IDName.User+" = ?", userID).
		Find(&examResults).
		Error
	return examResults, err
}

func (r userRepository) InsertUserBadgeTransaction(tx database.ITransaction, userBadge models.UserBadgeDB) (models.UserBadgeDB, error) {
	err := tx.GetDB().
		Table(models.TableName.UserBadge).
		Create(&userBadge).
		Error
	return userBadge, err
}
