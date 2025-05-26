package di

import (
	"slices"

	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di/app"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di/recipe"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di/shared"
)

var defaultDefs = slices.Concat(app.Defs, recipe.Defs, shared.Defs)
