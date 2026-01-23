package elevenlabs

import (
	"context"
	"encoding/json"
	"fmt"
	"store/pkg/sdk/conv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	APIKey  string
	BaseURL string
}

type Client struct {
	restClient *resty.Client
	conf       Config
}

func NewClient() *Client {

	conf := Config{
		APIKey: "sk_014353b30a4310e42c8d584b83e08a5b1286771fe944eb30",
		//APIKey: "sk_ae62c02c8e3e63f0d214085014e9cc09d9d91867509f585f",
	}

	if conf.BaseURL == "" {
		conf.BaseURL = "https://api.elevenlabs.io"
	}

	client := resty.New().
		SetBaseURL(conf.BaseURL).
		SetHeader("xi-api-key", conf.APIKey).
		SetTimeout(30 * time.Second)

	return &Client{
		restClient: client,
		conf:       conf,
	}
}

func (c *Client) request(ctx context.Context) *resty.Request {
	return c.restClient.R().SetContext(ctx)
}

// ErrorResp represents an error response from ElevenLabs API
type ErrorResp struct {
	Detail json.RawMessage `json:"detail"`
}

type ValidationError struct {
	Type  string `json:"type"`
	Loc   []any  `json:"loc"`
	Msg   string `json:"msg"`
	Input any    `json:"input"`
}

func (e *ErrorResp) Error() string {
	if len(e.Detail) == 0 {
		return "elevenlabs error: unknown"
	}

	// Try to parse as the array of validation errors
	var validationErrors []ValidationError
	if err := json.Unmarshal(e.Detail, &validationErrors); err == nil && len(validationErrors) > 0 {
		// If it's a slice but the first element isn't a valid ValidationError, unmarshal might still "succeed" but with empty fields
		// We check if the first element has some identifying field if possible, or just check the length
		var messages []string
		for _, ve := range validationErrors {
			if ve.Msg == "" {
				continue
			}
			locStr := []string{}
			for _, l := range ve.Loc {
				locStr = append(locStr, fmt.Sprint(l))
			}
			if len(locStr) > 0 {
				messages = append(messages, fmt.Sprintf("[%s]: %s", strings.Join(locStr, " -> "), ve.Msg))
			} else {
				messages = append(messages, ve.Msg)
			}
		}
		if len(messages) > 0 {
			return "elevenlabs validation error: " + strings.Join(messages, "; ")
		}
	}

	// Try to parse as the simple object error { "message": "...", "status": "..." }
	var simpleDetail struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal(e.Detail, &simpleDetail); err == nil && simpleDetail.Message != "" {
		return fmt.Sprintf("elevenlabs error: %s (status: %s)", simpleDetail.Message, simpleDetail.Status)
	}

	// Fallback to raw string or raw json
	var strDetail string
	if err := json.Unmarshal(e.Detail, &strDetail); err == nil {
		return "elevenlabs error: " + strDetail
	}

	return "elevenlabs error: " + string(e.Detail)
}

func handleResponse(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.IsError() {
		var errResp ErrorResp
		if e := conv.J2S(resp.Body(), &errResp); e == nil && len(errResp.Detail) > 0 {
			return &errResp
		}
		return fmt.Errorf("elevenlabs api error: %s (status: %d)", resp.Status(), resp.StatusCode())
	}
	return nil
}
