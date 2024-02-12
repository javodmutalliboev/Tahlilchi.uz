package admin

import (
	"net/http"
	"strconv"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
)

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
	// get page, limit query parameters from the request url
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
	vncListRes, err := vnc.GetVideoNewsCommentListResponse(true, id, page, limit)
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
	response.Res(w, "success", http.StatusOK, vncListRes)
}

// approveVideoNewsComment is a route handler function to approve a video news comment
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func approveVideoNewsComment(w http.ResponseWriter, r *http.Request) {
	// get the id of the video news comment from the request url
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

	comment_id_str := mux.Vars(r)["comment_id"]
	comment_id, err := strconv.Atoi(comment_id_str)
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
	// approve the video news comment
	err = vnc.ApproveVideoNewsComment(id, comment_id)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}

	// respond with the success message
	response.Res(w, "success", http.StatusOK, "video news comment approved successfully")
}
