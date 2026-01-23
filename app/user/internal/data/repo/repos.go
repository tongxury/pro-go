package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"store/app/user/internal/data/repo/ent"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"time"
)

func NewEntClient(conf confcenter.Mysql) *ent.Client {

	client, err := ent.Open(conf.Driver, conf.Source)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}

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
	EntClient       *ent.Client
	RedisClient     *rediz.RedisClient
	Member          *MemberRepo
	MemberSubscribe *MemberSubscribeRepo
	User            *UserRepo
	Promotion       *PromotionRepo
}

func NewRepos(conf confcenter.Database) *Repos {

	entClient := NewEntClient(conf.Mysql)
	redisClient := rediz.NewRedisClient(conf.Rediz)
	lcClient := cache.New(10*time.Minute, 10*time.Minute)

	return &Repos{
		EntClient:   NewEntClient(conf.Mysql),
		RedisClient: redisClient,
		User:        NewUserRepo(entClient, redisClient, lcClient),
	}
}
