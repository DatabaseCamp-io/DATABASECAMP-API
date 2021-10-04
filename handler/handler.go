package handler

import (
	"DatabaseCamp/errs"
	"net/http"

	"github.com/labstack/echo/v4"
)

type message struct {
	Th string
	En string
}

func HandleError(c echo.Context, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		return c.JSON(e.Code, message{Th: e.ThMessage, En: e.EnMessage})
	case error:
		return c.JSON(http.StatusInternalServerError, message{Th: "เกิดข้อผิดพลาด", En: "Internal Server Error"})
	}
	return nil
}
