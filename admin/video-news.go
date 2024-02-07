package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
)

// addVideoNews is a handler function for the admin router to handle the request to add a video news
func addVideoNews(w http.ResponseWriter, r *http.Request) {
	// get the video news from the request body
	var videoNews model.VideoNews
	if err := json.NewDecoder(r.Body).Decode(&videoNews); err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error decoding the request body, return a bad request response
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}
	// add the video news to the database
	if err := videoNews.AddVideoNews(); err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error adding the video news to the database, return an internal server error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// return a success response using response package
	response.Res(w, "success", http.StatusCreated, "Video news added successfully")
}

// updateVideoNews is a handler function for the admin router to handle the request to update a video news
func updateVideoNews(w http.ResponseWriter, r *http.Request) {
	// get id from the request url
	id := r.URL.Query().Get("id")
	// if the id is empty, return a bad request response
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	// get the video news from the request body
	var videoNews model.VideoNews
	if err := json.NewDecoder(r.Body).Decode(&videoNews); err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error decoding the request body, return a bad request response
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}
	// update the video news in the database
	if err := videoNews.UpdateVideoNews(id); err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error updating the video news in the database, return an internal server error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// return a success response using response package
	response.Res(w, "success", http.StatusOK, "Video news updated successfully")
}

// deleteVideoNews is a handler function for the admin router to handle the request to delete a video news
func deleteVideoNews(w http.ResponseWriter, r *http.Request) {
	// get id from the request url
	id := r.URL.Query().Get("id")
	// if the id is empty, return a bad request response
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	// convert the id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error converting the id to int, return a bad request response
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}
	// create a new video news
	videoNews := model.VideoNews{ID: idInt}
	// delete the video news from the database
	if err := videoNews.DeleteVideoNews(); err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error deleting the video news from the database, return an internal server error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// return a success response using response package
	response.Res(w, "success", http.StatusOK, "Video news deleted successfully")
}

// getVideoNewsList is a handler function for the admin router to handle the request to get a list of video news
func getVideoNewsList(w http.ResponseWriter, r *http.Request) {
	// get page and limit from query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	// create int page and limit
	var page, limit int
	// if page is not empty, convert it to int
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	} else {
		page = 1
	}
	// if limit is not empty, convert it to int
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	} else {
		limit = 10
	}
	// calculate the offset
	offset := (page - 1) * limit

	// get the list of video news from the database
	videoNewsListResponse, err := model.GetVideoNewsList(limit, offset)
	if err != nil {
		// log the error
		log.Printf("%v: error: %v", r.URL, err)
		// if there is an error getting the list of video news from the database, return an internal server error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return the list of video news using response package
	response.Res(w, "success", http.StatusOK, *videoNewsListResponse)
}
