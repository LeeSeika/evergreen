package controller

import (
	"evergreen/biz"
	"evergreen/biz/logic"
	"evergreen/model"
	"evergreen/util"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageValue = 1
	defaultSizeValue = 10
)

func CreatePostHandler(c *gin.Context) {
	var p model.Post
	uid, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("get current user error", zap.Error(err))
		ResponseError(c, biz.CodeUserNotLogin)
		return
	}
	p.AuthorID = uid

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, biz.CodeInvalidParam)
		} else {
			errMsg := util.RemoveTopStruct(errors.Translate(util.Trans))
			ResponseErrorWithMsg(c, biz.CodeInvalidParam, errMsg)
		}
		return
	}

	err = logic.CreatePost(&p)
	if err != nil {
		zap.L().Error("create post error", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, biz.CodeInvalidParam)
		return
	}
	postDetail, err := logic.GetPostDetailByID(postID)
	if err != nil {
		zap.L().Error("get post detail error", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, postDetail)
}

func GetPostListHandler(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = defaultPageValue
	}
	size, err := strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil {
		size = defaultSizeValue
	}
	postDetailList, err := logic.GetPostDetailList(page, size)
	if err != nil {
		zap.L().Error("get post list error", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, postDetailList)
}

func GetPostListInOrderHandler(c *gin.Context) {
	// default value
	p := model.ParamPostListInOrder{
		Page:  1,
		Size:  10,
		Order: model.OrderByScore,
	}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		zap.L().Error("get post list in order with invalid params", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	apiPostList, err := logic.GetPostListInOrder(&p)
	if err != nil {
		zap.L().Error("get post list in order error", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, apiPostList)
}

func GetCommunityPostListHandler(c *gin.Context) {
	p := model.ParamCommunityPostList{
		ParamPostListInOrder: &model.ParamPostListInOrder{
			Page:  1,
			Size:  10,
			Order: model.OrderByTime,
		},
		CommunityID: 0,
	}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		zap.L().Error("get post list from community in order with invalid params", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	postList, err := logic.GetCommunityPostList(&p)
	if err != nil {
		zap.L().Error("get community post list in order error", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, postList)
}
