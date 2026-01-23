package gmgnai

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

type SubmitSignedTransactionResult struct {
	Code    int
	Message string
	Data    struct {
		Hash string `json:"hash"`
	}
}

func (t *Client) SubmitSignedTransaction(ctx context.Context, txBase64 string) (*SubmitSignedTransactionResult, error) {

	post, err := resty.New().R().SetContext(ctx).
		SetBody(map[string]string{"signed_tx": txBase64}).
		Post("https://gmgn.ai/defi/router/v1/sol/tx/submit_signed_transaction")
	if err != nil {
		return nil, err
	}

	var result SubmitSignedTransactionResult
	if err := json.Unmarshal(post.Body(), &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, errors.New(result.Message)
	}

	return &result, nil
}
