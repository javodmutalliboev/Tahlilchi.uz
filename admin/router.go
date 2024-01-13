package admin

import (
	"github.com/gorilla/mux"
)

func AdminRouter(r *mux.Router) *mux.Router {
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/login", Login).Methods("POST") // .Schemes(os.Getenv("SCHEMES"))

	forgotPasswordRouter := adminRouter.PathPrefix("/forgot-password").Subrouter()
	forgotPasswordRouter.HandleFunc("/email", ForgotPasswordEmail).Methods("POST")
	forgotPasswordRouter.HandleFunc("/i-code", ForgotPasswordICode).Methods("POST")
	forgotPasswordRouter.HandleFunc("/new-password", ForgotPasswordNewPassword).Methods("POST")

	return adminRouter
}
