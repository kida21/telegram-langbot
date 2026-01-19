package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/kida21/telegram-langbot/config"
	"github.com/kida21/telegram-langbot/db"
	"github.com/kida21/telegram-langbot/internal/bot"
	"github.com/kida21/telegram-langbot/internal/handlers"
	"github.com/kida21/telegram-langbot/internal/repositories"
	"github.com/kida21/telegram-langbot/internal/services"
)

func main() {
	
	cfg := config.LoadConfig()
    db.InitDatabase(cfg)

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	
	userRepo := repositories.NewUserRepository(db.DB)
	vocabRepo := repositories.NewVocabularyRepository(db.DB)
	quizRepo := repositories.NewQuizRepository(db.DB)
	progressRepo := repositories.NewProgressRepository(db.DB)

	
	userService := services.NewUserService(userRepo)
	vocabService := services.NewVocabularyService(vocabRepo)
	quizService := services.NewQuizService(quizRepo)
	progressService := services.NewProgressService(progressRepo)

	
	handler := handlers.NewHandler(userService, vocabService, quizService, progressService)

	
	b := bot.NewBot(api, handler)
	b.Start()
}