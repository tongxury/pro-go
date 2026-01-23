package repo

import (
	"github.com/patrickmn/go-cache"
	"store/app/bi/internal/data/mongodb"
	"store/pkg/clients"
	"store/pkg/confcenter"
	"store/pkg/rediz"
	"time"
)

type Repos struct {
	ClickhouseClient *clients.ClickHouseClient
	RedisClient      *rediz.RedisClient
	JuGuangEvent     *JuGuangEventRepo
	EventLog         *EventLogRepo
	UserDevice       *UserDeviceRepo
	IpDetail         *IpDetailRepo
	ChatRecord       *ChatRecordRepo
	UserFunnelRecord *UserFunnelRecordRepo
	Stripe           *StripeRepo
}

func NewRepos(dbConf confcenter.Database) *Repos {

	ck := clients.NewClickHouseClient(dbConf.Clickhouse)
	lc := cache.New(24*time.Hour, 24*time.Hour)

	mongo := mongodb.NewCollections(dbConf.Mongo)

	return &Repos{
		ClickhouseClient: ck,
		RedisClient:      rediz.NewRedisClient(dbConf.Rediz),
		EventLog:         NewEventLogRepo(ck),
		JuGuangEvent:     NewJuGuangEventRepo(mongo),
		UserDevice:       NewUserDeviceRepo(ck, lc),
		IpDetail:         NewIpDetailRepo(ck, lc),
		ChatRecord:       NewChatRecordRepo(ck),
		UserFunnelRecord: NewUserFunnelRecordRepo(ck),
		Stripe:           NewStripeRepo(ck),
	}
}
