package openCloud_v1_imageVersion_create_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenCloudV1ImageVersionCreateResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OpenCloudV1ImageVersionCreateData `json:"data"`
}
type OpenCloudV1ImageVersionCreateData struct {
	// 镜像版本id
	VersionId string `json:"version_id"`
}
