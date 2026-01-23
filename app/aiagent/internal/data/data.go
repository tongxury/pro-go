package data

import (
	"store/app/aiagent/internal/conf"
	"store/app/aiagent/internal/data/repo/mongodb"
	"store/pkg/clients/grpcz"
	"store/pkg/clients/xhs"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/tikhub"
	"time"

	"store/confs"

	"github.com/patrickmn/go-cache"
	claude "github.com/potproject/claude-sdk-go"
)

type Repos struct {
	Uploader *Uploader
}

type Data struct {
	Conf  confcenter.Config[conf.BizConfig]
	Repos Repos
	Mongo *mongodb.Collections
	Redis *rediz.RedisClient
	//GenaiClient *genai.Client
	//GenaiClient2         *genai.Client
	GenaiFactory *gemini.GenaiFactory
	//Uploader     *Uploader
	//S3Client     *s3.Client
	//Alioss       *alioss.Client
	//Qiniu  *qiniu.Client
	Claude    *claude.Client
	TosClient *tos.Client
	//DeepSeek    *openai.Client
	//ThirdAI     *openai.Client
	GrpcClients *grpcz.Clients
	XhsClient   *xhs.Client
	LocalCache  *cache.Cache
	Tikhub      *tikhub.Client
}

func NewData(c confcenter.Config[conf.BizConfig]) (*Data, func(), error) {

	grpcClients, err := grpcz.NewClients(c.Component.Grpc)
	if err != nil {
		panic(err)
	}

	d := &Data{
		Conf:         c,
		Mongo:        mongodb.NewCollections(c.Database.Mongo),
		Redis:        rediz.NewRedisClient(c.Database.Rediz),
		GenaiFactory: gemini.NewGenaiFactory(&c.Component.Genai),
		//Uploader:     NewUploader(alioss.NewClient(c.Database.Oss)),
		//S3Client:     s3.NewS3Client(c.Database.S3),
		//Alioss:       alioss.NewClient(c.Database.Oss),
		//Qiniu: qiniu.NewClient(c.Database.Qiniu),
		Claude:      claude.NewClient(confs.ClaudeKey),
		GrpcClients: grpcClients,
		XhsClient:   xhs.NewClient(),
		LocalCache:  cache.New(time.Minute*10, time.Minute*5),
		Tikhub:      tikhub.NewClient(),
	}

	cleanup := func() {
	}
	return d, cleanup, nil
}
