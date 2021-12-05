package middleware

// interface.go
/**
 * 	This file used to be a interface of middleware
 */

import "github.com/gofiber/fiber/v2"

/**
 * 	 Interface to show function in JWT middleware that others can use
 */
type IJwt interface {

	/**
	 * Sign user for verification
	 *
	 * @param 	id  User ID to sign
	 *
	 * @return signing token
	 * @return the error of signing
	 */
	JwtSign(id int) (string, error)

	/**
	 * Verify request by token
	 *
	 * @param 	c  Context of the web framework
	 *
	 * @return the error of getting exam
	 */
	JwtVerify(c *fiber.Ctx) error
}
