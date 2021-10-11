package handler

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"net/http"

	"github.com/labstack/echo/v4"
)

type message struct {
	Th string `json:"th_message"`
	En string `json:"en_message"`
}

func handleError(c echo.Context, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		return c.JSON(e.Code, message{Th: e.ThMessage, En: e.EnMessage})
	case error:
		return c.JSON(http.StatusInternalServerError, message{Th: "เกิดข้อผิดพลาด", En: "Internal Server Error"})
	}
	return nil
}

func bindRequest(c echo.Context, request interface{}) error {
	err := c.Bind(&request)
	if err != nil {
		logs.New().Error(err)
		return errs.NewBadRequestError("คำร้องขอไม่ถูกต้อง", "Bad Request")
	}
	return nil
}
