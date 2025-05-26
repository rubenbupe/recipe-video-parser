package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di"
	handlers "github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server/handler"
	middleware "github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server/middleware"
	userssql "github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/sql"
)

func Register(router *gin.RouterGroup) {
	diContainer := di.Instance()

	extractController := diContainer.Container.Get("recipes.infrastructure.controller.extract").(handlers.Handler)
	userRepo := diContainer.Container.Get("users.domain.repository").(*userssql.UserRepository)

	router.GET("/extract", middleware.AuthMiddleware(*userRepo), extractController)
}
