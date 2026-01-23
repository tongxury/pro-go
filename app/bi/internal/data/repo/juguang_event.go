package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"store/app/bi/internal/data/mongodb"
	"store/app/bi/internal/data/types"
)

type JuGuangEventRepo struct {
	//db    *clients.ClickHouseClient
	mongo *mongodb.Collections
}

func NewJuGuangEventRepo(mongo *mongodb.Collections) *JuGuangEventRepo {
	return &JuGuangEventRepo{mongo: mongo}
}

//func (t *JuGuangEventRepo) FindByClickId(ctx context.Context, clickId string) (*types.JuGuangEvent, error) {
//
//	var result []types.JuGuangEvent
//	err := t.db.Select(ctx, &result, `select * from juguang_events where clickId = ?`, clickId)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(result) == 0 {
//		return nil, nil
//	}
//
//	return &result[0], nil
//}

func (t *JuGuangEventRepo) FindByDeviceId(ctx context.Context, deviceId string) (*types.JuGuangEvent, error) {

	list, err := t.mongo.JuGuang.List(ctx, bson.M{"deviceId": deviceId}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil

	//var result []types.JuGuangEvent
	//err := t.db.Select(ctx, &result, `select * from juguang_events where deviceId = ? order by createdAt desc`, deviceId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(result) == 0 {
	//	return nil, nil
	//}
	//
	//return &result[0], nil
}

func (t *JuGuangEventRepo) FindByIdfaMd5(ctx context.Context, idfaMd5 string) (*types.JuGuangEvent, error) {

	list, err := t.mongo.JuGuang.List(ctx, bson.M{"idfaMd5": idfaMd5})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
	//var result []types.JuGuangEvent
	//err := t.db.Select(ctx, &result, `select * from juguang_events where idfaMd5 = ? order by createdAt desc`, idfaMd5)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(result) == 0 {
	//	return nil, nil
	//}
	//
	//return &result[0], nil

	//return nil, nil
}

//func (t *JuGuangEventRepo) UpdateDeviceId(ctx context.Context, keyword, deviceId string) error {
//
//	//t.mongo.JuGuang.UpdateByIDXX()
//
//	//var result []JuGuangEvent
//	//err := t.db.Select(ctx, &result, `select * from juguang_events where idfaMd5 = ?`, keyword)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//if len(result) == 0 {
//	//	return nil
//	//}
//	//
//	//for _, x := range result {
//	//	if x.DeviceId == deviceId {
//	//		return nil
//	//	}
//	//}
//	//
//	//result[0].DeviceId = deviceId
//	//err = t.AsyncInsert(ctx, result[0])
//	//if err != nil {
//	//	return err
//	//}
//
//	return nil
//}

func (t *JuGuangEventRepo) AsyncInsert(ctx context.Context, e types.JuGuangEvent) error {
	return t.mongo.JuGuang.Insert(ctx, &e)

	//err := t.db.Exec(ctx, `
	//    insert into juguang_events (
	//        clickId,
	//        caidMd5,
	//        paid,
	//        caid,
	//        ua,
	//        ua1,
	//        idfaMd5,
	//        os,
	//        ip,
	//        ts,
	//        appId,
	//        advertiserId,
	//        unitId,
	//        campaignId,
	//        creativityId,
	//        placement,
	//        androidId,
	//        oaid,
	//        oaidMd5,
	//        imei,
	//        deviceId,
	//        createdAt
	//    )
	//    values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	//	e.ClickId,
	//	e.CaidMd5,
	//	e.Paid,
	//	e.Caid,
	//	e.Ua,
	//	e.Ua1,
	//	e.IdfaMd5,
	//	e.Os,
	//	helper.OrString(e.Ip, "0.0.0.0"),
	//	e.Ts.Format(time.DateTime),
	//	e.AppId,
	//	e.AdvertiserId,
	//	e.UnitId,
	//	e.CampaignId,
	//	e.CreativityId,
	//	e.Placement,
	//	e.AndroidId,
	//	e.Oaid,
	//	e.OaidMd5,
	//	e.Imei,
	//	e.DeviceId,
	//	time.Now().Unix(),
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}
