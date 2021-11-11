package models

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/utils"
)

// Model for receiving request
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Method for validate register request
func (r UserRequest) ValidateRegister() error {
	var err error
	if r.Name == "" {
		err = errs.NewBadRequestError("ไม่พบชื่อในคำร้องขอ", "Name Not Found")
	} else {
		err = r.ValidateLogin()
	}
	return err
}

// Method for validate login request
func (r UserRequest) ValidateLogin() error {
	var err error
	if r.Email == "" {
		err = errs.NewBadRequestError("ไม่พบอีเมลในคำร้องขอ", "Email Not Found")
	} else if !utils.NewHelper().IsEmailValid(r.Email) {
		err = errs.NewBadRequestError("รูปแบบ email ไม่ถูกต้อง", "Email Invalid")
	} else if r.Password == "" {
		err = errs.NewBadRequestError("ไม่พบรหัสผ่านในคำร้องขอ", "Password Not Found")
	} else if len(r.Password) < 8 {
		err = errs.NewBadRequestError("ความยาวของรหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร", "Password length must be at least 8 characters")
	}
	return err
}

// Method for validate edit request
func (r UserRequest) ValidateEdit() error {
	var err error
	if r.Name == "" {
		err = errs.NewBadRequestError("ไม่พบชื่อในคำร้องขอ", "Name Not Found")
	}
	return err
}
