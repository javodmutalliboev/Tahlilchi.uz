package admin

import (
	"net/http"
	"strconv"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
)

// getNewsPostCommentList is a route handler function to get a news post comment list response
func getNewsPostCommentList(w http.ResponseWriter, r *http.Request) {
	// get the news post id from the request url
	id, err := toolkit.GetID(r)

	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send a response with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// get the page and limit from the request url
	page, limit, err := toolkit.GetPageLimit(r)

	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send a response with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// create a new news post comment
	npc := model.NewsPostComment{} // Go file path: model/news_post_comment.go
	// get the news post comment list response from the database
	npcListResponse, err := npc.GetNewsPostCommentListResponse(true, id, page, limit)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send a response with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send a response with the news post comment list response
	response.Res(w, "success", http.StatusOK, npcListResponse)
}

// approveNewsPostComment is a route handler function to approve a news post comment
func approveNewsPostComment(w http.ResponseWriter, r *http.Request) {
	id, err := toolkit.GetID(r)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// respond with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}

	comment_id_str := mux.Vars(r)["comment_id"]
	comment_id, err := strconv.Atoi(comment_id_str)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// respond with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	npc := model.NewsPostComment{}
	err = npc.ApproveNewsPostComment(id, comment_id)
	if err != nil {
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "news post comment approved successfully")
}
