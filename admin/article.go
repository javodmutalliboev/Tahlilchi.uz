package admin

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func addArticleCategory(w http.ResponseWriter, r *http.Request) {
	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	var c Category
	c.TitleLatin = r.FormValue("title_latin")
	c.DescriptionLatin = r.FormValue("description_latin")
	c.TitleCyrillic = r.FormValue("title_cyrillic")
	c.DescriptionCyrillic = r.FormValue("description_cyrillic")

	if c.TitleLatin == "" {
		response.Res(w, "error", http.StatusBadRequest, "Title latin is required")
		return
	}

	if c.TitleCyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "Title cyrillic is required")
		return
	}

	_, err = db.Exec("INSERT INTO article_category(title_latin, description_latin, title_cyrillic, description_cyrillic) VALUES($1, $2, $3, $4)", c.TitleLatin, c.DescriptionLatin, c.TitleCyrillic, c.DescriptionCyrillic)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Category Added")
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	err := r.ParseMultipartForm(30 << 20) // 30 MB
	if err != nil {
		log.Printf("%v: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	titleLatin := r.FormValue("title_latin")
	titleCyrillic := r.FormValue("title_cyrillic")

	if titleLatin == "" || titleCyrillic == "" {
		log.Printf("%v: title_latin: %v; title_cyrillic: %v", r.URL, titleLatin, titleCyrillic)
		response.Res(w, "error", http.StatusBadRequest, "title_latin and title_cyrillic are required fields")
		return
	}

	descriptionLatin := r.FormValue("description_latin")
	descriptionCyrillic := r.FormValue("description_cyrillic")

	// Get the photos files
	photos := r.MultipartForm.File["photos"]
	if len(photos) == 0 {
		photos = nil
	}

	var photosForDb bytes.Buffer

	for _, fh := range photos {
		if fh.Size > 2<<20 {
			log.Printf("%v: photo size exceeds 2MB limit: %v", r.URL, fh.Size)
			response.Res(w, "error", http.StatusBadRequest, "photo size exceeds 2MB limit")
			return
		} else {
			file, _ := fh.Open()
			io.Copy(&photosForDb, file)
			file.Close()
		}
	}

	// Get the videos
	videos, ok := r.Form["videos"]
	if !ok {
		videos = []string{}
	}
	videosString := toolkit.SliceToString(videos)

	// Get the cover_image file
	coverImage, coverImageHeader, err := r.FormFile("cover_image")
	if err != nil {
		if err == http.ErrMissingFile {
			coverImage = nil
		} else {
			log.Printf("%v: cover_image error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, "cover_image error")
			return
		}
	}

	var coverImageForDB []byte = nil

	if coverImage != nil {
		if coverImageHeader.Size > 1<<20 {
			log.Printf("%v: cover_image size exceeds 1MB limit: %v", r.URL, coverImageHeader.Size)
			response.Res(w, "error", http.StatusBadRequest, "cover_image size exceeds 1MB limit")
			return
		}
		coverImageForDB, _ = io.ReadAll(coverImage)
		coverImage.Close()
	}

	tags, ok := r.Form["tags"]
	if !ok {
		// If tags don't exist, use an empty array
		tags = []string{}
	}

	// Convert tags to PostgreSQL array format
	tagsString := "{" + strings.Join(tags, ",") + "}"

	category := r.FormValue("category")
	var categoryInt64 sql.NullInt64
	if category != "" {
		categoryInt64.Int64, err = strconv.ParseInt(category, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, "invalid category value")
			return
		}
		categoryInt64.Valid = true
	}

	related := r.FormValue("related")
	var relatedInt64 sql.NullInt64
	if related != "" {
		relatedInt64.Int64, err = strconv.ParseInt(related, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, "invalid related value")
			return
		}
		relatedInt64.Valid = true
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO articles (title_latin, description_latin, title_cyrillic, description_cyrillic, photos, videos, cover_image, tags, category, related) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		titleLatin, descriptionLatin, titleCyrillic, descriptionCyrillic, pq.Array([][]byte{photosForDb.Bytes()}), videosString, coverImageForDB, tagsString, categoryInt64, relatedInt64)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Article Added")
}

func articleExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func articleIsArchived(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT archived FROM articles WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	var archived bool
	err = stmt.QueryRow(id).Scan(&archived)
	if err != nil {
		return nil, err
	}

	return &archived, nil
}

func editArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: edit article articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: edit article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: edit article articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: edit article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit archived article")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(15 << 20)
	if err != nil {
		log.Printf("%v: edit article: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	title_latin := r.FormValue("title_latin")
	if title_latin != "" {
		sqlStatement := `
			UPDATE articles
			SET title_latin = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, title_latin, id)
		if err != nil {
			log.Printf("%v: writing title_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	description_latin := r.FormValue("description_latin")
	if description_latin != "" {
		sqlStatement := `
			UPDATE articles
			SET description_latin = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, description_latin, id)
		if err != nil {
			log.Printf("%v: writing description_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	title_cyrillic := r.FormValue("title_cyrillic")
	if title_cyrillic != "" {
		sqlStatement := `
			UPDATE articles
			SET title_cyrillic = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, title_cyrillic, id)
		if err != nil {
			log.Printf("%v: writing title_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	description_cyrillic := r.FormValue("description_cyrillic")
	if description_cyrillic != "" {
		sqlStatement := `
			UPDATE articles
			SET description_cyrillic = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, description_cyrillic, id)
		if err != nil {
			log.Printf("%v: writing description_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	photos := r.MultipartForm.File["photos"]
	if len(photos) == 0 {
		photos = nil
	}

	var photosForDb bytes.Buffer

	for _, fh := range photos {
		if fh.Size > 2<<20 {
			log.Printf("%v: photo size exceeds 2MB limit: %v", r.URL, fh.Size)
			response.Res(w, "error", http.StatusBadRequest, "photo size exceeds 2MB limit")
			return
		} else {
			file, _ := fh.Open()
			io.Copy(&photosForDb, file)
			file.Close()
		}
	}

	if photosForDb.Len() > 0 {
		sqlStatement := `
			UPDATE articles
			SET photos = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, pq.Array([][]byte{photosForDb.Bytes()}), id)
		if err != nil {
			log.Printf("%v: writing photos into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if videos, ok := r.Form["videos"]; ok {
		videosString := "{" + strings.Join(videos, ",") + "}"
		sqlStatement := `
			UPDATE articles
			SET videos = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, videosString, id)
		if err != nil {
			log.Printf("%v: writing videos into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	coverImage, coverImageHeader, err := r.FormFile("cover_image")
	if err != nil {
		if err == http.ErrMissingFile {
			coverImage = nil
		} else {
			log.Printf("%v: cover_image error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, "cover_image error")
			return
		}
	}

	var coverImageForDB []byte = nil

	if coverImage != nil {
		if coverImageHeader.Size > 1<<20 {
			log.Printf("%v: cover_image size exceeds 1MB limit: %v", r.URL, coverImageHeader.Size)
			response.Res(w, "error", http.StatusBadRequest, "cover_image size exceeds 1MB limit")
			return
		}
		coverImageForDB, _ = io.ReadAll(coverImage)
		coverImage.Close()
		sqlStatement := `
			UPDATE articles
			SET cover_image = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, coverImageForDB, id)
		if err != nil {
			log.Printf("%v: writing cover_image into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if tags, ok := r.Form["tags"]; ok {
		tagsString := "{" + strings.Join(tags, ",") + "}"
		sqlStatement := `
			UPDATE articles
			SET tags = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, tagsString, id)
		if err != nil {
			log.Printf("%v: writing tags into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if category := r.FormValue("category"); category != "" {
		categoryInt, err := strconv.ParseInt(category, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE articles
			SET category = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, categoryInt, id)
		if err != nil {
			log.Printf("%v: writing category into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if related := r.FormValue("related"); related != "" {
		relatedInt, err := strconv.ParseInt(related, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE articles
			SET related = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, relatedInt, id)
		if err != nil {
			log.Printf("%v: writing related into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "Article edited")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: delete article articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: delete article articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*archived {
		log.Printf("%v: delete article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete not archived article")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM articles WHERE id=$1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "deleted")
}

func archiveArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: archive article articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: archive article articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: archive article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive already archived article")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE articles SET archived = true WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "archived")
}

type ArticleCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

func getArticleCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	period := vars["period"]

	if period != "week" && period != "month" && period != "year" {
		response.Res(w, "error", http.StatusBadRequest, "invalid period value")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM articles WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, ArticleCount{Period: period, Count: count})
}
