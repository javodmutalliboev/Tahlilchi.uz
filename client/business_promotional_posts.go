package client

import (
	"fmt"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
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
