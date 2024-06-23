package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/config"
	_ "github.com/yogenyslav/ldt-2024/chat/docs"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth"
	ac "github.com/yogenyslav/ldt-2024/chat/internal/auth/controller"
	ah "github.com/yogenyslav/ldt-2024/chat/internal/auth/handler"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/middleware"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat"
	cc "github.com/yogenyslav/ldt-2024/chat/internal/chat/controller"
	ch "github.com/yogenyslav/ldt-2024/chat/internal/chat/handler"
	cr "github.com/yogenyslav/ldt-2024/chat/internal/chat/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite"
	fc "github.com/yogenyslav/ldt-2024/chat/internal/favorite/controller"
	fh "github.com/yogenyslav/ldt-2024/chat/internal/favorite/handler"
	fr "github.com/yogenyslav/ldt-2024/chat/internal/favorite/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/notification"
	nc "github.com/yogenyslav/ldt-2024/chat/internal/notification/controller"
	nh "github.com/yogenyslav/ldt-2024/chat/internal/notification/handler"
	nr "github.com/yogenyslav/ldt-2024/chat/internal/notification/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/session"
	sc "github.com/yogenyslav/ldt-2024/chat/internal/session/controller"
	sh "github.com/yogenyslav/ldt-2024/chat/internal/session/handler"
	sr "github.com/yogenyslav/ldt-2024/chat/internal/session/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/stock"
	stockh "github.com/yogenyslav/ldt-2024/chat/internal/stock/handler"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
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

// Server структура сервера со всеми зависимостями.
type Server struct {
	cfg      *config.Config
	app      *fiber.App
	pg       storage.SQLDatabase
	tracer   trace.Tracer
	exporter sdktrace.SpanExporter
	kc       *gocloak.GoCloak
}

// NewServer создает новый сервер.
func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: srvresp.NewErrorHandler(errStatus).Handler,
		AppName:      "GKU Chat",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.CorsOrigins, ","),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Access-Control-Allow-Origin,Authorization,Origin,Accept,Content-Type,ngrok-skip-browser-warning",
		AllowCredentials: true,
	}))
	app.Use(otelfiber.Middleware())
	app.Use(recovermw.New())
	app.Get("/chat/swagger/*", swagger.HandlerDefault)

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
		kc:       gocloak.NewClient(cfg.KeyCloak.URL),
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
			log.Error().Err(err).Msg("failed to shutdown server")
		}
	}()

	apiClient, err := client.NewGrpcClient(s.cfg.API)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create api grpc client")
	}
	defer func() {
		if err = apiClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close api grpc client")
		}
	}()

	authController := ac.New(pb.NewAuthServiceClient(apiClient.GetConn()), s.cfg.Server.CipherKey, s.tracer)
	authHandler := ah.New(authController)
	auth.SetupAuthRoutes(s.app, authHandler)

	g := s.app.Group("/chat")
	g.Use(middleware.JWT(s.kc, s.cfg.KeyCloak.Realm, s.cfg.Server.CipherKey))

	sessionRepo := sr.New(s.pg)
	sessionController := sc.New(sessionRepo, s.tracer)
	sessionHandler := sh.New(sessionController)
	session.SetupSessionRoutes(g, sessionHandler)

	chatRepo := cr.New(s.pg)
	chatController := cc.New(
		chatRepo,
		sessionRepo,
		pb.NewPrompterClient(apiClient.GetConn()),
		pb.NewPredictorClient(apiClient.GetConn()),
		s.kc,
		s.cfg.KeyCloak.Realm,
		s.cfg.Server.CipherKey,
		s.tracer,
	)
	chatHandler := ch.New(chatController, s.tracer)
	chat.SetupChatRoutes(s.app, chatHandler, s.getWsConfig())

	favoriteRepo := fr.New(s.pg)
	favoriteController := fc.New(favoriteRepo, s.tracer)
	favoriteHandler := fh.New(favoriteController)
	favorite.SetupFavoriteRoutes(g, favoriteHandler)

	stockHandler := stockh.New(pb.NewPredictorClient(apiClient.GetConn()))
	stock.SetupStockRoutes(g, stockHandler)

	notificationRepo := nr.New(s.pg)
	notificationController := nc.New(notificationRepo, s.tracer)
	notificationHandler := nh.New(notificationController)
	notification.SetupNotificationRoutes(g, notificationHandler)

	go s.listen()
	go prom.HandlePrometheus(s.cfg.Prom)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
	log.Info().Msg("shutting down the server")
}

func (s *Server) listen() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	if err := s.app.Listen(addr); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}

func (s *Server) getWsConfig() websocket.Config {
	return websocket.Config{
		RecoverHandler: func(conn *websocket.Conn) {
			if e := recover(); e != nil {
				err := conn.WriteJSON(chatresp.Response{
					Msg: "internal error",
					Err: fmt.Errorf("can't handle error: %v", e).Error(),
				})
				if err != nil {
					log.Warn().Err(err).Msg("failed to recover ws panic")
				}
			}
		},
	}
}
