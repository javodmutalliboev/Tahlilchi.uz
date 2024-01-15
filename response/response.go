package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func Res(w http.ResponseWriter, status string, statusCode int, message string) {
	res := Response{
		Status:     status,
		StatusCode: statusCode,
		Message:    message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}
