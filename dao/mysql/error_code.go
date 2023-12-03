package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("user exists")
	ErrorUserNotFound    = errors.New("user not found")
	ErrorInvalidPassword = errors.New("invalid password")
	ErrorInvalidID       = errors.New("invalid id")
	ErrorCommentDeleted  = errors.New("comment has been deleted")
)
