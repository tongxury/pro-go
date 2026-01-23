package elastics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/deletebyquery"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"store/pkg/clients/elastics"
)

type Client struct {
	conf elastics.Config
	c    *elasticsearch.TypedClient
}

func NewClient(conf elastics.Config) *Client {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: conf.Addresses,
		Username:  conf.Username,
		Password:  conf.Password,
		//CloudID:                  "",
		//APIKey:                   "",
		//ServiceToken:             "",
		//CertificateFingerprint:   "",
		//Header:                   nil,
		//CACert:                   nil,
		//RetryOnStatus:            nil,
		//DisableRetry:             false,
		//MaxRetries:               0,
		//RetryOnError:             nil,
		//CompressRequestBody:      false,
		//CompressRequestBodyLevel: 0,
		//PoolCompressor:           false,
		//DiscoverNodesOnStart:     false,
		//DiscoverNodesInterval:    0,
		//EnableMetrics:            false,
		//EnableDebugLogger:        false,
		//EnableCompatibilityMode:  false,
		//DisableMetaHeader:        false,
		//RetryBackoff:             nil,
		//Transport:                nil,
		//Logger:                   nil,
		//Selector:                 nil,
		//ConnectionPoolFunc:       nil,
		//Instrumentation:          nil,
	})

	if err != nil {
		panic(err)
	}

	//do, err := client.Ping().Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//
	//if !do {
	//	panic("failed to ping elasticsearch")
	//}

	return &Client{
		conf: conf,
		c:    client,
	}
}

func (t *Client) GetElasticClient() *elasticsearch.TypedClient {
	return t.c
}

func (t *Client) Create(ctx context.Context, index string, doc interface{}) error {

	_, err := t.c.Index(index).Document(doc).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *Client) CreateBulk(ctx context.Context, index string, docs []interface{}) error {
	if len(docs) == 0 {
		return nil
	}

	// 构建 NDJSON 格式的请求体
	var buf bytes.Buffer
	for _, doc := range docs {
		// 元数据行
		buf.WriteString(fmt.Sprintf(`{"create":{"_index":"%s"}}`, index))
		buf.WriteString("\n")

		// 文档数据行
		json.NewEncoder(&buf).Encode(doc)
	}

	// 执行批量请求
	res, err := t.c.Bulk().
		Index(index).
		Raw(bytes.NewReader(buf.Bytes())).
		Do(ctx)

	if err != nil {
		return err
	}

	// 检查是否有错误
	if res.Errors {
		return fmt.Errorf("批量创建部分失败")
	}

	return nil
}

//func (t *Client) Get(ctx context.Context, index string, id string) (*types.Hit, error) {
//
//	v2, err := t.SearchV2(ctx, index, search.Request{
//		Query: NewTermQuery("id", id),
//	})
//	if err != nil {
//		return err
//	}
//	for _, hit := range v2.Hits {
//
//
//	}
//
//
//	fmt.Println(v2)
//
//	//d, err := t.c.Index(index).Id(id).Do(ctx)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//log.Debugw(d)
//
//	return nil
//}

func (t *Client) UpdateFields(ctx context.Context, index, id string, fields map[string]any) error {
	marshal, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	_, err = t.c.Update(index, id).
		Request(&update.Request{
			Doc: marshal,
		}).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *Client) Replace(ctx context.Context, index, id string, doc any) error {
	marshal, _ := json.Marshal(doc)

	_, err := t.c.Update(index, id).
		Request(&update.Request{
			Doc: marshal,
		}).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (t *Client) Delete(ctx context.Context, index string, id string) error {
	_, err := t.c.Delete(index, id).Do(ctx)

	return err
}

func (t *Client) DeleteByRequest(ctx context.Context, index string, params deletebyquery.Request) error {
	_, err := t.c.DeleteByQuery(index).Request(&params).Do(ctx)

	return err
}

type SearchParams struct {
	Index        string
	MatchQueries map[string]types.MatchQuery
	MinMatches   int
	Size         int
}

func (t *Client) Search(ctx context.Context, index string, params search.Request) (*types.HitsMetadata, error) {
	s := t.c.Search().Index(index)

	s = s.Request(&params)

	res, err := s.Do(ctx)
	if err != nil {
		return nil, err
	}

	return &res.Hits, nil
}
