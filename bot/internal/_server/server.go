package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/bot/config"
	"github.com/yogenyslav/ldt-2024/bot/internal/auth"
	"github.com/yogenyslav/ldt-2024/bot/internal/auth/handler"
	authmw "github.com/yogenyslav/ldt-2024/bot/internal/auth/middleware"
	mw "github.com/yogenyslav/ldt-2024/bot/internal/middleware"
	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	"github.com/yogenyslav/ldt-2024/bot/internal/state"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage"
	"github.com/yogenyslav/pkg/storage/postgres"
	rediscache "github.com/yogenyslav/pkg/storage/redis_cache"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

// Server is a struct that contains the necessary fields for the server.
type Server struct {
	cfg      *config.Config
	bot      *tele.Bot
	pg       storage.SQLDatabase
	redis    storage.Cache
	exporter sdktrace.SpanExporter
	tracer   trace.Tracer
}

// New creates a new Server.
func New(cfg *config.Config) *Server {
	bot, err := tele.NewBot(tele.Settings{
		Verbose:   cfg.Server.DebugMode,
		Token:     cfg.Server.BotToken,
		ParseMode: tele.ModeHTML,
		Poller:    &tele.LongPoller{Timeout: 60 * time.Second},
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // no need to verify the certificate
			},
		},
		OnError: func(err error, c tele.Context) {
			lg := log.Error().Err(err)
			if c != nil {
				lg.Any(shared.UserIDKey, c.Get(shared.UserIDKey))
				e := c.Send(shared.ErrorMessage)
				if e != nil {
					log.Warn().Err(e).Msg("can't send error response")
				}
			}
			lg.Msg("failed")
		},
	})
	if err != nil {
		log.Panic().Err(err).Msg("can't create bot")
	}

	bot.Use(middleware.Recover())
	bot.Use(middleware.AutoRespond())

	exporter := tracing.MustNewExporter(context.Background(), cfg.Jaeger.URL())
	provider := tracing.MustNewTraceProvider(exporter, "bot")
	otel.SetTracerProvider(provider)

	tracer := otel.Tracer("bot")

	bot.Use(mw.Tracing(tracer))

	return &Server{
		cfg:      cfg,
		bot:      bot,
		pg:       postgres.MustNew(cfg.Postgres, tracer),
		redis:    rediscache.MustNew(cfg.Redis, tracer),
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
	defer s.bot.Stop()

	machine := state.New(s.redis)

	authProvider := auth.Listen(machine, s.bot, s.cfg.Server.Port)
	defer func() {
		if err := authProvider.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown auth provider")
		}
	}()

	authHandler := handler.New(s.tracer)
	auth.SetupAuthRoutes(s.bot, authHandler)

	s.bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send(c.Text())
	}, authmw.JWT(machine, gocloak.NewClient(s.cfg.KeyCloak.URL), s.cfg.KeyCloak.Realm))

	go s.bot.Start()
	go prom.HandlePrometheus(s.cfg.Prometheus)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Info().Msg("shutting down the server")
}
