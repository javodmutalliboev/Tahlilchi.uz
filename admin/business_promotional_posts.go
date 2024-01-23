package admin

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/lib/pq"
)

func addBusinessPromotionalPost(w http.ResponseWriter, r *http.Request) {
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

	expiration := r.FormValue("expiration")
	expirationValid := checkExpiration(expiration)
	if !expirationValid {
		log.Printf("%v: FormValue(\"expiration\") valid: %v", r.URL, expirationValid)
		response.Res(w, "error", http.StatusBadRequest, "expiration value is invalid")
		return
	}
	var expirationForDB time.Time
	switch expiration {
	case "1 day":
		expirationForDB = time.Now().UTC().Add(24 * time.Hour)
	case "1 week":
		expirationForDB = time.Now().UTC().Add(7 * 24 * time.Hour)
	case "1 month":
		expirationForDB = time.Now().UTC().AddDate(0, 1, 0)
	}

	partner := r.FormValue("partner")
	if partner == "" {
		log.Printf("%v: partner: %v", r.URL, partner)
		response.Res(w, "error", http.StatusBadRequest, "partner is required field")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO business_promotional_posts (title_latin, description_latin, title_cyrillic, description_cyrillic, photos, videos, cover_image, expiration, partner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		titleLatin, descriptionLatin, titleCyrillic, descriptionCyrillic, pq.Array([][]byte{photosForDb.Bytes()}), pq.Array([][]byte{videosForDB.Bytes()}), coverImageForDB, expirationForDB, partner)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Article Added")
}

func checkExpiration(expiration string) bool {
	pattern := regexp.MustCompile(`^1 (day|week|month)$`)
	return pattern.MatchString(expiration)
}
