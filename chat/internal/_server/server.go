package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/config"
	_ "github.com/yogenyslav/ldt-2024/chat/docs"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth"
	ac "github.com/yogenyslav/ldt-2024/chat/internal/auth/controller"
	ah "github.com/yogenyslav/ldt-2024/chat/internal/auth/handler"
	"github.com/yogenyslav/ldt-2024/chat/pkg/client"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	srvresp "github.com/yogenyslav/pkg/response"
	"github.com/yogenyslav/pkg/storage"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Server main struct that holds dependencies.
type Server struct {
	cfg      *config.Config
	app      *fiber.App
	pg       storage.SQLDatabase
	tracer   trace.Tracer
	exporter sdktrace.SpanExporter
}

// New creates a new Server instance.
func New(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: srvresp.NewErrorHandler(errStatus).Handler,
		AppName:      "GKU Chat",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Server.CorsOrigins, ","),
	}))
	app.Use(otelfiber.Middleware())
	app.Use(recovermw.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	exporter := tracing.MustNewExporter(context.Background(), cfg.Jaeger.URL())
	provider := tracing.MustNewTraceProvider(exporter, "chat")
	otel.SetTracerProvider(provider)

	tracer := otel.Tracer("chat")

	return &Server{
		cfg:      cfg,
		app:      app,
		pg:       postgres.MustNew(cfg.Postgres, tracer),
		exporter: exporter,
		tracer:   tracer,
	}
}

// Run setups the server and starts it.
func (s *Server) Run() {
	defer s.pg.Close()
	defer func() {
		if err := s.exporter.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown exporter")
		}
	}()
	defer func() {
		if err := s.app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server")
		}
	}()

	apiClient, err := client.NewGrpcClient(s.cfg.API)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create grpc client")
	}
	defer func() {
		if err = apiClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close grpc client")
		}
	}()

	authController := ac.New(apiClient.GetConn(), s.tracer)
	authHandler := ah.New(authController, s.tracer)
	auth.SetupAuthRoutes(s.app, authHandler)

	go s.listen()
	go prom.HandlePrometheus(s.cfg.Prometheus)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	log.Info().Msg("shutting down the server")
}

func (s *Server) listen() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	if err := s.app.Listen(addr); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
