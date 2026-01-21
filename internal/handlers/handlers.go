package handlers

import (
	"fmt"
	
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kida21/telegram-langbot/internal/services"
)

type Handler struct {
	userService     *services.UserService
	vocabService    *services.VocabService
	
}

func NewHandler(us *services.UserService, vs *services.VocabService) *Handler {
	return &Handler{
		userService:     us,
		vocabService:    vs,
		
	}
}

func (h *Handler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    if update.Message == nil {
        return
    }

    cmd := update.Message.Command()
    if cmd != "" {
        switch cmd {
        case "start":
            h.handleStart(bot, update)
        case "setlanguage":
            h.handleSetLanguage(bot, update)
		case "importword":
			h.handleImportWord(bot,update)
        default:
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command."))
        }
        return
    }

    
    h.handleText(bot, update)
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






func (h *Handler) handleText(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    lang := strings.TrimSpace(update.Message.Text)
    validLangs := []string{"Spanish", "French", "German", "Italian", "Japanese"}

    for _, l := range validLangs {
        if lang == l {
            tgID := update.Message.From.ID
            username := update.Message.From.UserName
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

    bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please choose a valid language option."))
}

func (h *Handler) handleSetLanguage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose your new learning language:")
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
}


func (h *Handler) handleImportWord(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    text := update.Message.CommandArguments()
    if text == "" {
        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /importword <text>"))
        return
    }

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Which language would you like to translate into?")
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Spanish", "translate:Spanish:"+text),
            tgbotapi.NewInlineKeyboardButtonData("Japanese", "translate:Japanese:"+text),
            tgbotapi.NewInlineKeyboardButtonData("French", "translate:French:"+text),
            tgbotapi.NewInlineKeyboardButtonData("German", "translate:German:"+text),
        ),
    )
    msg.ReplyMarkup = keyboard
    bot.Send(msg)
}

