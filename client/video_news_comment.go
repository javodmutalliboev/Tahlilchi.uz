package client

import (
	"encoding/json"
	"net/http"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
)

// addVideoNewsComment is a route handler function to add a video news comment
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func addVideoNewsComment(w http.ResponseWriter, r *http.Request) {
	// get the id of the video news from the request url
	id, err := toolkit.GetID(r) // Go file path: toolkit/toolkit.go
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}
	// create a new video news comment
	vnc := model.VideoNewsComment{} // Go file path: model/video_news_comment.go
	// decode the request body into the video news comment
	err = json.NewDecoder(r.Body).Decode(&vnc)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error using response package Res function
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}
	// add the video news comment to the database
	err = vnc.AddVideoNewsComment(id)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the video news comment
	response.Res(w, "success", http.StatusCreated, "video news comment added successfully")
}

// getVideoNewsCommentList is a route handler function to get the video news comment list
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func getVideoNewsCommentList(w http.ResponseWriter, r *http.Request) {
	// get the id of the video news from the request url
	id, err := toolkit.GetID(r) // Go file path: toolkit/toolkit.go
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}

	// get the page, limit query parameters from the request url
	page, limit, err := toolkit.GetPageLimit(r) // Go file path: toolkit/toolkit.go

	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}

	// create a new video news comment
	vnc := model.VideoNewsComment{} // Go file path: model/video_news_comment.go
	// get the video news comment list response from the database
	vncListResponse, err := vnc.GetVideoNewsCommentListResponse(false, id, page, limit)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the video news comment list response
	response.Res(w, "success", http.StatusOK, vncListResponse)
}
