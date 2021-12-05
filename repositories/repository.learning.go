package repositories

// repository.learning.go
/**
 * 	This file is a part of repositories, used to do data manipulation of learning
 */

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/services"
)

/**
 * 	This class manipulation learning data to other application
 */
type learningRepository struct {
	Database database.IDatabase   // Database to do database manipulation
	Service  services.IAwsService // Service to do database manipulation
}

/**
 * Constructor creates a new learningRepository instance
 *
 * @param   db    		Database to data manipulation
 * @param   service    	service to data manipulation
 *
 * @return 	instance of learningRepository
 */
func NewLearningRepository(db database.IDatabase, service services.IAwsService) learningRepository {
	return learningRepository{Database: db, Service: service}
}

/**
 * Get learning content from the database
 *
 * @param 	contentID  Content ID for getting content data
 *
 * @return learning content
 * @return the error of getting content
 */
func (r learningRepository) GetContent(id int) (*storages.ContentDB, error) {
	content := storages.ContentDB{}
	err := r.Database.GetDB().
		Table(storages.TableName.Content).
		Where(storages.IDName.Content+" = ?", id).
		Find(&content).
		Error
	return &content, err
}

/**
 * Get content overview from the database
 *
 * @return all content in the application
 * @return the error of getting content
 */
func (r learningRepository) GetOverview() ([]storages.OverviewDB, error) {
	overview := make([]storages.OverviewDB, 0)
	err := r.Database.GetDB().
		Table(storages.TableName.ContentGroup).
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

/**
 * Get content of the exam from the database
 *
 * @param 	examType  Type of the exam for getting content of the exam
 *
 * @return all content of the exam in the application
 * @return the error of getting content
 */
func (r learningRepository) GetContentExam(examType string) ([]storages.ContentExamDB, error) {
	contentExam := make([]storages.ContentExamDB, 0)
	db := r.Database.GetDB()
	examSubquery := db.Table(storages.TableName.Exam).
		Select("exam_id").
		Where("type = ?", string(examType)).
		Order("created_timestamp desc").
		Limit(1)
	err := r.Database.GetDB().
		Table(storages.TableName.ContentExam).
		Where("exam_id = (?)", examSubquery).
		Find(&contentExam).
		Error
	return contentExam, err
}

/**
 * Get activities of the content from the database
 *
 * @param 	contentID  Content ID for getting activities of the content
 *
 * @return activities of the content
 * @return the error of getting activities
 */
func (r learningRepository) GetContentActivity(contentID int) ([]storages.ActivityDB, error) {
	activity := make([]storages.ActivityDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.Activity).
		Where(storages.IDName.Content+" = ?", contentID).
		Find(&activity).
		Error

	return activity, err
}

/**
 * Get activity from the database
 *
 * @param 	id  Activity ID for getting activity data
 *
 * @return activity data
 * @return the error of getting activity
 */
func (r learningRepository) GetActivity(id int) (*storages.ActivityDB, error) {
	activity := storages.ActivityDB{}

	err := r.Database.GetDB().
		Table(storages.TableName.Activity).
		Where(storages.IDName.Activity+" = ?", id).
		Find(&activity).
		Error

	return &activity, err
}

/**
 * Get matching choices of the activity from the database
 *
 * @param 	activityID  Activity ID for getting matching choices of the activity
 *
 * @return matching choice of the activity
 * @return the error of getting choices
 */
func (r learningRepository) GetMatchingChoice(activityID int) ([]storages.MatchingChoiceDB, error) {
	matchingChoice := make([]storages.MatchingChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.MatchingChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&matchingChoice).
		Error

	return matchingChoice, err
}

/**
 * Get multiple choices of the activity from the database
 *
 * @param 	activityID  Activity ID for getting multiple choices of the activity
 *
 * @return multiple choice of the activity
 * @return the error of getting choices
 */
func (r learningRepository) GetMultipleChoice(activityID int) ([]storages.MultipleChoiceDB, error) {
	multipleChoice := make([]storages.MultipleChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.MultipleChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&multipleChoice).
		Error

	return multipleChoice, err
}

/**
 * Get completion choices of the activity from the database
 *
 * @param 	activityID  Activity ID for getting completion choices of the activity
 *
 * @return completion choice of the activity
 * @return the error of getting choices
 */
func (r learningRepository) GetCompletionChoice(activityID int) ([]storages.CompletionChoiceDB, error) {
	completionChoice := make([]storages.CompletionChoiceDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.CompletionChoice).
		Where(storages.IDName.Activity+" = ?", activityID).
		Find(&completionChoice).
		Error

	return completionChoice, err
}

/**
 * Get hints of the activity from the database
 *
 * @param 	activityID  Activity ID for getting hints of the activity
 *
 * @return hints of the activity
 * @return the error of getting hints
 */
func (r learningRepository) GetActivityHints(activityID int) ([]storages.HintDB, error) {
	hints := make([]storages.HintDB, 0)

	err := r.Database.GetDB().
		Table(storages.TableName.Hint).
		Where(storages.IDName.Activity+" = ?", activityID).
		Order("level ASC").
		Find(&hints).
		Error

	return hints, err
}

/**
 * Get video file link from the AWS service
 *
 * @param 	imagekey  Image key of the file in amazon s3
 *
 * @return file link
 * @return the error of getting file link
 */
func (r learningRepository) GetVideoFileLink(imagekey string) (string, error) {
	return r.Service.GetFileLink(imagekey)
}
