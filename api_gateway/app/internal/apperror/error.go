package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	ErrNotFound = NewAppError("not found", "NS-000010", "", nil)
)

type AppError struct {
	Err              error    `json:"-"`
	Message          string   `json:"message,omitempty"`
	DeveloperMessage string   `json:"developer_message,omitempty"`
	Code             string   `json:"code,omitempty"`
	Errors           []string `json:"errors,omitempty"`
}

func NewAppError(message, code, developerMessage string, errors []string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
		Errors:           errors,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(message, "NS-000003", "", nil)
}

func BadRequestError(message string) *AppError {
	return NewAppError(message, "NS-000002", "some thing wrong with your data", nil)
}

func ValidationError(errors []string) *AppError {
	return NewAppError("validation error", "NS-000004", "some thing wrong with your data", errors)
}

func systemError(developerMessage string) *AppError {
	return NewAppError("system error", "NS-000001", developerMessage, nil)
}

func APIError(code, message, developerMessage string) *AppError {
	return NewAppError(message, code, developerMessage, nil)
}
