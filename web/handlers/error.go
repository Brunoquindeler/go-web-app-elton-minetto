package handlers

import (
	"encoding/json"
	"net/http"
)

func formatJSONError(message string) []byte {
	appError := struct {
		Message string `json:"message"`
	}{
		message,
	}

	response, err := json.Marshal(appError)
	if err != nil {
		return []byte(err.Error())
	}

	return response
}

func formatJSONValidationError(messages []string) []byte {
	appError := struct {
		Messages []string `json:"messages"`
	}{
		messages,
	}

	response, err := json.Marshal(appError)
	if err != nil {
		return []byte(err.Error())
	}

	return response
}

func writeResponseError(w http.ResponseWriter, errMessage string, status int) {
	w.WriteHeader(status)
	w.Write(formatJSONError(errMessage))
}

func writeResponseValidationError(w http.ResponseWriter, errMessages []string, status int) {
	w.WriteHeader(status)
	w.Write(formatJSONValidationError(errMessages))
}
