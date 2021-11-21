package utils

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"

	"github.com/gofiber/fiber/v2"
)

type message struct {
	Th string `json:"th_message"`
	En string `json:"en_message"`
}

type handle struct{}

func NewHandle() handle {
	return handle{}
}

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

func (h *handle) BindRequest(c *fiber.Ctx, request interface{}) error {
	err := c.BodyParser(&request)
	if err != nil {
		logs.New().Error(err)
		return errs.ErrBadRequestError
	}
	return nil
}
