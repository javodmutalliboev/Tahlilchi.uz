package admin

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

type PhotoGallery struct {
	TitleLatin    string
	TitleCyrillic string
}

func addPhotoGallery(w http.ResponseWriter, r *http.Request) {
	var p PhotoGallery
	err := r.ParseForm()
	if err != nil {
		log.Printf("%v: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	p.TitleLatin = r.FormValue("title_latin")
	p.TitleCyrillic = r.FormValue("title_cyrillic")

	if p.TitleLatin == "" || p.TitleCyrillic == "" {
		log.Printf("%v: title_latin: %v; title_cyrillic: %v", r.URL, p.TitleLatin, p.TitleCyrillic)
		response.Res(w, "error", http.StatusBadRequest, "Both title_latin and title_cyrillic are required")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO photo_gallery (title_latin, title_cyrillic) VALUES ($1, $2)", p.TitleLatin, p.TitleCyrillic)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "photo gallery added")
}

func photoGalleryExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM photo_gallery WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func photoGalleryAddPhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := photoGalleryExists(id)
	if err != nil {
		log.Printf("%v: photoGalleryAddPhotos photoGalleryExists(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: photoGalleryAddPhotos photoGalleryExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot add photos to non existent photo gallery")
		return
	}

	// Parse the multipart form in the request
	err = r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		log.Printf("%v: Could not parse multipart form: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called.
	files := r.MultipartForm.File["photos"] // "photos" is the key of the input form

	for _, fileHeader := range files {
		if fileHeader.Size > 2<<20 {
			message := fmt.Sprintf("%v: photo %v size exceeds 2MB limit: %v", r.URL, fileHeader.Filename, fileHeader.Size)
			log.Println(message)
			response.Res(w, "error", http.StatusBadRequest, message)
			return
		}
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			message := fmt.Sprintf("%v: Could not open multipart file %v: %v", r.URL, fileHeader.Filename, err)
			log.Println(message)
			response.Res(w, "error", http.StatusBadRequest, message)
			return
		}
		fileByteA, _ := io.ReadAll(file)
		file.Close()

		_, err = db.Exec("INSERT INTO photo_gallery_photos (photo_gallery, file_name, file) VALUES ($1, $2, $3)", id, fileHeader.Filename, fileByteA)
		if err != nil {
			message := fmt.Sprintf("%v: db.Exec INSERT INTO photo_gallery_photos %v error: %v", r.URL, fileHeader.Filename, err)
			log.Println(message)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		_, err = db.Exec("UPDATE photo_gallery SET edited_at = CURRENT_TIMESTAMP WHERE id = $1", id)
		if err != nil {
			message := fmt.Sprintf("%v: %v: UPDATE photo_gallery SET edited_at = CURRENT_TIMESTAMP error: %v", r.URL, fileHeader.Filename, err)
			log.Println(message)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "photo gallery photos added")
}
