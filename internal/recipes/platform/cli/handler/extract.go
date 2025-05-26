package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/ai"
	gallerydl "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/gallery"
	sharedai "github.com/rubenbupe/recipe-video-parser/internal/shared/platform/ai"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/gallery"
)

type ExtractRecipeInput struct {
	Url string
}

type ExtractRecipeHandler func(context.Context, ExtractRecipeInput) error

// Lógica compartida para extracción de receta
func ExtractRecipe(url string, galleryConfig *gallery.Galleryconfig, aiConfig *sharedai.Aiconfig) (ai.AiResponse, string, error) {
	if url == "" {
		return ai.AiResponse{}, "", fmt.Errorf("url is required")
	}
	id := uuid.New().String()
	var res ai.AiResponse
	var err error
	if needsDownload(url) {
		downloaded, errDownload := gallerydl.DownloadFile(url, id, galleryConfig.DownloadDir, galleryConfig.PublicUrl)
		if errDownload != nil {
			return ai.AiResponse{}, id, fmt.Errorf("failed to download file: %w", errDownload)
		}
		res, err = ai.AskModelWithFile(downloaded, *aiConfig)
		gallerydl.RemoveFile(downloaded.FilePath)
	} else {
		res, err = ai.AskModelWithUrl(url, *aiConfig)
	}
	if err != nil {
		return ai.AiResponse{}, id, fmt.Errorf("failed to extract recipe: %w", err)
	}
	return res, id, nil
}

func NewExtractRecipeHandler(galleryConfig *gallery.Galleryconfig, aiConfig *sharedai.Aiconfig) ExtractRecipeHandler {
	return func(ctx context.Context, input ExtractRecipeInput) error {
		res, _, err := ExtractRecipe(input.Url, galleryConfig, aiConfig)
		if err != nil {
			return err
		}
		jsonRecipe, err := toJSONString(res.Recipe)
		if err != nil {
			return fmt.Errorf("failed to serialize recipe to JSON: %w", err)
		}
		fmt.Println(jsonRecipe)
		return nil
	}
}

func needsDownload(url string) bool {
	return !(contains(url, "youtube.com") || contains(url, "youtu.be"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func toJSONString(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
