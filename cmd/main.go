package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    token := os.Getenv("TELEGRAM_BOT_TOKEN")
	log.Printf("token:%s",token)
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Printf("error occured : %v",err)
    }

    bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, welcome to LanguageBot!")
            bot.Send(msg)
        }
    }
}