package data

import (
	"store/app/bi/internal/conf"
	"store/app/bi/internal/data/mongodb"
	"store/app/bi/internal/data/repo"
	"store/pkg/clients"
	"store/pkg/clients/grpcz"
	"store/pkg/confcenter"
	"store/pkg/consts"
	"store/pkg/sdk/third/feishuo"
	"store/pkg/sdk/third/stripe"
	"store/pkg/sdk/third/xhs"
)

// Data .
type Data struct {
	//BizConfig    conf.BizConfig
	Conf        confcenter.Config[conf.BizConfig]
	Meta        confcenter.Meta
	Repos       *repo.Repos
	GrpcClients *grpcz.Clients
	//SyncClient     *clients.CanalClient
	StripeClient   *stripe.Client
	StripeClientV1 *stripe.Client
	KafkaClient    *clients.KafkaClient
	FeishuClient   *feishuo.Client
	XhsClient      *xhs.Client
	Mongo          *mongodb.Collections
}

func NewData(
	conf confcenter.Config[conf.BizConfig],
) (*Data, func(), error) {

	grpcClients, err := grpcz.NewClients(conf.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Conf:        conf,
		Meta:        conf.Meta,
		Repos:       repo.NewRepos(conf.Database),
		GrpcClients: grpcClients,
		//SyncClient:     clients.NewCanalClient(conf.Biz.Sync),
		//StripeClientV1: stripe.NewStripeClient(conf.Component.StripePay.ProdV1),
		StripeClient: stripe.NewStripeClient(conf.Component.StripePay.Prod),
		KafkaClient:  clients.NewKafkaClient(conf.Component.Kafka),
		FeishuClient: feishuo.NewClient(consts.FEISHU_ENDPOINT2),
		XhsClient:    xhs.NewClient(),
		Mongo:        mongodb.NewCollections(conf.Database.Mongo),
	}
	return d, func() {
		//if err := d.DB.Close(); err != nil {
		//	log.Error(err)
		//}
		//if err := d.RedisCluster.Close(); err != nil {
		//	log.Error(err)
		//}
	}, nil
}
