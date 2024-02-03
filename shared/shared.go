package shared

import (
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

func GetNewsPostPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var photo []byte
	err = database.QueryRow("SELECT photo FROM news_posts WHERE id = $1", id).Scan(&photo)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(photo)
}

func GetNewsPostAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var audio []byte
	err = database.QueryRow("SELECT audio FROM news_posts WHERE id = $1", id).Scan(&audio)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(audio)
}

func GetNewsPostCoverImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var CoverImage []byte
	err = database.QueryRow("SELECT cover_image FROM news_posts WHERE id = $1", id).Scan(&CoverImage)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(CoverImage)
}
