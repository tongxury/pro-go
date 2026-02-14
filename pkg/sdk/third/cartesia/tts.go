package cartesia

import (
	"context"
	"fmt"
)

type TTSRequest struct {
	ModelId      string       `json:"model_id"`
	Transcript   string       `json:"transcript"`
	Voice        TTSVoice     `json:"voice"`
	OutputFormat OutputFormat `json:"output_format"`
}

type TTSVoice struct {
	Mode string `json:"mode"`
	Id   string `json:"id"`
}

type OutputFormat struct {
	Container  string `json:"container"`
	Encoding   string `json:"encoding"`
	SampleRate int    `json:"sample_rate"`
}

func (c *Client) TextToSpeechBytes(ctx context.Context, voiceId string, modelId string, transcript string) ([]byte, error) {
	url := fmt.Sprintf("%s/tts/bytes", BaseURL)

	if modelId == "" {
		modelId = "sonic-multilingual"
	}

	reqBody := TTSRequest{
		ModelId:    modelId,
		Transcript: transcript,
		Voice: TTSVoice{
			Mode: "id",
			Id:   voiceId,
		},
		OutputFormat: OutputFormat{
			Container:  "wav",
			Encoding:   "pcm_s16le",
			SampleRate: 44100,
		},
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(reqBody).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("cartesia tts api error: status %d body %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
