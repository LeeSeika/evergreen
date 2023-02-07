package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code RespCode    `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseError(c *gin.Context, code RespCode) {
	respData := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, respData)
}

func responseErrorWithMsg(c *gin.Context, code RespCode, msg interface{}) {
	respData := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, respData)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	respData := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, respData)
}
