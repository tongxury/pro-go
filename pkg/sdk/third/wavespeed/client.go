package wavespeed

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	conf Config
}

type Config struct {
	APIKey  string
	BaseURL string
}

func NewClient() *Client {
	return &Client{
		conf: Config{
			APIKey:  "d3630c4815fbdc65821dce615689e92b688c2e585ce34ab0faad43fc1360ef12",
			BaseURL: "https://api.wavespeed.ai",
		},
	}
}

func (t *Client) invoke(ctx context.Context, url string, body any, result any) error {

	post, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+t.conf.APIKey).
		SetBody(body).
		Post(t.conf.BaseURL + url)

	if err != nil {
		return err
	}

	err = json.Unmarshal(post.Body(), &result)
	if err != nil {
		return err
	}

	//post.IsSuccess()

	return nil
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Id      string   `json:"id"`
		Model   string   `json:"model"`
		Outputs []string `json:"outputs"`
		Urls    struct {
			Get string `json:"get"`
		} `json:"urls"`
		HasNsfwContents []interface{} `json:"has_nsfw_contents"`
		Status          string        `json:"status"`
		CreatedAt       time.Time     `json:"created_at"`
		Error           string        `json:"error"`
		ExecutionTime   int           `json:"executionTime"`
		Timings         struct {
			Inference int `json:"inference"`
		} `json:"timings"`
	} `json:"data"`
}

/*

{
  "code" : 200,
  "message" : "success",
  "data" : {
    "id" : "7d8d7372150044d283b4c66285bc8154",
    "model" : "google/gemini-3-pro-image/edit",
    "outputs" : [ ],
    "urls" : {
      "get" : "https://api.wavespeed.ai/api/v3/predictions/7d8d7372150044d283b4c66285bc8154/result"
    },
    "has_nsfw_contents" : [ ],
    "status" : "created",
    "created_at" : "2026-01-05T07:56:10.979Z",
    "error" : "",
    "executionTime" : 0,
    "timings" : {
      "inference" : 0
    }
  }
}
*/
