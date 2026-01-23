package data

import (
	"store/app/proj-admin/configs"
	"store/app/proj-admin/internal/data/repo"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/bytedance/arkr"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"store/pkg/sdk/third/gemini"
)

type Data struct {
	Mongo        *repo.Collections
	GenaiFactory *gemini.GenaiFactory
	//Elastics     *elastics.Client
	//Alioss   *alioss.Client
	TOS         *tos.Client
	Arkr        *arkr.Client
	VikingDB    *vikingdb.Client
	GrpcClients *grpcz.Clients
}

func NewData(c confcenter.Config[configs.BizConfig]) (*Data, func(), error) {

	//ck := ch.NewClient(c.Database.Ch)
	mg, err := mgz.Database(c.Database.Mongo)
	if err != nil {
		return nil, nil, err
	}

	grpcClients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		panic(err)
	}

	rz := rediz.NewRedisClient(c.Database.Rediz)

	//
	rc := NewRedisCache(rz)
	gf := &c.Component.Genai
	for _, x := range gf.Configs {
		x.Cache = rc
	}

	d := &Data{
		Mongo:        repo.NewCollections(mg),
		GenaiFactory: gemini.NewGenaiFactory(gf),
		//Elastics:     elastics.NewClient(c.Database.Elastics),
		//Alioss:   alioss.NewClient(c.Database.Oss),
		TOS:         tos.NewClient(c.Database.Tos),
		Arkr:        arkr.NewClient(),
		VikingDB:    vikingdb.NewClient(),
		GrpcClients: grpcClients,
	}

	cleanup := func() {
	}
	return d, cleanup, nil
}
