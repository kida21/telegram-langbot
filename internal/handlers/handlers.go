package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kida21/telegram-langbot/internal/services"
)

type Handler struct {
	userService     *services.UserService
	vocabService    *services.VocabularyService
	quizService     *services.QuizService
	progressService *services.ProgressService
}

func NewHandler(us *services.UserService, vs *services.VocabularyService, qs *services.QuizService, ps *services.ProgressService) *Handler {
	return &Handler{
		userService:     us,
		vocabService:    vs,
		quizService:     qs,
		progressService: ps,
	}
}

func (h *Handler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	switch update.Message.Command() {
	case "start":
		h.handleStart(bot, update)
	case "word":
		h.handleWord(bot, update)
	case "quiz":
		h.handleQuiz(bot, update)
	case "progress":
		h.handleProgress(bot, update)
	default:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Try /start, /word, /quiz, /progress"))
	}
}


func (h *Handler) handleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	tgID := update.Message.From.ID
	username := update.Message.From.UserName

	if username == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"You don’t have a Telegram username. Please reply with a unique username using /setusername <name>."))
		return
	}

	
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose the language you want to learn:")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Spanish"),
			tgbotapi.NewKeyboardButton("French"),
			tgbotapi.NewKeyboardButton("German"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Italian"),
			tgbotapi.NewKeyboardButton("Japanese"),
		),
	)
	bot.Send(msg)

	
	lang := strings.TrimSpace(update.Message.Text)
	validLangs := []string{"Spanish", "French", "German", "Italian", "Japanese"}
	for _, l := range validLangs {
		if lang == l {
			user, err := h.userService.RegisterOrGet(tgID, username, lang)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error registering: "+err.Error()))
				return
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				fmt.Sprintf("Welcome %s! You’ll now learn %s. Try /word or /quiz to begin.", user.Username, user.LanguagePref)))
			return
		}
	}
}

func (h *Handler) handleWord(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	vocab, err := h.vocabService.GetWordOfTheDay()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No vocabulary found."))
		return
	}
	text := fmt.Sprintf("Word: %s\nTranslation: %s\nExample: %s", vocab.Word, vocab.Translation, vocab.Example)
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
}

func (h *Handler) handleQuiz(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	quiz, err := h.quizService.GetRandomQuiz()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No quizzes available."))
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, quiz.Question+"\nOptions: "+quiz.Options)
	bot.Send(msg)
}

func (h *Handler) handleProgress(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	attempts, correct, accuracy, err := h.progressService.GetUserStats(update.Message.From.ID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error fetching progress."))
		return
	}
	text := fmt.Sprintf("Attempts: %d\nCorrect: %d\nAccuracy: %.2f%%", attempts, correct, accuracy)
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
}
