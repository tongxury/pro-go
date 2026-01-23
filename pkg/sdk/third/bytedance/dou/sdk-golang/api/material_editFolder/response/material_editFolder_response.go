package material_editFolder_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type MaterialEditFolderResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *MaterialEditFolderData `json:"data"`
}
type MaterialEditFolderData struct {
}
