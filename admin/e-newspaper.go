package admin

import (
	"io"
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
)

func addENewspaper(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20) // Max memory 20MB
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	title_latin := r.FormValue("title_latin")
	title_cyrillic := r.FormValue("title_cyrillic")

	file_latin, file_latin_header, err := r.FormFile("file_latin")
	var fileLatinForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: file_latin error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		fileLatinForDB = nil
	} else {
		// Check size limits
		if file_latin_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_latin exceeds 6MB limit")
			return
		}
		fileLatinForDB, _ = io.ReadAll(file_latin)
		file_latin.Close()
	}

	file_cyrillic, file_cyrillic_header, err := r.FormFile("file_cyrillic")
	var fileCyrillicForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: file_cyrillic error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		fileCyrillicForDB = nil
	} else {
		if file_cyrillic_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_cyrillic exceeds 6MB limit")
			return
		}
		fileCyrillicForDB, _ = io.ReadAll(file_cyrillic)
		file_cyrillic.Close()
	}

	cover_image, cover_image_header, err := r.FormFile("cover_image")
	var coverImageForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: cover_image error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		coverImageForDB = nil
	} else {
		if cover_image_header.Size > int64(3<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Cover image exceeds 3MB limit")
			return
		}
		coverImageForDB, _ = io.ReadAll(cover_image)
		cover_image.Close()
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO e_newspapers (title_latin, title_cyrillic, file_latin, file_cyrillic, cover_image) VALUES ($1, $2, $3, $4, $5)`,
		title_latin, title_cyrillic, fileLatinForDB, fileCyrillicForDB, coverImageForDB)
	if err != nil {
		log.Printf("%v: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "e-newspaper has been added successfully.")
}
