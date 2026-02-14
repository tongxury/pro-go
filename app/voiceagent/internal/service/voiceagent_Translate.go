package service

import (
	"context"
	"fmt"
	"store/pkg/sdk/third/gemini"
	"strings"

	"google.golang.org/genai"
)

// Translate uses Gemini to translate text to the target language.
func (s *VoiceAgentService) Translate(ctx context.Context, text string, targetLang string) string {
	if text == "" {
		return ""
	}

	// Simple mapping for common target languages to prompt friendly names
	langName := targetLang
	switch strings.ToLower(targetLang) {
	case "zh", "cn", "zh-cn":
		langName = "Chinese"
	case "en":
		langName = "English"
	case "jp", "ja":
		langName = "Japanese"
	case "kr", "ko":
		langName = "Korean"
	}

	prompt := fmt.Sprintf("Translate the following text to %s. Return ONLY the translated text, no quotes, no explanations:\n\n%s", langName, text)

	res, err := s.Data.Gemini.Get().GenerateContent(ctx, gemini.GenerateContentRequest{
		Parts: []*genai.Part{{Text: prompt}},
	})
	if err != nil {
		return text // Fallback to original text on error
	}

	// Clean up the response
	translated := strings.TrimSpace(res)
	translated = strings.Trim(translated, "\"")
	translated = strings.Trim(translated, "'")

	if translated == "" {
		return text
	}

	return translated
}
