package client

type NewsPost struct {
	ID                  int
	TitleLatin          string
	DescriptionLatin    string
	TitleCyrillic       string
	DescriptionCyrillic string
	Photo               []byte
	Video               string
	Audio               []byte
	CoverImage          []byte
	Tags                []string
}

/*

func fetNewsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	// Parse the page number from the query parameters
	keys, ok := r.URL.Query()["page"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'page' is missing")
		return
	}
	page, _ := strconv.Atoi(keys[0])

	// Define the number of posts per page
	var limit = 10

	// Calculate the starting index
	var start = (page - 1) * limit
}
*/
