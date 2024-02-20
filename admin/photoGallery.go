package admin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
)

type PhotoGallery struct {
	ID            int    `json:"id"`
	TitleLatin    string `json:"title_latin"`
	TitleCyrillic string `json:"title_cyrillic"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
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
	err = r.ParseMultipartForm(500 << 20) // Max memory 500 MB
	if err != nil {
		log.Printf("%v: Could not parse multipart form: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called.
	files := r.MultipartForm.File["photo"] // "photo" is the key of the input form

	for _, fileHeader := range files {
		// Check if the file is an image
		if fileHeader.Header.Get("Content-Type")[:5] != "image" {
			err := fmt.Errorf("photo %v is not an image: %v", fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}

		// Check if the file size is greater than 10MB
		if fileHeader.Size > 10<<20 {
			err := fmt.Errorf("photo %v size exceeds 10MB limit: %v", fileHeader.Filename, fileHeader.Size)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
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

		_, err = db.Exec("UPDATE photo_gallery SET updated_at = CURRENT_TIMESTAMP WHERE id = $1", id)
		if err != nil {
			message := fmt.Sprintf("%v: %v: UPDATE photo_gallery SET updated_at = CURRENT_TIMESTAMP error: %v", r.URL, fileHeader.Filename, err)
			log.Println(message)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "photo gallery photos added")
}

func getPhotoGalleryList(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT * FROM photo_gallery ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleryList []PhotoGalleryList
	for rows.Next() {
		var photoGallery PhotoGalleryList
		if err := rows.Scan(&photoGallery.ID, &photoGallery.TitleLatin, &photoGallery.TitleCyrillic, &photoGallery.CreatedAt, &photoGallery.UpdatedAt); err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photoGalleryList = append(photoGalleryList, photoGallery)
	}

	response.Res(w, "success", http.StatusOK, photoGalleryList)
}

type PhotoGalleryList struct {
	ID            int    `json:"id"`
	TitleLatin    string `json:"title_latin"`
	TitleCyrillic string `json:"title_cyrillic"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func getPhotoGalleryPhotos(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := photoGalleryExists(id)
	if err != nil {
		log.Printf("%v: getPhotoGalleryPhotos photoGalleryExists(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: getPhotoGalleryPhotos photoGalleryExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "photo gallery not found")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, photo_gallery, file_name, created_at FROM photo_gallery_photos WHERE photo_gallery = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleryPhotos []PhotoGalleryPhoto
	for rows.Next() {
		var photoGalleryPhoto PhotoGalleryPhoto
		if err := rows.Scan(&photoGalleryPhoto.ID, &photoGalleryPhoto.PhotoGallery, &photoGalleryPhoto.FileName, &photoGalleryPhoto.CreatedAt); err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photoGalleryPhotos = append(photoGalleryPhotos, photoGalleryPhoto)
	}

	response.Res(w, "success", http.StatusOK, photoGalleryPhotos)
}

type PhotoGalleryPhoto struct {
	ID           int       `json:"id"`
	PhotoGallery int       `json:"photo_gallery"`
	FileName     string    `json:"file_name"`
	CreatedAt    time.Time `json:"created_at"`
	File         []byte    `json:"file"`
}

// getPhotoGalleryPhoto is a handler to get photo gallery photo
func getPhotoGalleryPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Var(r)
	photo_gallery := vars["id"]
	id := vars["photo_id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var photoGalleryPhoto PhotoGalleryPhoto
	err = database.QueryRow("SELECT id, photo_gallery, file_name, created_at, file FROM photo_gallery_photos WHERE photo_gallery = $1 AND id = $2", photo_gallery, id).Scan(&photoGalleryPhoto.ID, &photoGalleryPhoto.PhotoGallery, &photoGalleryPhoto.FileName, &photoGalleryPhoto.CreatedAt, &photoGalleryPhoto.File)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	contentType := http.DetectContentType(photoGalleryPhoto.File)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(photoGalleryPhoto.File)))
	w.Write(photoGalleryPhoto.File)
}
