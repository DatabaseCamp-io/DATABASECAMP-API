package handlers

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/middleware"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	Controller controllers.IUserController
	Jwt        middleware.IJwt
}

type IUserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	GetOwnProfile(c *fiber.Ctx) error
	GetUserRanking(c *fiber.Ctx) error
	Edit(c *fiber.Ctx) error
}

// Crate user handler instance
func NewUserHandler(controller controllers.IUserController, jwt middleware.IJwt) userHandler {
	return userHandler{Controller: controller, Jwt: jwt}
}

// Get token
// Validaate register
// Register user's id
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

// Log user into website
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

// Convert id to parseint type
func (h userHandler) GetProfile(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Params("id")
	response, err := h.Controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

// Get user's profile
func (h userHandler) GetOwnProfile(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Locals("id")
	response, err := h.Controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

// Get user ranking
func (h userHandler) GetUserRanking(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Locals("id")
	response, err := h.Controller.GetRanking(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

// Edit user's profile
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
