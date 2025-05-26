package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	extractionhandlers "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/cli/handler"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server"
	userhandlers "github.com/rubenbupe/recipe-video-parser/internal/users/platform/cli/handler"
)

var diContainer = di.Instance()

func main() {
	server.ConfigureCommandBus()
	server.ConfigureQueryBus()
	server.ConfigureEventBus()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancelar el contexto si se recibe una señal de interrupción
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nCancelando operación...")
		cancel()
		os.Exit(1)
	}()

	if len(os.Args) < 2 {
		fmt.Println("Se requiere un comando: create-user, update-api-key, get-user, get-user-summary")
		fmt.Println("Uso: cli <comando> [opciones]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create-user":
		createUserCmd(ctx, os.Args[2:])
	case "update-api-key":
		updateApiKeyCmd(ctx, os.Args[2:])
	case "get-user":
		getUserCmd(ctx, os.Args[2:])
	case "get-user-summary":
		getExtractionsSummaryCmd(ctx, os.Args[2:])
	case "extract-recipe":
		extractRecipeCmd(ctx, os.Args[2:])
	default:
		fmt.Printf("Comando desconocido: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func createUserCmd(ctx context.Context, args []string) {
	createHandler := diContainer.Container.Get("users.infrastructure.cli.create").(userhandlers.CreateUserHandler)

	if len(args) < 1 {
		fmt.Println("Uso: cli create-user <username>")
		os.Exit(1)
	}
	userName := args[0]
	userId := uuid.New().String()
	apiKey := uuid.New().String()

	err := createHandler(
		ctx,
		userhandlers.CreateUserInput{
			ID:        userId,
			Name:      userName,
			ApiKey:    apiKey,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	)
	if err != nil {
		fmt.Printf("Error al crear usuario: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Usuario creado: %s (%s)\n", userName, userId)
	os.Exit(0)
}

func updateApiKeyCmd(ctx context.Context, args []string) {
	updateHandler := diContainer.Container.Get("users.infrastructure.cli.updateapikey").(userhandlers.UpdateApiKeyHandler)

	if len(args) < 1 {
		fmt.Println("Uso: cli update-api-key <username>")
		os.Exit(1)
	}
	username := args[0]
	apikey := uuid.New().String()

	err := updateHandler(
		ctx,
		userhandlers.UpdateApiKeyInput{
			Name:   username,
			ApiKey: apikey,
		},
	)
	if err != nil {
		fmt.Printf("Error al actualizar apiKey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("apiKey actualizada para %s: %s\n", username, apikey)
	os.Exit(0)
}

func getUserCmd(ctx context.Context, args []string) {
	getHandler := diContainer.Container.Get("users.infrastructure.cli.get").(userhandlers.GetUserHandler)

	if len(args) < 1 {
		fmt.Println("Uso: cli get-user <username>")
		os.Exit(1)
	}
	userName := args[0]

	result, err := getHandler(
		ctx,
		userhandlers.GetUserInput{
			Name: userName,
		},
	)
	if err != nil {
		fmt.Printf("Error al obtener usuario: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Usuario: %s\nNombre: %s\nApiKey: %s\nCreado: %s\n", result.ID, result.Name, result.ApiKey, result.CreatedAt)
	os.Exit(0)
}

func getExtractionsSummaryCmd(ctx context.Context, args []string) {
	userGetHandler := diContainer.Container.Get("users.infrastructure.cli.get").(userhandlers.GetUserHandler)
	extractionHandler := diContainer.Container.Get("recipes.infrastructure.cli.get").(extractionhandlers.GetExtractionHandler)

	if len(args) < 1 {
		fmt.Println("Uso: cli get-user-summary <user_name>")
		os.Exit(1)
	}
	userName := args[0]

	// Buscar usuario por nombre
	userResult, err := userGetHandler(ctx, userhandlers.GetUserInput{Name: userName})
	if err != nil {
		fmt.Printf("Error al buscar usuario '%s': %v\n", userName, err)
		os.Exit(1)
	}
	userId := userResult.ID

	result, err := extractionHandler(ctx, extractionhandlers.GetExtractionInput{UserID: userId})
	if err != nil {
		fmt.Printf("Error al obtener extracciones: %v\n", err)
		os.Exit(1)
	}

	type summary struct {
		Month                string
		Count                int
		PromptTokenCount     int
		CandidatesTokenCount int
		TotalTokenCount      int
	}
	summaries := make(map[string]*summary)

	for _, extraction := range result {
		createdAt, err := time.Parse(time.RFC3339, extraction.CreatedAt)
		if err != nil {
			continue // ignora fechas mal formateadas
		}
		monthYear := createdAt.Format("01-2006")

		promptTokens := 0
		candidateTokens := 0
		if extraction.Metadata != "" {
			var meta struct {
				PromptTokenCount     int `json:"promptTokenCount"`
				CandidatesTokenCount int `json:"candidatesTokenCount"`
			}
			if err := json.Unmarshal([]byte(extraction.Metadata), &meta); err == nil {
				promptTokens = meta.PromptTokenCount
				candidateTokens = meta.CandidatesTokenCount
			}
		}

		sum, ok := summaries[monthYear]
		if !ok {
			sum = &summary{Month: monthYear}
			summaries[monthYear] = sum
		}
		sum.Count++
		sum.PromptTokenCount += promptTokens
		sum.CandidatesTokenCount += candidateTokens
		sum.TotalTokenCount += promptTokens + candidateTokens
	}

	fmt.Printf("Resumen de extracciones para el usuario '%s' (ID: %s) por mes-año:\n", userName, userId)
	fmt.Printf("%-10s | %-12s | %-15s | %-18s | %-15s\n", "Mes-Año", "Extracciones", "PromptTokens", "CandidateTokens", "TotalTokens")
	fmt.Println("------------------------------------------------------------------------------------------")
	for _, sum := range summaries {
		fmt.Printf("%-10s | %-12d | %-15d | %-18d | %-15d\n",
			sum.Month, sum.Count, sum.PromptTokenCount, sum.CandidatesTokenCount, sum.TotalTokenCount)
	}
	os.Exit(0)
}

func extractRecipeCmd(ctx context.Context, args []string) {
	extractHandler := diContainer.Container.Get("recipes.infrastructure.cli.extract").(extractionhandlers.ExtractRecipeHandler)

	if len(args) < 1 {
		fmt.Println("Uso: cli extract-recipe <url>")
		os.Exit(1)
	}
	url := args[0]

	err := extractHandler(ctx, extractionhandlers.ExtractRecipeInput{Url: url})
	if err != nil {
		fmt.Printf("Error al extraer receta: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
