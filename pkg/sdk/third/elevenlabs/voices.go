package elevenlabs

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// ListVoices lists all voices.
func (c *Client) ListVoices(ctx context.Context) (*ListVoicesResp, error) {
	var respData ListVoicesResp
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/voices")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// GetVoice retrieves a specific voice.
func (c *Client) GetVoice(ctx context.Context, voiceID string) (*Voice, error) {
	var respData Voice
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/voices/" + voiceID)
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// AddVoice creates a new voice (Instant Voice Cloning).
func (c *Client) AddVoice(ctx context.Context, name string, description string, files []string) (string, error) {
	var result struct {
		VoiceID string `json:"voice_id"`
	}
	req := c.request(ctx).
		SetFormData(map[string]string{
			"name":        name,
			"description": description,
		})

	for _, file := range files {
		if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
			// If it's a URL, download the content first
			resp, err := c.restClient.R().SetContext(ctx).Get(file)
			if err != nil {
				return "", fmt.Errorf("failed to download voice sample from %s: %v", file, err)
			}
			if resp.IsError() {
				return "", fmt.Errorf("failed to download voice sample from %s: status %d", file, resp.StatusCode())
			}
			// Use the filename from URL or a default one
			fileName := filepath.Base(file)
			if fileName == "" || fileName == "." {
				fileName = "sample.mp3"
			}
			req.SetFileReader("files", fileName, bytes.NewReader(resp.Body()))
		} else {
			// Local file
			req.SetFile("files", file)
		}
	}

	resp, err := req.SetResult(&result).Post("/v1/voices/add")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	return result.VoiceID, nil
}

// EditVoice updates an existing voice.
func (c *Client) EditVoice(ctx context.Context, voiceID string, name string, description string) error {
	req := c.request(ctx).
		SetFormData(map[string]string{
			"name":        name,
			"description": description,
		})
	resp, err := req.Post("/v1/voices/" + voiceID + "/edit")
	return handleResponse(resp, err)
}

// GetVoiceSettings retrieves settings for a specific voice.
func (c *Client) GetVoiceSettings(ctx context.Context, voiceID string) (*VoiceSettings, error) {
	var respData VoiceSettings
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/voices/" + voiceID + "/settings")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// DeleteVoice deletes a voice.
func (c *Client) DeleteVoice(ctx context.Context, voiceID string) error {
	resp, err := c.request(ctx).Delete("/v1/voices/" + voiceID)
	return handleResponse(resp, err)
}
