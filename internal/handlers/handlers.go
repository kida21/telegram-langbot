package handlers

import "github.com/kida21/telegram-langbot/internal/services"

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

