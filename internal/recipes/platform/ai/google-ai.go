package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/gallery"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/ai"
)

func uploadFileToGoogleAI(filePath string, config ai.Aiconfig) (string, error) {
	// Obtener el mime type y tamaño del archivo
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("could not stat file: %w", err)
	}
	numBytes := fileInfo.Size()

	// Obtener el mime type usando el contenido del archivo
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	mimeType := http.DetectContentType(buf[:n])
	file.Seek(0, 0)

	displayName := filepath.Base(filePath)

	// Paso 1: Iniciar la subida (obtener upload URL)
	startPayload := map[string]interface{}{
		"file": map[string]interface{}{
			"display_name": displayName,
		},
	}
	startPayloadBytes, _ := json.Marshal(startPayload)

	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/upload/v1beta/files?key="+config.ApiKey, bytes.NewBuffer(startPayloadBytes))
	if err != nil {
		return "", fmt.Errorf("could not create start upload request: %w", err)
	}
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Upload-Header-Content-Length", fmt.Sprintf("%d", numBytes))
	req.Header.Set("X-Goog-Upload-Header-Content-Type", mimeType)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not start upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("start upload failed: %s, body: %s", resp.Status, string(body))
	}

	uploadURL := resp.Header.Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		return "", fmt.Errorf("upload URL not found in response headers")
	}

	// Paso 2: Subir el archivo binario
	uploadReq, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return "", fmt.Errorf("could not create upload request: %w", err)
	}
	uploadReq.Header.Set("Content-Length", fmt.Sprintf("%d", numBytes))
	uploadReq.Header.Set("X-Goog-Upload-Offset", "0")
	uploadReq.Header.Set("X-Goog-Upload-Command", "upload, finalize")

	resp2, err := client.Do(uploadReq)
	if err != nil {
		return "", fmt.Errorf("could not upload file: %w", err)
	}
	defer resp2.Body.Close()

	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return "", fmt.Errorf("could not read upload response: %w", err)
	}

	if resp2.StatusCode != 200 {
		return "", fmt.Errorf("upload failed: %s, body: %s", resp2.Status, string(body))
	}

	var fileInfoResp map[string]interface{}
	if err := json.Unmarshal(body, &fileInfoResp); err != nil {
		return "", fmt.Errorf("could not parse upload response: %w", err)
	}

	fileObj, ok := fileInfoResp["file"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("file object not found in response")
	}
	fileURI, ok := fileObj["uri"].(string)
	if !ok {
		return "", fmt.Errorf("file uri not found in response")
	}

	return fileURI, nil
}

func ensureFileActive(fileUrl string, config ai.Aiconfig) error {
	client := &http.Client{}
	for i := 0; i < 120; i++ { // hasta 120 intentos (~60s)
		req, err := http.NewRequest("GET", fileUrl+"?key="+config.ApiKey, nil)
		if err != nil {
			return fmt.Errorf("could not create request: %w", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("could not request file state: %w", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("could not read response: %w", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status: %s, body: %s", resp.Status, string(body))
		}
		var res map[string]interface{}
		if err := json.Unmarshal(body, &res); err != nil {
			return fmt.Errorf("could not parse response: %w", err)
		}
		if state, ok := res["state"].(string); ok && state == "ACTIVE" {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("file did not become ACTIVE after waiting")
}

type AiResponse struct {
	Recipe   Recipe `json:"recipe"`
	Metadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
	} `json:"metadata"`
}

type Recipe struct {
	Title           string          `json:"title" validate:"required"`
	Description     string          `json:"description" validate:"required"`
	Servings        int             `json:"servings"`
	PrepTime        int             `json:"prep_time"`
	CookTime        int             `json:"cook_time"`
	TotalTime       int             `json:"total_time"`
	Ingredients     []Ingredient    `json:"ingredients" validate:"required"`
	Sections        []Section       `json:"sections" validate:"required"`
	Notes           string          `json:"notes"`
	NutritionalInfo NutritionalInfo `json:"nutritional_info"`
	Url             string          `json:"-"`
}
type Ingredient struct {
	Name     string `json:"name" validate:"required"`
	Quantity string `json:"quantity" validate:"required"`
	Unit     string `json:"unit" validate:"required"`
}
type Section struct {
	Instructions []Instruction `json:"instructions" validate:"required"`
}
type Instruction struct {
	Optional bool   `json:"optional" validate:"required"`
	Text     string `json:"text" validate:"required"`
}
type NutritionalInfo struct {
	Calories      int `json:"calories"`
	Protein       int `json:"protein"`
	Carbohydrates int `json:"carbohydrates"`
	Fats          int `json:"fats"`
	Fiber         int `json:"fiber"`
	Sugar         int `json:"sugar"`
}

func parseGoogleAIResponse(apiResponse string) (AiResponse, error) {
	var res map[string]interface{}
	if err := json.Unmarshal([]byte(apiResponse), &res); err != nil {
		return AiResponse{}, fmt.Errorf("could not parse response: %w", err)
	}
	// Extraer el texto JSON de la receta
	candidates, ok := res["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return AiResponse{}, fmt.Errorf("no candidates in response")
	}
	candidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		return AiResponse{}, fmt.Errorf("invalid candidate format")
	}
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return AiResponse{}, fmt.Errorf("invalid content format")
	}
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return AiResponse{}, fmt.Errorf("no parts in content")
	}
	var recipeJson string
	for _, p := range parts {
		if partMap, ok := p.(map[string]interface{}); ok {
			if text, ok := partMap["text"].(string); ok {
				recipeJson = text
				break
			}
		}
	}
	if recipeJson == "" {
		return AiResponse{}, fmt.Errorf("no recipe JSON found in response")
	}

	// Extraer metadatos
	metadata := make(map[string]interface{})
	if usage, ok := res["usageMetadata"].(map[string]interface{}); ok {
		metadata = usage
	}

	// Parsear y validar el JSON de la receta y los metadatos
	var aiResponse AiResponse
	if err := json.Unmarshal([]byte(recipeJson), &aiResponse.Recipe); err != nil {
		return AiResponse{}, fmt.Errorf("error parsing AI response JSON: %w", err)
	}
	if metadata != nil {
		if v, ok := metadata["promptTokenCount"].(float64); ok {
			aiResponse.Metadata.PromptTokenCount = int(v)
		}
		if v, ok := metadata["candidatesTokenCount"].(float64); ok {
			aiResponse.Metadata.CandidatesTokenCount = int(v)
		}
	}
	validate := validator.New()
	if err := validate.Struct(aiResponse); err != nil {
		return AiResponse{}, fmt.Errorf("validation error: %w", err)
	}
	for i, sec := range aiResponse.Recipe.Sections {
		if len(sec.Instructions) == 0 {
			return AiResponse{}, fmt.Errorf("la sección %d no tiene instrucciones", i+1)
		}
		for j, inst := range sec.Instructions {
			if inst.Text == "" {
				return AiResponse{}, fmt.Errorf("la instrucción %d de la sección %d está vacía", j+1, i+1)
			}
		}
	}
	return aiResponse, nil
}

func FormatToMarkdown(aiResponse AiResponse) string {
	recipe := aiResponse.Recipe
	var b strings.Builder

	b.WriteString("# " + recipe.Title + "\n\n")
	if recipe.Description != "" {
		b.WriteString("## Descripción\n")
		b.WriteString(recipe.Description + "\n\n")
	}
	b.WriteString(fmt.Sprintf("**Porciones:** %d\n", recipe.Servings))
	b.WriteString(fmt.Sprintf("**Tiempo de preparación:** %d min\n", recipe.PrepTime))
	b.WriteString(fmt.Sprintf("**Tiempo de cocción:** %d min\n", recipe.CookTime))
	b.WriteString(fmt.Sprintf("**Tiempo total:** %d min\n\n", recipe.TotalTime))

	b.WriteString("## Ingredientes\n")
	for _, ing := range recipe.Ingredients {
		b.WriteString(fmt.Sprintf("- %s %s %s\n", ing.Quantity, ing.Unit, ing.Name))
	}
	b.WriteString("\n")

	b.WriteString("## Instrucciones\n")
	for i, sec := range recipe.Sections {
		if len(recipe.Sections) > 1 {
			b.WriteString(fmt.Sprintf("### Sección %d\n", i+1))
		}
		for j, inst := range sec.Instructions {
			b.WriteString(fmt.Sprintf("%d. %s%s\n", j+1, inst.Text, func() string {
				if inst.Optional {
					return " _(opcional)_"
				} else {
					return ""
				}
			}()))
		}
		b.WriteString("\n")
	}

	if recipe.Notes != "" {
		b.WriteString("**Notas:** " + recipe.Notes + "\n\n")
	}

	b.WriteString("## Información nutricional (por cada 100g)\n")
	b.WriteString(fmt.Sprintf("- Calorías: %d kcal\n", recipe.NutritionalInfo.Calories))
	b.WriteString(fmt.Sprintf("- Proteínas: %d g\n", recipe.NutritionalInfo.Protein))
	b.WriteString(fmt.Sprintf("- Carbohidratos: %d g\n", recipe.NutritionalInfo.Carbohydrates))
	b.WriteString(fmt.Sprintf("- Grasas: %d g\n", recipe.NutritionalInfo.Fats))
	b.WriteString(fmt.Sprintf("- Fibra: %d g\n", recipe.NutritionalInfo.Fiber))
	b.WriteString(fmt.Sprintf("- Azúcares: %d g\n", recipe.NutritionalInfo.Sugar))

	if recipe.Url != "" {
		b.WriteString("\n[Ver receta original](" + recipe.Url + ")\n")
	}

	return b.String()
}

func AskModel(download gallery.DownloadResult, config ai.Aiconfig) (AiResponse, error) {
	filePath := download.FilePath
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join("tmp/dl", filePath)
	}

	fileURI, err := uploadFileToGoogleAI(filePath, config)
	if err != nil {
		return AiResponse{}, fmt.Errorf("could not upload file: %w", err)
	}

	if err := ensureFileActive(fileURI, config); err != nil {
		return AiResponse{}, fmt.Errorf("file not ACTIVE: %w", err)
	}

	prompt := ExtractRecipePrompt()

	payload := map[string]interface{}{
		"generation_config": map[string]interface{}{
			"response_mime_type": "application/json",
			"temperature":        config.Temperature,
		},
		"system_instruction": map[string]interface{}{
			"parts": []interface{}{
				map[string]interface{}{
					"text": prompt,
				},
			},
		},
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"file_data": map[string]interface{}{
							"mime_type": download.MimeType,
							"file_uri":  fileURI,
						},
					},
					map[string]interface{}{"text": download.Description},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return AiResponse{}, fmt.Errorf("could not marshal payload: %w", err)
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/" + config.Model + ":generateContent?key=" + config.ApiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return AiResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AiResponse{}, fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return AiResponse{}, fmt.Errorf("error: %s, body: %s", resp.Status, string(body))
	}

	parsedResponse, err := parseGoogleAIResponse(string(body))
	if err != nil {
		return AiResponse{}, fmt.Errorf("error parsing AI response: %w", err)
	}
	parsedResponse.Recipe.Url = download.Url
	return parsedResponse, nil
}
