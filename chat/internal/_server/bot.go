package server

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

// Bot структура бота со всеми зависимостями.
type Bot struct {
	cfg *config.Config
	bot *tele.Bot
}

// NewBot создает новый Bot.
func NewBot(cfg *config.Config) *Bot {
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
	bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Это бот для прогнозирования закупок. Чтобы начать работу, нажмите на кнопку \"Чат\" в меню.")
	})
	if err := bot.SetMenuButton(bot.Me, &tele.MenuButton{
		Type:   tele.MenuButtonWebApp,
		Text:   "Чат",
		WebApp: &tele.WebApp{URL: cfg.Server.AppURL},
	}); err != nil {
		log.Panic().Err(err).Msg("can't set menu button")
	}

	return &Bot{
		cfg: cfg,
		bot: bot,
	}
}

// Run запускает бота.
func (b *Bot) Run() {
	defer b.bot.Stop()

	go b.bot.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Info().Msg("shutting down the server")
}
