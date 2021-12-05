package middleware

// middleware.jwt.go
/**
 * 	This file is a part of middleware, used to verificate user
 */

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/repositories"
	"DatabaseCamp/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

/**
 * This class verificate user
 */
type jwtMiddleware struct {
	Repo repositories.IUserRepository // User repository for getting user data
}

/**
 * Constructor creates a new jwtMiddleware instance
 *
 * @param 	repo User repository for getting user data
 *
 * @return 	instance of jwtMiddleware
 */
func NewJwtMiddleware(repo repositories.IUserRepository) jwtMiddleware {
	return jwtMiddleware{Repo: repo}
}

/**
 * Sign user for verification
 *
 * @param 	id  User ID to sign
 *
 * @return signing token
 * @return the error of signing
 */
func (j jwtMiddleware) JwtSign(id int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["id"] = id
	atClaims["secret"] = uuid.NewV4()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logs.New().Error(err)
		return "", errs.ErrInternalServerError
	}

	err = j.updateToken(id, token)
	if err != nil {
		logs.New().Error(err)
		return "", errs.ErrInternalServerError
	}

	return token, nil
}

/**
 * Verify request by token
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
func (j jwtMiddleware) JwtVerify(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	bearer, err := j.jwtFromHeader(c)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			EnMessage := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			ThMessage := fmt.Sprintf("วิธีการลงนามที่ไม่คาดคิด: %v", token.Header["alg"])
			return nil, handleUtil.HandleError(c, errs.NewForbiddenError(ThMessage, EnMessage))
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		logs.New().Error(err)
		return handleUtil.HandleError(c, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid"))
	}

	claims, err := j.getClaims(token)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	id := utils.NewType().ParseInt(claims["id"])

	if !j.validUser(bearer, id) {
		return handleUtil.HandleError(c, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid"))
	}

	err = j.updateToken(id, bearer)
	if err != nil {
		logs.New().Error(err)
		return handleUtil.HandleError(c, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error"))
	}

	j.setClaims(c, claims)
	return c.Next()
}

/**
 * Update user token into the database
 *
 * @param 	id  	User ID to update token
 * @param 	token  	token to be updated
 *
 * @return the error of updating token
 */
func (j jwtMiddleware) updateToken(id int, token string) error {
	tokenExpireHour := time.Hour * utils.NewType().ParseDuration(os.Getenv("TOKEN_EXPIRE_HOUR"))
	expiredTokenTimestamp := time.Now().Local().Add(tokenExpireHour)
	err := j.Repo.UpdatesByID(id, map[string]interface{}{
		"access_token":            token,
		"expired_token_timestamp": expiredTokenTimestamp,
	})
	return err
}

/**
 * Get claims of the token
 *
 * @param 	token  	JWT token
 *
 * @return claims of the token
 * @return the error of getting claims
 */
func (j jwtMiddleware) getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claims, errs.NewForbiddenError("โทเค็นไม่ถูกต้อง", "Token Invalid")
	} else {
		return claims, nil
	}
}

/**
 * Set claims to the header of the request
 *
 * @param	c  			context of the web framework
 * @param 	claims  	claims of the token to set
 */
func (j jwtMiddleware) setClaims(c *fiber.Ctx, claims jwt.MapClaims) {
	for k, v := range claims {
		if k != "secret" {
			c.Locals(k, utils.NewType().ParseString(v))
		}
	}
}

/**
 * Check user token
 *
 * @param	token  		Token of the request
 * @param 	id  		User ID to check token
 *
 * @return 	true if user token is valid, false otherwise
 */
func (j jwtMiddleware) validUser(token string, id int) bool {
	userDB, err := j.Repo.GetUserByID(id)
	if err != nil || userDB == nil {
		return false
	}

	if userDB.AccessToken != token || userDB.ExpiredTokenTimestamp.Before(time.Now().Local()) {
		return false
	}

	return true
}

/**
 * Get token from the header of the request
 *
 * @param 	c  context of the web framework
 *
 * @return 	token from the header
 * @return 	the error of the getting token
 */
func (j jwtMiddleware) jwtFromHeader(c *fiber.Ctx) (string, error) {
	auth := c.Get("Authorization")
	l := len("Bearer")
	if len(auth) > l+1 && strings.EqualFold(auth[:l], "Bearer") {
		return auth[l+1:], nil
	}
	return "", errs.NewBadRequestError("ไม่พบ JWT Token ในส่วนหัวของคำร้องขอ", "JWT Token Not found")
}
