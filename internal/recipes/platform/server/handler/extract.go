package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rubenbupe/recipe-video-parser/internal/recipes/application/create"
	googleai "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/ai"
	clihandlers "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/cli/handler"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/ai"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/gallery"
	"github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

func handleExtractionResult(ctx *gin.Context, res googleai.AiResponse, id string, commandBus command.Bus) {
	jsonRecipe, err := toJSONString(res.Recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize recipe to JSON"})
		return
	}
	jsonMetadata, err := toJSONString(res.Metadata)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize metadata to JSON"})
		return
	}

	// Obtener usuario autenticado del contexto
	userVal, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	commandBus.Dispatch(
		ctx,
		create.NewExtractionCommand(
			id,
			user.Id.String(),
			jsonRecipe,
			jsonMetadata,
			time.Now().Format(time.RFC3339),
		),
	)

	if ctx.GetHeader("Accept") == "text/markdown" {
		ctx.Header("Content-Type", "text/markdown")
		markdownResponse := googleai.FormatToMarkdown(res)
		ctx.String(http.StatusOK, markdownResponse)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func extractWithUrl(galleryconfig *gallery.Galleryconfig, aiconfig *ai.Aiconfig, commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Query("url")
		res, id, err := clihandlers.ExtractRecipe(url, galleryconfig, aiconfig)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		handleExtractionResult(ctx, res, id, commandBus)
	}
}

func extractWithDownloadedFile(galleryconfig *gallery.Galleryconfig, aiconfig *ai.Aiconfig, commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Query("url")
		res, id, err := clihandlers.ExtractRecipe(url, galleryconfig, aiconfig)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		handleExtractionResult(ctx, res, id, commandBus)
	}
}

func needsDownload(url string) bool {
	// Considera que si la url contiene "youtube.com" o "youtu.be" no necesita descarga
	return !(strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be"))
}

func ExtractHandler(galleryconfig *gallery.Galleryconfig, aiconfig *ai.Aiconfig, commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Query("url")
		if needsDownload(url) {
			extractWithDownloadedFile(galleryconfig, aiconfig, commandBus)(ctx)
		} else {
			extractWithUrl(galleryconfig, aiconfig, commandBus)(ctx)
		}
	}
}

// Helper para serializar a string JSON
func toJSONString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
