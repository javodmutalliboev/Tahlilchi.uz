package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data"`
}

func Res(w http.ResponseWriter, status string, statusCode int, data any) {
	res := Response{
		Status:     status,
		StatusCode: statusCode,
		Data:       data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}
