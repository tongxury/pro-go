package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"store/app/payment/internal/data/repo/ent"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"time"
)

func NewEntClient(conf confcenter.Mysql) *ent.Client {

	client, err := ent.Open(conf.Driver, conf.Source)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}

	// todo 容易把字段清掉 还没看原因
	//Run the auto migration tool.
	err = client.Schema.Create(
		context.Background(),
		//migrate.WithDropIndex(true),
		//migrate.WithDropColumn(true),
		//migrate.WithForeignKeys(false),
	)

	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

type Repos struct {
	EntClient   *ent.Client
	RedisClient *rediz.RedisClient
	LocalCache  *cache.Cache
	Payment     *PaymentRepo
}

func NewRepos(entClient *ent.Client, redisClient *rediz.RedisClient, payment *PaymentRepo) *Repos {
	return &Repos{
		EntClient:   entClient,
		RedisClient: redisClient,
		Payment:     payment,
		LocalCache:  cache.New(5*time.Minute, 10*time.Minute),
	}
}
