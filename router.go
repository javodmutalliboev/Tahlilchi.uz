package main

import (
	"net/http"
	"os"

	"Tahlilchi.uz/routerFuncs"
	"github.com/gorilla/mux"
)

func Router() {
	r := mux.NewRouter()

	r.HandleFunc("/", routerFuncs.Root).Methods("GET").Host(os.Getenv("HOST")).Schemes(os.Getenv("SCHEMES"))

	http.ListenAndServe(":8080", r)
}
