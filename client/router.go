package client

import (
	"github.com/gorilla/mux"
)

func ClientRouter(r *mux.Router) {
	clientRouter := r.PathPrefix("/client").Subrouter()
	clientRouter.HandleFunc("/appeal", Appeal).Methods("POST")

	newsRouter := clientRouter.PathPrefix("/news").Subrouter()
	newsRouter.HandleFunc("/all", getAllNewsPosts)
}
