package model

import "time"

type Post struct {
	ID           int64     `json:"id,string" db:"post_id"`
	AuthorID     int64     `json:"author_id,string" db:"author_id" binding:"required"`
	CommunityID  int64     `json:"community_id" db:"community_id"`
	Status       int32     `json:"status" db:"status"`
	Title        string    `json:"title" db:"title"`
	Content      string    `json:"content" db:"content"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	CommentCount int       `json:"comment_count" db:"comment_count"`
}

type ApiPostDetail struct {
	AuthorName  string `json:"author_name"`
	VoteNumbers int64  `json:"vote_numbers"`
	*Post
	*CommunityDetail `json:"community"`
}
