package controller

import (
	"evergreen/biz"
	"evergreen/biz/logic"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @Summary return an array of all communities
// @Param Authorization header string true "Bearer token"
// @Produce application/json
// @Success 200 {array} []model.Community
// @Router /community/list [get]
func CommunityListHandler(c *gin.Context) {
	list, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("get community list error:", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, list)
}

// @Summary return a specified community detail by id
// @Param Authorization header string true "Bearer token"
// @Param id path int true "community id"
// @Produce application/json
// @Success 200 {object} model.CommunityDetail
// @Router /community/{id} [get]
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
