package model

import "Tahlilchi.uz/db"

// VideoNewsList is a struct to map the video news list data
type VideoNewsListResponse struct {
	VideoNewsList []VideoNews `json:"video_news_list"`
	Previous      bool        `json:"previous"`
	Next          bool        `json:"next"`
}

// VideoNews is a struct to map the video news data
type VideoNews struct {
	ID           int    `json:"id"`
	Video        string `json:"video"`
	TextLatin    string `json:"text_latin"`
	TextCyrillic string `json:"text_cyrillic"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Archived     bool   `json:"archived"`
	Completed    bool   `json:"completed"`
}

// AddVideoNews is a method to add a video news to the database
func (vn *VideoNews) AddVideoNews() error {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return err
	}
	// defer the close of the database connection
	defer database.Close()

	// create a new transaction
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	// prepare the insert statement
	stmt, err := tx.Prepare("INSERT INTO video_news (video, text_latin, text_cyrillic) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	// defer the close of the statement
	defer stmt.Close()

	// execute the insert statement
	_, err = stmt.Exec(vn.Video, vn.TextLatin, vn.TextCyrillic)
	if err != nil {
		return err
	}
	// commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	// return nil
	return nil
}

// UpdateVideoNews is a method to update a video news in the database
func (vn *VideoNews) UpdateVideoNews(id string) error {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return err
	}
	// defer the close of the database connection
	defer database.Close()

	// if VideoNews Video is not empty, update the video news video, set current date time to updated_at
	if vn.Video != "" {
		_, err := database.Exec("UPDATE video_news SET video = $1, updated_at = now() WHERE id = $2", vn.Video, id)
		if err != nil {
			return err
		}
	}

	// if VideoNews TextLatin is not empty, update the video news text_latin
	if vn.TextLatin != "" {
		_, err := database.Exec("UPDATE video_news SET text_latin = $1, updated_at = now() WHERE id = $2", vn.TextLatin, id)
		if err != nil {
			return err
		}
	}
	// if VideoNews TextCyrillic is not empty, update the video news text_cyrillic
	if vn.TextCyrillic != "" {
		_, err := database.Exec("UPDATE video_news SET text_cyrillic = $1, updated_at = now() WHERE id = $2", vn.TextCyrillic, id)
		if err != nil {
			return err
		}
	}
	// return nil
	return nil
}

// DeleteVideoNews is a method to delete a video news from the database
func (vn *VideoNews) DeleteVideoNews() error {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return err
	}
	// defer the close of the database connection
	defer database.Close()

	// execute the delete statement
	_, err = database.Exec("DELETE FROM video_news WHERE id = $1", vn.ID)
	if err != nil {
		return err
	}
	// return nil
	return nil
}

// GetVideoNewsList is a function to get a list of video news from the database. The function receives limit and offset as parameters. it return pointer to VideoNewsListResponse and error.
func GetVideoNewsList(limit, offset int) (*VideoNewsListResponse, error) {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return nil, err
	}
	// defer the close of the database connection
	defer database.Close()

	// create a new VideoNewsListResponse
	vnList := VideoNewsListResponse{}

	// execute the select statement to get the video news list
	rows, err := database.Query("SELECT id, video, text_latin, text_cyrillic, created_at, updated_at, archived, completed FROM video_news ORDER BY id DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	// defer the close of the rows
	defer rows.Close()

	// iterate through the rows
	for rows.Next() {
		// create a new VideoNews
		vn := VideoNews{}
		// scan the rows into the VideoNews
		err := rows.Scan(&vn.ID, &vn.Video, &vn.TextLatin, &vn.TextCyrillic, &vn.CreatedAt, &vn.UpdatedAt, &vn.Archived, &vn.Completed)
		if err != nil {
			return nil, err
		}
		// append the VideoNews to the VideoNewsListResponse
		vnList.VideoNewsList = append(vnList.VideoNewsList, vn)
	}
	// if there is an error iterating through the rows, return the error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// set appropriate value for previous field
	vnList.Previous = offset > 0

	// set appropriate value for next field
	// total count of video news
	var total int
	// execute the select statement to get the count of video news
	err = database.QueryRow("SELECT COUNT(*) FROM video_news").Scan(&total)
	if err != nil {
		return nil, err
	}
	// if the count is greater than the sum of the offset and limit, set next to true
	vnList.Next = total > offset+limit

	// return the VideoNewsListResponse and nil
	return &vnList, nil
}
