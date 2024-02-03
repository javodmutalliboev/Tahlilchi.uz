package client

import (
	"Tahlilchi.uz/shared"
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

	newsPostRouter := newsRouter.PathPrefix("/post").Subrouter()
	newsPostRouter.HandleFunc("/{id}/photo", shared.GetNewsPostPhoto).Methods("GET")
	newsPostRouter.HandleFunc("/{id}/audio", shared.GetNewsPostAudio).Methods("GET")
	newsPostRouter.HandleFunc("/{id}/cover_image", shared.GetNewsPostCoverImage).Methods("GET")

	articleRouter := clientRouter.PathPrefix("/article").Subrouter()
	articleRouter.HandleFunc("/category/{category}", getArticleListByCategory).Methods("GET")
	articleRouter.HandleFunc("/related/{related}", getArticleListByRelated).Methods("GET")
	articleRouter.HandleFunc("/list", getAllArticles).Methods("GET")
	articleRouter.HandleFunc("/category", getArticleCategory).Methods("GET")

	bpPostRouter := clientRouter.PathPrefix("/business-promotional/post").Subrouter()
	bpPostRouter.HandleFunc("/list", getBusinessPromotionalPosts).Methods("GET")

	eNewspaperRouter := clientRouter.PathPrefix("/e-newspaper").Subrouter()
	eNewspaperRouter.HandleFunc("/list", getENewspaperList).Methods("GET")
	eNewspaperRouter.HandleFunc("/{alphabet}/{id}", getENewspaperByID).Methods("GET")

	photoGalleryRouter := clientRouter.PathPrefix("/photo-gallery").Subrouter()
	photoGalleryRouter.HandleFunc("/list", getPhotoGalleryList).Methods("GET")
	photoGalleryRouter.HandleFunc("/{id}/photos", getPhotoGalleryPhotos).Methods("GET")

	contactRouter := clientRouter.PathPrefix("/contact").Subrouter()
	contactRouter.HandleFunc("", getAdminContact).Methods("GET")
}

// contact.go
