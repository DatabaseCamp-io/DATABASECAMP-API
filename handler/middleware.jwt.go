package handler

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/repository"
	"DatabaseCamp/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type jwtMiddleware struct {
	repo repository.IUserRepository
}

type IJwt interface {
	JwtSign(id int) (string, error)
	JwtVerify(c *fiber.Ctx) error
}

func NewJwtMiddleware(repo repository.IUserRepository) jwtMiddleware {
	return jwtMiddleware{repo: repo}
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

func (j jwtMiddleware) JwtVerify(c *fiber.Ctx) error {
	bearer, err := j.jwtFromHeader(c)
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
		return handleError(c, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid"))
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
	return c.Next()
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

func (j jwtMiddleware) getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claims, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid")
	} else {
		return claims, nil
	}
}

func (j jwtMiddleware) setClaims(c *fiber.Ctx, claims jwt.MapClaims) {
	for k, v := range claims {
		if k != "secret" {
			c.Locals(k, utils.NewType().ParseString(v))
		}
	}
}

func (j jwtMiddleware) validUser(token string, id int) bool {
	userDB, err := j.repo.GetUserByID(id)
	if err != nil || userDB == nil {
		return false
	}

	if userDB.AccessToken != token || userDB.ExpiredTokenTimestamp.Before(time.Now().Local()) {
		return false
	}

	return true
}

func (j jwtMiddleware) jwtFromHeader(c *fiber.Ctx) (string, error) {
	auth := c.Get("Authorization")
	l := len("Bearer")
	if len(auth) > l+1 && strings.EqualFold(auth[:l], "Bearer") {
		return auth[l+1:], nil
	}
	return "", errs.NewBadRequestError("ไม่พบ JWT Token ในส่วนหัวของคำร้องขอ", "JWT Token Not found")
}