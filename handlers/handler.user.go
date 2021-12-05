package handlers

// handler.user.go
/**
 * 	This file is a part of handler, used to handle request of the user
 */

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/middleware"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/**
 * This class handle request of the user
 */
type userHandler struct {
	Controller controllers.IUserController // User controller for doing business logic of the user
	Jwt        middleware.IJwt             // Jwt middleware for user verification
}

/**
 * Constructor creates a new userHandler instance
 *
 * @param   controller    	User controller for doing business logic of the user
 * @param   jwt    			Jwt middleware for user verification
 *
 * @return 	instance of userHandler
 */
func NewUserHandler(controller controllers.IUserController, jwt middleware.IJwt) userHandler {
	return userHandler{Controller: controller, Jwt: jwt}
}

/**
 * Register
 *
 * @param 	c  Context of the web framework
 *
 * @return the error of getting exam
 */
func (h userHandler) Register(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	request := request.UserRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.ValidateRegister()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.Register(request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	token, err := h.Jwt.JwtSign(response.ID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response.AccessToken = token

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Login
 *
 * @param 	c  Context of the web framework
 *
 * @return the error of getting exam
 */
func (h userHandler) Login(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	request := request.UserRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.ValidateLogin()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.Login(request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	token, err := h.Jwt.JwtSign(response.ID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response.AccessToken = token

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Get user profile
 *
 * @param 	c  Context of the web framework
 *
 * @return the error of getting exam
 */
func (h userHandler) GetProfile(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	id := c.Params("id")

	response, err := h.Controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Get own profile
 *
 * @param 	c  Context of the web framework
 *
 * @return the error of getting exam
 */
func (h userHandler) GetOwnProfile(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	id := c.Locals("id")

	response, err := h.Controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Get ranking
 *
 * @param 	c  Context of the web framework
 *
 * @return the error of getting exam
 */
func (h userHandler) GetUserRanking(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	id := c.Locals("id")

	response, err := h.Controller.GetRanking(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Edit user profile
 *
 * @param 	c  Context of the web framework
 *
 * @return the error of getting exam
 */
func (h userHandler) Edit(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	request := request.UserRequest{}
	userID := utils.NewType().ParseInt(c.Locals("id"))

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.ValidateEdit()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.EditProfile(userID, request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}
