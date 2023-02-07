package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"evergreen/model"
)

const secret = "leeseika/evergreen"

var (
	ErrorUserExist       = errors.New("user exists")
	ErrorUserNotFound    = errors.New("user not found")
	ErrorInvalidPassword = errors.New("invalid password")
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 想数据库中插入一条新的用户记录
func InsertUser(user *model.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *model.User) error {
	originalPwd := user.Password
	sqlStr := "select user_id, username, password from user where username = ?"

	if err := db.Get(user, sqlStr, user.Username); err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotFound
		}
		return err
	}
	if user.Password != encryptPassword(originalPwd) {
		return ErrorInvalidPassword
	}

	return nil
}
