package mysql

import (
	"database/sql"
	"evergreen/model"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

func AddComment(tx *sqlx.Tx, p *model.Comment) (sql.Result, error) {
	sqlStr := "insert into comment (content, root_comment_id, to_comment_id, user_id, post_id) values (?, ?, ?, ?, ?)"
	rs, err := tx.Exec(sqlStr, p.Content, p.RootCommentId, p.ToCommentId, p.UserId, p.PostId)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		return nil, err
	}
	return rs, nil
}

func DeleteComment(tx *sqlx.Tx, commentId int64) (sql.Result, error) {
	sqlStr := "update comment set is_delete = b'1' where comment_id = ?"
	rs, err := tx.Exec(sqlStr, commentId)
	if err != nil {
		zap.L().Error("update comment is_delete failed", zap.Error(err), zap.Int64("comment_id", commentId))
		return nil, err
	}

	return rs, nil
}

func CheckCommentIsDeleted(tx *sqlx.Tx, commentId int64) (*model.Comment, error) {
	sqlStr := "select * from comment where comment_id = ?"
	comment := model.Comment{}
	err := tx.Get(&comment, sqlStr, commentId)
	if err != nil {
		zap.L().Error("check comment is_delete failed", zap.Error(err), zap.Int64("comment_id", commentId))
		return nil, err
	}
	if comment.IsDelete {
		err = ErrorCommentDeleted
		zap.L().Error("comment has been deleted", zap.Int64("comment_id", commentId))
		return nil, err
	}
	return &comment, nil
}

func GetRootCommentsFromPost(p *model.ParamComments) ([]*model.Comment, error) {
	var sqlStr string
	commentList := make([]*model.Comment, 0, p.Size)
	if p.Order == model.OrderByScore {
		// todo 与常规的limit m, n 比较优化后的提速
		sqlStr = "select * from comment where post_id = ? and root_comment_id = -1 and is_delete == b'0' and comment_like_count <= (select comment_like_count from comment where post_id = ? and root_comment_id = -1 and is_delete == b'0' order by comment_like_count DESC limit ?, 1) order by comment_like_count DESC limit ?"
	} else {
		sqlStr = "select * from comment where post_id = ? and root_comment_id = -1 and is_delete == b'0' and comment_id <= (select comment_id from comment where post_id = ? and root_comment_id = -1 and is_delete == b'0' order by comment_id DESC limit ?, 1) order by comment_id DESC limit ?"
	}
	err := db.Select(&commentList, sqlStr, p.PostID, p.PostID, (p.Page-1)*p.Size, p.Size)
	if err != nil {
		zap.L().Error("get comments failed", zap.String("order", p.Order), zap.Int64("page", p.Page), zap.Int64("size", p.Size), zap.Error(err))
		return nil, err
	}
	return commentList, nil
}
