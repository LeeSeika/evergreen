package model

type User struct {
	UserID   int64  `db:"user_id,string"`
	Username string `db:"username"`
	Password string `db:"password"`
}
