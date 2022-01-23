package repositories

import (
	"database-camp/internal/infrastructure/cache"
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/badge"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/utils"
	"encoding/json"
	"fmt"
	"time"
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
	GetPreExamID(userID int) (*int, error)
	GetPreTestResults(userID int) (user.PreTestResults, error)
	InsertUser(user user.User) (*user.User, error)
	InsertUserHint(userHint activity.UserHint) (*activity.UserHint, error)
	InsertBadge(userBadge badge.UserBadge) (*badge.UserBadge, error)
	InsertLearningProgression(userID int, activityID int, point int, isCorrect bool) error
	UpdatesByID(id int, updateData map[string]interface{}) error
}

type userRepository struct {
	db    database.MysqlDB
	cache cache.Cache
}

func NewUserRepository(db database.MysqlDB, cache cache.Cache) *userRepository {
	return &userRepository{db: db, cache: cache}
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
		Where(IDName.User+" = ? AND is_correct = 1", id).
		Order("created_timestamp desc").
		Group(IDName.Activity).
		Find(&progresstion).
		Error

	return progresstion, err
}

func (r userRepository) GetAllBadge() ([]badge.Badge, error) {
	badge := make([]badge.Badge, 0)

	key := "userRepository::GetAllBadge"

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &badge); err == nil {
			return badge, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.Badge).
		Find(&badge).
		Error

	if data, err := json.Marshal(badge); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*10); err != nil {
			return nil, err
		}
	}

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

	key := "userRepository::GetPointRanking::" + utils.ParseString(id)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &ranking); err == nil {
			return &ranking, nil
		}
	}

	err := r.db.GetDB().
		Table(ViewName.Ranking).
		Where(IDName.User+" = ?", id).
		Find(&ranking).
		Error

	if data, err := json.Marshal(ranking); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*10); err != nil {
			return nil, err
		}
	}

	return &ranking, err
}

func (r userRepository) GetRankingLeaderBoard() ([]user.Ranking, error) {
	ranking := make([]user.Ranking, 0)

	key := "userRepository::GetRankingLeaderBoard"

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &ranking); err == nil {
			return ranking, nil
		}
	}

	err := r.db.GetDB().
		Table(ViewName.Ranking).
		Limit(20).
		Find(&ranking).
		Error

	if data, err := json.Marshal(ranking); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*10); err != nil {
			return nil, err
		}
	}

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

func (r userRepository) InsertLearningProgression(userID int, activityID int, point int, isCorrect bool) error {
	tx := r.db.GetDB().Begin()

	progression := content.LearningProgression{
		UserID:           userID,
		ActivityID:       activityID,
		IsCorrect:        isCorrect,
		CreatedTimestamp: time.Now().Local(),
	}

	err := tx.Table(TableName.LearningProgression).Create(&progression).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if isCorrect {
		statement := fmt.Sprintf("UPDATE %s SET point = point + %d WHERE %s = %d",
			TableName.User,
			point,
			IDName.User,
			userID,
		)
		temp := map[string]interface{}{}
		err = tx.Raw(statement).Find(&temp).Error
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
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

func (r userRepository) GetPreExamID(userID int) (*int, error) {
	data := struct {
		UserID int `gorm:"column:user_id"`
		ExamID int `gorm:"column:exam_id"`
	}{}

	err := r.db.GetDB().
		Table(ViewName.UserPreTest).
		Where(IDName.User+"=?", userID).
		Find(&data).
		Error

	return &data.ExamID, err
}

func (r userRepository) GetPreTestResults(userID int) (user.PreTestResults, error) {
	var results user.PreTestResults

	err := r.db.GetDB().
		Table(ViewName.UserPreTestResult).
		Where(IDName.User+"=?", userID).
		Order("created_timestamp ASC").
		Limit(1).
		Find(&results).
		Error

	return results, err
}
