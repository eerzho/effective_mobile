package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Error struct {
	Status string   `json:"status"`
	Error  string   `json:"error"`
	Errors []string `json:"errors,omitempty"`
}

type Success struct {
	Status string `json:"status"`
}

func OK() Success {
	return Success{Status: StatusOK}
}

func Err(msg string) Error {
	return Error{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationErr(errs validator.ValidationErrors) Error {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", strings.ToLower(err.Field())))
		}
	}

	return Error{
		Status: StatusError,
		Error:  "validation error",
		Errors: errMsgs,
	}
}
