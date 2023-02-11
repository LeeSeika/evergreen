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

	group := engine.Group("/api/v1")

	group.POST("/signup", controller.SingUpHandler)
	group.POST("/login", controller.LoginHandler)

	group.GET("/ping", controller.PingController)
	group.Use(middleware.JWTAuthMiddleware())
	group.GET("/community/list", controller.CommunityListHandler)
	group.GET("/community/:id", controller.CommunityDetailHandler)
	group.GET("/post/:id", controller.GetPostDetailHandler)
	group.GET("/post/list", controller.GetPostListHandler)
	group.POST("/post/create", controller.CreatePostHandler)

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "page not found",
		})
	})
	return engine
}
