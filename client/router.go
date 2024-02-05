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

	newsPostRouter := newsRouter.PathPrefix("/post").Subrouter()
	newsPostRouter.HandleFunc("/{id}/photo", getNewsPostPhoto).Methods("GET")
	newsPostRouter.HandleFunc("/{id}/audio", getNewsPostAudio).Methods("GET")
	newsPostRouter.HandleFunc("/{id}/cover_image", getNewsPostCoverImage).Methods("GET")

	articleRouter := clientRouter.PathPrefix("/article").Subrouter()
	articleRouter.HandleFunc("/category/{category}", getArticleListByCategory).Methods("GET")
	articleRouter.HandleFunc("/related/{related}", getArticleListByRelated).Methods("GET")
	articleRouter.HandleFunc("/list", getAllArticles).Methods("GET")
	articleRouter.HandleFunc("/category", getArticleCategory).Methods("GET")

	bpPostRouter := clientRouter.PathPrefix("/business-promotional/post").Subrouter()
	bpPostRouter.HandleFunc("/list", getBusinessPromotionalPosts).Methods("GET")

	eNewspaperRouter := clientRouter.PathPrefix("/e-newspaper").Subrouter()
	eNewspaperRouter.HandleFunc("/list", getENewspaperList).Methods("GET")
	// route to get /e-newspaper/{id}/file/{alphabet} where file is pdf, alphabet is latin or cyrillic
	eNewspaperRouter.HandleFunc("/{id}/file/{alphabet}", getENewspaperFile).Methods("GET")

	photoGalleryRouter := clientRouter.PathPrefix("/photo-gallery").Subrouter()
	photoGalleryRouter.HandleFunc("/list", getPhotoGalleryList).Methods("GET")
	photoGalleryRouter.HandleFunc("/{id}/photos", getPhotoGalleryPhotos).Methods("GET")

	contactRouter := clientRouter.PathPrefix("/contact").Subrouter()
	contactRouter.HandleFunc("", getAdminContact).Methods("GET")

	// search router
	searchRouter := clientRouter.PathPrefix("/search").Subrouter()

	// route to search in articles
	searchRouter.HandleFunc("/article", searchArticle).Methods("GET")
	// route to search in e-newspapers
	searchRouter.HandleFunc("/e-newspaper", searchENewspaper).Methods("GET")
	// route to search news_posts
	searchRouter.HandleFunc("/news", searchNews).Methods("GET")
	// route to search photo_gallery
	searchRouter.HandleFunc("/photo-gallery", searchPhotoGallery).Methods("GET")
	// route to search photo_gallery_photos
	searchRouter.HandleFunc("/photo-gallery/{id}/photos", searchPhotoGalleryPhotos).Methods("GET")
}
