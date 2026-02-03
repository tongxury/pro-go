package data

import (
	"store/app/voiceagent/configs"
	"store/app/voiceagent/internal/data/repo/mongodb"
	"store/confs"
	"store/pkg/clients"
	"store/pkg/clients/grpcz"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/bytedance/sms"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/cartesia"
	"store/pkg/sdk/third/elevenlabs"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/tikhub"

	lksdk "github.com/livekit/server-sdk-go/v2"
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
	Cartesia      *cartesia.Client
	Tikhub        *tikhub.Client
	Gemini        *gemini.GenaiFactory
	KafkaClient   *clients.KafkaClient
	RoomClient    *lksdk.RoomServiceClient

	Conf confcenter.Config[configs.BizConfig]
}

func NewData(c confcenter.Config[configs.BizConfig]) (*Data, func(), error) {

	grpcClients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		Mongo:       mongodb.NewCollections(c.Database.Mongo),
		Redis:       rediz.NewRedisClient(c.Database.Rediz),
		GrpcClients: grpcClients,
		//Alioss:      alioss.NewClient(c.Database.Oss),
		TOS: tos.NewClient(c.Database.Tos),
		//Alisms:      client,
		VolcSmsClient: sms.NewClient(),
		Tikhub:        tikhub.NewClient(),
		ElevenLabs:    elevenlabs.NewClient(),
		Cartesia:      cartesia.NewClient(confs.CartesiaKey),
		Gemini:        gemini.NewGenaiFactory(&c.Component.Genai),
		KafkaClient:   clients.NewKafkaClient(c.Component.Kafka),
		RoomClient:    lksdk.NewRoomServiceClient(confs.LiveKitUrl, confs.LiveKitApiKey, confs.LiveKitApiSecret),

		Conf: c,
	}

	cleanup := func() {
	}
	return d, cleanup, nil
}
