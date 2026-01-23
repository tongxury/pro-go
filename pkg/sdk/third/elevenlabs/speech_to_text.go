package elevenlabs

import (
	"context"
)

// TranscribeFile transcribes an audio file.
func (c *Client) TranscribeFile(ctx context.Context, reqBody *TranscribeRequest) (*TranscriptionResult, error) {
	var respData TranscriptionResult
	req := c.request(ctx)
	
	if reqBody.File != "" {
		req.SetFile("file", reqBody.File)
	}
	if reqBody.URL != "" {
		req.SetFormData(map[string]string{"url": reqBody.URL})
	}
	req.SetFormData(map[string]string{"model_id": reqBody.ModelID})

	resp, err := req.SetResult(&respData).Post("/v1/speech-to-text")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}
