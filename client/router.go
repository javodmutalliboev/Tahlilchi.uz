package client

import (
	"github.com/gorilla/mux"
)

func ClientRouter(r *mux.Router) {
	clientRouter := r.PathPrefix("/client").Subrouter()
	clientRouter.HandleFunc("/appeal", addAppeal).Methods("POST")

	newsRouter := clientRouter.PathPrefix("/news").Subrouter()
	newsRouter.HandleFunc("/category/{category}", getNewsByCategory).Methods("GET")
	newsRouter.HandleFunc("/subcategory/{subcategory}", getNewsBySubCategory).Methods("GET")
	// get news region list
	newsRouter.HandleFunc("/region/list", getNewsRegionList).Methods("GET")
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

	// article router
	articleRouter := clientRouter.PathPrefix("/article").Subrouter()
	// route to get article list by category
	articleRouter.HandleFunc("/category/{category}", getArticleListByCategory).Methods("GET")
	// route to get article list by related
	articleRouter.HandleFunc("/related/{related}", getArticleListByRelated).Methods("GET")
	// route to get all articles
	articleRouter.HandleFunc("/list", getAllArticles).Methods("GET")
	// route to get article category list
	articleRouter.HandleFunc("/category/list", getArticleCategoryList).Methods("GET")

	// article photos router
	articlePhotosRouter := articleRouter.PathPrefix("/{id}/photos").Subrouter()
	// route to get article photos
	articlePhotosRouter.HandleFunc("", getArticlePhotos).Methods("GET")
	// route to get article photo
	articlePhotosRouter.HandleFunc("/{photo_id}", getArticlePhoto).Methods("GET")
	// article photos router

	// route to get article cover image
	articleRouter.HandleFunc("/{id}/cover_image", getArticleCoverImage).Methods("GET")

	bpPostRouter := clientRouter.PathPrefix("/business-promotional/post").Subrouter()
	bpPostRouter.HandleFunc("/list", getBusinessPromotionalPosts).Methods("GET")
	// business promotional post photo router
	bpPostPhotoRouter := bpPostRouter.PathPrefix("/{id}/photo").Subrouter()
	// route to get business promotional post photo list
	bpPostPhotoRouter.HandleFunc("", getBusinessPromotionalPostPhotoList).Methods("GET")
	// route to get business promotional post photo
	bpPostPhotoRouter.HandleFunc("/{photo_id}", getBusinessPromotionalPostPhoto).Methods("GET")
	// route to get business promotional post cover image
	bpPostRouter.HandleFunc("/{id}/cover_image", getBusinessPromotionalPostCoverImage).Methods("GET")

	eNewspaperRouter := clientRouter.PathPrefix("/e-newspaper").Subrouter()
	// route to get e-newspaper category list
	eNewspaperRouter.HandleFunc("/category/list", getENewspaperCategoryList).Methods("GET")
	// route to get e-newspaper list by category
	eNewspaperRouter.HandleFunc("/category/{category}", getENewspaperListByCategory).Methods("GET")
	// route to get /e-newspaper/list
	eNewspaperRouter.HandleFunc("/list", getENewspaperList).Methods("GET")
	// route to get /e-newspaper/{id}/cover_image
	eNewspaperRouter.HandleFunc("/{id}/cover_image", getENewspaperCoverImage).Methods("GET")
	// route to get /e-newspaper/{id}/file/{alphabet} where file is pdf, alphabet is latin or cyrillic
	eNewspaperRouter.HandleFunc("/{id}/file/{alphabet}", getENewspaperFile).Methods("GET")

	photoGalleryRouter := clientRouter.PathPrefix("/photo-gallery").Subrouter()
	photoGalleryRouter.HandleFunc("/list", getPhotoGalleryList).Methods("GET")
	photoGalleryRouter.HandleFunc("/{id}/photos", getPhotoGalleryPhotos).Methods("GET")
	// route to get photo gallery photo
	photoGalleryRouter.HandleFunc("/{id}/photo/{photo_id}", getPhotoGalleryPhoto).Methods("GET")

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

	// video news router
	videoNewsRouter := clientRouter.PathPrefix("/video-news").Subrouter()
	// route to get video news list
	videoNewsRouter.HandleFunc("/list", getVideoNewsList).Methods("GET") // Go file path: client/video_news.go

	// article comment router
	articleCommentRouter := articleRouter.PathPrefix("/{id}/comment").Subrouter()
	// route to add article comment
	articleCommentRouter.HandleFunc("", addArticleComment).Methods("POST") // Go file path: client/article_comment.go
	// route to get article comment list
	articleCommentRouter.HandleFunc("/list", getArticleCommentList).Methods("GET") // Go file path: client/article_comment.go

	// e_newspaper comment router
	eNewspaperCommentRouter := eNewspaperRouter.PathPrefix("/{id}/comment").Subrouter()
	// route to add e-newspaper comment
	eNewspaperCommentRouter.HandleFunc("", addENewspaperComment).Methods("POST") // Go file path: client/e_newspaper_comment.go
	// route to get e-newspaper comment list
	eNewspaperCommentRouter.HandleFunc("/list", getENewspaperCommentList).Methods("GET") // Go file path: client/e_newspaper_comment.go

	// news post comment router
	newsPostCommentRouter := newsPostRouter.PathPrefix("/{id}/comment").Subrouter()
	// route to add news post comment
	newsPostCommentRouter.HandleFunc("", addNewsPostComment).Methods("POST") // Go file path: client/news_post_comment.go
	// route to get news post comment list
	newsPostCommentRouter.HandleFunc("/list", getNewsPostCommentList).Methods("GET") // Go file path: client/news_post_comment.go

	// video news comment router
	videoNewsCommentRouter := videoNewsRouter.PathPrefix("/{id}/comment").Subrouter() // Go file path: client/video_news_comment.go
	// route to add video news comment
	videoNewsCommentRouter.HandleFunc("", addVideoNewsComment).Methods("POST") // Go file path: client/video_news_comment.go
	// route to get video news comment list
	videoNewsCommentRouter.HandleFunc("/list", getVideoNewsCommentList).Methods("GET") // Go file path: client/video_news_comment.go
}
