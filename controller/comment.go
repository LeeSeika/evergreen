package controller

import (
	"errors"
	"evergreen/biz"
	"evergreen/biz/logic"
	"evergreen/dao/mysql"
	"evergreen/model"
	"evergreen/util"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddCommentHandler(c *gin.Context) {
	p := model.Comment{}
	uid, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("get current user error", zap.Error(err))
		ResponseError(c, biz.CodeUserNotLogin)
		return
	}
	p.UserId = uid
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("add comment with invalid params", zap.Error(err))
		errorsV, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
		} else {
			errMsg := util.RemoveTopStruct(errorsV.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeInvalidParam, errMsg)
		}
		return
	}
	err = logic.AddComment(&p)
	if err != nil {
		errorsV, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeServerBusy)
		} else {
			errMsg := util.RemoveTopStruct(errorsV.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeServerBusy, errMsg)
		}
		return
	}

	ResponseSuccess(c, p)
}

func DeleteCommentHandler(c *gin.Context) {
	commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	if err != nil {
		zap.L().Error("delete comment with invalid param", zap.Error(err))
		errorsV, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
		} else {
			errMsg := util.RemoveTopStruct(errorsV.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeInvalidParam, errMsg)
		}
		return
	}

	err = logic.DeleteComment(commentId)
	if err != nil {
		errorsV, ok := err.(validator.ValidationErrors)
		code := biz.CodeServerBusy
		if errors.Is(err, mysql.ErrorCommentDeleted) {
			code = biz.CodeCommentDeleted
		}
		if !ok {
			ResponseError(c, code)
		} else {
			errMsg := util.RemoveTopStruct(errorsV.Translate(util.Trans))
			ResponseErrorWithMsg(c, code, errMsg)
		}
		return
	}
	ResponseSuccess(c, nil)
}

func GetCommentInOrder(c *gin.Context) {
	p := model.ParamComments{
		ParamCommentsInOrder: &model.ParamCommentsInOrder{
			Page:  model.DefaultCommentPageValue,
			Size:  model.DefaultPostSizeValue,
			Order: model.OrderByScore,
		},
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("get comments with invalid params", zap.Error(err))
		errorsV, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
		} else {
			errMsg := util.RemoveTopStruct(errorsV.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeInvalidParam, errMsg)
		}
		return
	}
	comments, err := logic.GetComments(&p)
	if err != nil {
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, comments)
}
