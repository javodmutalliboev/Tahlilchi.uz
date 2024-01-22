package admin

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
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

	// Get the videos files
	videos := r.MultipartForm.File["videos"]
	if len(videos) == 0 {
		videos = nil
	}

	var videosForDB bytes.Buffer

	for _, fh := range videos {
		if fh.Size > 6<<20 {
			log.Printf("%v: video size exceeds 6MB limit: %v", r.URL, fh.Size)
			response.Res(w, "error", http.StatusBadRequest, "video size exceeds 6MB limit")
			return
		} else {
			file, _ := fh.Open()
			io.Copy(&videosForDB, file)
			file.Close()
		}
	}

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

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO articles (title_latin, description_latin, title_cyrillic, description_cyrillic, photos, videos, cover_image, tags) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		titleLatin, descriptionLatin, titleCyrillic, descriptionCyrillic, pq.Array([][]byte{photosForDb.Bytes()}), pq.Array([][]byte{videosForDB.Bytes()}), coverImageForDB, tagsString)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Article Added")
}

func editArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse multipart form
	err := r.ParseMultipartForm(15 << 20)
	if err != nil {
		log.Printf("edit news post: %v", err)
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
			SET title_latin = $1
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
			SET description_latin = $1
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
			SET title_cyrillic = $1
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
			SET description_cyrillic = $1
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
			SET photos = $1
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, pq.Array([][]byte{photosForDb.Bytes()}), id)
		if err != nil {
			log.Printf("%v: writing photos into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	videos := r.MultipartForm.File["videos"]
	if len(videos) == 0 {
		videos = nil
	}

	var videosForDB bytes.Buffer

	for _, fh := range videos {
		if fh.Size > 6<<20 {
			log.Printf("%v: video size exceeds 6MB limit: %v", r.URL, fh.Size)
			response.Res(w, "error", http.StatusBadRequest, "video size exceeds 6MB limit")
			return
		} else {
			file, _ := fh.Open()
			io.Copy(&videosForDB, file)
			file.Close()
		}
	}

	if videosForDB.Len() > 0 {
		sqlStatement := `
			UPDATE articles
			SET videos = $1
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, pq.Array([][]byte{videosForDB.Bytes()}), id)
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
			SET cover_image = $1
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
			SET tags = $1
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, tagsString, id)
		if err != nil {
			log.Printf("%v: writing tags into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "Article edited")
}
