package confcenter

import (
	"store/pkg/clients/grpcz"
	"store/pkg/sdk/third/airwallex"
	"store/pkg/sdk/third/aliyun/alisms"
	"store/pkg/sdk/third/aws/awsmail"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/openaiz"
	"store/pkg/sdk/third/stripe"
)

type Component struct {
	AwsMail   awsmail.Config
	Alisms    alisms.Config
	AirWallex airwallex.Config
	Stripe    stripe.Config
	StripePay StripePay
	Kafka     KafkaConfig
	Solana    SolanaConfig
	Genai     gemini.FactoryConfig
	Grpc      grpcz.Configs
	OpenAI    openaiz.Config
	GetGoAPI  openaiz.Config
}

type StripePay struct {
	Prod stripe.Config
	Test stripe.Config
	//ProdV1     stripe.Config
	//TestEmails []string
}

type SolanaConfig struct {
	Endpoint   string
	WSEndpoint string
}

type KafkaConfig struct {
	Brokers  []string
	Username string
	Password string
}
