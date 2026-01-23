package restyd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Params struct {
	Url     string
	Params  map[string]string
	Body    any
	Headers map[string]string
}

func (t Params) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func Post[T any](ctx context.Context, params Params) (*T, error) {
	response, err := resty.New().R().
		SetContext(ctx).
		SetBody(params.Body).
		SetHeaders(params.Headers).
		Post(params.Url)

	if err != nil {
		return nil, err
	}

	var result T
	err = parseResult(params, response, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Get[T any](ctx context.Context, params Params) (*T, error) {
	response, err := resty.New().R().
		SetContext(ctx).
		SetQueryParams(params.Params).
		SetHeaders(params.Headers).
		Get(params.Url)

	if err != nil {
		return nil, err
	}

	var result T
	err = parseResult(params, response, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func parseResult[T any](params Params, response *resty.Response, result T) error {

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("resty response status code: %d, body: %v, params: %v", response.StatusCode(), response.String(), params)
	}

	if err := json.Unmarshal(response.Body(), result); err != nil {
		return fmt.Errorf("resty response json unmarshalling error: %v, body: %v, params: %v", err, response.String(), params)
	}

	return nil
}

func ParseResult[T any](response *resty.Response, result T) error {

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("resty response status code: %d, %v", response.StatusCode(), response.String())
	}

	if err := json.Unmarshal(response.Body(), result); err != nil {
		return fmt.Errorf("resty response json unmarshalling error: %v, %v", err, response.String())
	}

	return nil
}
