package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	bipb "store/api/bi"
	"store/app/bi/internal/data/types"
	"store/pkg/sdk/conv"
	"strings"
	"time"
)

func (t BiService) TrackJuGuangEvent(ctx context.Context, ev *bipb.TrackJuGuangEventParams) (*empty.Empty, error) {

	log.Debugw("TrackJuGuangEvent ", ev)

	empty_ := func(s string) string {
		if strings.HasPrefix(s, "__") {
			return ""
		}

		return s
	}

	if ev.ClickId == "" {
		return nil, errors.BadRequest("", "")
	}

	//if empty_(ev.IdfaMd5) != "" {
	//	t.data.Repos.RedisClient.Set(ctx, "clickId:"+ev.IdfaMd5, ev.ClickId, time.Hour)
	//}
	//
	//if empty_(ev.OaidMd5) != "" {
	//	t.data.Repos.RedisClient.Set(ctx, "clickId:"+ev.OaidMd5, ev.ClickId, time.Hour)
	//}

	err := t.data.Repos.JuGuangEvent.AsyncInsert(ctx, types.JuGuangEvent{
		ClickId:      empty_(ev.ClickId),
		CaidMd5:      empty_(ev.CaidMd5),              // 或 helper.MD5(ev.Caid)
		Paid:         empty_(ev.Paid),                 // 付费用户标识
		Caid:         empty_(ev.Caid),                 // 客户端广告ID
		Ua:           empty_(ev.Ua),                   // User-Agent
		Ua1:          empty_(ev.Ua),                   // 备用User-Agent
		IdfaMd5:      empty_(ev.IdfaMd5),              // 或 helper.MD5(ev.Idfa)
		Os:           empty_(ev.Os),                   // 操作系统
		Ip:           empty_(ev.Ip),                   // 客户端IP
		Ts:           time.Unix(conv.Int64(ev.Ts), 0), // 或 time.Now()
		AppId:        empty_(ev.AppId),                // 应用ID
		AdvertiserId: empty_(ev.AdvertiserId),         // 广告主ID
		UnitId:       empty_(ev.UnitId),               // 广告单元ID
		CampaignId:   empty_(ev.CampaignId),           // 活动ID
		CreativityId: empty_(ev.CreativityId),         // 创意ID
		Placement:    empty_(ev.Placement),            // 广告位置
		AndroidId:    empty_(ev.AndroidId),            // Android设备ID
		Oaid:         empty_(ev.Oaid),                 // 开放匿名设备标识符
		OaidMd5:      empty_(ev.OaidMd5),              // 或 helper.MD5(ev.Oaid)
		Imei:         empty_(ev.Imei),                 // 设备IMEI
	})
	if err != nil {
		log.Errorw("TrackJuGuangEvent err", err, "event", ev)
		return nil, err
	}

	//APP激活	411	APP激活数量
	//APP注册	412	APP注册数量
	//APP付费	413	APP付费数量
	//关键行为	414	用户在推广的应用/落地页场景下发生的关键行为/行为集合，具体行为取决于广告主业务模式
	//次留	450	次留量
	//三日留存	453	三日留存量
	//七日留存	457	七日留存量

	//getEventType := func(cid string) int {
	//	switch cid {
	//	case "":
	//		return 413
	//	default:
	//		return 413
	//	}
	//}
	//
	//err = t.data.XhsClient.SendV2(ctx, xhs.SendV2Params{
	//	Platform:     "veogo",
	//	Timestamp:    conv.Int64(ev.Ts),
	//	Scene:        "701",
	//	AdvertiserId: ev.AdvertiserId,
	//	Os:           conv.Int(ev.Os),
	//	Caid1Md5:     ev.CaidMd5,
	//	ClickId:      ev.ClickId,
	//	EventType:    getEventType(ev.CreativityId),
	//	Context: xhs.ConversionContext{
	//		Properties: xhs.ConversionProperties{
	//			ConvCnt:           1,
	//			Pay:               0,
	//			ActivateTimestamp: 0,
	//		},
	//	},
	//})
	//if err != nil {
	//	log.Errorw("TrackJuGuangEvent err", err, "event", ev)
	//	return nil, err
	//}

	return &empty.Empty{}, nil
}
