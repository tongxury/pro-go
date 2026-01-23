package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"store/pkg/events"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/third/xhs"
)

func (t *EventLogBiz) OnRegisterEvent(ctx context.Context, event events.AuthEvent) error {

	log.Debug("OnRegisterEvent ing", event)

	e, err := t.data.Repos.JuGuangEvent.FindByDeviceId(ctx, event.DeviceID)
	if err != nil {
		return err
	}

	if e == nil {
		return nil
	}

	//APP激活	411	APP激活数量
	//APP注册	412	APP注册数量
	//APP付费	413	APP付费数量
	//关键行为	414	用户在推广的应用/落地页场景下发生的关键行为/行为集合，具体行为取决于广告主业务模式
	//次留	450	次留量
	//三日留存	453	三日留存量
	//七日留存	457	七日留存量

	params := xhs.SendV2Params{
		Platform:     "veogo",
		Timestamp:    event.TS * 1000,
		Scene:        "701",
		AdvertiserId: e.AdvertiserId,
		Os:           conv.Int(e.Os),
		Caid1Md5:     e.CaidMd5,
		ClickId:      e.ClickId,
		EventType:    412,
		Context: xhs.ConversionContext{
			Properties: xhs.ConversionProperties{
				ConvCnt:           1,
				ActivateTimestamp: event.TS * 1000,
			},
		},
	}

	log.Debug("XhsClient SendV2 ing", params)

	err = t.data.XhsClient.SendV2(ctx, params)
	if err != nil {
		log.Errorw("SendV2 err", err, "event", event)
		return err
	}

	return nil
}
