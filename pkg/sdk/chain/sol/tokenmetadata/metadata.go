package tokenmetadata

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"strings"
)

type TokenMetadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ShowName    any    `json:"showName"` // true or "true"
	CreatedOn   string `json:"createdOn"`
	Twitter     string `json:"twitter"`
	Website     string `json:"website"`

	Websites []struct {
		Label string `json:"label"`
		Url   string `json:"url"`
	} `json:"websites"`
	Socials []struct {
		Url  string `json:"url"`
		Type string `json:"type"`
	} `json:"socials"`
}

func GetTokenMetadata(ctx context.Context, url string) (*TokenMetadata, error) {

	rsp, err := resty.New().R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}

	if !rsp.IsSuccess() {
		return nil, errors.New(rsp.Status())
	}

	if !strings.Contains(rsp.Header().Get("Content-Type"), "application/json") {
		return nil, errors.New(rsp.Status())
	}

	var m TokenMetadata
	err = json.Unmarshal(rsp.Body(), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil

}
