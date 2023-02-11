package controller

import (
	"errors"
	"evergreen/biz"
	"evergreen/biz/logic"
	"evergreen/dao/mysql"
	"evergreen/model"
	"evergreen/util"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const ContextUserIDKey = "userID"

func SingUpHandler(c *gin.Context) {
	var p model.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Sign up with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
		} else {
			errMsg := util.RemoveTopStruct(errors.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeInvalidParam, errMsg)
		}
		return
	}

	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("Sign up failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, biz.CodeUserExists)
		} else {
			ResponseError(c, biz.CodeServerBusy)
		}
		return
	}

	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	var p model.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Log in with invalid params", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
		} else {
			msg := util.RemoveTopStruct(errors.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeInvalidParam, msg)
		}
		return
	}

	user, token, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("Log in failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotFound) {
			ResponseError(c, biz.CodeUserNotFound)
		} else {
			ResponseError(c, biz.CodeServerBusy)
		}
		return
	}

	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     token,
	})
}
