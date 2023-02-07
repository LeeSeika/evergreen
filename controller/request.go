package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("user not login")

func getCurrentUser(c *gin.Context) (int64, error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	userID, ok := uid.(int64)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	return userID, nil
}
