package domain

import "regexp"

var mailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type UserError struct {
	Code int
}

const (
	UserErrorEmailEmptyCode = iota
	UserErrorInvalidEmail
	UserErrorInvalidPasswordLengthCode
	UserErrorNameEmptyCode
	UserErrorPasswordEmptyCode
	UserErrorWrongPasswordCode
	UserErrorEmailExistsCode
	UserErrorPasswordConfirmationNotMatchCode
	UserChangePasswordErrorWrongOldPasswordCode
)

func (e UserError) Error() string {
	switch e.Code {
	case UserErrorEmailEmptyCode:
		return "Email cannot be empty"
	case UserErrorInvalidEmail:
		return "Email is invalid"
	case UserErrorNameEmptyCode:
		return "Name cannot be empty"
	case UserErrorPasswordEmptyCode:
		return "Password cannot be empty"
	case UserErrorWrongPasswordCode:
		return "Wrong password"
	case UserErrorInvalidPasswordLengthCode:
		return "Password must be at least 6 characters long"
	case UserErrorEmailExistsCode:
		return "Email already exists"
	case UserErrorPasswordConfirmationNotMatchCode:
		return "Password confirmation didn't match"
	case UserChangePasswordErrorWrongOldPasswordCode:
		return "Invalid old password"
	default:
		return "Unrecognized user error code"
	}
}
