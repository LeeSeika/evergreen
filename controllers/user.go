package controllers

import (
	"evergreen/logic"
	"evergreen/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

func SingUpHandler(c *gin.Context) {
	var p model.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Sign up with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		} else {
			translate := errors.Translate(trans)
			c.JSON(http.StatusOK, gin.H{
				"msg": removeTopStruct(translate),
			})
		}
		return
	}

	if err := logic.SignUp(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
