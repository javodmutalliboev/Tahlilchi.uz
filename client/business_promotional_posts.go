package client

import (
	"fmt"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
)

func getBusinessPromotionalPosts(w http.ResponseWriter, r *http.Request) {
	// get page, limit
	page, limit, err := toolkit.GetPageLimit(r)
	if err != nil {
		err := fmt.Errorf("error getting page and limit: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// open a database connection
	database, err := db.DB()
	if err != nil {
		err := fmt.Errorf("error opening a database connection: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}
	defer database.Close()

	// get business promotional list response
	var bppListResponse model.BusinessPromotionalPostListResponse

	// get business promotional posts from the database: select id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, updated_at where archived is false and completed is true order by id
	// perform a database query: table is business_promotional_posts
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, updated_at FROM business_promotional_posts WHERE archived = false AND completed = true ORDER BY id LIMIT $1 OFFSET $2", limit, (page-1)*limit)
	if err != nil {
		err := fmt.Errorf("error querying the database: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	// get business promotional posts
	var bppList []model.BusinessPromotionalPost
	for rows.Next() {
		var bpp model.BusinessPromotionalPost
		err := rows.Scan(&bpp.ID, &bpp.TitleLatin, &bpp.DescriptionLatin, &bpp.TitleCyrillic, &bpp.DescriptionCyrillic, &bpp.Videos, &bpp.UpdatedAt)
		if err != nil {
			err := fmt.Errorf("error scanning the database: %v", err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusInternalServerError, err.Error())
			return
		}
		bppList = append(bppList, bpp)
	}

	if page > 1 {
		bppListResponse.Previous = true
	}

	var count int
	// count by id
	err = database.QueryRow("SELECT COUNT(id) FROM business_promotional_posts WHERE archived = false AND completed = true").Scan(&count)
	if err != nil {
		err := fmt.Errorf("error querying the database: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}

	if count > page*limit {
		bppListResponse.Next = true
	}

	bppListResponse.BPPList = bppList
	response.Res(w, "success", http.StatusOK, bppListResponse)
}

func getBusinessPromotionalPostPhotoList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// open a database connection
	database, err := db.DB()
	if err != nil {
		err := fmt.Errorf("error opening a database connection: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}
	defer database.Close()

	// first check if the business promotional post where id is $1, archived is false and completed is true exists
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM business_promotional_posts WHERE id = $1 AND archived = false AND completed = true)", id).Scan(&exists)
	if err != nil {
		err := fmt.Errorf("error querying the database: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}
	if !exists {
		err := fmt.Errorf("business promotional post where id is %s, archived is false and completed is true does not exist", id)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusNotFound, "business promotional post does not exist")
		return
	}

	// get business promotional post photo list
	var bppPhotoList []model.BusinessPromotionalPostPhoto

	// get business promotional post photo list from the database: select id, file_name, created_at from bpp_photos where bpp = $1
	// perform a database query: table is bpp_photos
	rows, err := database.Query("SELECT id, file_name, created_at FROM bpp_photos WHERE bpp = $1 ORDER BY id", id)
	if err != nil {
		err := fmt.Errorf("error querying the database: %v", err)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	// get business promotional post photo list
	for rows.Next() {
		var bppPhoto model.BusinessPromotionalPostPhoto
		err := rows.Scan(&bppPhoto.ID, &bppPhoto.FileName, &bppPhoto.CreatedAt)
		if err != nil {
			err := fmt.Errorf("error scanning the database: %v", err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusInternalServerError, err.Error())
			return
		}
		bppPhotoList = append(bppPhotoList, bppPhoto)
	}

	response.Res(w, "success", http.StatusOK, bppPhotoList)
}
