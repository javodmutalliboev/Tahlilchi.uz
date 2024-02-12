package admin

import (
	"net/http"
	"strconv"

	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
)

// getENewspaperCommentList is a route handler function to get the e-newspaper comment list response
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func getENewspaperCommentList(w http.ResponseWriter, r *http.Request) {
	// get the id of the e-newspaper from the request url
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
	// create a new e-newspaper comment
	enc := model.ENewspaperComment{}
	// get the e-newspaper comment list response from the database
	encListRes, err := enc.GetENewspaperCommentListResponse(true, id, page, limit)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the e-newspaper comment list response
	response.Res(w, "success", http.StatusOK, encListRes)
}

// approveENewspaperComment is a route handler function to approve an e-newspaper comment
// It takes a http.ResponseWriter and a http.Request as its parameters
// It returns nothing
func approveENewspaperComment(w http.ResponseWriter, r *http.Request) {
	// get the id of the e-newspaper comment from the request url
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
		// return
		return
	}

	// create a new e-newspaper comment
	enc := model.ENewspaperComment{}
	// approve the e-newspaper comment
	err = enc.ApproveENewspaperComment(id, comment_id)
	// check if there is an error
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// respond with the error
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		// return
		return
	}
	// respond with the e-newspaper comment
	response.Res(w, "success", http.StatusOK, "e-newspaper comment approved successfully")
}
