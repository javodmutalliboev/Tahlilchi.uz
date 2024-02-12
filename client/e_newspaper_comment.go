package client

import (
	"encoding/json"
	"net/http"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
)

// addENewspaperComment is a route handler function to add an e-newspaper comment
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func addENewspaperComment(w http.ResponseWriter, r *http.Request) {
	// get the id of the e-newspaper from the request url
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
	// create a new e-newspaper comment
	enc := model.ENewspaperComment{} // Go file path: model/e_newspaper_comment.go
	// decode the request body into the e-newspaper comment
	err = json.NewDecoder(r.Body).Decode(&enc)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error using response package Res function
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		// return
		return
	}
	// add the e-newspaper comment to the database
	err = enc.AddENewspaperComment(id)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the e-newspaper comment
	response.Res(w, "success", http.StatusCreated, "e-newspaper comment added successfully")
}

// getENewspaperCommentList is a route handler function to get the e-newspaper comment list
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func getENewspaperCommentList(w http.ResponseWriter, r *http.Request) {
	// get the id of the e-newspaper from the request url
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
	// create a new e-newspaper comment
	enc := model.ENewspaperComment{} // Go file path: model/e_newspaper_comment.go
	// get the e-newspaper comment list from the database
	encListRes, err := enc.GetENewspaperCommentListResponse(false, id, page, limit)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err) // Go file path: toolkit/log.go
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the e-newspaper comment list
	response.Res(w, "success", http.StatusOK, encListRes)
}
