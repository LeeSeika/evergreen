package logic

import (
	"evergreen/dao/mysql"
	"evergreen/model"
	"evergreen/pkg/jwt"
	"evergreen/pkg/snowflake"
)

func SignUp(p *model.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &model.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *model.ParamLogin) (*model.User, string, error) {
	user := &model.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, "", err
	}
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}
