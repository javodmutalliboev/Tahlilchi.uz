package model

import (
	"database/sql"

	"Tahlilchi.uz/db"
)

type NewsPostCommentListResponse struct {
	NewsPostCommentList []NewsPostComment `json:"news_post_comment_list"`
	Previous            bool              `json:"previous"`
	Next                bool              `json:"next"`
}

type NewsPostComment struct {
	ID        int    `json:"id"`
	NewsPost  int    `json:"news_post"`
	Text      string `json:"text" validate:"required"`
	Contact   string `json:"contact"`
	CreatedAt string `json:"created_at"`
	Approved  bool   `json:"approved"`
}

func (npc *NewsPostComment) AddNewsPostComment(newsPost int) error {
	database, err := db.DB()
	if err != nil {
		return err
	}
	defer database.Close()

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO news_post_comments (news_post, text) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(newsPost, npc.Text).Scan(&npc.ID)
	if err != nil {
		return err
	}

	if npc.Contact != "" {
		_, err = tx.Exec("UPDATE news_post_comments SET contact = $1 WHERE id = $2", npc.Contact, npc.ID)
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

func (npc *NewsPostComment) GetNewsPostCommentListResponse(admin bool, newsPost int, page, limit int) (NewsPostCommentListResponse, error) {
	database, err := db.DB()
	if err != nil {
		return NewsPostCommentListResponse{}, err
	}
	defer database.Close()

	var rows *sql.Rows
	if admin {
		// select id, news_post, text, contact, created_at, approved from news_post_comments where news_post = $1 order by id desc limit $2 offset $3
		rows, err = database.Query("SELECT id, news_post, text, contact, created_at, approved FROM news_post_comments WHERE news_post = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", newsPost, limit, (page-1)*limit)
	} else {
		// select id, text, created_at from news_post_comments where news_post = $1 and approved = true order by id desc limit $2 offset $3
		rows, err = database.Query("SELECT id, text, created_at FROM news_post_comments WHERE news_post = $1 AND approved = true ORDER BY id DESC LIMIT $2 OFFSET $3", newsPost, limit, (page-1)*limit)
	}
	if err != nil {
		return NewsPostCommentListResponse{}, err
	}
	defer rows.Close()

	var npcList []NewsPostComment
	for rows.Next() {
		var npc NewsPostComment
		if admin {
			err = rows.Scan(&npc.ID, &npc.NewsPost, &npc.Text, &npc.Contact, &npc.CreatedAt, &npc.Approved)
		} else {
			err = rows.Scan(&npc.ID, &npc.Text, &npc.CreatedAt)
		}
		if err != nil {
			return NewsPostCommentListResponse{}, err
		}
		npcList = append(npcList, npc)
	}

	if err = rows.Err(); err != nil {
		return NewsPostCommentListResponse{}, err
	}

	var npcListResponse NewsPostCommentListResponse
	var count int
	if admin {
		// select count(id) from news_post_comments where news_post = $1
		err = database.QueryRow("SELECT COUNT(id) FROM news_post_comments WHERE news_post = $1", newsPost).Scan(&count)
	} else {
		// select count(id) from news_post_comments where news_post = $1 and approved = true
		err = database.QueryRow("SELECT COUNT(id) FROM news_post_comments WHERE news_post = $1 AND approved = true", newsPost).Scan(&count)
	}
	if err != nil {
		return NewsPostCommentListResponse{}, err
	}

	if page > 1 {
		npcListResponse.Previous = true
	}
	if count > page*limit {
		npcListResponse.Next = true
	}

	npcListResponse.NewsPostCommentList = npcList
	return npcListResponse, nil
}

func (npc *NewsPostComment) ApproveNewsPostComment(newsPost, commentID int) error {
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
	_, err = tx.Exec("UPDATE news_post_comments SET approved = NOT approved WHERE news_post = $1 AND id = $2", newsPost, commentID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
