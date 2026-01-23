package data

import (
	"store/app/usercenter/configs"
	"store/app/usercenter/internal/data/repo"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/aliyun/alioss"
	"store/pkg/sdk/third/aliyun/alisms"
	"store/pkg/sdk/third/bytedance/sms"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/tikhub"
)

type Data struct {
	Mongo         *repo.Collections
	Redis         *rediz.RedisClient
	GrpcClients   *grpcz.Clients
	Alioss        *alioss.Client
	Alisms        *alisms.Client
	VolcSmsClient *sms.Client

	Tikhub *tikhub.Client
	TOS    *tos.Client

	Conf confcenter.Config[configs.BizConfig]
}

func NewData(c confcenter.Config[configs.BizConfig]) (*Data, func(), error) {

	//ck := ch.NewClient(c.Database.Ch)
	mg, err := mgz.Database(c.Database.Mongo)
	if err != nil {
		return nil, nil, err
	}

	client, err := alisms.NewClient(c.Component.Alisms)
	if err != nil {
		return nil, nil, err
	}

	clients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Mongo:         repo.NewCollections(mg),
		Redis:         rediz.NewRedisClient(c.Database.Rediz),
		GrpcClients:   clients,
		Alioss:        alioss.NewClient(c.Database.Oss),
		Alisms:        client,
		VolcSmsClient: sms.NewClient(),
		Tikhub:        tikhub.NewClient(),
		TOS:           tos.NewClient(c.Database.Tos),
		Conf:          c,
	}

	cleanup := func() {
	}
	return d, cleanup, nil
}
