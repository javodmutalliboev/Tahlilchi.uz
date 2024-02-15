package admin

import (
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
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

	var businessPromotionalPosts []model.BusinessPromotionalPost
	for rows.Next() {
		var businessPromotionalPost model.BusinessPromotionalPost
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

// searchENewspaper is the handler for the /admin/search/e-newspaper endpoint.
// It searches the e_newspapers table for the given query.
// search columns: title_latin, title_cyrillic.
func searchENewspaper(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic, created_at, updated_at, archived, completed FROM e_newspapers WHERE title_latin ILIKE $1 OR title_cyrillic ILIKE $1", "%"+search+"%")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var eNewspapers []ENewspaper
	for rows.Next() {
		var eNewspaper ENewspaper
		err := rows.Scan(&eNewspaper.ID, &eNewspaper.TitleLatin, &eNewspaper.TitleCyrillic, &eNewspaper.CreatedAt, &eNewspaper.UpdatedAt, &eNewspaper.Archived, &eNewspaper.Completed)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		eNewspapers = append(eNewspapers, eNewspaper)
	}

	response.Res(w, "success", http.StatusOK, eNewspapers)
}

// searchNews is the handler for the /admin/search/news endpoint.
// It searches the news_posts table for the given query.
// search columns: title_latin, description_latin, title_cyrillic, description_cyrillic, tags.
// tags is text[].
func searchNews(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, archived, created_at, updated_at, category, subcategory, region, top, latest, related, completed FROM news_posts WHERE title_latin ILIKE $1 OR description_latin ILIKE $1 OR title_cyrillic ILIKE $1 OR description_cyrillic ILIKE $1 OR tags @> ARRAY[$2]", "%"+search+"%", search)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var newsPosts []NewsPost
	for rows.Next() {
		var newsPost NewsPost
		var tags pq.StringArray
		err := rows.Scan(&newsPost.ID, &newsPost.TitleLatin, &newsPost.DescriptionLatin, &newsPost.TitleCyrillic, &newsPost.DescriptionCyrillic, &newsPost.Video, &tags, &newsPost.Archived, &newsPost.CreatedAt, &newsPost.UpdatedAt, &newsPost.Category, &newsPost.Subcategory, &newsPost.Region, &newsPost.Top, &newsPost.Latest, &newsPost.Related, &newsPost.Completed)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		newsPost.Tags = tags
		newsPosts = append(newsPosts, newsPost)
	}

	response.Res(w, "success", http.StatusOK, newsPosts)
}

// searchPhotoGallery is the handler for the /admin/search/photo-gallery endpoint.
// It searches the photo_gallery table for the given query.
// search columns: title_latin, title_cyrillic.
func searchPhotoGallery(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic, created_at, updated_at FROM photo_gallery WHERE title_latin ILIKE $1 OR title_cyrillic ILIKE $1", "%"+search+"%")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleries []PhotoGallery
	for rows.Next() {
		var photoGallery PhotoGallery
		err := rows.Scan(&photoGallery.ID, &photoGallery.TitleLatin, &photoGallery.TitleCyrillic, &photoGallery.CreatedAt, &photoGallery.UpdatedAt)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photoGalleries = append(photoGalleries, photoGallery)
	}

	response.Res(w, "success", http.StatusOK, photoGalleries)
}

// searchPhotoGalleryPhotos is the handler for the /admin/search/photo-gallery/photos endpoint.
// It searches the photo_gallery_photos table for the given query.
// search columns: file_name.
func searchPhotoGalleryPhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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

	rows, err := database.Query("SELECT id, photo_gallery, file_name, created_at from photo_gallery_photos WHERE file_name ILIKE $1 AND photo_gallery = $2", "%"+search+"%", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleryPhotos []PhotoGalleryPhoto
	for rows.Next() {
		var photoGalleryPhoto PhotoGalleryPhoto
		err := rows.Scan(&photoGalleryPhoto.ID, &photoGalleryPhoto.PhotoGallery, &photoGalleryPhoto.FileName, &photoGalleryPhoto.CreatedAt)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photoGalleryPhotos = append(photoGalleryPhotos, photoGalleryPhoto)
	}

	response.Res(w, "success", http.StatusOK, photoGalleryPhotos)
}
