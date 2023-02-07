package controllers

import (
	"errors"
	"evergreen/dao/mysql"
	"evergreen/logic"
	"evergreen/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SingUpHandler(c *gin.Context) {
	var p model.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Sign up with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
		} else {
			errMsg := removeTopStruct(errors.Translate(trans))
			responseErrorWithMsg(c, CodeInvalidParam, errMsg)
		}
		return
	}

	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("Sign up failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExists)
		} else {
			ResponseError(c, CodeServerBusy)
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
			ResponseError(c, CodeInvalidParam)
		} else {
			msg := removeTopStruct(errors.Translate(trans))
			responseErrorWithMsg(c, CodeInvalidParam, msg)
		}
		return
	}

	if err := logic.Login(&p); err != nil {
		zap.L().Error("Log in failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotFound) {
			ResponseError(c, CodeUserNotFound)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}

	ResponseSuccess(c, nil)
}
