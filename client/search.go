package client

import (
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
)

// searchArticle is a handler function for the /search/article route.
// It is used to search for articles.
// search columns: title_latin, description_latin, title_cyrillic, description_cyrillic, tags.
// tags is text[]. use @> to search for a tag in the tags column.
func searchArticle(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		log.Printf("%v: search query is empty", r.URL)
		response.Res(w, "error", http.StatusBadRequest, "search query is empty")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags FROM articles WHERE (title_latin ILIKE '%' || $1 || '%' OR description_latin ILIKE '%' || $1 || '%' OR title_cyrillic ILIKE '%' || $1 || '%' OR description_cyrillic ILIKE '%' || $1 || '%' OR tags @> ARRAY[$1]) AND archived = false AND completed = true", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.TitleCyrillic, &article.DescriptionLatin, &article.DescriptionCyrillic, &article.Videos, &article.Tags)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articles = append(articles, article)
	}

	response.Res(w, "success", http.StatusOK, articles)
}
