package demo_spi_response

import doudian_sdk "store/pkg/sdk/third/bytedance/dou/sdk-golang/core"

type DemoSpiResponse struct {
	doudian_sdk.BaseDoudianOpSpiResponse
	Data *DemoSpiData `json:"data"`
}

func (d *DemoSpiResponse) GetData() interface{} {
	return d.Data
}

type DemoSpiData struct {
	Data1 string
}
