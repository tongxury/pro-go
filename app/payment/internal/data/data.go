package data

import (
	"store/app/payment/internal/conf"
	"store/app/payment/internal/data/repo"
	"store/app/payment/internal/data/stripefactory"
	"store/pkg/clients"
	"store/pkg/clients/grpcz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/airwallex"
)

// Data .
type Data struct {
	Conf  confcenter.Config[bizconf.BizConfig]
	Repos *repo.Repos
	//Client  *stripe.StripeClient
	StripeFactory *stripefactory.StripeFactory
	GrpcClients   grpcz.Clients
	KafkaClient   *clients.KafkaClient
	AirWallex     *airwallex.Client
}

func NewData(conf confcenter.Config[bizconf.BizConfig]) (*Data, func(), error) {

	grpcClients, err := grpcz.NewClients(conf.Component.Grpc)

	entClient := repo.NewEntClient(conf.Database.Mysql)
	redisClient := rediz.NewRedisClient(conf.Database.Rediz)

	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Conf: conf,
		Repos: repo.NewRepos(entClient, redisClient,
			repo.NewPaymentRepo(entClient, redisClient),
		),
		//Client:  stripe.NewStripeClient(conf.Component.Stripe),
		StripeFactory: stripefactory.NewStripeFactory(conf.Component.StripePay),
		AirWallex:     airwallex.NewClient(conf.Component.AirWallex),
		GrpcClients:   *grpcClients,
		KafkaClient:   clients.NewKafkaClient(conf.Component.Kafka),
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
