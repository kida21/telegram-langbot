package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kida21/telegram-langbot/internal/models"
	"github.com/kida21/telegram-langbot/internal/repositories"
)

type VocabularyService struct {
	repo *repositories.VocabularyRepository
	apiBaseURL string

}
type libreTranslateRequest struct {
    Q      string `json:"q"`
    Source string `json:"source"`
    Target string `json:"target"`
    Format string `json:"format"`
}

type libreTranslateResponse struct {
    TranslatedText string `json:"translatedText"`
}


func NewVocabularyService(repo *repositories.VocabularyRepository,apiBaseURL string) *VocabularyService {
	return &VocabularyService{repo: repo,apiBaseURL: apiBaseURL}
}

func (s *VocabularyService) GetWordOfTheDay() (*models.Vocabulary, error) {
	return s.repo.GetRandom()
}

func (s *VocabularyService) FetchAndStore(word, sourceLang, targetLang string) error {
    reqBody := libreTranslateRequest{
        Q:      word,
        Source: sourceLang, // e.g. "en"
        Target: targetLang, // e.g. "es", "fr", "de", "it", "ja"
        Format: "text",
    }

    bodyBytes, _ := json.Marshal(reqBody)
    resp, err := http.Post(s.apiBaseURL+"/translate", "application/json", bytes.NewBuffer(bodyBytes))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var ltResp libreTranslateResponse
    if err := json.NewDecoder(resp.Body).Decode(&ltResp); err != nil {
        return err
    }

    translation := ltResp.TranslatedText
    example := fmt.Sprintf("Example: %s → %s", word, translation)

    return s.repo.Insert(word, translation, example)
}
