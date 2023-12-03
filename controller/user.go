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

// @Summary sign up a user
// @Param request body model.ParamSignUp true "param sign up"
// @Success 200
// @Router /signup [post]
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

// @Summary sign up a user
// @Param request body model.ParamLogin true "param sign up"
// @Accept application/json
// @Produce application/json
// @Success 200 {object} controller.LoginHandler.loginResponseMsg
// @Router /login [post]
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

	user, aToken, rToken, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("Log in failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotFound) {
			ResponseError(c, biz.CodeUserNotFound)
		} else {
			ResponseError(c, biz.CodeServerBusy)
		}
		return
	}

	type loginResponseMsg struct {
		UserID       string `json:"user_id"`
		Username     string `json:"username"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh"`
	}

	loginMsg := loginResponseMsg{
		UserID:       fmt.Sprintf("%d", user.UserID),
		Username:     user.Username,
		Token:        aToken,
		RefreshToken: rToken,
	}

	ResponseSuccess(c, loginMsg)
}
