package jwt

import (
	"database-camp/internal/errs"
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/logs"
	"database-camp/internal/repositories"
	"database-camp/internal/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type Jwt interface {
	Sign(id int) (string, error)
	Verify(application.Context)
}
type jwtMiddleware struct {
	repo repositories.UserRepository
}

func New(repo repositories.UserRepository) jwtMiddleware {
	return jwtMiddleware{repo: repo}
}

func (j jwtMiddleware) Sign(id int) (string, error) {
	atClaims := jwt.MapClaims{
		"id":     id,
		"secret": uuid.NewV4(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	err = j.updateToken(id, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j jwtMiddleware) Verify(c application.Context) {
	bearer, err := j.jwtFromHeader(c)
	if err != nil {
		logs.GetInstance().Error(err)
		c.Error(err)
		return
	}

	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			EnMessage := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			ThMessage := fmt.Sprintf("วิธีการลงนามที่ไม่คาดคิด: %v", token.Header["alg"])
			return nil, c.Error(errs.NewForbiddenError(ThMessage, EnMessage))
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		logs.GetInstance().Error(err)
		c.Error(errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid"))
		return
	}

	claims, err := j.getClaims(token)
	if err != nil {
		logs.GetInstance().Error(err)
		c.Error(err)
		return
	}

	id := utils.ParseInt(claims["id"])

	if !j.validUser(bearer, id) {
		c.Error(errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid"))
		return
	}

	err = j.updateToken(id, bearer)
	if err != nil {
		logs.GetInstance().Error(err)
		c.Error(errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error"))
		return
	}

	j.setClaims(c, claims)
	c.Next()
}

func (j jwtMiddleware) updateToken(id int, token string) error {
	tokenExpireHour := time.Hour * utils.ParseDuration(os.Getenv("TOKEN_EXPIRE_HOUR"))
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

func (j jwtMiddleware) setClaims(c application.Context, claims jwt.MapClaims) {
	for k, v := range claims {
		if k != "secret" {
			c.Locals(k, utils.ParseString(v))
		}
	}
}

func (j jwtMiddleware) validUser(token string, id int) bool {
	user, err := j.repo.GetUserByID(id)
	if err != nil || user == nil {
		return false
	}

	if user.AccessToken != token || user.ExpiredTokenTimestamp.Before(time.Now().Local()) {
		return false
	}

	return true
}

func (j jwtMiddleware) jwtFromHeader(c application.Context) (string, error) {
	auth := c.GetHeader("Authorization")
	l := len("Bearer")
	if len(auth) > l+1 && strings.EqualFold(auth[:l], "Bearer") {
		return auth[l+1:], nil
	}
	return "", errs.NewBadRequestError("ไม่พบ JWT Token ในส่วนหัวของคำร้องขอ", "JWT Token Not found")
}
