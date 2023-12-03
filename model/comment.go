package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Comment struct {
	CommentId        int64         `json:"comment_id,string" db:"comment_id"`
	Content          string        `json:"content" db:"content" binding:"required"`
	RootCommentId    int64         `json:"root_comment_id,string" db:"root_comment_id"`
	ToCommentId      int64         `json:"to_comment_id,string" db:"to_comment_id"`
	UserId           int64         `json:"user_id,string" db:"user_id"`
	PostId           int64         `json:"post_id,string" db:"post_id" binding:"required"`
	CreateTime       time.Time     `json:"create_time" db:"create_time"`
	UpdateTime       time.Time     `json:"update_time" db:"update_time"`
	CommentLikeCount int           `json:"comment_like_count" db:"comment_like_count"`
	IsDelete         types.BitBool `json:"is_delete" db:"is_delete"`
}

type ApiCommentDetail struct {
	*Comment
	AuthorName string `json:"author_name"`
}
