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

// searchBusinessPromotional is the handler for the /admin/search/business-promotional endpoint.
// It searches the business_promotional_posts table for the given query.
// search columns: title_latin, description_latin, title_cyrillic, description_cyrillic, videos, partner.
// videos is text[].
func searchBusinessPromotional(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, expiration, created_at, updated_at, archived, partner, completed FROM business_promotional_posts WHERE title_latin ILIKE $1 OR description_latin ILIKE $1 OR title_cyrillic ILIKE $1 OR description_cyrillic ILIKE $1 OR videos @> ARRAY[$2] OR partner ILIKE $1", "%"+search+"%", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var businessPromotionalPosts []BusinessPromotionalPost
	for rows.Next() {
		var businessPromotionalPost BusinessPromotionalPost
		var videos pq.StringArray
		err := rows.Scan(&businessPromotionalPost.ID, &businessPromotionalPost.TitleLatin, &businessPromotionalPost.DescriptionLatin, &businessPromotionalPost.TitleCyrillic, &businessPromotionalPost.DescriptionCyrillic, &videos, &businessPromotionalPost.Expiration, &businessPromotionalPost.CreatedAt, &businessPromotionalPost.UpdatedAt, &businessPromotionalPost.Archived, &businessPromotionalPost.Partner, &businessPromotionalPost.Completed)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		businessPromotionalPost.Videos = videos
		businessPromotionalPosts = append(businessPromotionalPosts, businessPromotionalPost)
	}

	response.Res(w, "success", http.StatusOK, businessPromotionalPosts)
}

// BusinessPromotionalPost is the model for the business_promotional_posts table.
type BusinessPromotionalPost struct {
	ID                  int            `json:"id"`
	TitleLatin          string         `json:"title_latin"`
	DescriptionLatin    string         `json:"description_latin"`
	TitleCyrillic       string         `json:"title_cyrillic"`
	DescriptionCyrillic string         `json:"description_cyrillic"`
	Videos              pq.StringArray `json:"videos"`
	Expiration          string         `json:"expiration"`
	CreatedAt           string         `json:"created_at"`
	UpdatedAt           string         `json:"updated_at"`
	Archived            bool           `json:"archived"`
	Partner             string         `json:"partner"`
	Completed           bool           `json:"completed"`
}
