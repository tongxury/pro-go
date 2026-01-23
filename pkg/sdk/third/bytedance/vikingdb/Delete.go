package vikingdb

import (
	"context"
	"encoding/json"
	"errors"
	"store/pkg/sdk/conv"

	"github.com/go-resty/resty/v2"
)

type DeleteRequest struct {
	CollectionName string
	IDs            []string
	DeleteAll      bool
}

func (t *Client) Delete(ctx context.Context, req DeleteRequest) error {

	body := map[string]interface{}{
		"collection_name": req.CollectionName,
		"ids":             req.IDs,
		"del_all":         req.DeleteAll,
	}

	headers := sign(t.conf.AccessKeyID, t.conf.AccessKeySecret, t.conf.Service,
		t.conf.Region, t.conf.Host,
		"POST",
		"/api/vikingdb/data/delete",
		nil,
		conv.M2B(body),
	)

	post, err := resty.New().R().SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Post(t.conf.BaseURL + "/api/vikingdb/data/delete")
	if err != nil {
		return err
	}

	var res Response[DeleteResponse]
	err = json.Unmarshal(post.Body(), &res)
	if err != nil {
		return err
	}

	if res.Code != "Success" {
		return errors.New(res.Code)
	}

	return nil
}

type DeleteResponse struct {
}
