package utils

// util.handle.go
/**
 * 	This file is a part of utilities, used to help handle module
 */

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"

	"github.com/gofiber/fiber/v2"
)

/**
 * 	This class is data model for message in Thai and English
 */
type message struct {
	Th string `json:"th_message"`
	En string `json:"en_message"`
}

/**
 * 	This class help handle module
 */
type handle struct{}

/**
 * Constructor creates a new handle instance
 *
 * @return 	instance of handle
 */
func NewHandle() handle {
	return handle{}
}

/**
 * Handle error
 *
 * @param	c  		Context of the web framework
 * @param	err  	Error in type App error
 *
 * @return 	error response
 */
func (h *handle) HandleError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		return c.Status(e.Code).JSON(message{Th: e.ThMessage, En: e.EnMessage})
	case error:
		return c.Status(fiber.StatusInternalServerError).JSON(message{
			Th: errs.INTERNAL_SERVER_ERROR_TH,
			En: errs.INTERNAL_SERVER_ERROR_EN,
		})
	}
	return nil
}

/**
 * Bind user request in format json to struct model
 *
 * @param	c  			Context of the web framework
 * @param	request  	Struct model
 *
 * @return 	the error of response
 */
func (h *handle) BindRequest(c *fiber.Ctx, request interface{}) error {
	err := c.BodyParser(&request)
	if err != nil {
		logs.New().Error(err)
		return errs.ErrBadRequestError
	}
	return nil
}
