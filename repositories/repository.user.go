package repositories

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/utils"
	"fmt"
)

type userRepository struct {
	Database database.IDatabase
}
//Interface that show how others function call and use function in this module
type IUserRepository interface {
	GetUserByEmail(email string) (*storages.UserDB, error)
	GetUserByID(id int) (*storages.UserDB, error)
	GetProfile(id int) (*storages.ProfileDB, error)
	GetLearningProgression(id int) ([]storages.LearningProgressionDB, error)
	GetAllBadge() ([]storages.BadgeDB, error)
	GetUserBadge(id int) ([]storages.UserBadgeDB, error)
	GetCollectedBadge(userID int) ([]storages.CorrectedBadgeDB, error)
	GetPointRanking(id int) (*storages.RankingDB, error)
	GetRankingLeaderBoard() ([]storages.RankingDB, error)
	GetUserHint(userID int, activityID int) ([]storages.UserHintDB, error)
	GetExamResult(userID int) ([]storages.ExamResultDB, error)
	GetExamResultByID(userID int, examResultID int) ([]storages.ExamResultDB, error)

	InsertUser(user storages.UserDB) (*storages.UserDB, error)
	InsertUserHint(userHint storages.UserHintDB) (*storages.UserHintDB, error)
	UpdatesByID(id int, updateData map[string]interface{}) error

	InsertUserHintTransaction(tx database.ITransaction, userHint storages.UserHintDB) (*storages.UserHintDB, error)
	InsertLearningProgressionTransaction(tx database.ITransaction, progression storages.LearningProgressionDB) (*storages.LearningProgressionDB, error)
	InsertUserBadgeTransaction(tx database.ITransaction, userBadge storages.UserBadgeDB) (*storages.UserBadgeDB, error)
	ChangePointTransaction(tx database.ITransaction, userID int, point int, mode entities.ChangePointMode) error
}
// 
func NewUserRepository(db database.IDatabase) userRepository {
	return userRepository{Database: db}
}
// Get user email from database
func (r userRepository) GetUserByEmail(email string) (*storages.UserDB, error) {
	user := storages.UserDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Where("email = ?", email).
		Find(&user).
		Error
	return &user, err
}
// Get user Id from database
func (r userRepository) GetUserByID(id int) (*storages.UserDB, error) {
	user := storages.UserDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Where(storages.IDName.User+" = ?", id).
		Find(&user).
		Error
	return &user, err
}
// Get user Profile from database
func (r userRepository) GetProfile(id int) (*storages.ProfileDB, error) {
	profile := storages.ProfileDB{}
	err := r.Database.GetDB().
		Table(storages.ViewName.Profile).
		Where(storages.IDName.User+" = ?", id).
		Find(&profile).
		Error
	if profile == (storages.ProfileDB{}) {
		return nil, nil
	}
	return &profile, err
}
// Get use learning progression from database
// use time stamp
func (r userRepository) GetLearningProgression(id int) ([]storages.LearningProgressionDB, error) {
	learningProgrogression := make([]storages.LearningProgressionDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.LearningProgression).
		Where(storages.IDName.User+" = ?", id).
		Order("created_timestamp desc").
		Find(&learningProgrogression).
		Error
	return learningProgrogression, err
}
// Get all badge from database
func (r userRepository) GetAllBadge() ([]storages.BadgeDB, error) {
	badge := make([]storages.BadgeDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.Badge).
		Find(&badge).
		Error
	return badge, err
}
// Get user badge from database
func (r userRepository) GetUserBadge(id int) ([]storages.UserBadgeDB, error) {
	badgePair := make([]storages.UserBadgeDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.UserBadge).
		Where(storages.IDName.User+" = ?", id).
		Find(&badgePair).
		Error
	return badgePair, err
}
// Get badge that user already collect from database
func (r userRepository) GetCollectedBadge(userID int) ([]storages.CorrectedBadgeDB, error) {
	correctedBadge := make([]storages.CorrectedBadgeDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.Badge).
		Select(
			storages.TableName.Badge+".badge_id AS badge_id",
			storages.TableName.Badge+".name AS badge_name",
			storages.TableName.UserBadge+".user_id AS user_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s AND %s.%s = %d",
			storages.TableName.UserBadge,
			storages.TableName.UserBadge,
			storages.IDName.Badge,
			storages.TableName.Badge,
			storages.IDName.Badge,
			storages.TableName.UserBadge,
			storages.IDName.User,
			userID,
		)).
		Find(&correctedBadge).
		Error
	return correctedBadge, err
}
// Get user point ranking from database
func (r userRepository) GetPointRanking(id int) (*storages.RankingDB, error) {
	ranking := storages.RankingDB{}
	err := r.Database.GetDB().
		Table(storages.ViewName.Ranking).
		Where(storages.IDName.User+" = ?", id).
		Find(&ranking).
		Error
	return &ranking, err
}
// Get leaderboard from database
// sort by point
func (r userRepository) GetRankingLeaderBoard() ([]storages.RankingDB, error) {
	ranking := make([]storages.RankingDB, 0)
	err := r.Database.GetDB().
		Table(storages.ViewName.Ranking).
		Limit(20).
		Order("ranking ASC").
		Order("name ASC").
		Order(storages.IDName.User + " ASC").
		Find(&ranking).
		Error
	return ranking, err
}
// Get hint that user already use from database
func (r userRepository) GetUserHint(userID int, activityID int) ([]storages.UserHintDB, error) {
	userhint := make([]storages.UserHintDB, 0)

	hintSubquery := r.Database.GetDB().
		Select("hint_id").
		Table(storages.TableName.Hint).
		Where(storages.IDName.Activity+" = ?", activityID)

	err := r.Database.GetDB().
		Table(storages.TableName.UserHint).
		Where(storages.IDName.Hint+" IN (?)", hintSubquery).
		Where(storages.IDName.User+" = ?", userID).
		Find(&userhint).
		Error

	return userhint, err
}
// Get user exam result from database
func (r userRepository) GetExamResult(userID int) ([]storages.ExamResultDB, error) {
	examResults := make([]storages.ExamResultDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.ExamResult).
		Select(
			storages.TableName.ExamResult+".exam_result_id AS exam_result_id",
			storages.TableName.ExamResult+".exam_id AS exam_id",
			storages.TableName.ExamResult+".user_id AS user_id",
			storages.TableName.ExamResult+".is_passed AS is_passed",
			storages.TableName.ExamResult+".created_timestamp AS created_timestamp",

			storages.TableName.ExamResultActivity+".score AS score",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.ExamResultActivity,
			storages.TableName.ExamResultActivity,
			storages.IDName.ExamResult,
			storages.TableName.ExamResult,
			storages.IDName.ExamResult,
		)).
		Where(storages.IDName.User+" = ?", userID).
		Find(&examResults).
		Error
	return examResults, err
}
// Get user exam result from database by exam ID 
func (r userRepository) GetExamResultByID(userID int, examResultID int) ([]storages.ExamResultDB, error) {
	examResults := make([]storages.ExamResultDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.ExamResult).
		Select(
			storages.TableName.ExamResult+".exam_result_id AS exam_result_id",
			storages.TableName.ExamResult+".exam_id AS exam_id",
			storages.TableName.ExamResult+".user_id AS user_id",
			storages.TableName.ExamResult+".is_passed AS is_passed",
			storages.TableName.ExamResult+".created_timestamp AS created_timestamp",
			storages.TableName.ExamResultActivity+".score AS score",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			storages.TableName.ExamResultActivity,
			storages.TableName.ExamResultActivity,
			storages.IDName.ExamResult,
			storages.TableName.ExamResult,
			storages.IDName.ExamResult,
		)).
		Where(storages.IDName.User+" = ?", userID).
		Where(storages.TableName.ExamResult+"."+storages.IDName.ExamResult+" = ?", examResultID).
		Find(&examResults).
		Error
	return examResults, err
}
// insert new user into database
func (r userRepository) InsertUser(user storages.UserDB) (*storages.UserDB, error) {
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Create(&user).
		Error
	return &user, err
}
// insert new user hint into database
func (r userRepository) InsertUserHint(userHint storages.UserHintDB) (*storages.UserHintDB, error) {
	err := r.Database.GetDB().
		Table(storages.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}
// Update user information into database
func (r userRepository) UpdatesByID(id int, updateData map[string]interface{}) error {
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Select("", utils.NewHelper().GetKeyList(updateData)).
		Where(storages.IDName.User+" = ?", id).
		Updates(updateData).
		Error
	return err
}
// Insert user hint trtansaction into database
func (r userRepository) InsertUserHintTransaction(tx database.ITransaction, userHint storages.UserHintDB) (*storages.UserHintDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}
// Insert user learning progression trtansaction into database
func (r userRepository) InsertLearningProgressionTransaction(tx database.ITransaction, progression storages.LearningProgressionDB) (*storages.LearningProgressionDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.LearningProgression).
		Create(&progression).
		Error
	return &progression, err
}
// Insert user badge trtansaction into database
func (r userRepository) InsertUserBadgeTransaction(tx database.ITransaction, userBadge storages.UserBadgeDB) (*storages.UserBadgeDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.UserBadge).
		Create(&userBadge).
		Error
	return &userBadge, err
}
// Update point transaction into database
func (r userRepository) ChangePointTransaction(tx database.ITransaction, userID int, point int, mode entities.ChangePointMode) error {
	statement := fmt.Sprintf("UPDATE %s SET point = point %s %d WHERE %s = %d",
		storages.TableName.User,
		mode,
		point,
		storages.IDName.User,
		userID,
	)
	temp := map[string]interface{}{}
	err := tx.GetDB().Raw(statement).Find(&temp).Error
	return err
}
