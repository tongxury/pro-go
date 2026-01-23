package data

import (
	"store/app/voiceagent/configs"
	"store/app/voiceagent/internal/data/repo/mongodb"
	"store/pkg/clients/grpcz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/bytedance/sms"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/elevenlabs"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/tikhub"
)

type Data struct {
	Mongo       *mongodb.Collections
	Redis       *rediz.RedisClient
	GrpcClients *grpcz.Clients
	//Alioss      *alioss.Client
	TOS *tos.Client
	//Alisms      *alisms.Client
	VolcSmsClient *sms.Client
	ElevenLabs    *elevenlabs.Client
	Tikhub        *tikhub.Client
	Gemini        *gemini.GenaiFactory

	Conf confcenter.Config[configs.BizConfig]
}

func NewData(c confcenter.Config[configs.BizConfig]) (*Data, func(), error) {

	clients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Mongo:       mongodb.NewCollections(c.Database.Mongo),
		Redis:       rediz.NewRedisClient(c.Database.Rediz),
		GrpcClients: clients,
		//Alioss:      alioss.NewClient(c.Database.Oss),
		TOS: tos.NewClient(c.Database.Tos),
		//Alisms:      client,
		VolcSmsClient: sms.NewClient(),
		Tikhub:        tikhub.NewClient(),
		ElevenLabs:    elevenlabs.NewClient(),
		Gemini:        gemini.NewGenaiFactory(&c.Component.Genai),
		Conf:          c,
	}

	cleanup := func() {
	}
	return d, cleanup, nil
}
