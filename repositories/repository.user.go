package repositories

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/general"
	"DatabaseCamp/utils"
	"fmt"
)

type userRepository struct {
	database database.IDatabase
}

type IUserReader interface {
	GetUserByEmail(email string) (*general.UserDB, error)
	GetUserByID(id int) (*general.UserDB, error)
	GetProfile(id int) (*general.ProfileDB, error)
	GetLearningProgression(id int) ([]general.LearningProgressionDB, error)
	GetAllBadge() ([]general.BadgeDB, error)
	GetUserBadge(id int) ([]general.UserBadgeDB, error)
	GetCollectedBadge(userID int) ([]general.CorrectedBadgeDB, error)
	GetPointRanking(id int) (*general.RankingDB, error)
	GetRankingLeaderBoard() ([]general.RankingDB, error)
	GetUserHint(userID int, activityID int) ([]general.UserHintDB, error)
	GetExamResult(userID int) ([]general.ExamResultDB, error)
	GetExamResultByID(userID int, examResultID int) ([]general.ExamResultDB, error)
}

type IUserWriter interface {
	InsertUser(user general.UserDB) (*general.UserDB, error)
	InsertUserHint(userHint general.UserHintDB) (*general.UserHintDB, error)
	UpdatesByID(id int, updateData map[string]interface{}) error
}

type IUserTransaction interface {
	InsertUserHintTransaction(tx database.ITransaction, userHint general.UserHintDB) (*general.UserHintDB, error)
	InsertLearningProgressionTransaction(tx database.ITransaction, progression general.LearningProgressionDB) (*general.LearningProgressionDB, error)
	InsertUserBadgeTransaction(tx database.ITransaction, userBadge general.UserBadgeDB) (*general.UserBadgeDB, error)
	ChangePointTransaction(tx database.ITransaction, userID int, point int, mode entities.ChangePointMode) error
}

type IUserRepository interface {
	IUserReader
	IUserWriter
	IUserTransaction
}

func NewUserRepository(db database.IDatabase) userRepository {
	return userRepository{database: db}
}

func (r userRepository) GetUserByEmail(email string) (*general.UserDB, error) {
	user := general.UserDB{}
	err := r.database.GetDB().
		Table(general.TableName.User).
		Where("email = ?", email).
		Find(&user).
		Error
	return &user, err
}

func (r userRepository) GetUserByID(id int) (*general.UserDB, error) {
	user := general.UserDB{}
	err := r.database.GetDB().
		Table(general.TableName.User).
		Where(general.IDName.User+" = ?", id).
		Find(&user).
		Error
	return &user, err
}

func (r userRepository) GetProfile(id int) (*general.ProfileDB, error) {
	profile := general.ProfileDB{}
	err := r.database.GetDB().
		Table(general.ViewName.Profile).
		Where(general.IDName.User+" = ?", id).
		Find(&profile).
		Error
	if profile == (general.ProfileDB{}) {
		return nil, nil
	}
	return &profile, err
}

func (r userRepository) GetLearningProgression(id int) ([]general.LearningProgressionDB, error) {
	learningProgrogression := make([]general.LearningProgressionDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.LearningProgression).
		Where(general.IDName.User+" = ?", id).
		Order("created_timestamp desc").
		Find(&learningProgrogression).
		Error
	return learningProgrogression, err
}

func (r userRepository) GetAllBadge() ([]general.BadgeDB, error) {
	badge := make([]general.BadgeDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.Badge).
		Find(&badge).
		Error
	return badge, err
}

func (r userRepository) GetUserBadge(id int) ([]general.UserBadgeDB, error) {
	badgePair := make([]general.UserBadgeDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.UserBadge).
		Where(general.IDName.User+" = ?", id).
		Find(&badgePair).
		Error
	return badgePair, err
}

func (r userRepository) GetCollectedBadge(userID int) ([]general.CorrectedBadgeDB, error) {
	correctedBadge := make([]general.CorrectedBadgeDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.Badge).
		Select(
			general.TableName.Badge+".badge_id AS badge_id",
			general.TableName.Badge+".name AS badge_name",
			general.TableName.UserBadge+".user_id AS user_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s AND %s.%s = %d",
			general.TableName.UserBadge,
			general.TableName.UserBadge,
			general.IDName.Badge,
			general.TableName.Badge,
			general.IDName.Badge,
			general.TableName.UserBadge,
			general.IDName.User,
			userID,
		)).
		Find(&correctedBadge).
		Error
	return correctedBadge, err
}

func (r userRepository) GetPointRanking(id int) (*general.RankingDB, error) {
	ranking := general.RankingDB{}
	err := r.database.GetDB().
		Table(general.ViewName.Ranking).
		Where(general.IDName.User+" = ?", id).
		Find(&ranking).
		Error
	return &ranking, err
}

func (r userRepository) GetRankingLeaderBoard() ([]general.RankingDB, error) {
	ranking := make([]general.RankingDB, 0)
	err := r.database.GetDB().
		Table(general.ViewName.Ranking).
		Limit(20).
		Order("ranking ASC").
		Order("name ASC").
		Order(models.IDName.User + " ASC").
		Find(&ranking).
		Error
	return ranking, err
}

func (r userRepository) GetUserHint(userID int, activityID int) ([]general.UserHintDB, error) {
	userhint := make([]general.UserHintDB, 0)

	hintSubquery := r.database.GetDB().
		Select("hint_id").
		Table(general.TableName.Hint).
		Where(general.IDName.Activity+" = ?", activityID)

	err := r.database.GetDB().
		Table(general.TableName.UserHint).
		Where(general.IDName.Hint+" IN (?)", hintSubquery).
		Where(general.IDName.User+" = ?", userID).
		Find(&userhint).
		Error

	return userhint, err
}

func (r userRepository) GetExamResult(userID int) ([]general.ExamResultDB, error) {
	examResults := make([]general.ExamResultDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.ExamResult).
		Select(
			general.TableName.ExamResult+".exam_result_id AS exam_result_id",
			general.TableName.ExamResult+".exam_id AS exam_id",
			general.TableName.ExamResult+".user_id AS user_id",
			general.TableName.ExamResult+".is_passed AS is_passed",
			general.TableName.ExamResult+".created_timestamp AS created_timestamp",

			general.TableName.ExamResultActivity+".score AS score",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.ExamResultActivity,
			general.TableName.ExamResultActivity,
			general.IDName.ExamResult,
			general.TableName.ExamResult,
			general.IDName.ExamResult,
		)).
		Where(general.IDName.User+" = ?", userID).
		Find(&examResults).
		Error
	return examResults, err
}

func (r userRepository) GetExamResultByID(userID int, examResultID int) ([]general.ExamResultDB, error) {
	examResults := make([]general.ExamResultDB, 0)
	err := r.database.GetDB().
		Table(general.TableName.ExamResult).
		Select(
			general.TableName.ExamResult+".exam_result_id AS exam_result_id",
			general.TableName.ExamResult+".exam_id AS exam_id",
			general.TableName.ExamResult+".user_id AS user_id",
			general.TableName.ExamResult+".is_passed AS is_passed",
			general.TableName.ExamResult+".created_timestamp AS created_timestamp",
			general.TableName.ExamResultActivity+".score AS score",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			general.TableName.ExamResultActivity,
			general.TableName.ExamResultActivity,
			general.IDName.ExamResult,
			general.TableName.ExamResult,
			general.IDName.ExamResult,
		)).
		Where(general.IDName.User+" = ?", userID).
		Where(general.TableName.ExamResult+"."+general.IDName.ExamResult+" = ?", examResultID).
		Find(&examResults).
		Error
	return examResults, err
}

func (r userRepository) InsertUser(user general.UserDB) (*general.UserDB, error) {
	err := r.database.GetDB().
		Table(general.TableName.User).
		Create(&user).
		Error
	return &user, err
}

func (r userRepository) InsertUserHint(userHint general.UserHintDB) (*general.UserHintDB, error) {
	err := r.database.GetDB().
		Table(general.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

func (r userRepository) UpdatesByID(id int, updateData map[string]interface{}) error {
	err := r.database.GetDB().
		Table(general.TableName.User).
		Select("", utils.NewHelper().GetKeyList(updateData)).
		Where(general.IDName.User+" = ?", id).
		Updates(updateData).
		Error
	return err
}

func (r userRepository) InsertUserHintTransaction(tx database.ITransaction, userHint general.UserHintDB) (*general.UserHintDB, error) {
	err := tx.GetDB().
		Table(general.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

func (r userRepository) InsertLearningProgressionTransaction(tx database.ITransaction, progression general.LearningProgressionDB) (*general.LearningProgressionDB, error) {
	err := tx.GetDB().
		Table(general.TableName.LearningProgression).
		Create(&progression).
		Error
	return &progression, err
}

func (r userRepository) InsertUserBadgeTransaction(tx database.ITransaction, userBadge general.UserBadgeDB) (*general.UserBadgeDB, error) {
	err := tx.GetDB().
		Table(general.TableName.UserBadge).
		Create(&userBadge).
		Error
	return &userBadge, err
}

func (r userRepository) ChangePointTransaction(tx database.ITransaction, userID int, point int, mode entities.ChangePointMode) error {
	statement := fmt.Sprintf("UPDATE %s SET point = point %s %d WHERE %s = %d",
		general.TableName.User,
		mode,
		point,
		general.IDName.User,
		userID,
	)
	temp := map[string]interface{}{}
	err := tx.GetDB().Raw(statement).Find(&temp).Error
	return err
}
