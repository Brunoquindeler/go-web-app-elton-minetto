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

func writeResponseError(w http.ResponseWriter, errMessage string, status int) {
	w.Write(formatJSONError(errMessage))
	w.WriteHeader(status)
}
