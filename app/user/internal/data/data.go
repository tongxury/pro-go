package data

import (
	"store/app/user/internal/conf"
	"store/app/user/internal/data/repo"
	"store/pkg/clients"
	"store/pkg/clients/grpcz"
	"store/pkg/confcenter"
	"store/pkg/sdk/mail"
	"store/pkg/sdk/third/aliyun/alisms"
	"store/pkg/sdk/third/aws/awsmail"
	"store/pkg/sdk/third/wxdev"
	"store/pkg/sdk/third/xhs"
)

// Data .
type Data struct {
	Config        confcenter.Config[conf.BizConfig]
	BizConfig     conf.BizConfig
	Repos         *repo.Repos
	GrpcClients   *grpcz.Clients
	KafkaClient   *clients.KafkaClient
	AwsMailClient *awsmail.Client
	MailClient    *mail.Client
	AlismsClient  *alisms.Client
	XhsClient     *xhs.Client
	WXDev         *wxdev.Client
}

func NewData(conf confcenter.Config[conf.BizConfig]) (*Data, func(), error) {

	grpcClients, err := grpcz.NewClients(conf.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	alismsClient, err := alisms.NewClient(conf.Component.Alisms)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Config:        conf,
		BizConfig:     conf.Biz,
		Repos:         repo.NewRepos(conf.Database),
		GrpcClients:   grpcClients,
		KafkaClient:   clients.NewKafkaClient(conf.Component.Kafka),
		AwsMailClient: awsmail.NewClient(conf.Component.AwsMail),
		MailClient:    mail.NewClient(mail.Config{}),
		AlismsClient:  alismsClient,
		XhsClient:     xhs.NewClient(),
		WXDev:         wxdev.NewClient(),
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
