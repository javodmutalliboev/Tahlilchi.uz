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
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		// return
		return
	}
	// respond with the article comment
	response.Res(w, "success", http.StatusCreated, "article comment added successfully")
}
