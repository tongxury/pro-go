package clients

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/conv"
	"time"
)

type LokiClient struct {
	endpoint string
}

func NewLokiClient(endpoint string) *LokiClient {
	return &LokiClient{
		endpoint: endpoint,
	}
}

func (t *LokiClient) Query(query string) (LokiResults, error) {

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"query": query,
		}).
		Get(t.endpoint + "/loki/api/v1/query")

	if err != nil {
		return nil, err
	}

	var result LokiResponse
	_ = conv.J2S(resp.Body(), &result)

	if result.Status != "success" {
		return nil, fmt.Errorf(resp.String())
	}

	return result.Data.Result, nil
}

// https://grafana.com/docs/loki/latest/reference/api/
// https://grafana.com/docs/loki/latest/query/log_queries/
func (t *LokiClient) QueryRange(query string, start, end time.Time) (LokiResults, error) {

	// sum by(container) (rate({namespace="prod", container="sg-llm"} |~ `(?i)error` [1m]))

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"query": query,
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		}).
		Get(t.endpoint + "/loki/api/v1/query_range")

	if err != nil {
		return nil, err
	}

	var result LokiResponse
	_ = conv.J2S(resp.Body(), &result)

	if result.Status != "success" {
		return nil, fmt.Errorf(resp.String())
	}

	return result.Data.Result, nil
}

type LokiResult struct {
	Metric struct {
		Container string `json:"container"`
	} `json:"metric"`
	Value []interface{} `json:"value"`
}

type LokiResults []LokiResult

type LokiResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string      `json:"resultType"`
		Result     LokiResults `json:"result"`
	} `json:"data"`
}
