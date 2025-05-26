package server

import (
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

func ConfigureCommandBus() {
	diContainer := di.Instance()
	commandBus := diContainer.Container.Get("shared.domain.commandbus").(command.Bus)
	commandHandlers := diContainer.GetByTag("command-handler")

	for _, handlerDef := range commandHandlers {
		handler := diContainer.Container.Get(handlerDef).(command.Handler)
		commandBus.Register(handler.SubscribedTo(), handler)
	}
}

func ConfigureQueryBus() {
	diContainer := di.Instance()
	queryBus := diContainer.Container.Get("shared.domain.querybus").(query.Bus)
	queryHandlers := diContainer.GetByTag("query-handler")

	for _, handlerDef := range queryHandlers {
		handler := diContainer.Container.Get(handlerDef).(query.Handler)
		queryBus.Register(handler.SubscribedTo(), handler)
	}
}

func ConfigureEventBus() {
	diContainer := di.Instance()
	eventBus := diContainer.Container.Get("shared.domain.eventbus").(event.Bus)
	eventHandlers := diContainer.GetByTag("event-handler")

	for _, handlerDef := range eventHandlers {
		handler := diContainer.Container.Get(handlerDef).(event.Handler)
		eventBus.Subscribe(handler.SubscribedTo(), handler)
	}
}
