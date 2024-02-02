package main

import (
	"net/http"

	"Tahlilchi.uz/admin"
	"Tahlilchi.uz/client"
	"Tahlilchi.uz/routerFuncs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Router() {
	r := mux.NewRouter()

	r.HandleFunc("/", routerFuncs.Root).Methods("GET") // .Schemes(os.Getenv("SCHEMES"))  add .Host(os.Getenv("HOST")) in the end

	admin.AdminRouter(r)
	client.ClientRouter(r)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PATCH", "DELETE", "OPTIONS"})
	handlers.AllowCredentials()

	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}
