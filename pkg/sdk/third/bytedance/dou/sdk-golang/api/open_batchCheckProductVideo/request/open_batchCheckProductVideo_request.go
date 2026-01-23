package open_batchCheckProductVideo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/open_batchCheckProductVideo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenBatchCheckProductVideoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *OpenBatchCheckProductVideoParam
}

func (c *OpenBatchCheckProductVideoRequest) GetUrlPath() string {
	return "/open/batchCheckProductVideo"
}

func New() *OpenBatchCheckProductVideoRequest {
	request := &OpenBatchCheckProductVideoRequest{
		Param: &OpenBatchCheckProductVideoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *OpenBatchCheckProductVideoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*open_batchCheckProductVideo_response.OpenBatchCheckProductVideoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &open_batchCheckProductVideo_response.OpenBatchCheckProductVideoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *OpenBatchCheckProductVideoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *OpenBatchCheckProductVideoRequest) GetParams() *OpenBatchCheckProductVideoParam {
	return c.Param
}

type ProductVideoCheckListItem struct {
	// 商品id
	ProductId *int64 `json:"product_id"`
	// 视频素材id
	Vid *string `json:"vid"`
}
type OpenBatchCheckProductVideoParam struct {
	// 入参
	ProductVideoCheckList []ProductVideoCheckListItem `json:"product_video_check_list"`
}
