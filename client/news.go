package client

import (
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/lib/pq"
)

type NewsPost struct {
	ID                  int64    `json:"id"`
	TitleLatin          string   `json:"title_latin"`
	DescriptionLatin    string   `json:"description_latin"`
	TitleCyrillic       string   `json:"title_cyrillic"`
	DescriptionCyrillic string   `json:"description_cyrillic"`
	Photo               []byte   `json:"photo"`
	Video               string   `json:"video"`
	Audio               []byte   `json:"audio"`
	CoverImage          []byte   `json:"cover_image"`
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
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, photo, video, audio, cover_image, tags FROM news_posts WHERE archived = false ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var posts []NewsPost
	for rows.Next() {
		var post NewsPost
		err := rows.Scan(&post.ID, &post.TitleLatin, &post.DescriptionLatin, &post.TitleCyrillic, &post.DescriptionCyrillic, &post.Photo, &post.Video, &post.Audio, &post.CoverImage, pq.Array(&post.Tags))
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
