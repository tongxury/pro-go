package material_editMaterial_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type MaterialEditMaterialResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *MaterialEditMaterialData `json:"data"`
}
type MaterialEditMaterialData struct {
}
