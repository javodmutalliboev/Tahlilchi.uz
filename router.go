package main

import (
	"net/http"

	"Tahlilchi.uz/admin"
	"Tahlilchi.uz/routerFuncs"
	"github.com/gorilla/mux"
)

func Router() {
	r := mux.NewRouter()

	r.HandleFunc("/", routerFuncs.Root).Methods("GET") // .Schemes(os.Getenv("SCHEMES"))  add .Host(os.Getenv("HOST")) in the end

	admin.AdminRouter(r)

	http.ListenAndServe(":8080", r)
}
