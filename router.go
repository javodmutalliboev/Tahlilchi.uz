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

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/login", admin.Login).Methods("POST") // .Schemes(os.Getenv("SCHEMES"))
	adminRouter.HandleFunc("/forgot-password/email", admin.ForgotPasswordEmail).Methods("POST")
	adminRouter.HandleFunc("/forgot-password/i_code", admin.ForgotPasswordICode).Methods("POST")

	http.ListenAndServe(":8080", r)
}
