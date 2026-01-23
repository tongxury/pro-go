package airwallex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

type Config struct {
	ClientId       string
	ClientSecret   string
	Endpoint       string
	AccountId      string
	CallbackSecret string
}

type Client struct {
	config Config
	c      *cache.Cache
}

func NewClient(conf Config) *Client {

	//config := Config{
	//	ClientId:     "sv9jws-mTFGv2iejM1C3HA",
	//	ClientSecret: "329f7a951023503d81337598f7af89b17b668343a761f3a3de4fe735a7788268fa36d5d3bac02f7ec2323ab9e7b97dde",
	//	//ClientId:     "9jtRbZaXSAy5GcPGgZVHeA",
	//	//ClientSecret: "9fd7512838865b660d0c4de708b8edb28282b93a633a335bfd05eeca22567e7e17eea9dc22487cb2d687be77d0184178",
	//	Endpoint:       "https://api.airwallex.com",
	//	AccountId:      "acct_Q6NUqVssMGu--2O8trlsVw",
	//}

	return &Client{
		config: conf,
		c:      cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (t *Client) getAuthToken(ctx context.Context) (string, error) {

	if token, ok := t.c.Get("authToken"); ok {
		return token.(string), nil
	}

	var at AuthToken

	_, err := resty.New().R().
		SetContext(ctx).
		SetHeader("x-client-id", t.config.ClientId).
		SetHeader("x-api-key", t.config.ClientSecret).
		//SetHeader("x-api-version", "2025-06-16").
		SetResult(&at).
		Post(t.config.Endpoint + "/api/v1/authentication/login")

	if err != nil {
		return "", nil
	}

	if at.Token == "" {
		return "", errors.New("empty auth token")
	}

	t.c.Set("authToken", at.Token, 5*time.Minute)

	return at.Token, nil

}

type AuthToken struct {
	ExpiresAt string `json:"expires_at"`
	Token     string `json:"token"`
}

func (t *Client) post(ctx context.Context, api string, body map[string]interface{}, result any) error {

	token, err := t.getAuthToken(ctx)
	if err != nil {
		return err
	}

	r, err := resty.New().R().
		SetContext(ctx).
		SetBody(body).
		SetHeader("Authorization", "Bearer "+token).
		//SetHeader("x-api-version", "2025-06-16").
		Post(t.config.Endpoint + api)

	if err != nil {
		return err
	}

	var em ErrorResponse
	err = json.Unmarshal(r.Body(), &em)
	if err != nil {
		return err
	}

	if em.Code != "" {
		return errors.New(em.Message)
	}

	err = json.Unmarshal(r.Body(), result)
	if err != nil {
		return err
	}

	return nil
}

func (t *Client) get(ctx context.Context, api string, params map[string]string, result any) error {

	token, err := t.getAuthToken(ctx)
	if err != nil {
		return err
	}

	response, err := resty.New().R().
		SetContext(ctx).
		SetQueryParams(params).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("x-api-version", "2025-06-16").
		SetResult(result).
		Get(t.config.Endpoint + api)

	fmt.Println(response)

	if err != nil {
		return err
	}

	return nil
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"trace_id"`
	Details struct {
		ResourceId string `json:"resource_id"`
	} `json:"details"`
}

type List[T any] struct {
	HasMore bool `json:"has_more"`
	Items   []T  `json:"items"`
}
