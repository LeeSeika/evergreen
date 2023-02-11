package mysql

import (
	"evergreen/model"

	"go.uber.org/zap"
)

func CreatePost(p *model.Post) error {
	sqlStr := "insert into post (post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)"
	if _, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID); err != nil {
		zap.L().Error("insert post error", zap.Error(err))
		return err
	}
	return nil
}

func GetPostDetailByID(postID int64) (*model.Post, error) {
	var post model.Post
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post where post_id = ?"
	err := db.Get(&post, sqlStr, postID)
	if err != nil {
		zap.L().Error("get post detail by id error", zap.Int64("postID", postID), zap.Error(err))
		return nil, err
	}
	return &post, nil
}

func GetPostList(page, size int64) ([]*model.Post, error) {
	postList := make([]*model.Post, 0, size)
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post limit ?, ?"
	err := db.Select(&postList, sqlStr, (page-1)*size, size)
	if err != nil {
		zap.L().Error("get post list error", zap.Int64("page", page), zap.Int64("size", size), zap.Error(err))
		return nil, err
	}
	return postList, nil
}
