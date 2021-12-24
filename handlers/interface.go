package handlers

// interface.go
/**
 * 	This file used to be a interface of handlers
 */

import "github.com/gofiber/fiber/v2"

/**
 * 	 Interface to show function in exam handler that others can use
 */
type IExamHandler interface {

	/**
	 * Get the exam to use for the test
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetExam(c *fiber.Ctx) error

	/**
	 * Check answer of the exam
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	CheckExam(c *fiber.Ctx) error

	/**
	 * Get overview of the exam
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetExamOverview(c *fiber.Ctx) error

	/**
	 * Get exam result of the user
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetExamResult(c *fiber.Ctx) error
}

/**
 * 	 Interface to show function in learning handler that others can use
 */
type ILearningHandler interface {

	/**
	 * Get content roadmap
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetContentRoadmap(c *fiber.Ctx) error

	/**
	 * Get video lecture of the content
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetVideo(c *fiber.Ctx) error

	/**
	 * Get content overview
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetOverview(c *fiber.Ctx) error

	/**
	 * Get activity for user to do
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetActivity(c *fiber.Ctx) error

	/**
	 * Use hint of the activity
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	UseHint(c *fiber.Ctx) error

	/**
	 * Check matching choice answer of the activity
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	CheckMatchingAnswer(c *fiber.Ctx) error

	/**
	 * Check multiple choice answer of the activity
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	CheckMultipleAnswer(c *fiber.Ctx) error

	/**
	 * Check completion choice answer of the activity
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	CheckCompletionAnswer(c *fiber.Ctx) error
}

/**
 * 	 Interface to show function in user handler that others can use
 */
type IUserHandler interface {

	/**
	 * Register
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	Register(c *fiber.Ctx) error

	/**
	 * Login
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	Login(c *fiber.Ctx) error

	/**
	 * Get user profile
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetProfile(c *fiber.Ctx) error

	/**
	 * Get own profile
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetOwnProfile(c *fiber.Ctx) error

	/**
	 * Get ranking
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	GetUserRanking(c *fiber.Ctx) error

	/**
	 * Edit user profile
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	Edit(c *fiber.Ctx) error
}
