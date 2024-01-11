package main

import (
	"net/http"
	"os"

	"Tahlilchi.uz/routerFuncs"
	"github.com/gorilla/mux"
)

func Router() {
	r := mux.NewRouter()

	r.HandleFunc("/", routerFuncs.Root).Methods("GET").Schemes(os.Getenv("SCHEMES")) // add .Host(os.Getenv("HOST")) in the end

	http.ListenAndServe(":8080", r)
}
