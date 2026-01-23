package volcengine

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

type StreamUploadDataRequest struct {
	Md5       string
	Size      int64
	Offset    int64
	OwnerId   int64
	OwnerType string
	Data      []byte
}

// CreateBytesMaterial 封装了从检查上传状态、上传数据到创建素材的完整流程
func (t *Client) CreateBytesMaterial(ctx context.Context, data []byte, info CreateMaterialInfo) (*CreateMaterialResult, error) {
	h := md5.New()
	h.Write(data)
	md5Str := hex.EncodeToString(h.Sum(nil))
	size := int64(len(data))

	owner := &Owner{
		Type: "PERSON",
		Id:   7541913281535950894,
	}

	// 1. 获取上传状态
	state, err := t.GetUploadState(ctx, GetUploadStateRequest{
		Owner: owner,
		Md5:   md5Str,
		Size:  size,
	})
	if err != nil {
		return nil, fmt.Errorf("get upload state failed: %w", err)
	}

	// 2. 如果未完成上传，则上传数据
	if !state.SkipDataComplete {
		err = t.StreamUploadData(ctx, StreamUploadDataRequest{
			Md5:       md5Str,
			Size:      size,
			Offset:    0,
			OwnerId:   owner.Id,
			OwnerType: owner.Type,
			Data:      data,
		})
		if err != nil {
			return nil, fmt.Errorf("stream upload data failed: %w", err)
		}
	}

	// 3. 创建素材
	res, err := t.CreateMaterial(ctx, CreateMaterialRequest{
		Owner: owner,
		StoreItem: &StoreItem{
			Md5:              md5Str,
			Size:             size,
			SkipDataComplete: true,
			FileExtension:    info.MediaExtension,
		},
		CreateMaterialInfo: &info,
	})
	if err != nil {
		return nil, fmt.Errorf("create material failed: %w", err)
	}

	return res, nil
}

func (t *Client) StreamUploadData(ctx context.Context, req StreamUploadDataRequest) error {
	_, err := t.doRequest(ctx,
		Req{
			Version: "2022-02-01",
			Action:  "StreamUploadData",
			Method:  "POST",
			Params: map[string]string{
				"Md5":       req.Md5,
				"Size":      fmt.Sprintf("%d", req.Size),
				"Offset":    fmt.Sprintf("%d", req.Offset),
				"OwnerId":   fmt.Sprintf("%d", req.OwnerId),
				"OwnerType": req.OwnerType,
			},
			Body: req.Data,
		},
	)
	return err
}

func (t *Client) GetUploadState(ctx context.Context, req GetUploadStateRequest) (*GetUploadStateResult, error) {
	body, _ := json.Marshal(req)
	data, err := t.doRequest(ctx,
		Req{
			Version: "2022-02-01",
			Action:  "GetUploadState",
			Method:  "POST",
			Body:    body,
		},
	)
	if err != nil {
		return nil, err
	}

	var res Response[GetUploadStateResult]
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if res.Code != 0 || res.ResponseMetadata.Code != 0 {
		return nil, errors.New(string(data))
	}

	return &res.Result, nil
}

func (t *Client) CreateMaterial(ctx context.Context, req CreateMaterialRequest) (*CreateMaterialResult, error) {
	body, _ := json.Marshal(req)
	data, err := t.doRequest(ctx,
		Req{
			Version: "2022-02-01",
			Action:  "CreateMaterial",
			Method:  "POST",
			Body:    body,
		},
	)
	if err != nil {
		return nil, err
	}

	var res Response[CreateMaterialResult]
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if res.Code != 0 || res.ResponseMetadata.Code != 0 {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
