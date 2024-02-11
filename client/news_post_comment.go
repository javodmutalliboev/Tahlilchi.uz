package client

import (
	"encoding/json"
	"net/http"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
)

// addNewsPostComment is a method to add a news post comment to the database
func addNewsPostComment(w http.ResponseWriter, r *http.Request) {
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

	// create a new news post comment
	npc := model.NewsPostComment{} // Go file path: model/news_post_comment.go
	// decode the request body to the news post comment
	err = json.NewDecoder(r.Body).Decode(&npc)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send a response with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// add the news post comment to the database
	err = npc.AddNewsPostComment(id)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send a response with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send a response with the news post comment
	response.Res(w, "success", http.StatusCreated, "news post comment added successfully")
}

// getNewsPostCommentList is a method to get the news post comment list from the database
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
	npcListResponse, err := npc.GetNewsPostCommentListResponse(false, id, page, limit)
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
