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
	
	if update.CallbackQuery != nil {
		h.handleCallback(bot, update)
		return
	}

	if update.Message == nil {
		return
	}

	cmd := update.Message.Command()
	if cmd != "" {
		switch cmd {
		case "start":
			h.handleStart(bot, update)
		case "importword":
			h.handleImportWord(bot, update)
		case "history":
			h.handleHistory(bot, update)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Try /importword or /history."))
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

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("Welcome %s! This bot helps you translate any word or sentence into your chosen language.\n\nTry `/importword hello, how are you?` to get started.\nYou can also use `/history` to see your past translations.", username))
	bot.Send(msg)

	
	_, _ = h.userService.RegisterOrGet(tgID, username, "")
}


func (h *Handler) handleText(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    text := strings.TrimSpace(update.Message.Text)
    if text == "" {
        return
    }

    
    msg := tgbotapi.NewMessage(update.Message.Chat.ID,
        fmt.Sprintf("You wrote:\n\n%s\n\nChoose a language to translate into:", text))

    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Spanish", "translate:Spanish:"+text),
            tgbotapi.NewInlineKeyboardButtonData("Japanese", "translate:Japanese:"+text),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("French", "translate:French:"+text),
            tgbotapi.NewInlineKeyboardButtonData("German", "translate:German:"+text),
        ),
    )

    msg.ReplyMarkup = keyboard
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

func (h *Handler) handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    data := update.CallbackQuery.Data
    parts := strings.SplitN(data, ":", 3)
    if len(parts) != 3 || parts[0] != "translate" {
        return
    }

    targetLang := parts[1]
    text := parts[2]
    userID := update.CallbackQuery.From.ID

    translation, err := h.vocabService.TranslateAndLog(userID, text, targetLang)
    if err != nil {
        bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error: "+err.Error()))
        return
    }

    msg := fmt.Sprintf("Translation into %s:\n%s", targetLang, translation)
    bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg))
    bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Translated!"))

}


func (h *Handler) handleHistory(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    userID := update.Message.From.ID
    history, err := h.vocabService.GetHistory(userID, 5)
    if err != nil || len(history) == 0 {
        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No history found. Try /importword first."))
        return
    }

    var sb strings.Builder
    sb.WriteString("Your last translations:\n\n")
    for _, entry := range history {
        sb.WriteString(fmt.Sprintf("'%s' → (%s) %s\n", entry.SourceText, entry.TargetLang, entry.TranslatedText))
    }
    bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, sb.String()))
}

