package data

import (
	"store/app/demo/configs"
	"store/app/demo/internal/data/repo"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
)

type Data struct {
	Mongo *repo.Collections
	Redis *rediz.RedisClient
	Grpc  *grpcz.Clients
	Conf  confcenter.Config[configs.BizConfig]
}

func NewData(c confcenter.Config[configs.BizConfig]) (*Data, func(), error) {
	mg, err := mgz.Database(c.Database.Mongo)
	if err != nil {
		return nil, nil, err
	}

	clients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Mongo: repo.NewCollections(mg),
		Redis: rediz.NewRedisClient(c.Database.Rediz),
		Grpc:  clients,
		Conf:  c,
	}

	cleanup := func() {}
	return d, cleanup, nil
}
