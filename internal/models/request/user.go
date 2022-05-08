package request

// request.user.go
/**
 * 	This file is a part of models, used to collect request of user
 */

import (
	"database-camp/internal/errs"
	"database-camp/internal/utils"
)

/**
 * 	This class represent user request
 */
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
 * Validate register request
 *
 * @return the error of validating request
 */
func (r UserRequest) ValidateRegister() error {
	var err error
	if r.Name == "" {
		err = errs.NewBadRequestError("ไม่พบชื่อในคำร้องขอ", "Name Not Found")
	} else {
		err = r.ValidateLogin()
	}
	return err
}

/**
 * Validate login request
 *
 * @return the error of validating request
 */
func (r UserRequest) ValidateLogin() error {
	var err error
	if r.Email == "" {
		err = errs.NewBadRequestError("ไม่พบอีเมลในคำร้องขอ", "Email Not Found")
	} else if !utils.IsEmailValid(r.Email) {
		err = errs.NewBadRequestError("รูปแบบ email ไม่ถูกต้อง", "Email Invalid")
	} else if r.Password == "" {
		err = errs.NewBadRequestError("ไม่พบรหัสผ่านในคำร้องขอ", "Password Not Found")
	} else if len(r.Password) < 8 {
		err = errs.NewBadRequestError("ความยาวของรหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร", "Password length must be at least 8 characters")
	}
	return err
}

/**
 * Validate edit user profile request
 *
 * @return the error of validating request
 */
func (r UserRequest) ValidateEdit() error {
	var err error
	if r.Name == "" {
		err = errs.NewBadRequestError("ไม่พบชื่อในคำร้องขอ", "Name Not Found")
	}
	return err
}
