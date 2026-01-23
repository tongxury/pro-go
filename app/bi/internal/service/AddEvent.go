package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	bipb "store/api/bi"
	responsepb "store/api/public/response"
	"store/app/bi/internal/biz"
	"store/app/bi/internal/data"
	"store/app/bi/internal/data/repo"
	"store/pkg/krathelper"
	"time"
)

type BiService struct {
	bipb.UnimplementedBiServer
	eventLog *biz.EventLogBiz
	data     *data.Data
}

func NewBiService(eventLogBiz *biz.EventLogBiz, data *data.Data) *BiService {
	return &BiService{
		eventLog: eventLogBiz,
		data:     data,
	}
}

func (t BiService) AddEventLog(ctx context.Context, params *bipb.AddEventLogParams) (*responsepb.RedirectResponse, error) {

	userId := krathelper.FindUserId(ctx)

	now := time.Now()

	deviceId := params.GetCommonParams().GetDeviceId()
	channel := params.GetCommonParams().GetChannel()
	platform := params.GetCommonParams().GetPlatform()

	log.Debugw("AddEventLog ", "", "params", params)

	t.eventLog.OnEvent(ctx, params)

	err := t.data.Repos.EventLog.AsyncInsert(ctx, repo.EventLog{
		EventName:   params.Name,
		Ip:          krathelper.ClientPublicIP(ctx),
		CountryCode: krathelper.ClientCountryCode(ctx),
		CreatedAt:   now,
		UserID:      userId,
		DeviceID:    deviceId,
		Channel:     channel,
		Platform:    platform,
		Idfa:        params.GetCommonParams().GetIdfa(),
		Imei:        params.GetCommonParams().GetImei(),
		Oaid:        params.GetCommonParams().GetOaid(),
	})

	if err != nil {
		log.Errorw("AddEventLog err", err, "params", params)
		return nil, err
	}

	//if params.Name == "visit" && deviceId != "" {
	//	set, err := t.data.Repos.RedisClient.SetNX(ctx, deviceId, "1", time.Hour*24*365).Result()
	//	if err != nil {
	//		log.Errorw("AddEventLog err", err, "params", params)
	//		return nil, err
	//	}
	//
	//	if set {
	//		err := t.data.XhsClient.Send(ctx, "7688267", clickId, "101")
	//		if err != nil {
	//			log.Errorw("AddEventLog XhsClient err", err, "params", params)
	//			return nil, err
	//		}
	//
	//		log.Debugw("AddEventLog XhsClient ", "done", "params", params)
	//
	//	}
	//}

	return &responsepb.RedirectResponse{
		Cookies: []*responsepb.Cookie{
			{
				Name:   "device_id", // todo 目前set-cookie不成功
				Value:  params.GetCommonParams().GetDeviceId(),
				Domain: t.data.Conf.Biz.Domain,
			},
		},
	}, nil
}
