package client

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

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

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic FROM photo_gallery ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleryList []PhotoGallery
	for rows.Next() {
		var photoGallery PhotoGallery
		if err := rows.Scan(&photoGallery.ID, &photoGallery.TitleLatin, &photoGallery.TitleCyrillic); err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photoGalleryList = append(photoGalleryList, photoGallery)
	}

	response.Res(w, "success", http.StatusOK, photoGalleryList)
}

type PhotoGallery struct {
	ID            int    `json:"id"`
	TitleLatin    string `json:"title_latin"`
	TitleCyrillic string `json:"title_cyrillic"`
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

	rows, err := database.Query("SELECT id, file_name FROM photo_gallery_photos WHERE photo_gallery = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photoGalleryPhotos []PhotoGalleryPhoto
	for rows.Next() {
		var photoGalleryPhoto PhotoGalleryPhoto
		if err := rows.Scan(&photoGalleryPhoto.ID, &photoGalleryPhoto.FileName); err != nil {
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
	vars := mux.Vars(r)
	photo_gallery := vars["id"]
	id := vars["photo_id"]

	// open a database connection
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

	// send file
	contentType := http.DetectContentType(photoGalleryPhoto.File)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(photoGalleryPhoto.File)))
	w.Write(photoGalleryPhoto.File)
}
