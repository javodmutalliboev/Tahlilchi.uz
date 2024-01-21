package admin

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

type Category struct {
	TitleLatin          string
	DescriptionLatin    string
	TitleCyrillic       string
	DescriptionCyrillic string
}

func addCategory(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec("INSERT INTO news_category(title_latin, description_latin, title_cyrillic, description_cyrillic) VALUES($1, $2, $3, $4)", c.TitleLatin, c.DescriptionLatin, c.TitleCyrillic, c.DescriptionCyrillic)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Category Added")
}

type Subcategory struct {
	CategoryID          int    `json:"category_id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

func addSubcategory(w http.ResponseWriter, r *http.Request) {
	var s Subcategory
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// Check if category_id and titles are provided
	if s.CategoryID == 0 || s.TitleLatin == "" || s.TitleCyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "category_id and titles are required")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO news_subcategory (category_id, title_latin, description_latin, title_cyrillic, description_cyrillic) VALUES ($1, $2, $3, $4, $5)", s.CategoryID, s.TitleLatin, s.DescriptionLatin, s.TitleCyrillic, s.DescriptionCyrillic)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "subcategory added")
}

type Region struct {
	NameLatin           string `json:"name_latin"`
	DescriptionLatin    string `json:"description_latin,omitempty"`
	NameCyrillic        string `json:"name_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

func addRegions(w http.ResponseWriter, r *http.Request) {
	var regions []Region

	err := json.NewDecoder(r.Body).Decode(&regions)
	if err != nil {
		log.Println(err)
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

	for _, region := range regions {
		if region.NameLatin == "" || region.NameCyrillic == "" {
			response.Res(w, "error", http.StatusBadRequest, "Name fields are required")
			return
		}

		_, err = db.Exec("INSERT INTO news_regions (name_latin, description_latin, name_cyrillic, description_cyrillic) VALUES ($1, $2, $3, $4)", region.NameLatin, region.DescriptionLatin, region.NameCyrillic, region.DescriptionCyrillic)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "Regions added successfully")
}

func addNewsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(15 << 20) // Max memory 15MB
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	title_latin := r.FormValue("title_latin")
	description_latin := r.FormValue("description_latin")
	title_cyrillic := r.FormValue("title_cyrillic")
	description_cyrillic := r.FormValue("description_cyrillic")

	// Check if required fields are not empty
	if title_latin == "" || description_latin == "" || title_cyrillic == "" || description_cyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "Required fields are missing")
		return
	}

	photo, photo_header, err := r.FormFile("photo")
	var photoForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("photo: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		photoForDB = nil
	} else {
		// Check size limits
		if photo_header.Size > int64(2<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Photo exceeds 2MB limit")
			return
		}
		photoForDB, _ = io.ReadAll(photo)
		photo.Close()
	}

	video, video_header, err := r.FormFile("video")
	var videoForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("video: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		videoForDB = nil
	} else {
		if video_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Video exceeds 6MB limit")
			return
		}
		videoForDB, _ = io.ReadAll(video)
		video.Close()
	}

	audio, audio_header, err := r.FormFile("audio")
	var audioForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("audio: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		audioForDB = nil
	} else {
		if audio_header.Size > int64(4<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Audio exceeds 4MB limit")
			return
		}
		audioForDB, _ = io.ReadAll(audio)
		audio.Close()
	}

	cover_image, cover_image_header, err := r.FormFile("cover_image")
	var coverImageForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("cover_image: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		coverImageForDB = nil
	} else {
		if cover_image_header.Size > int64(1<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Cover image exceeds 1MB limit")
			return
		}
		coverImageForDB, _ = io.ReadAll(cover_image)
		cover_image.Close()
	}

	// Get tags if they exist
	tags, ok := r.Form["tags"]
	if !ok {
		// If tags don't exist, use an empty array
		tags = []string{}
	}

	// Convert tags to PostgreSQL array format
	tagsString := "{" + strings.Join(tags, ",") + "}"

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO news_posts (title_latin, description_latin, title_cyrillic, description_cyrillic, photo, video, audio, cover_image, tags) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		title_latin, description_latin, title_cyrillic, description_cyrillic, photoForDB, videoForDB, audioForDB, coverImageForDB, tagsString)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "New post has been created successfully.")
}

func editNewsPost(w http.ResponseWriter, r *http.Request) {
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
			UPDATE news_posts
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
			UPDATE news_posts
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
			UPDATE news_posts
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
			UPDATE news_posts
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

	photo, photo_header, err := r.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("photo: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		// Check size limits
		if photo_header.Size > int64(2<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Photo exceeds 2MB limit")
			return
		}
		photoForDB, _ := io.ReadAll(photo)
		photo.Close()
		sqlStatement := `
			UPDATE news_posts
			SET photo = $1
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, photoForDB, id)
		if err != nil {
			log.Printf("%v: writing photo into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	video, video_header, err := r.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("video: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		if video_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Video exceeds 6MB limit")
			return
		}
		videoForDB, _ := io.ReadAll(video)
		video.Close()
		sqlStatement := `
			UPDATE news_posts
			SET video = $1
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, videoForDB, id)
		if err != nil {
			log.Printf("%v: writing video into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	audio, audio_header, err := r.FormFile("audio")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("audio: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		if audio_header.Size > int64(4<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Audio exceeds 4MB limit")
			return
		}
		audioForDB, _ := io.ReadAll(audio)
		audio.Close()
		sqlStatement := `
			UPDATE news_posts
			SET audio = $1
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, audioForDB, id)
		if err != nil {
			log.Printf("%v: writing audio into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	cover_image, cover_image_header, err := r.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("cover_image: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		if cover_image_header.Size > int64(1<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Cover image exceeds 1MB limit")
			return
		}
		coverImageForDB, _ := io.ReadAll(cover_image)
		cover_image.Close()
		sqlStatement := `
			UPDATE news_posts
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
			UPDATE news_posts
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

	response.Res(w, "success", http.StatusOK, "OK")
}
