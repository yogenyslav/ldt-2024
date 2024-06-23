package _server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/config"
	_ "github.com/yogenyslav/ldt-2024/admin/docs"
	"github.com/yogenyslav/ldt-2024/admin/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/admin/internal/auth"
	ac "github.com/yogenyslav/ldt-2024/admin/internal/auth/controller"
	ah "github.com/yogenyslav/ldt-2024/admin/internal/auth/handler"
	"github.com/yogenyslav/ldt-2024/admin/internal/auth/middleware"
	"github.com/yogenyslav/ldt-2024/admin/internal/notification"
	nc "github.com/yogenyslav/ldt-2024/admin/internal/notification/controller"
	nh "github.com/yogenyslav/ldt-2024/admin/internal/notification/handler"
	nr "github.com/yogenyslav/ldt-2024/admin/internal/notification/repo"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization"
	oc "github.com/yogenyslav/ldt-2024/admin/internal/organization/controller"
	oh "github.com/yogenyslav/ldt-2024/admin/internal/organization/handler"
	or "github.com/yogenyslav/ldt-2024/admin/internal/organization/repo"
	"github.com/yogenyslav/ldt-2024/admin/internal/user"
	uc "github.com/yogenyslav/ldt-2024/admin/internal/user/controller"
	uh "github.com/yogenyslav/ldt-2024/admin/internal/user/handler"
	ur "github.com/yogenyslav/ldt-2024/admin/internal/user/repo"
	"github.com/yogenyslav/ldt-2024/admin/pkg/client"
	"github.com/yogenyslav/ldt-2024/admin/pkg/metrics"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	srvresp "github.com/yogenyslav/pkg/response"
	"github.com/yogenyslav/pkg/storage"
	"github.com/yogenyslav/pkg/storage/minios3"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Server сервер.
type Server struct {
	cfg      *config.Config
	app      *fiber.App
	pg       storage.SQLDatabase
	kc       *gocloak.GoCloak
	s3       minios3.S3
	exporter sdktrace.SpanExporter
	tracer   trace.Tracer
}

// New создает новый Server.
func New(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit:    1024 * 1024 * 1024,
		ErrorHandler: srvresp.NewErrorHandler(errStatus).Handler,
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.CorsOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders:     "Access-Control-Allow-Origin,Authorization,Origin,Accept,Content-Type,ngrok-skip-browser-warning",
		AllowCredentials: true,
	}))
	app.Use(recovermw.New())
	app.Use(otelfiber.Middleware())
	app.Use("/admin/swagger/*", swagger.HandlerDefault)

	exporter := tracing.MustNewExporter(context.Background(), cfg.Jaeger.URL())
	provider := tracing.MustNewTraceProvider(exporter, "admin")
	otel.SetTracerProvider(provider)

	tracer := otel.Tracer("admin")

	s3 := minios3.MustNew(cfg.S3, tracer)
	if err := s3.CreateBuckets(context.Background()); err != nil {
		log.Panic().Err(err).Msg("failed to create buckets")
	}

	return &Server{
		cfg:      cfg,
		app:      app,
		pg:       postgres.MustNew(cfg.Postgres, tracer),
		kc:       gocloak.NewClient(cfg.KeyCloak.URL),
		s3:       s3,
		exporter: exporter,
		tracer:   tracer,
	}
}

// Run запускает сервер.
func (s *Server) Run() {
	defer s.pg.Close()
	defer func() {
		if err := s.exporter.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown exporter")
		}
	}()
	defer func() {
		if err := s.app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown app")
		}
	}()

	apiClient, err := client.NewGrpcClient(s.cfg.API)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create grpc client")
	}
	defer func() {
		if err := apiClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close grpc client")
		}
	}()

	m := metrics.New()
	m.Collect()

	authController := ac.New(pb.NewAuthServiceClient(apiClient.GetConn()), s.cfg.Server.CipherKey, s.tracer)
	authHandler := ah.New(authController)
	auth.SetupAuthRoutes(s.app, authHandler)

	g := s.app.Group("/admin")
	g.Use(middleware.JWT(s.kc, s.cfg.KeyCloak.Realm, s.cfg.Server.CipherKey))

	userRepo := ur.New(s.pg)
	userController := uc.New(userRepo, s.cfg.KeyCloak, s.kc, s.tracer)
	userHandler := uh.New(userController, m, s.tracer)
	user.SetupUserRoutes(g, userHandler)

	organizationRepo := or.New(s.pg)
	organizationController := oc.New(organizationRepo, s.s3, pb.NewPredictorClient(apiClient.GetConn()), s.tracer)
	organizationHandler := oh.New(organizationController, m, s.tracer)
	organization.SetupOrganizationRoutes(g, organizationHandler)

	notificationRepo := nr.New(s.pg)
	notificationController := nc.New(notificationRepo, s.tracer)
	notificationHandler := nh.New(notificationController)
	notification.SetupNotificationRoutes(g, notificationHandler)

	go s.listen()
	go prom.HandlePrometheus(s.cfg.Prom)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("shutting down the server")
}

func (s *Server) listen() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	if err := s.app.Listen(addr); err != nil {
		log.Error().Err(err).Msg("failed to listen")
	}
}
