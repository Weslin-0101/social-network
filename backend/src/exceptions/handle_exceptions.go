package exceptions

import (
	"backend/src/model"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	URL 	string 	`json:"url"`
	Type 	int 	`json:"type"`
	Message string 	`json:"message"`
}

type ValidationErrorResponse struct {
	URL     string `json:"url"`
	Type    int    `json:"type"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Code    string `json:"code"`
}

func HandleResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

func HandleError(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	err error,
) {
	if validationErr, ok := err.(model.ValidationError); ok {
		errorResponse := ValidationErrorResponse {
			URL:	 r.URL.Path,
			Type: 	 statusCode,
			Message: validationErr.Message,
			Field:   validationErr.Field,
			Code:    validationErr.Code,
		}
		HandleResponse(w, statusCode, errorResponse)
		return
	}

	errorResponse := ErrorResponse {
		URL:	 r.URL.Path,
		Type: 	 statusCode,
		Message: err.Error(),
	}

	HandleResponse(w, statusCode, errorResponse)
}

func HandleErrorWithCustomMessage(
	w http.ResponseWriter,
	r *http.Request,
	statuCode int,
	message string,
) {
	errorResponse := ErrorResponse {
		URL:	 r.URL.Path,
		Type: 	 statuCode,
		Message: message,
	}

	HandleResponse(w, statuCode, errorResponse)
}