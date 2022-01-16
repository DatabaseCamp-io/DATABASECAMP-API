package repositories

import (
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/badge"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/utils"
	"fmt"
)

type UserRepository interface {
	GetUserByEmail(email string) (*user.User, error)
	GetUserByID(id int) (*user.User, error)
	GetProfile(id int) (*user.Profile, error)
	GetLearningProgression(id int) ([]content.LearningProgression, error)
	GetAllBadge() ([]badge.Badge, error)
	GetUserBadge(id int) ([]badge.UserBadge, error)
	GetCollectedBadge(userID int) ([]user.CorrectedBadge, error)
	GetPointRanking(id int) (*user.Ranking, error)
	GetRankingLeaderBoard() ([]user.Ranking, error)
	GetUserHint(userID int, activityID int) ([]activity.UserHint, error)
	InsertUser(user user.User) (*user.User, error)
	InsertUserHint(userHint activity.UserHint) (*activity.UserHint, error)
	InsertBadge(userBadge badge.UserBadge) (*badge.UserBadge, error)
	UpdatesByID(id int, updateData map[string]interface{}) error
}

type userRepository struct {
	db database.MysqlDB
}

func NewUserRepository(db database.MysqlDB) *userRepository {
	return &userRepository{db: db}
}

func (r userRepository) GetUserByEmail(email string) (*user.User, error) {
	user := user.User{}
	err := r.db.GetDB().
		Table(TableName.User).
		Where("email = ?", email).
		Find(&user).
		Error
	return &user, err
}

func (r userRepository) GetUserByID(id int) (*user.User, error) {
	user := user.User{}
	err := r.db.GetDB().
		Table(TableName.User).
		Where(IDName.User+" = ?", id).
		Find(&user).
		Error
	return &user, err
}

func (r userRepository) GetProfile(id int) (*user.Profile, error) {
	profile := user.Profile{}
	err := r.db.GetDB().
		Table(ViewName.Profile).
		Where(IDName.User+" = ?", id).
		Find(&profile).
		Error
	return &profile, err
}

func (r userRepository) GetLearningProgression(id int) ([]content.LearningProgression, error) {
	progresstion := make([]content.LearningProgression, 0)
	err := r.db.GetDB().
		Table(TableName.LearningProgression).
		Where(IDName.User+" = ?", id).
		Order("created_timestamp desc").
		Find(&progresstion).
		Error
	return progresstion, err
}

func (r userRepository) GetAllBadge() ([]badge.Badge, error) {
	badge := make([]badge.Badge, 0)
	err := r.db.GetDB().
		Table(TableName.Badge).
		Find(&badge).
		Error
	return badge, err
}

func (r userRepository) GetUserBadge(id int) ([]badge.UserBadge, error) {
	userBadgeDB := make([]badge.UserBadge, 0)
	err := r.db.GetDB().
		Table(TableName.UserBadge).
		Where(IDName.User+" = ?", id).
		Find(&userBadgeDB).
		Error
	return userBadgeDB, err
}

func (r userRepository) GetCollectedBadge(userID int) ([]user.CorrectedBadge, error) {
	correctedBadge := make([]user.CorrectedBadge, 0)
	err := r.db.GetDB().
		Table(TableName.Badge).
		Select(
			TableName.Badge+".badge_id AS badge_id",
			TableName.Badge+".name AS badge_name",
			TableName.UserBadge+".user_id AS user_id",
		).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s AND %s.%s = %d",
			TableName.UserBadge,
			TableName.UserBadge,
			IDName.Badge,
			TableName.Badge,
			IDName.Badge,
			TableName.UserBadge,
			IDName.User,
			userID,
		)).
		Find(&correctedBadge).
		Error
	return correctedBadge, err
}

func (r userRepository) GetPointRanking(id int) (*user.Ranking, error) {
	ranking := user.Ranking{}
	err := r.db.GetDB().
		Table(ViewName.Ranking).
		Where(IDName.User+" = ?", id).
		Find(&ranking).
		Error
	return &ranking, err
}

func (r userRepository) GetRankingLeaderBoard() ([]user.Ranking, error) {
	ranking := make([]user.Ranking, 0)
	err := r.db.GetDB().
		Table(ViewName.Ranking).
		Limit(20).
		Find(&ranking).
		Error
	return ranking, err
}

func (r userRepository) GetUserHint(userID int, activityID int) ([]activity.UserHint, error) {
	userhint := make([]activity.UserHint, 0)

	hintSubquery := r.db.GetDB().
		Select("hint_id").
		Table(TableName.Hint).
		Where(IDName.Activity+" = ?", activityID)

	err := r.db.GetDB().
		Table(TableName.UserHint).
		Where(IDName.Hint+" IN (?)", hintSubquery).
		Where(IDName.User+" = ?", userID).
		Find(&userhint).
		Error

	return userhint, err
}

func (r userRepository) InsertUser(user user.User) (*user.User, error) {
	err := r.db.GetDB().
		Table(TableName.User).
		Create(&user).
		Error
	return &user, err
}

func (r userRepository) InsertUserHint(userHint activity.UserHint) (*activity.UserHint, error) {
	err := r.db.GetDB().
		Table(TableName.UserHint).
		Create(&userHint).
		Error
	return &userHint, err
}

func (r userRepository) InsertBadge(userBadge badge.UserBadge) (*badge.UserBadge, error) {
	err := r.db.GetDB().
		Table(TableName.UserBadge).
		Create(&userBadge).
		Error
	return &userBadge, err
}

func (r userRepository) UpdatesByID(id int, updateData map[string]interface{}) error {
	err := r.db.GetDB().
		Table(TableName.User).
		Select("", utils.GetKeyList(updateData)).
		Where(IDName.User+" = ?", id).
		Updates(updateData).
		Error
	return err
}
