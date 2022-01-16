package repositories

import (
	"database-camp/internal/errs"
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/infrastructure/storage"
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/utils"
	"fmt"
	"time"
)

type LearningRepository interface {
	GetContent(id int) (*content.Content, error)
	GetOverview() ([]content.Overview, error)
	GetActivity(id int) (*activity.Activity, error)
	GetActivityHints(activityID int) ([]activity.Hint, error)
	GetContentActivity(contentID int) ([]activity.Activity, error)
	GetVideoFileLink(imagekey string) (string, error)
	GetActivityChoices(activityID int, activityTypeID int) (activity.Choices, error)
	UseHint(userID int, reducePoint int, hintID int) error
	InsertActivityResult(userID int, activityID int, point int) error
}

type learningRepository struct {
	db database.MysqlDB
}

func NewLearningRepository(db database.MysqlDB) *learningRepository {
	return &learningRepository{db: db}
}

func (r learningRepository) GetContent(id int) (*content.Content, error) {
	content := content.Content{}
	err := r.db.GetDB().
		Table(TableName.Content).
		Where(IDName.Content+" = ?", id).
		Find(&content).
		Error
	return &content, err
}

func (r learningRepository) GetOverview() ([]content.Overview, error) {
	overview := make([]content.Overview, 0)
	err := r.db.GetDB().
		Table(TableName.ContentGroup).
		Select("ContentGroup.content_group_id AS content_group_id",
			"Content.content_id AS content_id",
			"Activity.activity_id AS activity_id",
			"ContentGroup.name AS group_name",
			"Content.name AS content_name",
		).
		Joins("LEFT JOIN Content ON ContentGroup.content_group_id = Content.content_group_id").
		Joins("LEFT JOIN Activity ON Content.content_id = Activity.content_id").
		Order("content_group_id ASC").
		Find(&overview).
		Error
	return overview, err
}

func (r learningRepository) GetContentActivity(contentID int) ([]activity.Activity, error) {
	activity := make([]activity.Activity, 0)

	err := r.db.GetDB().
		Table(TableName.Activity).
		Where(IDName.Content+" = ?", contentID).
		Find(&activity).
		Error

	return activity, err
}

func (r learningRepository) GetActivity(id int) (*activity.Activity, error) {
	activity := activity.Activity{}

	err := r.db.GetDB().
		Table(TableName.Activity).
		Where(IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	return &activity, err
}

func (r learningRepository) getMatchingChoice(activityID int) (activity.MatchingChoices, error) {
	matchingChoice := make([]activity.MatchingChoice, 0)

	err := r.db.GetDB().
		Table(TableName.MatchingChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	return matchingChoice, err
}

func (r learningRepository) getMultipleChoice(activityID int) (activity.MultipleChoices, error) {
	multipleChoice := make([]activity.MultipleChoice, 0)

	err := r.db.GetDB().
		Table(TableName.MultipleChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&multipleChoice).
		Error

	return multipleChoice, err
}

func (r learningRepository) getCompletionChoice(activityID int) (activity.CompletionChoices, error) {
	completionChoice := make([]activity.CompletionChoice, 0)

	err := r.db.GetDB().
		Table(TableName.CompletionChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&completionChoice).
		Error

	return completionChoice, err
}

func (r learningRepository) GetActivityHints(activityID int) ([]activity.Hint, error) {
	hints := make([]activity.Hint, 0)

	err := r.db.GetDB().
		Table(TableName.Hint).
		Where(IDName.Activity+" = ?", activityID).
		Order("level ASC").
		Find(&hints).
		Error

	return hints, err
}

func (r learningRepository) GetVideoFileLink(objectName string) (string, error) {
	storage := storage.GetCloudStorageServiceInstance()
	return storage.GetFileLink(objectName)
}

func (r learningRepository) GetActivityChoices(activityID int, activityTypeID int) (activity.Choices, error) {
	switch activityTypeID {
	case 1:
		return r.getMatchingChoice(activityID)
	case 2:
		return r.getMultipleChoice(activityID)
	case 3:
		return r.getCompletionChoice(activityID)
	default:
		return nil, errs.ErrActivityTypeInvalid
	}
}

func (r learningRepository) UseHint(userID int, reducePoint int, hintID int) error {
	routine := 0
	errs := make(chan error, 2)
	tx := r.db.GetDB().Begin()

	go func() {
		statement := fmt.Sprintf("UPDATE %s SET point = point - %d WHERE %s = %d",
			TableName.User,
			reducePoint,
			IDName.User,
			userID,
		)
		temp := map[string]interface{}{}
		errs <- tx.Raw(statement).Find(&temp).Error
	}()

	go func() {
		hint := activity.UserHint{
			UserID:           userID,
			HintID:           hintID,
			CreatedTimestamp: time.Now().Local(),
		}
		errs <- tx.Table(TableName.UserHint).Create(&hint).Error
	}()

	for err := range errs {
		routine++
		if err != nil {
			close(errs)
			tx.Rollback()
			return err
		}
		if routine == 2 {
			close(errs)
		}
	}

	tx.Commit()
	return nil
}

func (r learningRepository) InsertActivityResult(userID int, activityID int, point int) error {
	tx := r.db.GetDB().Begin()

	progression := content.LearningProgression{
		UserID:           userID,
		ActivityID:       activityID,
		CreatedTimestamp: time.Now().Local(),
	}

	err := tx.Table(TableName.LearningProgression).Create(&progression).Error
	if err != nil && !utils.IsSqlDuplicateError(err) {
		tx.Rollback()
		return err
	}

	if !utils.IsSqlDuplicateError(err) {
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
