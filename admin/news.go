package admin

import (
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
