package routes

import (
	"evergreen/controllers"
	"evergreen/logger"
	"evergreen/settings"

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

	engine.POST("/signup", controllers.SingUpHandler)
	return engine
}
