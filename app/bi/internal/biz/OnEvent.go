package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	bipb "store/api/bi"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/xhs"
	"time"
)

func (t *EventLogBiz) OnEvent(ctx context.Context, params *bipb.AddEventLogParams) {

	idfa := params.GetCommonParams().GetIdfa()
	//oaid := params.GetCommonParams().GetOaid()
	deviceId := params.GetCommonParams().GetDeviceId()

	helper.Go(ctx, func(ctx context.Context) {
		if params.Name != "visit" {
			return
		}

		event, err := t.data.Repos.JuGuangEvent.FindByIdfaMd5(ctx, helper.MD52(idfa))
		if err != nil {
			return
		}

		if event == nil || event.DeviceId != "" {
			log.Debugw("FindByIdfaMd5 empty ", "", "idfa", idfa)
			return
		}

		// 聚光event中没有给deviceId 补充一下
		event.DeviceId = deviceId
		err = t.data.Repos.JuGuangEvent.AsyncInsert(ctx, *event)
		if err != nil {
			log.Errorw("AsyncInsert err", err, "event", event)
			return
		}

		//APP激活	411	APP激活数量
		//APP注册	412	APP注册数量
		//APP付费	413	APP付费数量
		//关键行为	414	用户在推广的应用/落地页场景下发生的关键行为/行为集合，具体行为取决于广告主业务模式
		//次留	450	次留量
		//三日留存	453	三日留存量
		//七日留存	457	七日留存量
		now := time.Now().UnixMilli()

		xhsParams := xhs.SendV2Params{
			Platform:     "veogo",
			Timestamp:    time.Now().UnixMilli(),
			Scene:        "701",
			AdvertiserId: event.AdvertiserId,
			Os:           conv.Int(event.Os),
			Caid1Md5:     event.CaidMd5,
			ClickId:      event.ClickId,
			EventType:    411,
			Context: xhs.ConversionContext{
				Properties: xhs.ConversionProperties{
					ActivateTimestamp: now,
				},
			},
		}

		log.Debug("XhsClient SendV2 ing", xhsParams)

		err = t.data.XhsClient.SendV2(ctx, xhsParams)
		if err != nil {
			log.Errorw("SendV2 err", err, "event", event)
		}

	})

}
