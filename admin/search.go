package admin

import (
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/lib/pq"
)

// searchAppeal is the handler for the /admin/search/appeal endpoint.
// It searches the appeals table for the given query.
// search columns: name, surname, phone_number, message
func searchAppeal(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		response.Res(w, "error", http.StatusBadRequest, "search query is missing")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, name, surname, phone_number, message, created_at FROM appeals WHERE name ILIKE $1 OR surname ILIKE $1 OR phone_number ILIKE $1 OR message ILIKE $1", "%"+search+"%")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var appeals []Appeal
	for rows.Next() {
		var appeal Appeal
		err := rows.Scan(&appeal.ID, &appeal.Name, &appeal.Surname, &appeal.PhoneNumber, &appeal.Message, &appeal.CreatedAt)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		appeals = append(appeals, appeal)
	}

	response.Res(w, "success", http.StatusOK, appeals)
}

// searchArticle is the handler for the /admin/search/article endpoint.
// It searches the articles table for the given query.
// search columns: title_latin, description_latin, title_cyrillic, description_cyrillic, tags.
// tags is text[].
func searchArticle(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		response.Res(w, "error", http.StatusBadRequest, "search query is missing")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags, archived, created_at, updated_at, category, related, completed FROM articles WHERE title_latin ILIKE $1 OR description_latin ILIKE $1 OR title_cyrillic ILIKE $1 OR description_cyrillic ILIKE $1 OR tags @> ARRAY[$2]", "%"+search+"%", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		var tags pq.StringArray
		var videos pq.StringArray
		err := rows.Scan(&article.ID, &article.TitleLatin, &article.DescriptionLatin, &article.TitleCyrillic, &article.DescriptionCyrillic, &videos, &tags, &article.Archived, &article.CreatedAt, &article.UpdatedAt, &article.Category, &article.Related, &article.Completed)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		article.Videos = videos
		article.Tags = tags
		articles = append(articles, article)
	}

	response.Res(w, "success", http.StatusOK, articles)
}
