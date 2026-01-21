package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kida21/telegram-langbot/internal/handlers"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *handlers.Handler
}


func NewBot(api *tgbotapi.BotAPI, handler *handlers.Handler) *Bot {
	return &Bot{api: api, handler: handler}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
    u.AllowedUpdates = []string{"message", "callback_query"}
    updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		b.handler.HandleUpdate(b.api, update)
	 }
}