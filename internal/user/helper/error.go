package helper

import (
	"fmt"
	"github.com/SemmiDev/blog/internal/common/logger"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

const (
	ErrEmailEmptyCode = 1 << iota
	ErrInvalidEmailCode
	ErrInvalidPasswordLengthCode
	ErrNameEmptyCode
	ErrPasswordEmptyCode
	ErrWrongPasswordCode
	ErrEmailExistsCode
	ErrPasswordConfirmationNotMatchCode
	ErrWrongOldPasswordCode
	ErrInvalidCode
	ErrParseCode

	ErrAuthorizationHeaderKeyCode
	ErrAuthorizationHeaderFormatCode
	ErrAuthorizationTypeBearerCode
	ErrAuthorizationPayloadKeyCode
	ErrAuthorizationInvalidTokenCode
)

type Err struct {
	FieldName    string `json:"field_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func NewErr(errorCode int, fieldName string) Err {
	return Err{
		FieldName:    fieldName,
		ErrorCode:    strconv.Itoa(errorCode),
		ErrorMessage: Message(errorCode),
	}
}

func (r Err) Error() string {
	return fmt.Sprintf(
		"Field Name: %s, Error Code: %s, Error Message: %s",
		r.FieldName,
		r.ErrorCode,
		r.ErrorMessage,
	)
}

func Message(errorCode int) string {
	switch errorCode {
	case ErrEmailEmptyCode:
		return "Email is empty"
	case ErrInvalidEmailCode:
		return "Email is invalid"
	case ErrInvalidPasswordLengthCode:
		return "Password length must be greater than 6"
	case ErrNameEmptyCode:
		return "Name is empty"
	case ErrPasswordEmptyCode:
		return "Password is empty"
	case ErrWrongPasswordCode:
		return "Wrong password"
	case ErrEmailExistsCode:
		return "Email exists"
	case ErrPasswordConfirmationNotMatchCode:
		return "Password confirmation not match"
	case ErrWrongOldPasswordCode:
		return "Wrong old password"
	case ErrAuthorizationHeaderKeyCode:
		return "Authorization header key is invalid"
	case ErrAuthorizationHeaderFormatCode:
		return "Authorization header format is invalid"
	case ErrAuthorizationTypeBearerCode:
		return "Authorization type is invalid"
	case ErrAuthorizationPayloadKeyCode:
		return "Authorization payload key is invalid"
	case ErrAuthorizationInvalidTokenCode:
		return "Authorization token is invalid"
	case ErrInvalidCode:
		return "code is invalid"
	case ErrParseCode:
		return "parse error"
	default:
		return "Internal server error"
	}
}

func Error(c *fiber.Ctx, err error) error {
	errorResponse := map[string]string{
		"field_name":    "",
		"error_code":    "",
		"error_message": "",
	}

	file, line := getFileAndLineNumber()

	theError, ok := err.(Err)
	if ok {
		errorResponse["error_code"] = theError.ErrorCode
		errorResponse["error_message"] = theError.ErrorMessage
		errorResponse["field_name"] = theError.FieldName
		fields := fiber.Map{
			"IP":    c.IP(),
			"FILE":  file,
			"LINE":  line,
			"ERROR": errorResponse,
		}
		logger.Log.Error().Interface("err", fields).Send()
		return c.Status(http.StatusBadRequest).JSON(errorResponse)
	}

	logger.Log.Error().Interface("err", err.Error()).Send()
	errorResponse["error_message"] = err.Error()
	return c.Status(http.StatusInternalServerError).JSON(errorResponse)
}

func getFileAndLineNumber() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}
