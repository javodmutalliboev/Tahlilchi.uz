package client

import (
	"encoding/json"
	"net/http"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
)

// addArticleComment is a route handler function to add an article comment
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func addArticleComment(w http.ResponseWriter, r *http.Request) {
	// create a new article comment
	ac := model.ArticleComment{} // Go file path: model/article_comment.go
	// decode the request body into the article comment
	err := json.NewDecoder(r.Body).Decode(&ac)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error using response package Res function
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}
	// add the article comment to the database
	err = ac.AddArticleComment()
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the article comment
	response.Res(w, "success", http.StatusCreated, "article comment added successfully")
}

// getArticleCommentList is a route handler function to get the article comment list
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func getArticleCommentList(w http.ResponseWriter, r *http.Request) {
	// get the id of the article from the request url
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
	// create a new article comment
	ac := model.ArticleComment{} // Go file path: model/article_comment.go
	// get the article comment list from the database
	acs, err := ac.GetArticleCommentList(false, id, page, limit)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the article comment list
	response.Res(w, "success", http.StatusOK, acs)
}
