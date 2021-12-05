package repositories

// repository.user.go
/**
 * 	This file is a part of repositories, used to do data manipulation of user
 */

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/utils"
	"fmt"
)

/**
 * 	This class manipulation user data to other application
 */
type userRepository struct {
	Database database.IDatabase // Database to do database manipulation
}

/**
 * Constructor creates a new userRepository instance
 *
 * @param   db    Database to data manipulation
 *
 * @return 	instance of userRepository
 */
func NewUserRepository(db database.IDatabase) userRepository {
	return userRepository{Database: db}
}

/**
 * Get user by email from the database
 *
 * @param 	email  Email for getting user data
 *
 * @return user data
 * @return the error of getting user data
 */
func (r userRepository) GetUserByEmail(email string) (*storages.UserDB, error) {
	user := storages.UserDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Where("email = ?", email).
		Find(&user).
		Error
	return &user, err
}

/**
 * Get user by id from the database
 *
 * @param 	id  User ID for getting user data
 *
 * @return user data
 * @return the error of getting user data
 */
func (r userRepository) GetUserByID(id int) (*storages.UserDB, error) {
	user := storages.UserDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Where(storages.IDName.User+" = ?", id).
		Find(&user).
		Error
	return &user, err
}

/**
 * Get user profile from the database
 *
 * @param 	id  User ID for getting user profile
 *
 * @return user profile
 * @return the error of getting user profile
 */
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

/**
 * Get learning progression of the user from the database
 *
 * @param 	id  User ID for getting learning progression of the user
 *
 * @return learning progression of the user
 * @return the error of getting learning progression of the user
 */
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

/**
 * Get all badges of the application from the database
 *
 * @return all badges of the application
 * @return the error of getting all badges of the application
 */
func (r userRepository) GetAllBadge() ([]storages.BadgeDB, error) {
	badge := make([]storages.BadgeDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.Badge).
		Find(&badge).
		Error
	return badge, err
}

/**
 * Get user badge data by user badge id from the database
 *
 * @param 	id  User Badge ID for getting user badge data
 *
 * @return user badge data
 * @return the error of getting user badge data
 */
func (r userRepository) GetUserBadge(id int) ([]storages.UserBadgeDB, error) {
	userBadgeDB := make([]storages.UserBadgeDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.UserBadge).
		Where(storages.IDName.User+" = ?", id).
		Find(&userBadgeDB).
		Error
	return userBadgeDB, err
}

/**
 * Get collected badges of the user from the database
 *
 * @param 	userID  User ID for getting collected badges of the user
 *
 * @return collected badges of the user
 * @return the error of getting collected badges of the user
 */
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

/**
 * Get point ranking of the user from the database
 *
 * @param 	id  User ID for getting point ranking of the user
 *
 * @return point ranking of the user
 * @return the error of getting point ranking of the user
 */
func (r userRepository) GetPointRanking(id int) (*storages.RankingDB, error) {
	ranking := storages.RankingDB{}
	err := r.Database.GetDB().
		Table(storages.ViewName.Ranking).
		Where(storages.IDName.User+" = ?", id).
		Find(&ranking).
		Error
	return &ranking, err
}

/**
 * Get all user point ranking from the database
 *
 * @return all user point ranking
 * @return the error of getting all user point ranking
 */
func (r userRepository) GetRankingLeaderBoard() ([]storages.RankingDB, error) {
	ranking := make([]storages.RankingDB, 0)
	err := r.Database.GetDB().
		Table(storages.ViewName.Ranking).
		Limit(20).
		Find(&ranking).
		Error
	return ranking, err
}

/**
 * Get User Hint data from the database
 *
 * @param 	userID  	User ID for getting User Hint data
 * @param 	activityID  Activity ID for getting User Hint data
 *
 * @return User Hint data
 * @return the error of getting User Hint data
 */
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

/**
 * Get exam results of the user from the database
 *
 * @param 	userID  	User ID for getting exam results of the user
 *
 * @return exam results of the user
 * @return the error of getting exam results of the user
 */
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

/**
 * Get exam results of the user by exam result id from the database
 *
 * @param 	userID  	User ID for getting exam results of the user
 *
 * @return exam results of the user
 * @return the error of getting exam results of the user
 */
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

/**
 * Insert user data into the database
 *
 * @param 	user  	User model for insert into the database
 *
 * @return inserted user
 * @return the error of inserting user
 */
func (r userRepository) InsertUser(user storages.UserDB) (*storages.UserDB, error) {
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Create(&user).
		Error
	return &user, err
}

/**
 * Insert user hint into the database
 *
 * @param 	userHint  	User hint model for insert into the database
 *
 * @return inserted user hint
 * @return the error of inserting user hint
 */
func (r userRepository) InsertUserHint(userHint storages.UserHintDB) (*storages.UserHintDB, error) {
	err := r.Database.GetDB().
		Table(storages.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

/**
 * Update user data into the database
 *
 * @param 	id  			User id to update into the database
 * @param 	updateData  	User data to update into the database
 *
 * @return inserted user hint
 * @return the error of updating user data
 */
func (r userRepository) UpdatesByID(id int, updateData map[string]interface{}) error {
	err := r.Database.GetDB().
		Table(storages.TableName.User).
		Select("", utils.NewHelper().GetKeyList(updateData)).
		Where(storages.IDName.User+" = ?", id).
		Updates(updateData).
		Error
	return err
}

/**
 * Insert user hint into the database by database transaction
 *
 * @param 	tx  		Transaction model to do database transaction
 * @param 	userHint  	User hint model for insert into the database
 *
 * @return inserted user hint
 * @return the error of inserting user hint
 */
func (r userRepository) InsertUserHintTransaction(tx database.ITransaction, userHint storages.UserHintDB) (*storages.UserHintDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

/**
 * Insert learning progression of the user into the database by database transaction
 *
 * @param 	tx  			Transaction model to do database transaction
 * @param 	progression  	Learning progession model for insert into the database
 *
 * @return inserted learning progression
 * @return the error of inserting learning progression
 */
func (r userRepository) InsertLearningProgressionTransaction(tx database.ITransaction, progression storages.LearningProgressionDB) (*storages.LearningProgressionDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.LearningProgression).
		Create(&progression).
		Error
	return &progression, err
}

/**
 * Insert user badge into the database by database transaction
 *
 * @param 	tx  			Transaction model to do database transaction
 * @param 	userBadge  		User badge model for insert into the database
 *
 * @return inserted user badge
 * @return the error of inserting user badge
 */
func (r userRepository) InsertUserBadgeTransaction(tx database.ITransaction, userBadge storages.UserBadgeDB) (*storages.UserBadgeDB, error) {
	err := tx.GetDB().
		Table(storages.TableName.UserBadge).
		Create(&userBadge).
		Error
	return &userBadge, err
}

/**
 * Update user point into the database by database transaction
 *
 * @param 	tx  		Transaction model to do database transaction
 * @param 	userID  	User ID to update point
 * @param 	point  		Point to change
 * @param 	mode  		Mode to change user point
 *
 * @return the error of updating user point
 */
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
