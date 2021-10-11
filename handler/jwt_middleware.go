package handler

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type jwtMiddleware struct {
	repo repository.IUserRepository
}

type IJwt interface {
	JwtSign(id int) (string, error)
}

func NewJwtMiddleware(repo repository.IUserRepository) jwtMiddleware {
	return jwtMiddleware{repo: repo}
}

func (j jwtMiddleware) updateToken(id int, token string) error {
	tokenExpireHour := time.Hour * utils.NewType().ParseDuration(os.Getenv("TOKEN_EXPIRE_HOUR"))
	expiredTokenTimestamp := time.Now().Local().Add(tokenExpireHour)
	err := j.repo.UpdatesByID(id, map[string]interface{}{
		"access_token":            token,
		"expired_token_timestamp": expiredTokenTimestamp,
	})
	return err
}

func (j jwtMiddleware) JwtSign(id int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["id"] = id
	atClaims["secret"] = uuid.NewV4()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logs.New().Error(err)
		return "", errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}

	err = j.updateToken(id, token)
	if err != nil {
		logs.New().Error(err)
		return "", errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}

	return token, nil
}

func (j jwtMiddleware) getBearer(header []string) (string, error) {
	if len(header) == 0 {
		return "", errs.NewForbiddenError("ส่วนหัวของคำร้องขอไม่ถูกต้อง", "Bad header")
	}

	splitHeader := strings.Split(header[0], " ")
	if len(splitHeader) != 2 {
		return "", errs.NewForbiddenError("ส่วนหัวของคำร้องขอไม่ถูกต้อง", "Bad header")
	}

	bearer := splitHeader[1]
	return bearer, nil
}

func (j jwtMiddleware) getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claims, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid")
	} else {
		return claims, nil
	}
}

func (j jwtMiddleware) setClaims(c echo.Context, claims jwt.MapClaims) {
	for k, v := range claims {
		c.Set(k, v)
	}
}

func (j jwtMiddleware) validUser(token string, id int) bool {
	user, err := j.repo.GetUserByID(id)
	if err != nil || user == (models.User{}) {
		return false
	}

	if user.AccessToken != token || user.ExpiredTokenTimestamp.Before(time.Now().Local()) {
		return false
	}

	return true
}

func (j jwtMiddleware) JwtVerify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearer, err := j.getBearer(c.Request().Header["Authorization"])
		if err != nil {
			return handleError(c, err)
		}

		token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				EnMessage := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
				ThMessage := fmt.Sprintf("วิธีการลงนามที่ไม่คาดคิด: %v", token.Header["alg"])
				return nil, handleError(c, errs.NewForbiddenError(ThMessage, EnMessage))
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			logs.New().Error(err)
			return handleError(c, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error"))
		}

		claims, err := j.getClaims(token)
		if err != nil {
			return handleError(c, err)
		}

		id := utils.NewType().ParseInt(claims["id"])

		if !j.validUser(bearer, id) {
			return handleError(c, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid"))
		}

		err = j.updateToken(id, bearer)
		if err != nil {
			logs.New().Error(err)
			return handleError(c, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error"))
		}

		j.setClaims(c, claims)
		return next(c)
	}
}
