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
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/auth"
	ah "github.com/yogenyslav/ldt-2024/chat/internal/bot/auth/handler"
	authmw "github.com/yogenyslav/ldt-2024/chat/internal/bot/auth/middleware"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/chat"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/chat/handler"
	mw "github.com/yogenyslav/ldt-2024/chat/internal/bot/middleware"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	cc "github.com/yogenyslav/ldt-2024/chat/internal/chat/controller"
	cr "github.com/yogenyslav/ldt-2024/chat/internal/chat/repo"
	sc "github.com/yogenyslav/ldt-2024/chat/internal/session/controller"
	sr "github.com/yogenyslav/ldt-2024/chat/internal/session/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg/client"
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

// Bot структура бота со всеми зависимостями.
type Bot struct {
	cfg      *config.Config
	bot      *tele.Bot
	pg       storage.SQLDatabase
	redis    storage.Cache
	exporter sdktrace.SpanExporter
	tracer   trace.Tracer
}

// NewBot создает новый Bot.
func NewBot(cfg *config.Config) *Bot {
	bot, err := tele.NewBot(tele.Settings{
		Verbose:   cfg.Server.DebugMode,
		Token:     cfg.Server.BotToken,
		ParseMode: tele.ModeMarkdown,
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

	return &Bot{
		cfg:      cfg,
		bot:      bot,
		pg:       postgres.MustNew(cfg.Postgres, tracer),
		redis:    rediscache.MustNew(cfg.Redis, tracer),
		exporter: exporter,
		tracer:   tracer,
	}
}

// Run запускает бота.
func (b *Bot) Run() {
	defer b.pg.Close()
	defer func() {
		if err := b.exporter.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown exporter")
		}
	}()
	defer b.bot.Stop()

	machine := state.New(b.redis)
	b.bot.Use(mw.Tracing(machine, b.tracer))

	apiClient, err := client.NewGrpcClient(b.cfg.API)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create api grpc client")
	}
	defer func() {
		if err = apiClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close api grpc client")
		}
	}()

	authHandler := ah.New(b.tracer)
	auth.SetupAuthRoutes(b.bot, authHandler)

	g := b.bot.Group()
	g.Use(authmw.JWT(machine, gocloak.NewClient(b.cfg.KeyCloak.URL), b.cfg.KeyCloak.Realm))

	chatRepo := cr.New(b.pg)
	sessionRepo := sr.New(b.pg)
	chatController := cc.New(
		chatRepo,
		sessionRepo,
		pb.NewPrompterClient(apiClient.GetConn()),
		pb.NewPredictorClient(apiClient.GetConn()),
		gocloak.NewClient(b.cfg.KeyCloak.URL),
		b.cfg.KeyCloak.Realm,
		b.cfg.Server.CipherKey,
		b.tracer,
	)
	sessionController := sc.New(sessionRepo, b.tracer)
	chatHandler := handler.New(chatController, sessionController, machine, b.tracer)
	chat.SetupChatRoutes(g, chatHandler, machine)

	go b.bot.Start()
	go prom.HandlePrometheus(b.cfg.BotProm)

	authProvider := auth.Listen(machine, b.bot, b.cfg.Server.BotPort)
	defer func() {
		if err := authProvider.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown auth provider")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Info().Msg("shutting down the server")
}
