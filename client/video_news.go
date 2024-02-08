package client

import (
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
)

// VideoNewsListResponse is a struct to represent a video news list response.
type VideoNewsListResponse struct {
	Previous      bool        `json:"previous"`
	VideoNewsList []VideoNews `json:"video_news_list"`
	Next          bool        `json:"next"`
}

// VideoNews is a struct to represent a video news.
type VideoNews struct {
	ID           int    `json:"id"`
	Video        string `json:"video"`
	TextLatin    string `json:"text_latin"`
	TextCyrillic string `json:"text_cyrillic"`
	CreatedAt    string `json:"created_at"`
}

// getVideoNewsList is a handler function for the /video-news/list route.
// It is used to get a list of video news.
func getVideoNewsList(w http.ResponseWriter, r *http.Request) {
	// page, limit query parameters
	// page
	pageStr := r.URL.Query().Get("page")
	var page int
	var err error
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			response.Res(w, "error", http.StatusBadRequest, "Invalid page parameter")
			return
		}
	}
	// limit
	limitStr := r.URL.Query().Get("limit")
	var limit int
	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			response.Res(w, "error", http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	// offset
	offset := (page - 1) * limit

	// Get video news list
	// connect to the database
	database, err := db.DB()
	if err != nil {
		// log error
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "Internal server error")
		return
	}
	// defer closing the database
	defer database.Close()

	// get video news list from the database with limit and offset parameters querying the database
	rows, err := database.Query("SELECT id, video, text_latin, text_cyrillic, created_at FROM video_news ORDER BY id DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		// log error
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "Internal server error")
		return
	}
	// defer closing the rows
	defer rows.Close()

	// create a slice of video news
	var videoNewsList []VideoNews
	// iterate over the rows
	for rows.Next() {
		// create a video news
		var videoNews VideoNews
		// scan the row
		err := rows.Scan(&videoNews.ID, &videoNews.Video, &videoNews.TextLatin, &videoNews.TextCyrillic, &videoNews.CreatedAt)
		if err != nil {
			// log error
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "Internal server error")
			return
		}
		// append the video news to the video news list
		videoNewsList = append(videoNewsList, videoNews)
	}
	// check for errors
	err = rows.Err()
	if err != nil {
		// log error
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "Internal server error")
		return
	}

	// get the total number of video news
	var total int
	err = database.QueryRow("SELECT COUNT(*) FROM video_news").Scan(&total)
	if err != nil {
		// log error
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "Internal server error")
		return
	}

	// create a video news list response
	videoNewsListResponse := VideoNewsListResponse{
		Previous:      page > 1,
		VideoNewsList: videoNewsList,
		Next:          total > page*limit,
	}
	// send the response
	response.Res(w, "success", http.StatusOK, videoNewsListResponse)
}
