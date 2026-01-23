package biz

import (
	"store/app/aiagent/internal/data"
)

type TrackingEvent struct {
	data *data.Data
}

func NewTrackingEvent(data *data.Data) *TrackingEvent {
	return &TrackingEvent{
		data: data,
	}
}

//func (t *TrackingEvent) InsertMany(ctx context.Context, params []*trackerpb.TrackingEvent, extra *trackerpb.Extra) error {
//
//	parseSize := func(src string) (int, int) {
//		resParts := strings.Split(src, "x")
//		if len(resParts) != 2 {
//			return 0, 0
//		}
//		return conv.Int(resParts[0]), conv.Int(resParts[1])
//	}
//
//	for _, x := range params {
//
//		sw, sh := parseSize(x.Res)
//		vw, vh := parseSize(x.Vp)
//		dw, dh := parseSize(x.Ds)
//
//		cx, _ := base64.StdEncoding.DecodeString(x.Cx)
//		//if err != nil {
//		//	return err
//		//}
//
//		userEventContext, _ := base64.StdEncoding.DecodeString(x.UePx)
//		//if err != nil {
//		//	return err
//		//}
//
//		//log.Debug(string(cx))
//
//		//var ctxData trackerpb.Context
//		//err = conv.J2S(cx, &ctxData)
//		//if err != nil {
//		//	return err
//		//}
//
//		y := dbtypes.TrackingEvent{
//			Event:            x.E,
//			EventId:          x.Eid,
//			AppId:            x.Aid,
//			TrackerName:      x.Tna,
//			Platform:         x.P,
//			DomainSessionIdx: conv.Int(x.Vid),
//			SessionId:        x.Sid,
//			DeviceId:         x.Duid,
//			Page:             x.Page,
//			Referrer:         x.Refr,
//			TrackerVersion:   x.Tv,
//			CreatedAt:        time.UnixMilli(conv.Int64(x.Dtm)),
//			SendAt:           time.UnixMilli(conv.Int64(x.Stm)),
//			PermitCookie:     x.Cookie == "1",
//			DocCharset:       x.Cs,
//			Lang:             x.Lang,
//			ScreenWidth:      sw,
//			ScreenHeight:     sh,
//			ViewWidth:        vw,
//			ViewHeight:       vh,
//			DocWidth:         dw,
//			DocHeight:        dh,
//			ColorDepth:       conv.Int(x.Cd),
//			Context:          string(cx),
//			Ip:               extra.Ip,
//			UserEventContext: string(userEventContext),
//		}
//
//		err := t.data.Repos.TrackingEvent.InsertMany(ctx, &y)
//		if err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//}
