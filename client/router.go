package client

import (
	"github.com/gorilla/mux"
)

func ClientRouter(r *mux.Router) {
	clientRouter := r.PathPrefix("/client").Subrouter()
	clientRouter.HandleFunc("/appeal", Appeal).Methods("POST")

	newsRouter := clientRouter.PathPrefix("/news").Subrouter()
	newsRouter.HandleFunc("/category/{category}", getNewsByCategory).Methods("GET")
	newsRouter.HandleFunc("/subcategory/{subcategory}", getNewsBySubCategory).Methods("GET")
	newsRouter.HandleFunc("/region/{region}", getNewsByRegion).Methods("GET")
	newsRouter.HandleFunc("/top", getTopNews).Methods("GET")
	newsRouter.HandleFunc("/latest", getLatestNews).Methods("GET")
	newsRouter.HandleFunc("/related/{id}", getRelatedNewsPosts).Methods("GET")
	newsRouter.HandleFunc("", getAllNewsPosts).Methods("GET")
	newsRouter.HandleFunc("/category", getCategoryList).Methods("GET")

	articleRouter := clientRouter.PathPrefix("/article").Subrouter()
	articleRouter.HandleFunc("/category/{category}", getArticleListByCategory).Methods("GET")
	articleRouter.HandleFunc("/related/{related}", getArticleListByRelated).Methods("GET")
	articleRouter.HandleFunc("/list", getAllArticles).Methods("GET")

	bpPostRouter := clientRouter.PathPrefix("/business-promotional/post").Subrouter()
	bpPostRouter.HandleFunc("/list", getBusinessPromotionalPosts).Methods("GET")
}
