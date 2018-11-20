package responder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func (e ErrorResponse) Error() string {
	var b bytes.Buffer
	for _, err := range e.Errors {
		b.WriteString(fmt.Sprintf("message: %v, field: %v", err.Message, err.Field))
	}
	return b.String()
}

type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func Respond(w http.ResponseWriter, httpStatus int) error {
	w.WriteHeader(httpStatus)
	return nil
}

func RespondResult(w http.ResponseWriter, i interface{}, httpStatus int) error {
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(i)
}

func RespondError(w http.ResponseWriter, message, field string, httpStatus int) error {
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(&ErrorResponse{
		Errors: []Error{
			Error{
				Message: message,
				Field:   field,
			},
		}})
}

func RespondErrors(w http.ResponseWriter, errors []Error, httpStatus int) error {
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(&ErrorResponse{
		Errors: errors,
	})
}
