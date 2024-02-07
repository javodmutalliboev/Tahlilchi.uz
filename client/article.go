package client

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// ArticleCategory is a struct to represent article category
type ArticleCategory struct {
	ID                  int    `json:"id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

// getArticleCategoryList is a handler to get article category list
func getArticleCategoryList(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("select * from article_category order by id")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	var articleCategoryList []ArticleCategory
	for rows.Next() {
		var a ArticleCategory
		err := rows.Scan(&a.ID, &a.TitleLatin, &a.DescriptionLatin, &a.TitleCyrillic, &a.DescriptionCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articleCategoryList = append(articleCategoryList, a)
	}

	response.Res(w, "success", http.StatusOK, articleCategoryList)
}

// Article is a struct to represent article
type Article struct {
	ID                  int            `json:"id"`
	TitleLatin          string         `json:"title_latin"`
	DescriptionLatin    string         `json:"description_latin"`
	TitleCyrillic       string         `json:"title_cyrillic"`
	DescriptionCyrillic string         `json:"description_cyrillic"`
	Videos              pq.StringArray `json:"videos"`
	Tags                pq.StringArray `json:"tags"`
}

// getArticleListByCategory is a handler to get article list by category
func getArticleListByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	// Get 'page' and 'limit' query parameters
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		// If 'page' parameter is not provided or is not a number, default to 1
		page = 1
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		// If 'limit' parameter is not provided or is not a number, default to 10
		limit = 10
	}

	// Calculate the start index for the slice of posts
	start := (page - 1) * limit

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags FROM articles WHERE category = $1 AND archived = false AND completed = true ORDER BY id DESC LIMIT $2 OFFSET $3", category, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, &article.Videos, &article.Tags)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, articles)
}

// getArticleListByRelated is a handler to get article list by related
func getArticleListByRelated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	related := vars["related"]

	// Get 'page' and 'limit' query parameters
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		// If 'page' parameter is not provided or is not a number, default to 1
		page = 1
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		// If 'limit' parameter is not provided or is not a number, default to 10
		limit = 10
	}

	// Calculate the start index for the slice of posts
	start := (page - 1) * limit

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags FROM articles WHERE related = $1 AND archived = false AND completed = true ORDER BY id DESC LIMIT $2 OFFSET $3", related, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, &article.Videos, &article.Tags)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, articles)
}

// getAllArticles is a handler to get all articles
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	// Get 'page' and 'limit' query parameters
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		// If 'page' parameter is not provided or is not a number, default to 1
		page = 1
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		// If 'limit' parameter is not provided or is not a number, default to 10
		limit = 10
	}

	// Calculate the start index for the slice of posts
	start := (page - 1) * limit

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Get the slice of posts
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags FROM articles WHERE archived = false AND completed = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, &article.Videos, &article.Tags)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, articles)
}

// ArticlePhoto is a struct to represent article photo
type ArticlePhoto struct {
	ID        int    `json:"id"`
	Article   int    `json:"article"`
	FileName  string `json:"file_name"`
	File      []byte `json:"file"`
	CreatedAt string `json:"created_at"`
}

// getArticlePhotos is a handler to get article photos
func getArticlePhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check if the article exists, and archived is false, completed is true
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id = $1 AND archived = false AND completed = true)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	if !exists {
		response.Res(w, "error", http.StatusNotFound, "article not found")
		return
	}

	// get article photos from the database where article id is equal to the id
	rows, err := database.Query("SELECT id, article, file_name FROM article_photos WHERE article = $1", id)
	if err != nil {
		// check if the error is no rows in result set using sql.ErrNoRows
		if err == sql.ErrNoRows {
			response.Res(w, "error", http.StatusNotFound, "article photos not found")
			return
		}
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articlePhotos []ArticlePhoto
	for rows.Next() {
		var articlePhoto ArticlePhoto
		err := rows.Scan(&articlePhoto.ID, &articlePhoto.Article, &articlePhoto.FileName)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articlePhotos = append(articlePhotos, articlePhoto)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, articlePhotos)
}

// getArticlePhoto is a handler to get article photo
func getArticlePhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	photoID := vars["photo_id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check if the article exists, and archived is false, completed is true
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id = $1 AND archived = false AND completed = true)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	if !exists {
		response.Res(w, "error", http.StatusNotFound, "article not found")
		return
	}

	// get article photo from the database where article id is equal to the id and photo id is equal to the photoID
	var articlePhoto ArticlePhoto
	err = database.QueryRow("SELECT id, article, file_name, file, created_at FROM article_photos WHERE article = $1 AND id = $2", id, photoID).Scan(&articlePhoto.ID, &articlePhoto.Article, &articlePhoto.FileName, &articlePhoto.File, &articlePhoto.CreatedAt)
	if err != nil {
		// check if the error is no rows in result set using sql.ErrNoRows
		if err == sql.ErrNoRows {
			response.Res(w, "error", http.StatusNotFound, "article photo not found")
			return
		}
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send the article photo as photo
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(articlePhoto.File)))
	w.Write(articlePhoto.File)
}

// getArticleCoverImage is a handler to get article cover image
func getArticleCoverImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check if the article exists, and archived is false, completed is true
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id = $1 AND archived = false AND completed = true)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	if !exists {
		response.Res(w, "error", http.StatusNotFound, "article not found")
		return
	}

	// get article cover image from the database where article id is equal to the id
	var coverImage []byte
	err = database.QueryRow("SELECT cover_image FROM articles WHERE id = $1", id).Scan(&coverImage)
	if err != nil {
		// check if the error is no rows in result set using sql.ErrNoRows
		if err == sql.ErrNoRows {
			response.Res(w, "error", http.StatusNotFound, "article cover image not found")
			return
		}
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send the article cover image as photo
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(coverImage)))
	w.Write(coverImage)
}
