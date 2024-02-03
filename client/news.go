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

type NewsPost struct {
	ID                  int64    `json:"id"`
	TitleLatin          string   `json:"title_latin"`
	DescriptionLatin    string   `json:"description_latin"`
	TitleCyrillic       string   `json:"title_cyrillic"`
	DescriptionCyrillic string   `json:"description_cyrillic"`
	Video               string   `json:"video"`
	Tags                []string `json:"tags"`
}

func getAllNewsPosts(w http.ResponseWriter, r *http.Request) {
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
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE archived = false ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getNewsByCategory(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE category = $1 AND archived = false ORDER BY id DESC LIMIT $2 OFFSET $3", category, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getNewsBySubCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subcategory := vars["subcategory"]

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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE subcategory = $1 AND archived = false ORDER BY id DESC LIMIT $2 OFFSET $3", subcategory, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getNewsByRegion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := vars["region"]

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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE region = $1 AND archived = false ORDER BY id DESC LIMIT $2 OFFSET $3", region, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getTopNews(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE archived = false AND top = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getLatestNews(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE archived = false AND latest = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getRelatedNewsPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags FROM news_posts WHERE related = $1 AND archived = false ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, pq.Array(&post.Tags))
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, posts)
}

func getCategoryList(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT * FROM news_category ORDER BY id")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var categories []CategoryForGet
	for rows.Next() {
		var c CategoryForGet
		err := rows.Scan(&c.ID, &c.TitleLatin, &c.DescriptionLatin, &c.TitleCyrillic, &c.DescriptionCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		subRows, err := database.Query("SELECT * FROM news_subcategory WHERE category_id = $1", c.ID)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer subRows.Close()

		for subRows.Next() {
			var s Subcategory
			err := subRows.Scan(&s.ID, &s.CategoryID, &s.TitleLatin, &s.DescriptionLatin, &s.TitleCyrillic, &s.DescriptionCyrillic)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			c.Subcategories = append(c.Subcategories, s)
		}

		categories = append(categories, c)
	}

	response.Res(w, "success", http.StatusOK, categories)
}

type CategoryForGet struct {
	ID                  int           `json:"id"`
	TitleLatin          string        `json:"title_latin"`
	DescriptionLatin    string        `json:"description_latin"`
	TitleCyrillic       string        `json:"title_cyrillic"`
	DescriptionCyrillic string        `json:"description_cyrillic"`
	Subcategories       []Subcategory `json:"subcategories"`
}

type Subcategory struct {
	ID                  int    `json:"id"`
	CategoryID          int    `json:"category_id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}
