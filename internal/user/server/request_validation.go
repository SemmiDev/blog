package server

import (
	"fmt"
	. "github.com/SemmiDev/blog/internal/common/logger"
	"github.com/SemmiDev/blog/internal/user/domain"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	mailRegex   = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	numberRegex = regexp.MustCompile(`^[0-9]+$`)
)

const (
	Required      = "Required"
	Alphanumeric  = "Alphanumeric"
	Alpha         = "Alpha"
	Numeric       = "Numeric"
	Float         = "Float"
	ParseFailed   = "Parse Failed"
	InvalidOption = "Invalid Option"
	NotFound      = "Not Found"
	NotMatch      = "Not Match"
	Invalid       = "Invalid"

	ErrAuthorizationHeaderKey    = "Authorization Header Key"
	ErrAuthorizationHeaderFormat = "Authorization Header Format"
	ErrAuthorizationTypeBearer   = "Authorization Type Bearer"
	ErrAuthorizationPayloadKey   = "Authorization Payload Key"
	ErrAuthorizationInvalidToken = "Authorization Invalid Token"
)

type RequestValidationError struct {
	FieldName    string `json:"field_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (rve RequestValidationError) Error() string {
	return fmt.Sprintf(
		"Field Name: %s, Error Code: %s, Error Message: %s",
		rve.FieldName,
		rve.ErrorCode,
		rve.ErrorMessage,
	)
}

func NewRequestValidationError(errorCode, fieldName string) RequestValidationError {
	return RequestValidationError{
		FieldName:    fieldName,
		ErrorCode:    errorCode,
		ErrorMessage: Message(errorCode),
	}
}

type AuthorizationValidationError struct {
	FieldName    string `json:"field_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func NewAuthorizationValidationError(errorCode, fieldName string) RequestValidationError {
	return RequestValidationError{
		FieldName:    fieldName,
		ErrorCode:    errorCode,
		ErrorMessage: Message(errorCode),
	}
}

func (ave AuthorizationValidationError) Error() string {
	return fmt.Sprintf(
		"Field Name: %s, Error Code: %s, Error Message: %s",
		ave.FieldName,
		ave.ErrorCode,
		ave.ErrorMessage,
	)
}

func Message(errorCode string) string {
	switch errorCode {
	case Required:
		return "This field is required"
	case Alphanumeric:
		return "Alphanumeric only"
	case Alpha:
		return "Alphabet only"
	case Numeric:
		return "Number only"
	case Float:
		return "Float only"
	case ParseFailed:
		return "Parsing failed. Make sure the input is correct."
	case InvalidOption:
		return "This value is not available in options. Please give the correct options."
	case NotFound:
		return "Data not found."
	case NotMatch:
		return "Password didn't match with confirmation password"
	case Invalid:
		return "Invalid value"
	case ErrAuthorizationHeaderKey:
		return "authorization header is not provided"
	case ErrAuthorizationHeaderFormat:
		return "invalid authorization header format"
	case ErrAuthorizationTypeBearer:
		return "unsupported authorization type"
	case ErrAuthorizationPayloadKey:
		return "invalid authorization payload"
	case ErrAuthorizationInvalidToken:
		return "invalid token"
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
	fields := fiber.Map{
		"ip":   c.IP(),
		"file": file,
		"line": line,
	}

	if userErr, ok := err.(domain.UserError); ok {
		errorResponse["error_code"] = strconv.Itoa(userErr.Code)
		errorResponse["error_message"] = userErr.Error()

		// --------------------------------------------------------
		fields["error_code"] = userErr.Code
		fields["error_message"] = userErr.Error()
		Log.Error().Interface("user error", fields).Send()
		// --------------------------------------------------------

		return c.Status(http.StatusBadRequest).JSON(errorResponse)
	} else if reqValidationErr, ok := err.(RequestValidationError); ok {
		errorResponse["field_name"] = reqValidationErr.FieldName
		errorResponse["error_code"] = reqValidationErr.ErrorCode
		errorResponse["error_message"] = reqValidationErr.ErrorMessage

		// --------------------------------------------------------
		fields["error_code"] = reqValidationErr.ErrorCode
		fields["field_name"] = reqValidationErr.FieldName
		fields["error_message"] = reqValidationErr.ErrorMessage
		Log.Error().Interface("validation error", fields).Send()
		// --------------------------------------------------------

		return c.Status(http.StatusBadRequest).JSON(errorResponse)
	} else if authValidationErr, ok := err.(AuthorizationValidationError); ok {
		errorResponse["field_name"] = authValidationErr.FieldName
		errorResponse["error_code"] = authValidationErr.ErrorCode
		errorResponse["error_message"] = authValidationErr.ErrorMessage

		// --------------------------------------------------------
		fields["error_code"] = authValidationErr.ErrorCode
		fields["field_name"] = authValidationErr.FieldName
		fields["error_message"] = authValidationErr.ErrorMessage
		Log.Error().Interface("authorization error", fields).Send()
		// --------------------------------------------------------

		return c.Status(http.StatusUnauthorized).JSON(errorResponse)
	}

	// --------------------------------------------------------
	fields["field_name"] = "undefined"
	fields["error_code"] = -1
	fields["error_message"] = err.Error()
	Log.Error().Interface("internal error", fields).Send()
	// --------------------------------------------------------

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
