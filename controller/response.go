package controller

import (
	"evergreen/biz"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code biz.RespCode `json:"code"`
	Msg  interface{}  `json:"msg"`
	Data interface{}  `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code biz.RespCode) {
	respData := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, respData)
}

func ResponseErrorWithMsg(c *gin.Context, code biz.RespCode, msg interface{}) {
	respData := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, respData)
}

func ResponseErrorWithHttpCode(c *gin.Context, code biz.RespCode, httpCode int) {
	respData := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(httpCode, respData)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	respData := &ResponseData{
		Code: biz.CodeSuccess,
		Msg:  biz.CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, respData)
}
