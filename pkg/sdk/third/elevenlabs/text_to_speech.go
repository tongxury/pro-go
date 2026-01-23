package elevenlabs

import (
	"context"
	"fmt"
)

type TextToSpeechRequest struct {
	Text          string         `json:"text"`
	ModelID       string         `json:"model_id,omitempty"`
	VoiceSettings *VoiceSettings `json:"voice_settings,omitempty"`
}

// TextToSpeech generates speech from text using the specified voice ID.
// It returns the raw audio data (usually MP3).
func (c *Client) TextToSpeech(ctx context.Context, voiceID string, reqBody *TextToSpeechRequest) ([]byte, error) {
	if voiceID == "" {
		return nil, fmt.Errorf("voice_id is required")
	}

	url := fmt.Sprintf("/v1/text-to-speech/%s", voiceID)
	req := c.request(ctx).SetBody(reqBody)

	resp, err := req.Post(url)
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
