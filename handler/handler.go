package handler

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type message struct {
	Th string `json:"th_message"`
	En string `json:"en_message"`
}

func handleError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		return c.Status(e.Code).JSON(message{Th: e.ThMessage, En: e.EnMessage})
	case error:
		return c.Status(http.StatusInternalServerError).JSON(message{Th: "เกิดข้อผิดพลาด", En: "Internal Server Error"})
	}
	return nil
}

func bindRequest(c *fiber.Ctx, request interface{}) error {
	err := c.BodyParser(&request)
	if err != nil {
		logs.New().Error(err)
		return errs.NewBadRequestError("คำร้องขอไม่ถูกต้อง", "Bad Request")
	}
	return nil
}
