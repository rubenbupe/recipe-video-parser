package bootstrap

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/di"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

func Run() error {
	var cfg config
	err := envconfig.Process("APP", &cfg)
	if err != nil {
		return err
	}

	server.ConfigureCommandBus()
	server.ConfigureQueryBus()
	server.ConfigureEventBus()

	commandBus := di.Instance().Container.Get("shared.domain.commandbus").(command.Bus)

	ctx, srv := server.New(context.Background(), cfg.Host, cfg.Port, cfg.ShutdownTimeout, commandBus)
	return srv.Run(ctx)
}

type config struct {
	// Server configuration
	Host            string        `default:"0.0.0.0"`
	Port            uint          `default:"8080"`
	ShutdownTimeout time.Duration `default:"10s"`
}
