package router

import (
	"evergreen/controller"
	"evergreen/logger"
	"evergreen/middleware"
	"evergreen/settings"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	_ "evergreen/docs"
)

func Setup() *gin.Engine {
	if settings.Conf.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(logger.GinLogger(), logger.GinRecovery(true))

	engine.LoadHTMLFiles("./templates/index.html")
	engine.Static("/static", "./static")

	engine.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	engine.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "hello")
	})

	group := engine.Group("/api/v1")

	group.POST("/signup", controller.SingUpHandler)
	group.POST("/login", controller.LoginHandler)
	group.POST("/refresh", controller.RefreshTokenHandler)

	group.Use(middleware.JWTAuthMiddleware())
	group.GET("/ping", controller.PingController)

	group.GET("/community/list", controller.CommunityListHandler)
	group.GET("/community/:id", controller.CommunityDetailHandler)
	group.GET("/post/:id", controller.GetPostDetailHandler)
	group.GET("/post/list", controller.GetPostListHandler)
	group.GET("/post/list/order", controller.GetPostListInOrderHandler)

	group.POST("/post/vote", controller.PostVoteController)
	group.POST("/post/create", controller.CreatePostHandler)
	group.POST("/post/list/order", controller.GetCommunityPostListHandler)

	group.POST("/comment/add", controller.AddCommentHandler)
	group.GET("/comment/delete", controller.DeleteCommentHandler)

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "page not found",
		})
	})
	return engine
}
