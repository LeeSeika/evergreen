package controller

import (
	"evergreen/biz"
	"evergreen/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshTokenHandler(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseError(c, biz.CodeAuthNotFound)
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 3)
	if !(len(parts) == 3 && parts[0] == "Bearer") {
		ResponseError(c, biz.CodeInvalidAuth)
		c.Abort()
		return
	}

	aToken, rToken, err := jwt.RefreshToken(parts[1], parts[2])
	if err != nil {
		ResponseError(c, biz.CodeInvalidAuth)
		c.Abort()
		return
	}

	type Authorization struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh"`
	}

	authMsg := Authorization{
		Token:        aToken,
		RefreshToken: rToken,
	}

	ResponseSuccess(c, authMsg)
}
