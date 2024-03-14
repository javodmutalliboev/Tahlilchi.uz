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
	ID                  int64          `json:"id"`
	TitleLatin          string         `json:"title_latin"`
	DescriptionLatin    string         `json:"description_latin"`
	TitleCyrillic       string         `json:"title_cyrillic"`
	DescriptionCyrillic string         `json:"description_cyrillic"`
	Video               string         `json:"video"`
	Tags                pq.StringArray `json:"tags"`
	CreatedAt           string         `json:"created_at"`
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
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE archived = false AND completed = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE category = $1 AND archived = false AND completed = true ORDER BY id DESC LIMIT $2 OFFSET $3", category, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE subcategory = $1 AND archived = false AND completed = true ORDER BY id DESC LIMIT $2 OFFSET $3", subcategory, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE region = $1 AND archived = false AND completed = true ORDER BY id DESC LIMIT $2 OFFSET $3", region, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE archived = false AND top = true AND completed = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE archived = false AND latest = true AND completed = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, created_at FROM news_posts WHERE related = $1 AND archived = false AND completed = true ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Video, &post.Tags, &post.CreatedAt)
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

func getNewsPostPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var photo []byte
	err = database.QueryRow("SELECT photo FROM news_posts WHERE id = $1 AND archived = false AND completed = true", id).Scan(&photo)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	contentType := http.DetectContentType(photo)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(photo)))
	w.Write(photo)
}

func getNewsPostAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var audio []byte
	err = database.QueryRow("SELECT audio FROM news_posts WHERE id = $1 AND archived = false AND completed = true", id).Scan(&audio)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	contentType := http.DetectContentType(audio)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(audio)))
	w.Write(audio)
}

func getNewsPostCoverImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var CoverImage []byte
	err = database.QueryRow("SELECT cover_image FROM news_posts WHERE id = $1 AND archived = false AND completed = true", id).Scan(&CoverImage)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	contentType := http.DetectContentType(CoverImage)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(CoverImage)))
	w.Write(CoverImage)
}

// getNewsRegionList is a route handler function to get the news region list
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func getNewsRegionList(w http.ResponseWriter, r *http.Request) {
	// database connection
	database, err := db.DB()
	// check if there is an error
	if err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// close the database connection
	defer database.Close()
	// get the news region list from the database: table news_regions
	rows, err := database.Query("SELECT * FROM news_regions ORDER BY id")
	// check if there is an error
	if err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// close the rows
	defer rows.Close()
	// create a slice of news regions
	var regions []Region
	// loop through the rows
	for rows.Next() {
		// create a new region
		var region Region
		// scan the row into the region
		err := rows.Scan(&region.ID, &region.NameLatin, &region.DescriptionLatin, &region.NameCyrillic, &region.DescriptionCyrillic)
		// check if there is an error
		if err != nil {
			// log the error
			log.Printf("%v: error: %v", r.URL, err)
			// respond with the error
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			// return
			return
		}
		// append the region to the slice of regions
		regions = append(regions, region)
	}
	// respond with the success and the slice of regions
	response.Res(w, "success", http.StatusOK, regions)
}

// Region is a struct to represent a news region
type Region struct {
	ID                  int    `json:"id"`
	NameLatin           string `json:"name_latin"`
	DescriptionLatin    string `json:"description_latin"`
	NameCyrillic        string `json:"name_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}
