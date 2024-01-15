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

func Res(w http.ResponseWriter, Status string, StatusCode int, Message string) {
	res := Response{
		Status:     Status,
		StatusCode: StatusCode,
		Message:    Message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
