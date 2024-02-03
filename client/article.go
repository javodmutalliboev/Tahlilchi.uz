package client

import (
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, photos, videos, cover_image, tags FROM articles WHERE category = $1 AND archived = false ORDER BY id DESC LIMIT $2 OFFSET $3", category, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, pq.Array(&article.Photos), pq.Array(&article.Videos), &article.CoverImage, pq.Array(&article.Tags))
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

type Article struct {
	ID                  int      `json:"id"`
	TitleLatin          string   `json:"title_latin"`
	DescriptionLatin    string   `json:"description_latin"`
	TitleCyrillic       string   `json:"title_cyrillic"`
	DescriptionCyrillic string   `json:"description_cyrillic"`
	Photos              [][]byte `json:"photos"`
	Videos              []string `json:"videos"`
	CoverImage          []byte   `json:"cover_image"`
	Tags                []string `json:"tags"`
}

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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, photos, videos, cover_image, tags FROM articles WHERE related = $1 AND archived = false ORDER BY id DESC LIMIT $2 OFFSET $3", related, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, pq.Array(&article.Photos), pq.Array(&article.Videos), &article.CoverImage, pq.Array(&article.Tags))
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

func getArticleCategory(w http.ResponseWriter, r *http.Request) {
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

type ArticleCategory struct {
	ID                  int    `json:"id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

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
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, photos, videos, cover_image, tags FROM articles WHERE archived = false ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, pq.Array(&article.Photos), pq.Array(&article.Videos), &article.CoverImage, pq.Array(&article.Tags))
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
