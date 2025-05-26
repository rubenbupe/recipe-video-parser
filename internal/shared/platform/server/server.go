package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	recipesroutes "github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/server/routes"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server/middleware/logging"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/server/middleware/recovery"
	statusroutes "github.com/rubenbupe/recipe-video-parser/internal/status/platform/server/routes"

	// extractionsroutes "github.com/rubenbupe/recipe-video-parser/internal/extractions/platform/server/routes"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	shutdownTimeout time.Duration
}

func New(ctx context.Context, host string, port uint, shutdownTimeout time.Duration, commandBus command.Bus) (context.Context, Server) {
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", host, port),

		shutdownTimeout: shutdownTimeout,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

func (s *Server) registerRoutes() {
	s.engine.Use(recovery.Middleware(), logging.Middleware())

	status := s.engine.Group("/status")
	recipes := s.engine.Group("/recipes")
	// users := apiV1.Group("/recipes")

	statusroutes.Register(status)
	recipesroutes.Register(recipes)
	// extractionsroutes.Register(users)
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
