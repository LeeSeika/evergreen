package logic

import (
	"database/sql"
	"errors"
	"evergreen/dao/mysql"
	"evergreen/middleware/mq"
	"evergreen/model"
)

func AddComment(p *model.Comment) error {
	tx, err, txFunc := mysql.BeginTransaction()
	if err != nil {
		return err
	}
	defer txFunc(err)

	var rs sql.Result

	rs, err = mysql.AddComment(tx, p)
	if err != nil {
		return err
	}
	rowsAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("rows affected != 1")
	}
	commentId, _ := rs.LastInsertId()
	p.CommentId = commentId

	rs, err = mysql.IncrCommentCountAtPost(tx, p.PostId)
	if err != nil {
		return err
	}
	rowsAffected, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("rows affected != 1")
	}

	errMQ := mq.PublishComment(p)
	if errMQ != nil {
		// todo 消息冗余
	}

	return nil
}

func DeleteComment(commentId int64) error {
	tx, err, txFunc := mysql.BeginTransaction()
	if err != nil {
		return err
	}
	defer txFunc(err)

	var comment *model.Comment
	var rs sql.Result
	var rowsAffected int64

	comment, err = mysql.CheckCommentIsDeleted(tx, commentId)
	if err != nil {
		return err
	}

	rs, err = mysql.DeleteComment(tx, commentId)
	if err != nil {
		return err
	}
	rowsAffected, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("rows affected != 1")
	}

	rs, err = mysql.DecrCommentCountAtPost(tx, comment.PostId)
	if err != nil {
		return err
	}
	rowsAffected, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("rows affected != 1")
	}

	return nil
}

func GetComments(p *model.ParamComments) ([]*model.Comment, error) {
	// get root comments
	comments, err := mysql.GetRootCommentsFromPost(p)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
