package toolkit

import (
	"net/http"
	"strconv"
	"strings"
)

func SliceToString(data []string) string {
	return "{" + strings.Join(data, ",") + "}"
}

func GetID(r *http.Request) (int, error) {
	id := r.URL.Query().Get("id")
	// convert the id to int
	idInt, err := strconv.Atoi(id)
	// check if there is an error
	if err != nil {
		// return the error
		return 0, err
	}
	// return the id
	return idInt, nil
}

func GetPageLimit(r *http.Request) (int, int, error) {
	// get the page query parameter from the request url
	page := r.URL.Query().Get("page")
	// if page is empty set it to 1
	if page == "" {
		// log page is empty
		LogInfo(r, "query parameter page is empty, setting it to 1")
		page = "1"
	}
	// convert the page to int
	pageInt, err := strconv.Atoi(page)
	// check if there is an error
	if err != nil {
		// return the error
		return 0, 0, err
	}
	// get the limit query parameter from the request url
	limit := r.URL.Query().Get("limit")
	// if limit is empty set it to 10
	if limit == "" {
		// log limit is empty
		LogInfo(r, "query parameter limit is empty, setting it to 10")
		limit = "10"
	}
	// convert the limit to int
	limitInt, err := strconv.Atoi(limit)
	// check if there is an error
	if err != nil {
		// return the error
		return 0, 0, err
	}
	// return page and limit
	return pageInt, limitInt, nil
}
