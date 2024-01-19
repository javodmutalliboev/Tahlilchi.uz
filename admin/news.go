package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
)

type Category struct {
	Title       string
	Description string
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
	c.Title = r.FormValue("title")
	c.Description = r.FormValue("description")

	if c.Title == "" {
		response.Res(w, "error", http.StatusBadRequest, "Title is required")
		return
	}

	if c.Description == "" {
		_, err := db.Exec("INSERT INTO news_category(title) VALUES($1)", c.Title)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	} else {
		_, err := db.Exec("INSERT INTO news_category(title, description) VALUES($1, $2)", c.Title, c.Description)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "Category Added")
}

type Subcategory struct {
	CategoryID  int    `json:"category_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func addSubcategory(w http.ResponseWriter, r *http.Request) {
	var s Subcategory
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// Check if category_id and title are provided
	if s.CategoryID == 0 || s.Title == "" {
		response.Res(w, "error", http.StatusBadRequest, "category_id and title are required")
		return
	}

	// If description is not provided, set it to an empty string
	if s.Description == "" {
		s.Description = ""
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO news_subcategory (category_id, title, description) VALUES ($1, $2, $3)", s.CategoryID, s.Title, s.Description)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "subcategory added")
}

type Region struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
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
		if region.Name == "" {
			response.Res(w, "error", http.StatusBadRequest, "Name field is required")
			return
		}

		_, err = db.Exec("INSERT INTO news_regions (name, description) VALUES ($1, $2)", region.Name, region.Description)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "Regions added successfully")
}
