package middleware

import (
	"evergreen/biz"
	"evergreen/controller"
	"evergreen/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, biz.CodeAuthNotFound)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, biz.CodeInvalidAuth)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			// controller.ResponseError(c, biz.CodeInvalidToken)
			controller.ResponseErrorWithHttpCode(c, biz.CodeInvalidToken, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set(controller.ContextUserIDKey, mc.UserID)
		c.Next()
	}
}
