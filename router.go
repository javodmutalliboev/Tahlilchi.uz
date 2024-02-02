package main

import (
	"net/http"
	"os"

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
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ADMINCLIENT")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PATCH", "DELETE", "OPTIONS"})
	credentials := handlers.AllowCredentials()

	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk, credentials)(r))
}
