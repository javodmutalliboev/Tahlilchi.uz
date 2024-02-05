package client

import (
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

// searchArticle is a handler function for the /search/article route.
// It is used to search for articles.
// search columns: title_latin, description_latin, title_cyrillic, description_cyrillic, tags.
// tags is text[]. use @> to search for a tag in the tags column.
func searchArticle(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		log.Printf("%v: search query is empty", r.URL)
		response.Res(w, "error", http.StatusBadRequest, "search query is empty")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags FROM articles WHERE (title_latin ILIKE '%' || $1 || '%' OR description_latin ILIKE '%' || $1 || '%' OR title_cyrillic ILIKE '%' || $1 || '%' OR description_cyrillic ILIKE '%' || $1 || '%' OR tags @> ARRAY[$1]) AND archived = false AND completed = true", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.TitleCyrillic, &article.DescriptionLatin, &article.DescriptionCyrillic, &article.Videos, &article.Tags)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articles = append(articles, article)
	}

	response.Res(w, "success", http.StatusOK, articles)
}

// searchENewspaper is a handler function for the /search/e-newspaper route.
// It is used to search for e-newspapers.
// search columns: title_latin, title_cyrillic.
// where archived is false, completed is true.
func searchENewspaper(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		log.Printf("%v: search query is empty", r.URL)
		response.Res(w, "error", http.StatusBadRequest, "search query is empty")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic FROM e_newspapers WHERE (title_latin ILIKE '%' || $1 || '%' OR title_cyrillic ILIKE '%' || $1 || '%') AND archived = false AND completed = true", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var eNewspapers []ENewspaper
	for rows.Next() {
		var eNewspaper ENewspaper
		err := rows.Scan(&eNewspaper.ID, &eNewspaper.TitleLatin, &eNewspaper.TitleCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		eNewspapers = append(eNewspapers, eNewspaper)
	}

	response.Res(w, "success", http.StatusOK, eNewspapers)
}

// searchNews is a handler function for the /search/news route.
// It is used to search for news.
// search columns: title_latin, title_cyrillic, description_latin, description_cyrillic, tags.
// tags is text[].
// where archived is false, completed is true.
func searchNews(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		log.Printf("%v: search query is empty", r.URL)
		response.Res(w, "error", http.StatusBadRequest, "search query is empty")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE (title_latin ILIKE '%' || $1 || '%' OR title_cyrillic ILIKE '%' || $1 || '%' OR description_latin ILIKE '%' || $1 || '%' OR description_cyrillic ILIKE '%' || $1 || '%' OR tags @> ARRAY[$1]) AND archived = false AND completed = true", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var news []NewsPost
	for rows.Next() {
		var n NewsPost
		err := rows.Scan(&n.ID, &n.TitleLatin, &n.DescriptionLatin, &n.TitleCyrillic, &n.DescriptionCyrillic, &n.Video, &n.Tags)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		news = append(news, n)
	}

	response.Res(w, "success", http.StatusOK, news)
}

// searchPhotoGallery is a handler function for the /search/photo-gallery route.
// It is used to search for photo galleries.
// search columns: title_latin, title_cyrillic.
func searchPhotoGallery(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		log.Printf("%v: search query is empty", r.URL)
		response.Res(w, "error", http.StatusBadRequest, "search query is empty")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic FROM photo_gallery WHERE title_latin ILIKE '%' || $1 || '%' OR title_cyrillic ILIKE '%' || $1 || '%'", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleries []PhotoGallery
	for rows.Next() {
		var p PhotoGallery
		err := rows.Scan(&p.ID, &p.TitleLatin, &p.TitleCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photoGalleries = append(photoGalleries, p)
	}

	response.Res(w, "success", http.StatusOK, photoGalleries)
}

// searchPhotoGalleryPhotos is a handler function for the /search/photo-gallery/photos route.
// It is used to search for photos in a photo gallery.
// search columns: file_name.
func searchPhotoGalleryPhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	search := r.URL.Query().Get("search")
	if search == "" {
		log.Printf("%v: search query is empty", r.URL)
		response.Res(w, "error", http.StatusBadRequest, "search query is empty")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, file_name from photo_gallery_photos WHERE photo_gallery = $1 AND file_name ILIKE '%' || $2 || '%'", id, search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photos []PhotoGalleryPhoto
	for rows.Next() {
		var p PhotoGalleryPhoto
		err := rows.Scan(&p.ID, &p.FileName)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photos = append(photos, p)
	}

	response.Res(w, "success", http.StatusOK, photos)
}
