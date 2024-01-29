package copilotHistory

/*

 */

/*
package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/news/category/{category}", GetNewsByCategory).Methods("GET")
	r.HandleFunc("/news/subcategory/{subcategory}", GetNewsBySubcategory).Methods("GET")
	r.HandleFunc("/news/region/{region}", GetNewsByRegion).Methods("GET")
	r.HandleFunc("/news/top", GetTopNews).Methods("GET")
	r.HandleFunc("/news/latest", GetLatestNews).Methods("GET")
	r.HandleFunc("/news/related/{newsId}", GetRelatedNews).Methods("GET")
	r.HandleFunc("/news/all", GetAllNews).Methods("GET")

	http.ListenAndServe(":8000", r)
}

func GetNewsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	// Your code here to fetch news by category from your PostgreSQL database
}

func GetNewsBySubcategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subcategory := vars["subcategory"]

	// Your code here to fetch news by subcategory from your PostgreSQL database
}

func GetNewsByRegion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := vars["region"]

	// Your code here to fetch news by region from your PostgreSQL database
}

func GetTopNews(w http.ResponseWriter, r *http.Request) {
	// Your code here to fetch top news from your PostgreSQL database
}

func GetLatestNews(w http.ResponseWriter, r *http.Request) {
	// Your code here to fetch latest news from your PostgreSQL database
}

func GetRelatedNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsId := vars["newsId"]

	// Your code here to fetch related news from your PostgreSQL database
}

func GetAllNews(w http.ResponseWriter, r *http.Request) {
	// Your code here to fetch all news from your PostgreSQL database
}

*/

/*
package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/news/category/{category}", GetNewsByCategory).Methods("GET")
	r.HandleFunc("/news/subcategory/{subcategory}", GetNewsBySubcategory).Methods("GET")
	r.HandleFunc("/news/region/{region}", GetNewsByRegion).Methods("GET")
	r.HandleFunc("/news/top", GetTopNews).Methods("GET")
	r.HandleFunc("/news/latest", GetLatestNews).Methods("GET")
	r.HandleFunc("/news/related/{newsId}", GetRelatedNews).Methods("GET")
	r.HandleFunc("/news/all", GetAllNews).Methods("GET")

	http.ListenAndServe(":8000", r)
}

func GetNewsByCategory(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

func GetNewsBySubcategory(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

func GetNewsByRegion(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

func GetTopNews(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

func GetLatestNews(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

func GetRelatedNews(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

func GetAllNews(w http.ResponseWriter, r *http.Request) {
	// Your code here
}

*/
