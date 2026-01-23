package bitquery

import (
	"context"
	json "encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	//c *graphql.SubscriptionClient
	endpoint string
	apiKey   string
}

func New() *Client {
	return &Client{
		endpoint: "https://streaming.bitquery.io/eap",
		apiKey:   "ory_at_RlCavpqd7pLFaz6SXTa9ZpPVYCziBRwIw0TQOYaKZ7Q.A76LPN4ovOupXC2rSeki2Hiz7csyjRglrwcca5Ac0p0",
	}
}

func (t *Client) Query(ctx context.Context, sql string, params map[string]any, response any) error {

	r, err := resty.New().R().SetContext(ctx).
		SetHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", t.apiKey),
		}).
		SetBody(map[string]interface{}{
			"query":     sql,
			"variables": params,
		}).
		Post(t.endpoint)

	if err != nil {
		return err
	}

	if r.StatusCode() != 200 {
		return fmt.Errorf("query failed: %s", r.String())
	}

	err = json.Unmarshal(r.Body(), response)
	if err != nil {
		return err
	}

	return nil
}

//func NewClient(conf Config) *Client {
//
//	client := graphql.NewSubscriptionClient(
//		conf.URL,
//	).WithConnectionParams(map[string]interface{}{
//		"headers": map[string]string{
//			"Sec-WebSocket-Protocol": "graphql-ws",
//			"Content-Type":           "application/json",
//		},
//	}).
//		WithProtocol(graphql.GraphQLWS).
//		//WithLog(log.Println).
//		WithoutLogTypes(graphql.GQLData, graphql.GQLConnectionKeepAlive).
//		OnError(func(sc *graphql.SubscriptionClient, err error) error {
//			return err
//		})
//	return &Client{
//		c: client,
//	}
//}
//
//func (t *Client) Close() error {
//	return t.c.Close()
//}
//
//func (t *Client) Run() {
//	t.c.Run()
//}
//
//func (t *Client) Register(query string, variables map[string]interface{}, handler func(data []byte, err error) error) error {
//	_, err := t.c.Exec(query, variables, handler)
//	return err
//
//}
