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

// @Summary create a post
// @Param Authorization header string true "Bearer token"
// @Param request body model.Post true "post model"
// @Produce application/json
// @Success 200
// @Router /post/create [post]
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

// @Summary get a post
// @Param Authorization header string true "Bearer token"
// @Param id path int true "post id"
// @Produce application/json
// @Success 200 {object} model.ApiPostDetail
// @Router /post/{id} [get]
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

// @Summary get post list
// @Param Authorization header string true "Bearer token"
// @Param page query int false "page"
// @Param size query int false "size"
// @Produce application/json
// @Success 200 {array} []model.ApiPostDetail
// @Router /post/list [get]
func GetPostListHandler(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = model.DefaultPostPageValue
	}
	size, err := strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil {
		size = model.DefaultPostSizeValue
	}
	postDetailList, err := logic.GetPostDetailList(page, size)
	if err != nil {
		zap.L().Error("get post list error", zap.Error(err))
		ResponseError(c, biz.CodeServerBusy)
		return
	}
	ResponseSuccess(c, postDetailList)
}

// @Summary get a post by given order
// @Param Authorization header string true "Bearer token"
// @Param page query int false "page"
// @Param size query int flase "size"
// @Param order query string false "order"
// @Produce application/json
// @Success 200 {array} []model.ApiPostDetail
// @Router /post/list/order [get]
func GetPostListInOrderHandler(c *gin.Context) {
	// default value
	p := model.ParamPostListInOrder{
		Page:  model.DefaultPostPageValue,
		Size:  model.DefaultPostSizeValue,
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

// @Summary get post list from a specified community id by order
// @Param Authorization header string true "Bearer token"
// @Produce application/json
// @Accept application/json
// @Param request body model.ParamCommunityPostList true "param json"
// @Success 200 {array} []model.ApiPostDetail
// @Router /post/list/order [post]
func GetCommunityPostListHandler(c *gin.Context) {
	p := model.ParamCommunityPostList{
		ParamPostListInOrder: &model.ParamPostListInOrder{
			Page:  model.DefaultPostPageValue,
			Size:  model.DefaultPostSizeValue,
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
