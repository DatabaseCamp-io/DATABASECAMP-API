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
	"sync"
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
	GetPeerChoice(erAnswerID *int) (activity.ERAnswer, error)
	GetERChoice(activityID int) (activity.ERChoice, error)
	UseHint(userID int, reducePoint int, hintID int) error
	InsertERAnswer(answer activity.ERAnswer, userID int) error
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

func (r learningRepository) getVocabGroupChoice(activityID int) (activity.VocabGroupChoice, error) {
	vocalGroupChoice := activity.VocabGroupChoice{}

	key := "learningRepository::getVocabGroupChoice::" + utils.ParseString(activityID)
	if cacheData, err := r.cache.Get(key); err == nil {
		if err = json.Unmarshal([]byte(cacheData), &vocalGroupChoice); err == nil {
			return vocalGroupChoice, nil
		}
	}

	rows, err := r.db.GetDB().
		Select("name", "vocab").
		Table(TableName.VocabGroupChoice).
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.%s = %s.%s",
			TableName.VocabGroup,
			TableName.VocabGroup,
			IDName.VocabGroup,
			TableName.VocabGroupChoice,
			IDName.VocabGroup,
		)).
		Where(IDName.Activity+" = ?", activityID).
		Rows()

	groupMap := map[string]*activity.VocabGroup{}

	for rows.Next() {
		var name string
		var vocab string

		err = rows.Scan(&name, &vocab)
		if err != nil {
			return vocalGroupChoice, err
		}

		if _, ok := groupMap[name]; !ok {
			groupMap[name] = &activity.VocabGroup{
				GroupName: name,
				Vocabs:    []string{vocab},
			}
		} else {
			groupMap[name].Vocabs = append(groupMap[name].Vocabs, vocab)
		}
	}

	for _, v := range groupMap {
		vocalGroupChoice.Groups = append(vocalGroupChoice.Groups, *v)
	}

	if data, err := json.Marshal(vocalGroupChoice); err != nil {
		return vocalGroupChoice, err
	} else {
		if err = r.cache.Set(key, string(data), time.Minute*300); err != nil {
			return vocalGroupChoice, err
		}
	}

	return vocalGroupChoice, err
}

func (r learningRepository) getDependencyChoice(activityID int) (activity.DependencyChoice, error) {
	choice := activity.DependencyChoice{}

	rows, err := r.db.GetDB().
		Table(TableName.DependencyChoice).
		Select(
			TableName.DependencyChoice+"."+IDName.DependencyChoice,
			TableName.Dependency+"."+IDName.Dependency,
			TableName.Dependency+".dependent",
			TableName.Dependency+".fixed",
			TableName.Determinant+".value",
			TableName.Determinant+".fixed",
		).
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.%s = %s.%s",
			TableName.Dependency,
			TableName.Dependency,
			IDName.DependencyChoice,
			TableName.DependencyChoice,
			IDName.DependencyChoice,
		)).
		Joins(fmt.Sprintf("INNER JOIN %s ON %s.%s = %s.%s",
			TableName.Determinant,
			TableName.Determinant,
			IDName.Dependency,
			TableName.Dependency,
			IDName.Dependency,
		)).
		Where(IDName.Activity+" = ?", activityID).
		Rows()

	if err != nil {
		return choice, err
	}

	defer rows.Close()

	choice.Dependencies = make([]activity.Dependency, 0)

	dependencyMap := map[int]*activity.Dependency{}

	for rows.Next() {
		var id int
		dependency := activity.Dependency{}
		determinant := activity.Determinant{}

		err = rows.Scan(&choice.ID, &id, &dependency.Dependent, &dependency.Fixed, &determinant.Value, &determinant.Fixed)
		if err != nil {
			return choice, err
		}

		if _, ok := dependencyMap[id]; !ok {
			dependencyMap[id] = &dependency
			dependencyMap[id].Determinants = make([]activity.Determinant, 0)
		}

		dependencyMap[id].Determinants = append(dependencyMap[id].Determinants, determinant)
	}

	for _, v := range dependencyMap {
		choice.Dependencies = append(choice.Dependencies, *v)
	}

	return choice, err
}

func (r learningRepository) GetERChoice(activityID int) (activity.ERChoice, error) {
	choice := activity.ERChoice{}

	rows, err := r.db.GetDB().Debug().
		Select(
			TableName.ERChoice+".type",
			TableName.Tables+"."+IDName.Table,
			TableName.Tables+".title",
			TableName.Tables+".fixed",
			TableName.Attributes+"."+IDName.Attribute,
			TableName.Attributes+".value",
			TableName.Attributes+".key",
			TableName.Attributes+".fixed",
		).
		Table(TableName.ERChoice).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			TableName.ERChoiceTables,
			TableName.ERChoiceTables,
			IDName.ERChoice,
			TableName.ERChoice,
			IDName.ERChoice,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			TableName.Tables,
			TableName.Tables,
			IDName.Table,
			TableName.ERChoiceTables,
			IDName.Table,
		)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s",
			TableName.Attributes,
			TableName.Attributes,
			IDName.Table,
			TableName.Tables,
			IDName.Table,
		)).
		Where(IDName.Activity+" = ?", activityID).
		Rows()

	if err != nil {
		return choice, err
	}

	tablesMap := map[string]*activity.Table{}

	for rows.Next() {

		var attributeID *int
		var attributeValue *string
		var attributeKey *string
		var attributeFixed *bool

		var attribute activity.Attribute

		table := activity.Table{}

		err = rows.Scan(&choice.Type, &table.ID, &table.Title, &table.Fixed, &attributeID, &attributeValue, &attributeKey, &attributeFixed)
		if err != nil {
			return choice, err
		}

		if attributeID != nil {
			attribute = activity.Attribute{
				ID:      *attributeID,
				TableID: table.ID,
				Key:     attributeKey,
				Value:   *attributeValue,
				Fixed:   *attributeFixed,
			}

		}

		if _, ok := tablesMap[table.ID]; !ok {
			tablesMap[table.ID] = &table
			tablesMap[table.ID].Attributes = make(activity.Attributes, 0)
		}

		tablesMap[table.ID].Attributes = append(tablesMap[table.ID].Attributes, attribute)

	}

	tableIDs := make([]interface{}, 0)

	for _, v := range tablesMap {
		choice.Tables = append(choice.Tables, *v)
		tableIDs = append(tableIDs, v.ID)
	}

	relationships := make(activity.Relationships, 0)

	if len(tableIDs) > 0 {
		err = r.db.GetDB().
			Table(TableName.Relationship).
			Where("table1_id IN (" + utils.ToStrings(tableIDs) + ") OR table2_id IN (" + utils.ToStrings(tableIDs) + ")").
			Find(&relationships).
			Error

		if err != nil {
			return choice, err
		}
	}

	choice.Relationships = append(choice.Relationships, relationships...)

	return choice, nil
}

func (r learningRepository) GetPeerChoice(erAnswerID *int) (activity.ERAnswer, error) {

	answer := activity.ERAnswer{}

	query := r.db.GetDB().Debug().
		Select(IDName.Table, "title", IDName.Attribute, "value", "attribute_key", IDName.ERAnswer).
		Table(ViewName.RandomERAnswer)

	if erAnswerID != nil {
		query = query.Table(ViewName.AllERAnswer).Where(IDName.ERAnswer+" = ?", *erAnswerID)
	} else {
		query = query.Table(ViewName.RandomERAnswer)
	}

	rows, err := query.Rows()

	if err != nil {
		return answer, err
	}

	tablesMap := map[string]*activity.Table{}

	for rows.Next() {
		table := activity.Table{}

		var attributeID *int
		var attributeValue *string
		var attributeKey *string

		var attribute activity.Attribute

		err = rows.Scan(&table.ID, &table.Title, &attributeID, &attributeValue, &attributeKey, &answer.ID)
		if err != nil {
			return answer, err
		}

		if attributeID != nil {
			attribute = activity.Attribute{
				ID:      *attributeID,
				TableID: table.ID,
				Key:     attributeKey,
				Value:   *attributeValue,
				Fixed:   false,
			}

		}

		if _, ok := tablesMap[table.ID]; !ok {
			tablesMap[table.ID] = &table
			tablesMap[table.ID].Attributes = make(activity.Attributes, 0)
		}

		tablesMap[table.ID].Attributes = append(tablesMap[table.ID].Attributes, attribute)
	}

	tableIDs := make([]interface{}, 0)

	for _, v := range tablesMap {
		answer.Tables = append(answer.Tables, *v)
		tableIDs = append(tableIDs, v.ID)
	}

	if len(tableIDs) == 0 {
		return answer, errs.ErrNotFoundError
	}

	relationships := make(activity.Relationships, 0)

	err = r.db.GetDB().
		Table(TableName.Relationship).
		Where("table1_id IN (" + utils.ToStrings(tableIDs) + ") OR table2_id IN (" + utils.ToStrings(tableIDs) + ")").
		Find(&relationships).
		Error

	if err != nil {
		return answer, err
	}

	answer.Relationships = append(answer.Relationships, relationships...)

	return answer, err

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
	case 6:
		return r.GetERChoice(activityID)
	case 7:
		return r.GetPeerChoice(nil)
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

func (r learningRepository) InsertERAnswer(answer activity.ERAnswer, userID int) error {
	var wg sync.WaitGroup
	var err error

	a := activity.ERAnswer{}

	tx := r.db.GetDB().Debug().Begin()

	wg.Add(3)

	go func() {
		defer wg.Done()
		e := r.db.GetDB().Table(TableName.ERAnswer).Where(IDName.User+" = ?", userID).Find(&a).Error
		if e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()
		e := tx.Table(TableName.Tables).Create(&answer.Tables).Error
		if e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()
		e := tx.Table(TableName.ERAnswer).Create(&answer).Error
		if e != nil {
			err = e
		}
	}()

	wg.Wait()

	if a.UserID == userID {
		tx.Rollback()
		return nil
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	wg.Add(3)

	go func() {
		defer wg.Done()
		e := tx.Table(TableName.Relationship).Create(&answer.Relationships).Error
		if e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()

		attributes := make([]activity.Attribute, 0)

		for i, v := range answer.Tables {

			for j := range v.Attributes {
				answer.Tables[i].Attributes[j].TableID = v.ID
			}

			attributes = append(attributes, v.Attributes...)
		}

		if len(attributes) == 0 {
			return
		}

		e := tx.Table(TableName.Attributes).Create(&attributes).Error
		if e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()

		erAnswerTables := make([]activity.ERAnswerTables, 0)
		for _, v := range answer.Tables {
			erAnswerTables = append(erAnswerTables, activity.ERAnswerTables{
				ERAnswerID: answer.ID,
				TableID:    v.ID,
			})

		}

		e := tx.Table(TableName.ERAnswerTables).Create(&erAnswerTables).Error
		if e != nil {
			err = e
		}
	}()

	wg.Wait()

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
