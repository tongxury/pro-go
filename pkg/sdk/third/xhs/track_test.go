package xhs

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {

	c := NewClient()
	//c.Send(context.Background(), "7688267", "atest67594118e4b02347c3c4db6f", "102")
	c.SendV2(context.Background(), SendV2Params{
		Platform:     "veogo",
		Timestamp:    1753376341000,
		Scene:        "701",
		AdvertiserId: "7676953",
		Os:           1,
		Caid1Md5:     "atest68825e80e4b0234776db80e6",
		ClickId:      "atest6882663de4b023476dc41a80",
		EventType:    411,
		Context: ConversionContext{
			Properties: ConversionProperties{},
		},
	})
}
