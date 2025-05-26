package app

import (
	recipesclihandlers "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/cli/handler"
	recipeshandlers "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/server/handler"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/ai"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/gallery"
	statushandlers "github.com/rubenbupe/recipe-video-parser/internal/status/platform/server/handler"
	usershandlers "github.com/rubenbupe/recipe-video-parser/internal/users/platform/cli/handler"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
	"github.com/rubenbupe/recipe-video-parser/kit/query"
	"github.com/sarulabs/di/v2"
)

var Defs = []di.Def{
	// STATUS (HTTP)
	{
		Name: "status.infrastructure.controller.check",
		Build: func(ctn di.Container) (interface{}, error) {
			return statushandlers.CheckHandler(), nil
		},
	},
	// USERS (CLI)
	{
		Name: "users.infrastructure.cli.create",
		Build: func(ctn di.Container) (interface{}, error) {
			commandBus := ctn.Get("shared.domain.commandbus").(command.Bus)
			return usershandlers.CreateHandler(commandBus), nil
		},
	},
	{
		Name: "users.infrastructure.cli.get",
		Build: func(ctn di.Container) (interface{}, error) {
			queryBus := ctn.Get("shared.domain.querybus").(query.Bus)
			return usershandlers.CreateGetUserHandler(queryBus), nil
		},
	},
	{
		Name: "users.infrastructure.cli.updateapikey",
		Build: func(ctn di.Container) (interface{}, error) {
			commandBus := ctn.Get("shared.domain.commandbus").(command.Bus)
			return usershandlers.CreateUpdateApiKeyHandler(commandBus), nil
		},
	},

	// RECIPES (HTTP)
	{
		Name: "recipes.infrastructure.controller.extract",
		Build: func(ctn di.Container) (interface{}, error) {
			galleryConfig := ctn.Get("shared.infrastructure.galleryconfig").(*gallery.Galleryconfig)
			aiConfig := ctn.Get("shared.infrastructure.aiconfig").(*ai.Aiconfig)

			commandBus := ctn.Get("shared.domain.commandbus").(command.Bus)

			return recipeshandlers.ExtractHandler(galleryConfig, aiConfig, commandBus), nil
		},
	},

	// RECIPES (CLI)
	{
		Name: "recipes.infrastructure.cli.get",
		Build: func(ctn di.Container) (interface{}, error) {
			queryBus := ctn.Get("shared.domain.querybus").(query.Bus)
			return recipesclihandlers.CreateGetExtractionsHandler(queryBus), nil
		},
	},
}
