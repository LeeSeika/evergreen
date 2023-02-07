package router

import (
	"evergreen/controller"
	"evergreen/logger"
	"evergreen/middleware"
	"evergreen/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	if settings.Conf.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(logger.GinLogger(), logger.GinRecovery(true))

	engine.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "hello")
	})
	engine.POST("/signup", controller.SingUpHandler)
	engine.POST("/login", controller.LoginHandler)
	engine.GET("/ping", middleware.JWTAuthMiddleware(), controller.PingController)

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "page not found",
		})
	})
	return engine
}
