package controller

import (
	"evergreen/biz"
	"evergreen/biz/logic"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CommunityListHandler(c *gin.Context) {
	list, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("get community list error:", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, list)
}

func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, biz.CodeInvalidParam)
		return
	}
	detail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("get community detail error:", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, detail)
}
