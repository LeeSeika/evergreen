package model

import "time"

type Post struct {
	ID          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id" binding:"required"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	*Post
	*CommunityDetail `json:"community"`
}
