package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kida21/telegram-langbot/internal/repositories"
)

type VocabService struct {
    repo       *repositories.TranslationRepository
    apiBaseURL string
}

type geminiResponse struct {
    Candidates []struct {
        Content struct {
            Parts []struct {
                Text string `json:"text"`
            } `json:"parts"`
        } `json:"content"`
    } `json:"candidates"`
}

func NewVocabService(repo *repositories.TranslationRepository, apiBaseURL string) *VocabService {
    return &VocabService{repo: repo, apiBaseURL: apiBaseURL}
}

func (s *VocabService) TranslateAndLog(userID int64, sourceText, targetLang string) (string, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    url := fmt.Sprintf("%s/v1/models/gemini-2.5-flash:generateContent?key=%s", s.apiBaseURL, apiKey)

    prompt := fmt.Sprintf("Translate \"%s\" into %s. Respond with only the translated text.", sourceText, targetLang)
    reqBody := map[string]interface{}{
        "contents": []map[string]interface{}{
            {"role": "user", "parts": []map[string]string{{"text": prompt}}},
        },
    }

    bodyBytes, _ := json.Marshal(reqBody)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
    if err != nil {
        return "", fmt.Errorf("Gemini request failed: %w", err)
    }
    defer resp.Body.Close()

    raw, _ := io.ReadAll(resp.Body)
    var gResp geminiResponse
    if err := json.Unmarshal(raw, &gResp); err != nil {
        return "", fmt.Errorf("decode error: %w", err)
    }
    if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
        return "", fmt.Errorf("Gemini returned no translation")
    }
    translation := gResp.Candidates[0].Content.Parts[0].Text

    // Log into DB
    if err := s.repo.LogTranslation(userID, sourceText, targetLang, translation); err != nil {
        return translation, fmt.Errorf("failed to log translation: %w", err)
    }

    return translation, nil
}