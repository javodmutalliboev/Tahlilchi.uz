package admin

import (
	"net/http"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
)

// getArticleCommentList is a route handler function to get the article comment list response
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func getArticleCommentList(w http.ResponseWriter, r *http.Request) {
	// get the id of the article from the request url
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
	// get page, limit query parameters from the request url
	page, limit, err := toolkit.GetPageLimit(r)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// respond with the error
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}
	// create a new article comment
	ac := model.ArticleComment{}
	// get the article comment list from the database
	acs, err := ac.GetArticleCommentList(true, id, page, limit)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the article comment list
	response.Res(w, "success", http.StatusOK, acs)
}
