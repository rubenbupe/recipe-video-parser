package shared

import (
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/ai"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/bus/inmemory"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/gallery"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
	"github.com/sarulabs/di/v2"
)

var Defs = []di.Def{
	// BUSES
	{
		Name: "shared.domain.commandbus",
		Build: func(ctn di.Container) (interface{}, error) {
			return inmemory.NewCommandBus(), nil
		},
	},
	{
		Name: "shared.domain.querybus",
		Build: func(ctn di.Container) (interface{}, error) {
			return inmemory.NewQueryBus(), nil
		},
	},
	{
		Name: "shared.domain.eventbus",
		Build: func(ctn di.Container) (interface{}, error) {
			return inmemory.NewEventBus(), nil
		},
	},

	// DB
	{
		Name: "shared.infrastructure.sqlconfig",
		Build: func(ctn di.Container) (interface{}, error) {
			return storage.CreateConfig()
		},
	},
	{
		Name: "shared.infrastructure.sqlconnection",
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("shared.infrastructure.sqlconfig").(*storage.Dbconfig)
			return storage.CreateConnection("shared", cfg)
		},
	},

	// GALLERY
	{
		Name: "shared.infrastructure.galleryconfig",
		Build: func(ctn di.Container) (interface{}, error) {
			return gallery.CreateConfig()
		},
	},

	// AI
	{
		Name: "shared.infrastructure.aiconfig",
		Build: func(ctn di.Container) (interface{}, error) {
			return ai.CreateConfig()
		},
	},
}
