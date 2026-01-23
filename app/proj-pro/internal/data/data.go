package data

import (
	"store/app/proj-pro/configs"
	"store/app/proj-pro/internal/data/cache"
	"store/app/proj-pro/internal/data/repo"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/mgz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/bytedance/arkr"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"store/pkg/sdk/third/bytedance/volcengine"
	"store/pkg/sdk/third/douyin"
	dypenapi "store/pkg/sdk/third/douyin/openapi"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/openaiz"
	"store/pkg/sdk/third/qiniu"
	"store/pkg/sdk/third/wavespeed"
	"store/pkg/sdk/third/wxdev"
	"store/pkg/sdk/third/wxdev/wxpay"
)

type Data struct {
	Mongo        *repo.Collections
	Cache        *cache.Caches
	GenaiFactory *gemini.GenaiFactory
	//Elastics     *elastics.Client
	Douyin *douyin.Client
	//AlismsClient  *alisms.Client
	//VolcSmsClient *sms.Client
	WXDev       *wxdev.Client
	WXPay       *wxpay.Client
	Redis       *rediz.RedisClient
	DyOpenapi   *dypenapi.Client
	Qiniu       *qiniu.Client
	GrpcClients *grpcz.Clients
	Arkr        *arkr.Client
	//Alioss       *alioss.Client
	TOS        *tos.Client
	Volcengine *volcengine.Client
	OpenAI     *openaiz.Client
	GetGoAPI   *openaiz.Client
	VikingDB   *vikingdb.Client
	Wavespeed  *wavespeed.Client

	Conf confcenter.Config[configs.BizConfig]
}

func NewData(c confcenter.Config[configs.BizConfig]) (*Data, func(), error) {

	//ck := ch.NewClient(c.Database.Ch)
	mg, err := mgz.Database(c.Database.Mongo)
	if err != nil {
		return nil, nil, err
	}

	//client, err := alisms.NewClient(c.Component.Alisms)
	//if err != nil {
	//	return nil, nil, err
	//}

	dyopenai, err := dypenapi.NewClient(dypenapi.Config{
		ClientKey:    "awdihb5n6xqjs2lm",
		ClientSecret: "2ff1fe85e1840020d68ccd92b23fb844",
	})
	if err != nil {
		return nil, nil, err
	}

	clients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	rd := rediz.NewRedisClient(c.Database.Rediz)
	d := &Data{
		Mongo:        repo.NewCollections(mg),
		Cache:        cache.NewCaches(rd),
		GenaiFactory: gemini.NewGenaiFactory(&c.Component.Genai),
		//Elastics:     elastics.NewClient(c.Database.Elastics),
		Douyin: douyin.NewClient(),
		//AlismsClient:  client,
		//VolcSmsClient: sms.NewClient(),
		WXDev:       wxdev.NewClient(),
		WXPay:       wxpay.NewClient(),
		Redis:       rd,
		DyOpenapi:   dyopenai,
		Qiniu:       qiniu.NewClient(c.Database.Qiniu),
		GrpcClients: clients,
		Arkr:        arkr.NewClient(),
		//Alioss:       alioss.NewClient(c.Database.Oss),
		TOS:        tos.NewClient(c.Database.Tos),
		Volcengine: volcengine.NewClient(),
		OpenAI:     openaiz.NewClient(c.Component.OpenAI),
		VikingDB:   vikingdb.NewClient(),
		Wavespeed:  wavespeed.NewClient(),
		GetGoAPI:   openaiz.NewClient(c.Component.GetGoAPI),
		Conf:       c,
	}

	cleanup := func() {
	}
	return d, cleanup, nil
}
