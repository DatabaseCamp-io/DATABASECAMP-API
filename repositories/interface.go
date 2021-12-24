package repositories

// interface.go
/**
 * 	This file used to be a interface of repositories
 */

import (
	"DatabaseCamp/database"
	"DatabaseCamp/models/entities"
	"DatabaseCamp/models/storages"
)

/**
 * 	 Interface to show function in exam respository that others can use
 */
type IExamRepository interface {

	/**
	 * Get activities of the exam from the database
	 *
	 * @param 	examID  Exam ID for getting exam data
	 *
	 * @return activities of the exam
	 * @return the error of getting exam
	 */
	GetExamActivity(examID int) ([]storages.ExamActivityDB, error)

	/**
	 * Get exam overview from the database
	 *
	 * @return all exam in the application
	 * @return the error of getting exam
	 */
	GetExamOverview() ([]storages.ExamDB, error)

	/**
	 * Insert exam result into the database by database transaction
	 *
	 * @param 	tx  		Transaction model to do database transaction
	 * @param 	examResult  Exam result for insert into the database
	 *
	 * @return inserted exam result
	 * @return the error of inserting exam result
	 */
	InsertExamResultTransaction(tx database.ITransaction, examResult storages.ExamResultDB) (storages.ExamResultDB, error)

	/**
	 * Insert activities of the exam result into the database by database transaction
	 *
	 * @param 	tx  				Transaction model to do database transaction
	 * @param 	examResultActivity  Activities of the exam result for insert into the database
	 *
	 * @return inserted activities of the exam result
	 * @return the error of inserting activities of the exam result
	 */
	InsertExamResultActivityTransaction(tx database.ITransaction, examResultActivity []storages.ExamResultActivityDB) ([]storages.ExamResultActivityDB, error)
}

/**
 * 	 Interface to show function in learning respository that others can use
 */
type ILearningRepository interface {

	/**
	 * Get learning content from the database
	 *
	 * @param 	contentID  Content ID for getting content data
	 *
	 * @return learning content
	 * @return the error of getting content
	 */
	GetContent(id int) (*storages.ContentDB, error)

	/**
	 * Get content overview from the database
	 *
	 * @return all content in the application
	 * @return the error of getting content
	 */
	GetOverview() ([]storages.OverviewDB, error)

	/**
	 * Get content of the exam from the database
	 *
	 * @param 	examType  Type of the exam for getting content of the exam
	 *
	 * @return all content of the exam in the application
	 * @return the error of getting content
	 */
	GetContentExam(examType string) ([]storages.ContentExamDB, error)

	/**
	 * Get activity from the database
	 *
	 * @param 	id  Activity ID for getting activity data
	 *
	 * @return activity data
	 * @return the error of getting activity
	 */
	GetActivity(id int) (*storages.ActivityDB, error)

	/**
	 * Get matching choices of the activity from the database
	 *
	 * @param 	activityID  Activity ID for getting matching choices of the activity
	 *
	 * @return matching choice of the activity
	 * @return the error of getting choices
	 */
	GetMatchingChoice(activityID int) ([]storages.MatchingChoiceDB, error)

	/**
	 * Get multiple choices of the activity from the database
	 *
	 * @param 	activityID  Activity ID for getting multiple choices of the activity
	 *
	 * @return multiple choice of the activity
	 * @return the error of getting choices
	 */
	GetMultipleChoice(activityID int) ([]storages.MultipleChoiceDB, error)

	/**
	 * Get completion choices of the activity from the database
	 *
	 * @param 	activityID  Activity ID for getting completion choices of the activity
	 *
	 * @return completion choice of the activity
	 * @return the error of getting choices
	 */
	GetCompletionChoice(activityID int) ([]storages.CompletionChoiceDB, error)

	/**
	 * Get hints of the activity from the database
	 *
	 * @param 	activityID  Activity ID for getting hints of the activity
	 *
	 * @return hints of the activity
	 * @return the error of getting hints
	 */
	GetActivityHints(activityID int) ([]storages.HintDB, error)

	/**
	 * Get activities of the content from the database
	 *
	 * @param 	contentID  Content ID for getting activities of the content
	 *
	 * @return activities of the content
	 * @return the error of getting activities
	 */
	GetContentActivity(contentID int) ([]storages.ActivityDB, error)

	/**
	 * Get video file link from the AWS service
	 *
	 * @param 	imagekey  Image key of the file in amazon s3
	 *
	 * @return file link
	 * @return the error of getting file link
	 */
	GetVideoFileLink(imagekey string) (string, error)
}

//Interface that show how others function call and use function in module user respository
type IUserRepository interface {

	/**
	 * Get user by email from the database
	 *
	 * @param 	email  Email for getting user data
	 *
	 * @return user data
	 * @return the error of getting user data
	 */
	GetUserByEmail(email string) (*storages.UserDB, error)

	/**
	 * Get user by id from the database
	 *
	 * @param 	id  User ID for getting user data
	 *
	 * @return user data
	 * @return the error of getting user data
	 */
	GetUserByID(id int) (*storages.UserDB, error)

	/**
	 * Get user profile from the database
	 *
	 * @param 	id  User ID for getting user profile
	 *
	 * @return user profile
	 * @return the error of getting user profile
	 */
	GetProfile(id int) (*storages.ProfileDB, error)

	/**
	 * Get learning progression of the user from the database
	 *
	 * @param 	id  User ID for getting learning progression of the user
	 *
	 * @return learning progression of the user
	 * @return the error of getting learning progression of the user
	 */
	GetLearningProgression(id int) ([]storages.LearningProgressionDB, error)

	/**
	 * Get all badges of the application from the database
	 *
	 * @return all badges of the application
	 * @return the error of getting all badges of the application
	 */
	GetAllBadge() ([]storages.BadgeDB, error)

	/**
	 * Get user badge data by user badge id from the database
	 *
	 * @param 	id  User Badge ID for getting user badge data
	 *
	 * @return user badge data
	 * @return the error of getting user badge data
	 */
	GetUserBadge(id int) ([]storages.UserBadgeDB, error)

	/**
	 * Get collected badges of the user from the database
	 *
	 * @param 	userID  User ID for getting collected badges of the user
	 *
	 * @return collected badges of the user
	 * @return the error of getting collected badges of the user
	 */
	GetCollectedBadge(userID int) ([]storages.CorrectedBadgeDB, error)

	/**
	 * Get point ranking of the user from the database
	 *
	 * @param 	id  User ID for getting point ranking of the user
	 *
	 * @return point ranking of the user
	 * @return the error of getting point ranking of the user
	 */
	GetPointRanking(id int) (*storages.RankingDB, error)

	/**
	 * Get all user point ranking from the database
	 *
	 * @return all user point ranking
	 * @return the error of getting all user point ranking
	 */
	GetRankingLeaderBoard() ([]storages.RankingDB, error)

	/**
	 * Get User Hint data from the database
	 *
	 * @param 	userID  	User ID for getting User Hint data
	 * @param 	activityID  Activity ID for getting User Hint data
	 *
	 * @return User Hint data
	 * @return the error of getting User Hint data
	 */
	GetUserHint(userID int, activityID int) ([]storages.UserHintDB, error)

	/**
	 * Get exam results of the user from the database
	 *
	 * @param 	userID  	User ID for getting exam results of the user
	 *
	 * @return exam results of the user
	 * @return the error of getting exam results of the user
	 */
	GetExamResult(userID int) ([]storages.ExamResultDB, error)

	/**
	 * Get exam results of the user by exam result id from the database
	 *
	 * @param 	userID  	User ID for getting exam results of the user
	 *
	 * @return exam results of the user
	 * @return the error of getting exam results of the user
	 */
	GetExamResultByID(userID int, examResultID int) ([]storages.ExamResultDB, error)

	/**
	 * Insert user data into the database
	 *
	 * @param 	user  	User model for insert into the database
	 *
	 * @return inserted user
	 * @return the error of inserting user
	 */
	InsertUser(user storages.UserDB) (*storages.UserDB, error)

	/**
	 * Insert user hint into the database
	 *
	 * @param 	userHint  	User hint model for insert into the database
	 *
	 * @return inserted user hint
	 * @return the error of inserting user hint
	 */
	InsertUserHint(userHint storages.UserHintDB) (*storages.UserHintDB, error)

	/**
	 * Update user data into the database
	 *
	 * @param 	id  			User id to update into the database
	 * @param 	updateData  	User data to update into the database
	 *
	 * @return inserted user hint
	 * @return the error of updating user data
	 */
	UpdatesByID(id int, updateData map[string]interface{}) error

	/**
	 * Insert user hint into the database by database transaction
	 *
	 * @param 	tx  		Transaction model to do database transaction
	 * @param 	userHint  	User hint model for insert into the database
	 *
	 * @return inserted user hint
	 * @return the error of inserting user hint
	 */
	InsertUserHintTransaction(tx database.ITransaction, userHint storages.UserHintDB) (*storages.UserHintDB, error)

	/**
	 * Insert learning progression of the user into the database by database transaction
	 *
	 * @param 	tx  			Transaction model to do database transaction
	 * @param 	progression  	Learning progession model for insert into the database
	 *
	 * @return inserted learning progression
	 * @return the error of inserting learning progression
	 */
	InsertLearningProgressionTransaction(tx database.ITransaction, progression storages.LearningProgressionDB) (*storages.LearningProgressionDB, error)

	/**
	 * Insert user badge into the database by database transaction
	 *
	 * @param 	tx  			Transaction model to do database transaction
	 * @param 	userBadge  		User badge model for insert into the database
	 *
	 * @return inserted user badge
	 * @return the error of inserting user badge
	 */
	InsertUserBadgeTransaction(tx database.ITransaction, userBadge storages.UserBadgeDB) (*storages.UserBadgeDB, error)

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
	ChangePointTransaction(tx database.ITransaction, userID int, point int, mode entities.ChangePointMode) error
}
