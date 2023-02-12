package controller

import (
	"evergreen/biz"
	"evergreen/biz/logic"
	"evergreen/model"
	"evergreen/util"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteController(c *gin.Context) {
	p := model.ParamVoteData{}
	if err := c.ShouldBindJSON(&p); err != nil {
		errors, ok := err.(validator.ValidationErrors)
		zap.L().Error("vote for post failed", zap.Error(err))
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
			return
		}
		translate := errors.Translate(util.Trans)
		ResponseErrorWithMsg(c, biz.CodeInvalidParam, util.RemoveTopStruct(translate))
		return
	}
	uid, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("vote for post failed", zap.Error(err))
		ResponseError(c, biz.CodeInvalidParam)
		return
	}
	err = logic.VoteForPost(uid, &p)
	if err != nil {
		zap.L().Error("vote for post failed", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
