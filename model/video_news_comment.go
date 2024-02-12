package model

import (
	"database/sql"

	"Tahlilchi.uz/db"
)

// VideoNewsCommentListResponse is a struct to represent the video news comment list response
type VideoNewsCommentListResponse struct {
	VideoNewsCommentList []VideoNewsComment `json:"video_news_comment_list"`
	Previous             bool               `json:"previous"`
	Next                 bool               `json:"next"`
}

// VideoNewsComment is a struct to represent a video news comment
type VideoNewsComment struct {
	ID        int    `json:"id"`
	VideoNews int    `json:"video_news"`
	Text      string `json:"text" validate:"required"`
	Contact   string `json:"contact"`
	CreatedAt string `json:"created_at"`
	Approved  bool   `json:"approved"`
}

// AddVideoNewsComment is a method to add a video news comment to the database
func (vnc *VideoNewsComment) AddVideoNewsComment(videoNews int) error {
	database, err := db.DB()
	if err != nil {
		return err
	}
	defer database.Close()

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO video_news_comments (video_news, text) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(videoNews, vnc.Text).Scan(&vnc.ID)
	if err != nil {
		return err
	}

	if vnc.Contact != "" {
		_, err = tx.Exec("UPDATE video_news_comments SET contact = $1 WHERE id = $2", vnc.Contact, vnc.ID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetVideoNewsCommentListResponse is a method to get the video news comment list response
func (vnc *VideoNewsComment) GetVideoNewsCommentListResponse(admin bool, videoNews int, page, limit int) (VideoNewsCommentListResponse, error) {
	database, err := db.DB()
	if err != nil {
		return VideoNewsCommentListResponse{}, err
	}
	defer database.Close()

	var rows *sql.Rows
	if admin {
		// select id, video_news, text, contact, created_at, approved from video_news_comments where video_news = $1 order by id desc limit $2 offset $3
		rows, err = database.Query("SELECT id, video_news, text, contact, created_at, approved FROM video_news_comments WHERE video_news = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", videoNews, limit, (page-1)*limit)
	} else {
		// select id, text, created_at from video_news_comments where video_news = $1 and approved = true order by id desc limit $2 offset $3
		rows, err = database.Query("SELECT id, text, created_at FROM video_news_comments WHERE video_news = $1 AND approved = true ORDER BY id DESC LIMIT $2 OFFSET $3", videoNews, limit, (page-1)*limit)
	}
	if err != nil {
		return VideoNewsCommentListResponse{}, err
	}
	defer rows.Close()

	var vncList []VideoNewsComment
	for rows.Next() {
		var vnc VideoNewsComment
		if admin {
			err = rows.Scan(&vnc.ID, &vnc.VideoNews, &vnc.Text, &vnc.Contact, &vnc.CreatedAt, &vnc.Approved)
		} else {
			err = rows.Scan(&vnc.ID, &vnc.Text, &vnc.CreatedAt)
		}
		if err != nil {
			return VideoNewsCommentListResponse{}, err
		}
		vncList = append(vncList, vnc)
	}

	if err = rows.Err(); err != nil {
		return VideoNewsCommentListResponse{}, err
	}

	var vncListResponse VideoNewsCommentListResponse
	var count int
	if admin {
		err = database.QueryRow("SELECT COUNT(id) FROM video_news_comments WHERE video_news = $1", videoNews).Scan(&count)
	} else {
		err = database.QueryRow("SELECT COUNT(id) FROM video_news_comments WHERE video_news = $1 AND approved = true", videoNews).Scan(&count)
	}
	if err != nil {
		return VideoNewsCommentListResponse{}, err
	}

	if page > 1 {
		vncListResponse.Previous = true
	}

	if count > page*limit {
		vncListResponse.Next = true
	}

	vncListResponse.VideoNewsCommentList = vncList

	return vncListResponse, nil
}

// ApproveVideoNewsComment is a method to approve a video news comment
func (vnc *VideoNewsComment) ApproveVideoNewsComment(videoNews, commentID int) error {
	database, err := db.DB()
	if err != nil {
		return err
	}
	defer database.Close()

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	// toggle approved column
	_, err = tx.Exec("UPDATE video_news_comments SET approved = NOT approved WHERE id = $1 AND video_news = $2", commentID, videoNews)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
