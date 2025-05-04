package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type SuccessResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func ValidationError(errs validator.ValidationErrors) ErrorResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must have more than %s characters", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must have less than %s characters", err.Field(), err.Param()))
		case "gte":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must have greater than %s", err.Field(), err.Param()))
		case "lte":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must have less than %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return ErrorResponse{
		Status: http.StatusBadRequest,
		Error:  strings.Join(errMsgs, ", "),
	}
}
