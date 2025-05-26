package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rubenbupe/recipe-video-parser/internal/recipes/application/create"
	googleai "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/ai"
	gallerydl "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/gallery"
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

func extractWithUrl(aiconfig *ai.Aiconfig, commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Query("url")
		if url == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "url query parameter is required"})
			return
		}

		id := uuid.New().String()
		res, err := googleai.AskModelWithUrl(url, *aiconfig)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process url with AI model: " + err.Error()})
			return
		}

		handleExtractionResult(ctx, res, id, commandBus)
	}
}

func extractWithDownloadedFile(galleryconfig *gallery.Galleryconfig, aiconfig *ai.Aiconfig, commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Query("url")
		if url == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "url query parameter is required"})
			return
		}

		id := uuid.New().String()
		downloaded, err := gallerydl.DownloadFile(url, id, galleryconfig.DownloadDir, galleryconfig.PublicUrl)
		if err != nil {
			// Log the error for debugging purposes
			fmt.Printf("Error downloading file: %s\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download file: " + err.Error()})
			return
		}

		res, err := googleai.AskModelWithFile(downloaded, *aiconfig)

		gallerydl.RemoveFile(downloaded.FilePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process file with AI model: " + err.Error()})
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
			extractWithUrl(aiconfig, commandBus)(ctx)
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
