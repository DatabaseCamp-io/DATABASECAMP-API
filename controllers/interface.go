package controllers

// interface.go
/**
 * 	This file used to be a interface of controllers
 */

import (
	"DatabaseCamp/models/entities/activity"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/response"
	"DatabaseCamp/models/storages"
)

/**
 * 	 Interface to show function in exam controller that others can use
 */
type IExamController interface {

	/**
	 * Get the exam to use for the test
	 *
	 * @param 	examID 		   Exam ID for getting activities of the exam
	 * @param 	userID 		   User ID for getting detail of user about the exam
	 *
	 * @return response of the exam
	 * @return error of getting the exam
	 */
	GetExam(examID int, userID int) (*response.ExamResponse, error)

	/**
	 * Get overview of the exam
	 *
	 * @param 	userID 		   User ID for getting detail of user about the exam overview
	 *
	 * @return response of the exam overview
	 * @return error of getting exam overview
	 */
	GetOverview(userID int) (*response.ExamOverviewResponse, error)

	/**
	 * Check answer of the exam
	 *
	 * @param 	userID 		   	User ID for record user exam
	 * @param 	request 		Exam answer request for check answer of the exam
	 *
	 * @return response of the exam result overview
	 * @return error of checking exam
	 */
	CheckExam(userID int, request request.ExamAnswerRequest) (*response.ExamResultOverviewResponse, error)

	/**
	 * Get exam result of the user
	 *
	 * @param 	userID 		   	User ID for getting user data
	 * @param 	examResultID 	Exam result ID for getting exam results of the user
	 *
	 * @return response of the exam result overview
	 * @return error of getting exam result
	 */
	GetExamResult(userID int, examResultID int) (*response.ExamResultOverviewResponse, error)
}

/**
 * 	 Interface to show function in learning controller that others can use
 */
type ILearningController interface {

	/**
	 * Get video lecture of the content
	 *
	 * @param 	id 		   	Content ID for getting video lecture of the content
	 *
	 * @return response of the video lecture
	 * @return error of getting video lecture
	 */
	GetVideoLecture(id int) (*response.VideoLectureResponse, error)

	/**
	 * Get content overview of thew user
	 *
	 * @param 	userID 		   	User ID for getting content overview of the user
	 *
	 * @return response of the content overview
	 * @return error of getting content overview
	 */
	GetOverview(userID int) (*response.ContentOverviewResponse, error)

	/**
	 * Get activity for user to do
	 *
	 * @param 	userID 		   		User ID for getting user hints
	 * @param 	activityID 			Activity ID for getting activity data
	 *
	 * @return response of the activity
	 * @return error of getting activity
	 */
	GetActivity(userID int, activityID int) (*activity.Response, error)

	/**
	 * Use hint of the activity
	 *
	 * @param 	userID 		   		User ID for getting user hints
	 * @param 	activityID 			Activity ID for getting all hints of the activity
	 *
	 * @return hint of the activity that user can use
	 * @return error of using hint
	 */
	UseHint(userID int, activityID int) (*storages.HintDB, error)

	/**
	 * Get roadmap of the content
	 *
	 * @param 	userID 		   		User ID for getting learning progression of the user
	 * @param 	contentID 			Content ID for getting roadmap of the content
	 *
	 * @return response of the content roadmap
	 * @return error of getting content roadmap
	 */
	GetContentRoadmap(userID int, contentID int) (*response.ContentRoadmapResponse, error)

	/**
	 * Check activity answer
	 *
	 * @param 	userID 		   		User ID for record user activity
	 * @param 	activityID 			Activity ID for getting activity solution
	 * @param 	typeID 				Activity ID for indicate type of the activity
	 * @param 	answer 				Answer of the user
	 *
	 * @return response of the answer
	 * @return error of checking activity answer
	 */
	CheckAnswer(userID int, activityID int, typeID int, answer interface{}) (*response.AnswerResponse, error)
}

/**
 * 	 Interface to show function in user controller that others can use
 */
type IUserController interface {

	/**
	 * Register for learning in platform
	 *
	 * @param 	request 		User request for create account
	 *
	 * @return response of user
	 * @return error of register
	 */
	Register(request request.UserRequest) (*response.UserResponse, error)

	/**
	 * Login
	 *
	 * @param 	request 		User request for login to the platform
	 *
	 * @return response of user
	 * @return error of login
	 */
	Login(request request.UserRequest) (*response.UserResponse, error)

	/**
	 * Get user profile
	 *
	 * @param 	userID 		   	User ID for getting user profile
	 *
	 * @return response of get profile
	 * @return error of getting profile
	 */
	GetProfile(userID int) (*response.GetProfileResponse, error)

	/**
	 * Edit profile data
	 *
	 * @param 	userID 		   	User ID for getting user profile
	 * @param 	request 		User request for editing the user profile
	 *
	 * @return response of edit profile
	 * @return error of edit profile
	 */
	EditProfile(userID int, request request.UserRequest) (*response.EditProfileResponse, error)

	/**
	 * Get others and own ranking
	 *
	 * @param 	id		User ID for getting user ranking
	 *
	 * @return response of ranking
	 * @return error of getting ranking
	 */
	GetRanking(id int) (*response.RankingResponse, error)
}
