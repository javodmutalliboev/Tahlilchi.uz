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
	newsRouter.HandleFunc("/category", middleware.Chain(getCategoryList, authPackage.AdminAuth())).Methods("GET")
	newsRouter.HandleFunc("/subcategory", middleware.Chain(addSubcategory, authPackage.AdminAuth())).Methods("POST")
	newsRouter.HandleFunc("/subcategory", middleware.Chain(getSubCategoryList, authPackage.AdminAuth())).Methods("GET")
	newsRouter.HandleFunc("/regions", middleware.Chain(addRegions, authPackage.AdminAuth())).Methods("POST")
	newsRouter.HandleFunc("/post", middleware.Chain(addNewsPost, authPackage.AdminAuth())).Methods("POST")

	newsPostRouter := newsRouter.PathPrefix("/post").Subrouter()
	newsPostRouter.HandleFunc("/edit/{id}", middleware.Chain(editNewsPost, authPackage.AdminAuth())).Methods("PATCH")
	newsPostRouter.HandleFunc("/delete/{id}", middleware.Chain(deleteNewsPost, authPackage.AdminAuth())).Methods("DELETE")
	newsPostRouter.HandleFunc("/archive/{id}", middleware.Chain(archiveNewsPost, authPackage.AdminAuth())).Methods("PATCH")
	newsPostRouter.HandleFunc("/count/{period}", middleware.Chain(getNewsPostCount, authPackage.AdminAuth())).Methods("GET")
	newsPostRouter.HandleFunc("/count", middleware.Chain(getNewsPostCountAll, authPackage.AdminAuth())).Methods("GET")
	newsPostRouter.HandleFunc("/list", middleware.Chain(getNewsPosts, authPackage.AdminAuth())).Methods("GET")
	newsPostRouter.HandleFunc("/unarchive/{id}", middleware.Chain(unArchiveNewsPost, authPackage.AdminAuth())).Methods("PATCH")
	newsPostRouter.HandleFunc("/{id}/photo", middleware.Chain(getNewsPostPhoto, authPackage.AdminAuth())).Methods("GET")
	newsPostRouter.HandleFunc("/{id}/audio", middleware.Chain(getNewsPostAudio, authPackage.AdminAuth())).Methods("GET")
	newsPostRouter.HandleFunc("/{id}/cover_image", middleware.Chain(getNewsPostCoverImage, authPackage.AdminAuth())).Methods("GET")
	// route to make news post completed field true/false
	newsPostRouter.HandleFunc("/completed/{id}", middleware.Chain(newsPostCompleted, authPackage.AdminAuth())).Methods("PATCH")

	articleRouter := adminRouter.PathPrefix("/article").Subrouter()
	articleRouter.HandleFunc("/category", middleware.Chain(addArticleCategory, authPackage.AdminAuth())).Methods("POST")
	articleRouter.HandleFunc("/category", middleware.Chain(getArticleCategory, authPackage.AdminAuth())).Methods("GET")
	articleRouter.HandleFunc("", middleware.Chain(addArticle, authPackage.AdminAuth())).Methods("POST")
	articleRouter.HandleFunc("/edit/{id}", middleware.Chain(editArticle, authPackage.AdminAuth())).Methods("PATCH")
	articleRouter.HandleFunc("/delete/{id}", middleware.Chain(deleteArticle, authPackage.AdminAuth())).Methods("DELETE")
	articleRouter.HandleFunc("/archive/{id}", middleware.Chain(archiveArticle, authPackage.AdminAuth())).Methods("PATCH")
	articleRouter.HandleFunc("/count/{period}", middleware.Chain(getArticleCount, authPackage.AdminAuth())).Methods("GET")
	articleRouter.HandleFunc("/count", middleware.Chain(getArticleCountAll, authPackage.AdminAuth())).Methods("GET")
	articleRouter.HandleFunc("/list", middleware.Chain(getArticles, authPackage.AdminAuth())).Methods("GET")
	articleRouter.HandleFunc("/unarchive/{id}", middleware.Chain(unArchiveArticle, authPackage.AdminAuth())).Methods("PATCH")
	// route to make article completed field true/false
	articleRouter.HandleFunc("/completed/{id}", middleware.Chain(articleCompleted, authPackage.AdminAuth())).Methods("PATCH")

	businessPromotionalRouter := adminRouter.PathPrefix("/business-promotional").Subrouter()
	businessPromotionalRouter.HandleFunc("/post", middleware.Chain(addBusinessPromotionalPost, authPackage.AdminAuth())).Methods("POST")

	businessPromotionalPostRouter := businessPromotionalRouter.PathPrefix("/post").Subrouter()
	businessPromotionalPostRouter.HandleFunc("/edit/{id}", middleware.Chain(editBusinessPromotionalPost, authPackage.AdminAuth())).Methods("PATCH")
	businessPromotionalPostRouter.HandleFunc("/delete/{id}", middleware.Chain(deleteBPPost, authPackage.AdminAuth())).Methods("DELETE")
	businessPromotionalPostRouter.HandleFunc("/archive/{id}", middleware.Chain(archiveBPPost, authPackage.AdminAuth())).Methods("PATCH")
	businessPromotionalPostRouter.HandleFunc("/count/{period}", middleware.Chain(getBusinessPromotionalPostCount, authPackage.AdminAuth())).Methods("GET")
	businessPromotionalPostRouter.HandleFunc("/list", middleware.Chain(getBusinessPromotionalPosts, authPackage.AdminAuth())).Methods("GET")
	businessPromotionalPostRouter.HandleFunc("/unarchive/{id}", middleware.Chain(unArchiveBPPost, authPackage.AdminAuth())).Methods("PATCH")
	// route to make business promotional post completed field true/false
	businessPromotionalPostRouter.HandleFunc("/completed/{id}", middleware.Chain(businessPromotionalPostCompleted, authPackage.AdminAuth())).Methods("PATCH")

	eNewspaperRouter := adminRouter.PathPrefix("/e-newspaper").Subrouter()
	eNewspaperRouter.HandleFunc("/add", middleware.Chain(addENewspaper, authPackage.AdminAuth())).Methods("POST")
	eNewspaperRouter.HandleFunc("/edit/{id}", middleware.Chain(editENewspaper, authPackage.AdminAuth())).Methods("PATCH")
	eNewspaperRouter.HandleFunc("/delete/{id}", middleware.Chain(deleteENewspaper, authPackage.AdminAuth())).Methods("DELETE")
	eNewspaperRouter.HandleFunc("/archive/{id}", middleware.Chain(archiveENewspaper, authPackage.AdminAuth())).Methods("PATCH")
	eNewspaperRouter.HandleFunc("/count/{period}", middleware.Chain(getENewspaperCount, authPackage.AdminAuth())).Methods("GET")
	eNewspaperRouter.HandleFunc("/count", middleware.Chain(getENewspaperCountAll, authPackage.AdminAuth())).Methods("GET")
	eNewspaperRouter.HandleFunc("/list", middleware.Chain(getENewspaperList, authPackage.AdminAuth())).Methods("GET")
	eNewspaperRouter.HandleFunc("/unarchive/{id}", middleware.Chain(unArchiveENewspaper, authPackage.AdminAuth())).Methods("PATCH")
	// route to make e-newspaper completed field true/false
	eNewspaperRouter.HandleFunc("/completed/{id}", middleware.Chain(eNewspaperCompleted, authPackage.AdminAuth())).Methods("PATCH")
	// route to get /e-newspaper/{id}/file/{alphabet} where file is pdf, alphabet is latin or cyrillic
	eNewspaperRouter.HandleFunc("/{id}/file/{alphabet}", middleware.Chain(getENewspaperFile, authPackage.AdminAuth())).Methods("GET")

	photoGalleryRouter := adminRouter.PathPrefix("/photo-gallery").Subrouter()
	photoGalleryRouter.HandleFunc("/add", middleware.Chain(addPhotoGallery, authPackage.AdminAuth())).Methods("POST")
	photoGalleryRouter.HandleFunc("/list", middleware.Chain(getPhotoGalleryList, authPackage.AdminAuth())).Methods("GET")

	photoGalleryPhotosRouter := photoGalleryRouter.PathPrefix("/{id}/photos").Subrouter()
	photoGalleryPhotosRouter.HandleFunc("/add", middleware.Chain(photoGalleryAddPhotos, authPackage.AdminAuth())).Methods("POST")
	photoGalleryPhotosRouter.Handle("", middleware.Chain(getPhotoGalleryPhotos, authPackage.AdminAuth())).Methods("GET")

	contactRouter := adminRouter.PathPrefix("/contact").Subrouter()

	contactAppealRouter := contactRouter.PathPrefix("/appeal").Subrouter()
	contactAppealRouter.HandleFunc("/list", middleware.Chain(appealList, authPackage.AdminAuth())).Methods("GET")
	contactAppealRouter.HandleFunc("/{id}/picture", middleware.Chain(appealPicture, authPackage.AdminAuth())).Methods("GET")
	contactAppealRouter.HandleFunc("/{id}/video", middleware.Chain(appealVideo, authPackage.AdminAuth())).Methods("GET")
	contactAppealRouter.HandleFunc("/count/{period}", middleware.Chain(getAppealCount, authPackage.AdminAuth())).Methods("GET")
	contactAppealRouter.HandleFunc("/count", middleware.Chain(getAppealCountAll, authPackage.AdminAuth())).Methods("GET")

	contactRouter.HandleFunc("", middleware.Chain(createAdminContact, authPackage.AdminAuth())).Methods("POST")
	contactRouter.HandleFunc("", middleware.Chain(getAdminContact, authPackage.AdminAuth())).Methods("GET")
	contactRouter.HandleFunc("/{id}", middleware.Chain(updateAdminContact, authPackage.AdminAuth())).Methods("PATCH")

	searchRouter := adminRouter.PathPrefix("/search").Subrouter()
	// route to search in appeals table
	searchRouter.HandleFunc("/appeal", middleware.Chain(searchAppeal, authPackage.AdminAuth())).Methods("GET")
	// route to search in articles table
	searchRouter.HandleFunc("/article", middleware.Chain(searchArticle, authPackage.AdminAuth())).Methods("GET")
	// route to search in business_promotional_posts table
	searchRouter.HandleFunc("/business-promotional", middleware.Chain(searchBusinessPromotional, authPackage.AdminAuth())).Methods("GET")
	// route to search in e_newspapers table
	searchRouter.HandleFunc("/e-newspaper", middleware.Chain(searchENewspaper, authPackage.AdminAuth())).Methods("GET")

	return adminRouter
}
