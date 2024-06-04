package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/controller"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/handler"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	srvresp "github.com/yogenyslav/pkg/response"
	"github.com/yogenyslav/pkg/storage"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	orchestratorAddr := "api:9999"
	authConn, err := grpc.NewClient(orchestratorAddr, grpcOpts...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to orchestrator")
	}
	defer func() {
		if err = authConn.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to properly close grpc connection")
		}
	}()

	authController := controller.New(authConn, s.tracer)
	authHandler := handler.New(authController, s.tracer)
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
