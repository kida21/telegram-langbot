package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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





func (s *VocabularyService) FetchAndStore(word, sourceLang, targetLang string) (string, string, error) {
    // Translate the word itself
    reqBody := libreTranslateRequest{
        Q:      word,
        Source: sourceLang,
        Target: targetLang,
        Format: "text",
    }

    bodyBytes, _ := json.Marshal(reqBody)
    resp, err := http.Post(s.apiBaseURL+"/translate", "application/json", bytes.NewBuffer(bodyBytes))
    if err != nil {
        return "", "", fmt.Errorf("translation request failed: %w", err)
    }
    defer resp.Body.Close()

    raw, _ := io.ReadAll(resp.Body)
    log.Printf("Raw translation response: %s", string(raw))

    var ltResp libreTranslateResponse
    if err := json.Unmarshal(raw, &ltResp); err != nil {
        return "", "", fmt.Errorf("decode error: %w, body: %s", err, raw)
    }
    translation := ltResp.TranslatedText

    // Translate an example sentence
    exampleSentence := fmt.Sprintf("%s, how are you?", word)
    exReq := libreTranslateRequest{
        Q:      exampleSentence,
        Source: sourceLang,
        Target: targetLang,
        Format: "text",
    }
    exBody, _ := json.Marshal(exReq)
    exResp, err := http.Post(s.apiBaseURL+"/translate", "application/json", bytes.NewBuffer(exBody))
    if err != nil {
        return translation, "", fmt.Errorf("example request failed: %w", err)
    }
    defer exResp.Body.Close()

    exRaw, _ := io.ReadAll(exResp.Body)
    log.Printf("Raw example response: %s", string(exRaw))

    var exLtResp libreTranslateResponse
    if err := json.Unmarshal(exRaw, &exLtResp); err != nil {
        return translation, "", fmt.Errorf("decode error: %w, body: %s", err, exRaw)
    }
    example := exLtResp.TranslatedText

    // Store in DB
    if err := s.repo.Insert(word, translation, example); err != nil {
        return translation, example, fmt.Errorf("db insert failed: %w", err)
    }

    return translation, example, nil
}