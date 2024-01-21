package admin

import (
	"Tahlilchi.uz/authPackage"
	"Tahlilchi.uz/middleware"
	"github.com/gorilla/mux"
)

func AdminRouter(r *mux.Router) *mux.Router {
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/login", login).Methods("POST") // .Schemes(os.Getenv("SCHEMES"))

	forgotPasswordRouter := adminRouter.PathPrefix("/forgot-password").Subrouter()
	forgotPasswordRouter.HandleFunc("/email", forgotPasswordEmail).Methods("POST")
	forgotPasswordRouter.HandleFunc("/i-code", forgotPasswordICode).Methods("POST")
	forgotPasswordRouter.HandleFunc("/new-password", forgotPasswordNewPassword).Methods("POST")

	newsRouter := adminRouter.PathPrefix("/news").Subrouter()
	newsRouter.HandleFunc("/category", middleware.Chain(addCategory, authPackage.AdminAuth())).Methods("POST")
	newsRouter.HandleFunc("/subcategory", middleware.Chain(addSubcategory, authPackage.AdminAuth())).Methods("POST")
	newsRouter.HandleFunc("/regions", middleware.Chain(addRegions, authPackage.AdminAuth())).Methods("POST")
	newsRouter.HandleFunc("/post", middleware.Chain(addNewsPost, authPackage.AdminAuth())).Methods("POST")

	newsPostRouter := newsRouter.PathPrefix("/post").Subrouter()
	newsPostRouter.HandleFunc("/edit/{id}", middleware.Chain(editNewsPost, authPackage.AdminAuth())).Methods("PATCH")

	return adminRouter
}
