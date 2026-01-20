package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kida21/telegram-langbot/internal/models"
	"github.com/kida21/telegram-langbot/internal/repositories"
)

type VocabularyService struct {
    repo       *repositories.VocabularyRepository
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

func NewVocabularyService(repo *repositories.VocabularyRepository,apiBaseURL string) *VocabularyService {
	return &VocabularyService{repo: repo,apiBaseURL: apiBaseURL}
}

func (s *VocabularyService) GetWordOfTheDay() (*models.Vocabulary, error) {
	return s.repo.GetRandom()
}

func (s *VocabularyService) FetchAndStore(word, sourceLang, targetLang string) (string, string, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return "", "", fmt.Errorf("GEMINI_API_KEY not set")
    }

    
    url := fmt.Sprintf("%s/v1/models/gemini-2.5-flash:generateContent?key=%s", s.apiBaseURL, apiKey)

    
    prompt := fmt.Sprintf("Translate the word '%s' from %s to %s. Respond with only the translated word.", word, sourceLang, targetLang)
    reqBody := map[string]interface{}{
        "contents": []map[string]interface{}{
            {"role": "user", "parts": []map[string]string{{"text": prompt}}},
        },
    }

    bodyBytes, _ := json.Marshal(reqBody)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
    if err != nil {
        return "", "", fmt.Errorf("gemini request failed: %w", err)
    }
    defer resp.Body.Close()

    raw, _ := io.ReadAll(resp.Body)
    // log.Printf("Gemini translation raw response: %s", string(raw))

    var gResp geminiResponse
    if err := json.Unmarshal(raw, &gResp); err != nil {
        return "", "", fmt.Errorf("decode error: %w, body: %s", err, raw)
    }
    if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
        return "", "", fmt.Errorf("Gemini returned no translation for '%s'", word)
    }
    translation := gResp.Candidates[0].Content.Parts[0].Text

    // --- Example sentence ---
    exPrompt := fmt.Sprintf("Give me a simple example sentence using '%s' in %s.", translation, targetLang)
    exReq := map[string]interface{}{
        "contents": []map[string]interface{}{
            {"role": "user", "parts": []map[string]string{{"text": exPrompt}}},
        },
    }
    exBody, _ := json.Marshal(exReq)
    exResp, err := http.Post(url, "application/json", bytes.NewBuffer(exBody))
    if err != nil {
        return translation, "", fmt.Errorf("example request failed: %w", err)
    }
    defer exResp.Body.Close()

    exRaw, _ := io.ReadAll(exResp.Body)
    // log.Printf("Gemini example raw response: %s", string(exRaw))

    var exRespObj geminiResponse
    if err := json.Unmarshal(exRaw, &exRespObj); err != nil {
        return translation, "", fmt.Errorf("decode error: %w, body: %s", err, exRaw)
    }
    if len(exRespObj.Candidates) == 0 || len(exRespObj.Candidates[0].Content.Parts) == 0 {
        return translation, "", fmt.Errorf("Gemini returned no example for '%s'", translation)
    }
    example := exRespObj.Candidates[0].Content.Parts[0].Text

    // --- Store in DB ---
    if err := s.repo.Insert(word, translation, example); err != nil {
        return translation, example, fmt.Errorf("db insert failed: %w", err)
    }

    return translation, example, nil
}