package recipe

import (
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
	"github.com/sarulabs/di/v2"

	usercreate "github.com/rubenbupe/recipe-video-parser/internal/users/application/create"
	userget "github.com/rubenbupe/recipe-video-parser/internal/users/application/get"
	userupdateapikey "github.com/rubenbupe/recipe-video-parser/internal/users/application/updateapikey"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	userssql "github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/sql"

	extractioncreate "github.com/rubenbupe/recipe-video-parser/internal/recipes/application/create"
	extractionget "github.com/rubenbupe/recipe-video-parser/internal/recipes/application/get"
	extractionsdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	extractionsql "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/storage/sql"
)

var Defs = []di.Def{
	// REPOSITORIES
	{
		Name: "users.domain.repository",
		Build: func(ctn di.Container) (interface{}, error) {
			conn := ctn.Get("shared.infrastructure.sqlconnection").(*storage.Connection)
			dbconfig := ctn.Get("shared.infrastructure.sqlconfig").(*storage.Dbconfig)
			return userssql.NewUserRepository(conn, dbconfig), nil
		},
	},
	{
		Name: "extractions.domain.repository",
		Build: func(ctn di.Container) (interface{}, error) {
			conn := ctn.Get("shared.infrastructure.sqlconnection").(*storage.Connection)
			dbconfig := ctn.Get("shared.infrastructure.sqlconfig").(*storage.Dbconfig)
			return extractionsql.NewExtractionRepository(conn, dbconfig), nil
		},
	},
	// USE CASES, COMMAND HANDLERS, AND EVENT HANDLERS
	{
		Name: "users.domain.create",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("users.domain.repository").(usersdomain.UserRepository)
			eventBus := ctn.Get("shared.domain.eventbus").(event.Bus)
			return usercreate.NewUserService(repo, eventBus), nil
		},
	},
	{
		Name: "users.domain.createcommandhandler",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("users.domain.create").(usercreate.UserService)
			return usercreate.NewUserCommandHandler(service), nil
		},
		Tags: []di.Tag{
			{Name: "command-handler"},
		},
	},

	{
		Name: "users.domain.get",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("users.domain.repository").(usersdomain.UserRepository)
			return userget.NewUserService(repo), nil
		},
	},
	{
		Name: "users.domain.getqueryhandler",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("users.domain.get").(userget.UserService)
			return userget.NewUserQueryHandler(service), nil
		},
		Tags: []di.Tag{
			{Name: "query-handler"},
		},
	},
	{
		Name: "users.domain.updateapikey",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("users.domain.repository").(usersdomain.UserRepository)
			return userupdateapikey.NewUserApiKeyService(repo, nil), nil
		},
	},
	{
		Name: "users.domain.updateapikeycommandhandler",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("users.domain.updateapikey").(userupdateapikey.UserApiKeyService)
			return userupdateapikey.NewUserApiKeyCommandHandler(service), nil
		},
		Tags: []di.Tag{
			{Name: "command-handler"},
		},
	},

	{
		Name: "extractions.domain.create",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("extractions.domain.repository").(extractionsdomain.ExtractionRepository)
			eventBus := ctn.Get("shared.domain.eventbus").(event.Bus)
			return extractioncreate.NewExtractionService(repo, eventBus), nil
		},
	},
	{
		Name: "extractions.domain.createcommandhandler",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("extractions.domain.create").(extractioncreate.ExtractionService)
			return extractioncreate.NewExtractionCommandHandler(service), nil
		},
		Tags: []di.Tag{
			{Name: "command-handler"},
		},
	},
	{
		Name: "extractions.domain.get",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("extractions.domain.repository").(extractionsdomain.ExtractionRepository)
			return extractionget.NewExtractionService(repo), nil
		},
	},
	{
		Name: "extractions.domain.getqueryhandler",
		Build: func(ctn di.Container) (interface{}, error) {
			service := ctn.Get("extractions.domain.get").(extractionget.ExtractionService)
			return extractionget.NewExtractionQueryHandler(service), nil
		},
		Tags: []di.Tag{
			{Name: "query-handler"},
		},
	},
}
