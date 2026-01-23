package dou

import doudian_sdk "store/pkg/sdk/third/bytedance/dou/sdk-golang/core"

type Product struct {
}

func (t *Client) GetProductDetail() (*Product, error) {

	accessToken, err := doudian_sdk.BuildAccessToken(&doudian_sdk.BuildAccessTokenParam{
		AuthId: &shopId, AuthSubjectType: &subId,
	})
	if err != nil {
		panic(err)
	}

	return nil, nil
}
