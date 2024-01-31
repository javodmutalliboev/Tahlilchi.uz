package admin

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
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

func bpPostExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM business_promotional_posts WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func bpPostIsArchived(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT archived FROM business_promotional_posts WHERE id = $1")
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

func editBusinessPromotionalPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: edit business promotional post bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: edit business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: edit business promotional post bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: edit business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit archived business promotional post")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(15 << 20)
	if err != nil {
		log.Printf("%v: edit business promotional post: %v", r.URL, err)
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
			UPDATE business_promotional_posts
			SET title_latin = $1, edited_at = NOW()
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
			UPDATE business_promotional_posts
			SET description_latin = $1, edited_at = NOW()
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
			UPDATE business_promotional_posts
			SET title_cyrillic = $1, edited_at = NOW()
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
			UPDATE business_promotional_posts
			SET description_cyrillic = $1, edited_at = NOW()
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
			UPDATE business_promotional_posts
			SET photos = $1, edited_at = NOW()
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
			UPDATE business_promotional_posts
			SET videos = $1, edited_at = NOW()
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
			UPDATE business_promotional_posts
			SET cover_image = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, coverImageForDB, id)
		if err != nil {
			log.Printf("%v: writing cover_image into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	expiration := r.FormValue("expiration")
	if expiration != "" {
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
		sqlStatement := `
			UPDATE business_promotional_posts
			SET expiration = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, expirationForDB, id)
		if err != nil {
			log.Printf("%v: writing expiration into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	partner := r.FormValue("partner")
	if partner != "" {
		sqlStatement := `
			UPDATE business_promotional_posts
			SET partner = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, partner, id)
		if err != nil {
			log.Printf("%v: writing partner into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "business promotional post edited")
}

func deleteBPPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: delete business promotional post bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: delete business promotional post bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*archived {
		log.Printf("%v: delete business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete not archived business promotional post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM business_promotional_posts WHERE id=$1")
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

func archiveBPPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: archive business promotional post bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: archive business promotional post bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: archive business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive already archived business promotional post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE business_promotional_posts SET archived = true WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "archived")
}

func CheckAndArchiveExpiredBPPosts() {
	db, err := db.DB()
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): error: %v", err)
		return
	}
	defer db.Close()

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Start a new transaction: error: %v", err)
		return
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare("UPDATE business_promotional_posts SET archived = true WHERE expiration < $1")
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Prepare the SQL statement: error: %v", err)
		return
	}

	// Execute the SQL statement
	_, err = stmt.Exec(time.Now().UTC())
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Execute the SQL statement: error: %v", err)
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Commit the transaction: error: %v", err)
		return
	}
}

type businessPromotionalPostCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

func getBusinessPromotionalPostCount(w http.ResponseWriter, r *http.Request) {
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
	query := fmt.Sprintf("SELECT COUNT(*) FROM business_promotional_posts WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, businessPromotionalPostCount{Period: period, Count: count})
}

func getBusinessPromotionalPosts(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT * FROM business_promotional_posts ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var bpPosts []bpPost
	for rows.Next() {
		var bpp bpPost
		err := rows.Scan(&bpp.ID, &bpp.TitleLatin, &bpp.DescriptionLatin, &bpp.TitleCyrillic, &bpp.DescriptionCyrillic, pq.Array(&bpp.Videos), &bpp.CoverImage, &bpp.Expiration, &bpp.CreatedAt, &bpp.UpdatedAt, &bpp.Archived, &bpp.Partner)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		bpPosts = append(bpPosts, bpp)
	}

	response.Res(w, "success", http.StatusOK, bpPosts)
}

type bpPost struct {
	ID                  int      `json:"id"`
	TitleLatin          string   `json:"title_latin"`
	DescriptionLatin    string   `json:"description_latin"`
	TitleCyrillic       string   `json:"title_cyrillic"`
	DescriptionCyrillic string   `json:"description_cyrillic"`
	Videos              []string `json:"videos"`
	CoverImage          []byte   `json:"cover_image"`
	Expiration          string   `json:"expiration"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Archived            bool     `json:"archived"`
	Partner             string   `json:"partner"`
}
