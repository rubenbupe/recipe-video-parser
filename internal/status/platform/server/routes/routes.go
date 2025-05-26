package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di"
	handlers "github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server/handler"
)

func Register(router *gin.RouterGroup) {

	diContainer := di.Instance()

	getController := diContainer.Container.Get("status.infrastructure.controller.check").(handlers.Handler)
	print("Registering status routes")
	router.GET("/", getController)
}
