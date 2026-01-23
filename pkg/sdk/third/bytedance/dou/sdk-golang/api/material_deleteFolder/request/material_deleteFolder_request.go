package material_deleteFolder_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/material_deleteFolder/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type MaterialDeleteFolderRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *MaterialDeleteFolderParam
}

func (c *MaterialDeleteFolderRequest) GetUrlPath() string {
	return "/material/deleteFolder"
}

func New() *MaterialDeleteFolderRequest {
	request := &MaterialDeleteFolderRequest{
		Param: &MaterialDeleteFolderParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *MaterialDeleteFolderRequest) Execute(accessToken *doudian_sdk.AccessToken) (*material_deleteFolder_response.MaterialDeleteFolderResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &material_deleteFolder_response.MaterialDeleteFolderResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *MaterialDeleteFolderRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *MaterialDeleteFolderRequest) GetParams() *MaterialDeleteFolderParam {
	return c.Param
}

type MaterialDeleteFolderParam struct {
	// 文件夹id list，最多不超过20个
	FolderIds []string `json:"folder_ids"`
}
