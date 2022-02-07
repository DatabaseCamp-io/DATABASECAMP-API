package repositories

import (
	"database-camp/internal/errs"
	"database-camp/internal/infrastructure/cache"
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/infrastructure/storage"
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/utils"
	"encoding/json"
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
	GetContentGroups() (groups content.ContentGroups, err error)
	GetCorrectProgression(activityID int) (progression *content.LearningProgression, err error)
	UseHint(userID int, reducePoint int, hintID int) error
}

type learningRepository struct {
	db    database.MysqlDB
	cache cache.Cache
}

func NewLearningRepository(db database.MysqlDB, cache cache.Cache) *learningRepository {
	return &learningRepository{db: db, cache: cache}
}

func (r learningRepository) GetContent(id int) (*content.Content, error) {
	content := content.Content{}

	key := "learningRepository::GetContent::" + utils.ParseString(id)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &content); err == nil {
			return &content, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.Content).
		Where(IDName.Content+" = ?", id).
		Find(&content).
		Error

	if data, err := json.Marshal(content); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return &content, err
}

func (r learningRepository) GetOverview() ([]content.Overview, error) {
	overview := make([]content.Overview, 0)

	key := "learningRepository::GetOverview"

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &overview); err == nil {
			return overview, nil
		}
	}

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

	if data, err := json.Marshal(overview); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return overview, err
}

func (r learningRepository) GetContentActivity(contentID int) ([]activity.Activity, error) {
	activity := make([]activity.Activity, 0)

	key := "learningRepository::GetContentActivity::" + utils.ParseString(contentID)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &activity); err == nil {
			return activity, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.Activity).
		Where(IDName.Content+" = ?", contentID).
		Find(&activity).
		Error

	if data, err := json.Marshal(activity); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return activity, err
}

func (r learningRepository) GetActivity(id int) (*activity.Activity, error) {
	activity := activity.Activity{}

	key := "learningRepository::GetActivity::" + utils.ParseString(id)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &activity); err == nil {
			return &activity, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.Activity).
		Where(IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	if data, err := json.Marshal(activity); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return &activity, err
}

func (r learningRepository) getMatchingChoice(activityID int) (activity.MatchingChoices, error) {
	matchingChoice := make([]activity.MatchingChoice, 0)

	key := "learningRepository::getMatchingChoice::" + utils.ParseString(activityID)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &matchingChoice); err == nil {
			return matchingChoice, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.MatchingChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	if data, err := json.Marshal(matchingChoice); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return matchingChoice, err
}

func (r learningRepository) getMultipleChoice(activityID int) (activity.MultipleChoices, error) {
	multipleChoice := make([]activity.MultipleChoice, 0)

	key := "learningRepository::getMultipleChoice::" + utils.ParseString(activityID)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &multipleChoice); err == nil {
			return multipleChoice, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.MultipleChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&multipleChoice).
		Error

	if data, err := json.Marshal(multipleChoice); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return multipleChoice, err
}

func (r learningRepository) getCompletionChoice(activityID int) (activity.CompletionChoices, error) {
	completionChoice := make([]activity.CompletionChoice, 0)

	key := "learningRepository::getCompletionChoice::" + utils.ParseString(activityID)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &completionChoice); err == nil {
			return completionChoice, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.CompletionChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&completionChoice).
		Error

	if data, err := json.Marshal(completionChoice); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return completionChoice, err
}

func (r learningRepository) getVocabGroupChoice(activityID int) (activity.VocalGroupChoices, error) {
	vocalGroupChoices := make([]activity.VocalGroupChoice, 0)

	key := "learningRepository::getVocabGroupChoice::" + utils.ParseString(activityID)
	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &vocalGroupChoices); err == nil {
			return vocalGroupChoices, nil
		}
	}

	err := r.db.GetDB().
		Select("vocab_group_name", "vocab").
		Table(TableName.VocabGroupChoice).
		Joins("INNER JOIN %s ON %s.%s = %s.%s",
			TableName.VocabGroup,
			TableName.VocabGroup,
			IDName.VocabGroup,
			TableName.VocabGroupChoice,
			IDName.VocabGroup,
		).
		Where(IDName.Activity+" = ?", activityID).
		Find(&vocalGroupChoices).
		Error

	if data, err := json.Marshal(vocalGroupChoices); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return vocalGroupChoices, err
}

func (r learningRepository) getDependencyChoice(activityID int) (*activity.DependencyChoice, error) {
	choice := activity.DependencyChoice{}

	key := "learningRepository::getDependencyChoice::" + utils.ParseString(activityID)
	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &choice); err == nil {
			return &choice, nil
		}
	}

	err := r.db.GetDB().
		Preload(TableName.Dependency).
		Preload(TableName.Determinant).
		Table(TableName.DependencyChoice).
		Where(IDName.Activity+" = ?", activityID).
		Find(&choice).
		Error

	if data, err := json.Marshal(choice); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return &choice, err
}

func (r learningRepository) GetActivityHints(activityID int) ([]activity.Hint, error) {
	hints := make([]activity.Hint, 0)

	key := "learningRepository::GetActivityHints::" + utils.ParseString(activityID)

	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &hints); err == nil {
			return hints, nil
		}
	}

	err := r.db.GetDB().
		Table(TableName.Hint).
		Where(IDName.Activity+" = ?", activityID).
		Order("level ASC").
		Find(&hints).
		Error

	if data, err := json.Marshal(hints); err != nil {
		return nil, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return nil, err
		}
	}

	return hints, err
}

func (r learningRepository) GetVideoFileLink(objectName string) (string, error) {
	storage := storage.GetCloudStorageServiceInstance()

	var link string
	var err error

	key := "learningRepository::GetVideoFileLink::" + objectName

	if cacheData, err := r.cache.Get(key); err == nil {
		return cacheData, nil
	}

	link, err = storage.GetFileLink(objectName)
	if err != nil {
		return "", err
	}

	err = r.cache.Set(key, link, time.Minute*15)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (r learningRepository) GetActivityChoices(activityID int, activityTypeID int) (activity.Choices, error) {
	switch activityTypeID {
	case 1:
		return r.getMatchingChoice(activityID)
	case 2:
		return r.getMultipleChoice(activityID)
	case 3:
		return r.getCompletionChoice(activityID)
	case 4:
		return r.getVocabGroupChoice(activityID)
	case 5:
		return r.getDependencyChoice(activityID)
	default:
		return nil, errs.ErrActivityTypeInvalid
	}
}

func (r learningRepository) GetContentGroups() (groups content.ContentGroups, err error) {
	err = r.db.GetDB().
		Table(TableName.ContentGroup).
		Find(&groups).
		Error
	return
}

func (r learningRepository) GetCorrectProgression(activityID int) (progression *content.LearningProgression, err error) {
	err = r.db.GetDB().
		Table(TableName.LearningProgression).
		Where("is_correct = 1").
		Find(&progression).
		Error
	return
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
