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

	// news router
	newsRouter := adminRouter.PathPrefix("/news").Subrouter()
	// route to add news category
	newsRouter.HandleFunc("/category", middleware.Chain(addCategory, authPackage.AdminAuth())).Methods("POST")
	// route to get news category list
	newsRouter.HandleFunc("/category", middleware.Chain(getCategoryList, authPackage.AdminAuth())).Methods("GET")
	// route to update news category
	newsRouter.HandleFunc("/category/{id}", middleware.Chain(updateCategory, authPackage.AdminAuth())).Methods("PATCH") // Go file path: admin/news.go
	// route to delete news category
	newsRouter.HandleFunc("/category/{id}", middleware.Chain(deleteCategory, authPackage.AdminAuth())).Methods("DELETE") // Go file path: admin/news.go
	// route to add news subcategory
	newsRouter.HandleFunc("/subcategory", middleware.Chain(addSubcategory, authPackage.AdminAuth())).Methods("POST")
	// route to get news subcategory list
	newsRouter.HandleFunc("/subcategory", middleware.Chain(getSubCategoryList, authPackage.AdminAuth())).Methods("GET")
	// route news/category/{id}/subcategory/list
	newsRouter.HandleFunc("/category/{id}/subcategory/list", middleware.Chain(getSubCategoryListByCategory, authPackage.AdminAuth())).Methods("GET")
	// route to add news region
	newsRouter.HandleFunc("/regions", middleware.Chain(addRegions, authPackage.AdminAuth())).Methods("POST")
	newsRouter.HandleFunc("/post", middleware.Chain(addNewsPost, authPackage.AdminAuth())).Methods("POST")

	// news region router
	newsRegionRouter := newsRouter.PathPrefix("/regions").Subrouter()
	// route to get news region list
	newsRegionRouter.HandleFunc("", middleware.Chain(getRegions, authPackage.AdminAuth())).Methods("GET") // Go file path: admin/news.go
	// route to update news region
	newsRegionRouter.HandleFunc("/{id}", middleware.Chain(updateRegion, authPackage.AdminAuth())).Methods("PATCH") // Go file path: admin/news.go
	// route to delete news region
	newsRegionRouter.HandleFunc("/{id}", middleware.Chain(deleteRegion, authPackage.AdminAuth())).Methods("DELETE") // Go file path: admin/news.go

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

	// article router
	articleRouter := adminRouter.PathPrefix("/article").Subrouter()
	// route to add article category
	articleRouter.HandleFunc("/category", middleware.Chain(addArticleCategory, authPackage.AdminAuth())).Methods("POST")
	// route to get article category list
	articleRouter.HandleFunc("/category", middleware.Chain(getArticleCategory, authPackage.AdminAuth())).Methods("GET")
	// route to update article category
	articleRouter.HandleFunc("/category/{id}", middleware.Chain(updateArticleCategory, authPackage.AdminAuth())).Methods("PATCH")
	// route to delete article category
	articleRouter.HandleFunc("/category/{id}", middleware.Chain(deleteArticleCategory, authPackage.AdminAuth())).Methods("DELETE")
	// route to add article
	articleRouter.HandleFunc("", middleware.Chain(addArticle, authPackage.AdminAuth())).Methods("POST")

	// route to edit article
	articleRouter.HandleFunc("/edit/{id}", middleware.Chain(editArticle, authPackage.AdminAuth())).Methods("PATCH")
	// article photos router
	articlePhotosRouter := articleRouter.PathPrefix("/{id}/photos").Subrouter()
	// route to add article photos
	articlePhotosRouter.HandleFunc("/add", middleware.Chain(addArticlePhotos, authPackage.AdminAuth())).Methods("POST")
	// route to get article photos
	articlePhotosRouter.HandleFunc("", middleware.Chain(getArticlePhotos, authPackage.AdminAuth())).Methods("GET")
	// route to get article photo
	articlePhotosRouter.HandleFunc("/{photo_id}", middleware.Chain(getArticlePhoto, authPackage.AdminAuth())).Methods("GET")
	// route to delete article photo
	articlePhotosRouter.HandleFunc("/delete/{photo_id}", middleware.Chain(deleteArticlePhoto, authPackage.AdminAuth())).Methods("DELETE")
	// route to edit article

	// route to get article cover_image
	articleRouter.HandleFunc("/{id}/cover_image", middleware.Chain(getArticleCoverImage, authPackage.AdminAuth())).Methods("GET")
	// route to delete article
	articleRouter.HandleFunc("/delete/{id}", middleware.Chain(deleteArticle, authPackage.AdminAuth())).Methods("DELETE")
	// route to archive article
	articleRouter.HandleFunc("/archive/{id}", middleware.Chain(archiveArticle, authPackage.AdminAuth())).Methods("PATCH")
	// route to get article count
	articleRouter.HandleFunc("/count/{period}", middleware.Chain(getArticleCount, authPackage.AdminAuth())).Methods("GET")
	// route to get article count all
	articleRouter.HandleFunc("/count", middleware.Chain(getArticleCountAll, authPackage.AdminAuth())).Methods("GET")
	// route to get article list
	articleRouter.HandleFunc("/list", middleware.Chain(getArticles, authPackage.AdminAuth())).Methods("GET")
	// route to unarchive article
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

	// e-newspaper router
	eNewspaperRouter := adminRouter.PathPrefix("/e-newspaper").Subrouter()
	// route to add e-newspaper category
	eNewspaperRouter.HandleFunc("/category", middleware.Chain(addENewspaperCategory, authPackage.AdminAuth())).Methods("POST")
	// route to get e-newspaper category list
	eNewspaperRouter.HandleFunc("/category/list", middleware.Chain(getENewspaperCategoryList, authPackage.AdminAuth())).Methods("GET")
	// route to update e-newspaper category
	eNewspaperRouter.HandleFunc("/category/{id}", middleware.Chain(updateENewspaperCategory, authPackage.AdminAuth())).Methods("PATCH")
	// route to delete e-newspaper category
	eNewspaperRouter.HandleFunc("/category/{id}", middleware.Chain(deleteENewspaperCategory, authPackage.AdminAuth())).Methods("DELETE")
	// route to add e-newspaper
	eNewspaperRouter.HandleFunc("/add", middleware.Chain(addENewspaper, authPackage.AdminAuth())).Methods("POST")
	// route to edit e-newspaper
	eNewspaperRouter.HandleFunc("/edit/{id}", middleware.Chain(editENewspaper, authPackage.AdminAuth())).Methods("PATCH")
	// route to delete e-newspaper
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
	// route to search in news_posts table
	searchRouter.HandleFunc("/news", middleware.Chain(searchNews, authPackage.AdminAuth())).Methods("GET")
	// route to search in photo_gallery table
	searchRouter.HandleFunc("/photo-gallery", middleware.Chain(searchPhotoGallery, authPackage.AdminAuth())).Methods("GET")
	// route to search in photo_gallery_photos table
	searchRouter.HandleFunc("/photo-gallery/{id}/photos", middleware.Chain(searchPhotoGalleryPhotos, authPackage.AdminAuth())).Methods("GET")

	// video news router: location: admin/video-news.go
	videoNewsRouter := adminRouter.PathPrefix("/video-news").Subrouter()
	// route to add video news
	videoNewsRouter.HandleFunc("/add", middleware.Chain(addVideoNews, authPackage.AdminAuth())).Methods("POST")
	// route to update video news
	videoNewsRouter.HandleFunc("/update/{id}", middleware.Chain(updateVideoNews, authPackage.AdminAuth())).Methods("PATCH")
	// route to delete video news
	videoNewsRouter.HandleFunc("/delete/{id}", middleware.Chain(deleteVideoNews, authPackage.AdminAuth())).Methods("DELETE")
	// route to get a video news list
	videoNewsRouter.HandleFunc("/list", middleware.Chain(getVideoNewsList, authPackage.AdminAuth())).Methods("GET")

	// article comment router
	articleCommentRouter := articleRouter.PathPrefix("/{id}/comment").Subrouter()
	// route to get article comment list
	articleCommentRouter.HandleFunc("/list", middleware.Chain(getArticleCommentList, authPackage.AdminAuth())).Methods("GET") // Go file path: admin/article_comment.go
	// route approve/disapprove article comment
	articleCommentRouter.HandleFunc("/approve/{comment_id}", middleware.Chain(approveArticleComment, authPackage.AdminAuth())).Methods("PATCH") // Go file path: admin/article_comment.go

	// news post comment router
	newsPostCommentRouter := newsPostRouter.PathPrefix("/{id}/comment").Subrouter()
	// route to get news post comment list
	newsPostCommentRouter.HandleFunc("/list", middleware.Chain(getNewsPostCommentList, authPackage.AdminAuth())).Methods("GET") // Go file path: admin/news_post_comment.go
	// route approve/disapprove news post comment
	newsPostCommentRouter.HandleFunc("/approve/{comment_id}", middleware.Chain(approveNewsPostComment, authPackage.AdminAuth())).Methods("PATCH") // Go file path: admin/news_post_comment.go

	return adminRouter
}
