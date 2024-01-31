package client

import (
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/lib/pq"
)

func getBusinessPromotionalPosts(w http.ResponseWriter, r *http.Request) {
	// Parse the page number from the query parameters
	pageStr, ok := r.URL.Query()["page"]
	if !ok || len(pageStr[0]) < 1 {
		log.Printf("%v: Url Param 'page' is missing. Setting default value to 1.", r.URL)
		pageStr = []string{"1"}
	}
	page, _ := strconv.Atoi(pageStr[0])

	// Parse the limit from the query parameters
	limitStr, ok := r.URL.Query()["limit"]
	if !ok || len(limitStr[0]) < 1 {
		log.Printf("%v: Url Param 'limit' is missing. Setting default value to 10.", r.URL)
		limitStr = []string{"10"}
	}
	limit, _ := strconv.Atoi(limitStr[0])

	// Calculate the starting index
	var start = (page - 1) * limit

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT * FROM business_promotional_posts WHERE archived = false ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var bpPosts []bpPost
	for rows.Next() {
		var bpp bpPost
		err := rows.Scan(&bpp.ID, &bpp.TitleLatin, &bpp.DescriptionLatin, &bpp.TitleCyrillic, &bpp.DescriptionCyrillic, pq.Array(&bpp.Videos), &bpp.CoverImage)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		bpPosts = append(bpPosts, bpp)
	}

	response.Res(w, "success", http.StatusOK, bpPosts)
}

type bpPost struct {
	ID                  int      `json:"id"`
	TitleLatin          string   `json:"title_latin"`
	DescriptionLatin    string   `json:"description_latin"`
	TitleCyrillic       string   `json:"title_cyrillic"`
	DescriptionCyrillic string   `json:"description_cyrillic"`
	Videos              []string `json:"videos"`
	CoverImage          []byte   `json:"cover_image"`
}
