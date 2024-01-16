package client

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ClientRouter(r *mux.Router) {
	clientRouter := r.PathPrefix("/client").Subrouter()
	clientRouter.HandleFunc("/appeal", func(w http.ResponseWriter, r *http.Request) {})
}

type Appeal struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	PhoneNumber string    `json:"phone_number"`
	Message     string    `json:"message"`
	Picture     []byte    `json:"picture"`
	Video       []byte    `json:"video"`
	CreatedAt   time.Time `json:"created_at"`
}
