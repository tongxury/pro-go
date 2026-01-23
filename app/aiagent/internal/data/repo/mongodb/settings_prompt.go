package mongodb

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type PromptConfigCollection struct {
	*mgz.CoreCollection[aiagentpb.PromptConfig]
	c *cache.Cache
}

func NewPromptConfigCollection(db *mongo.Database) *PromptConfigCollection {

	cl := &PromptConfigCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.PromptConfig](db, "prompt_configs"),
		c:              cache.New(cache.NoExpiration, cache.NoExpiration),
	}

	cr := cron.New(cron.WithSeconds())

	_, _ = cr.AddFunc("@every 10s", func() {
		ctx := context.Background()

		settings, err := cl.List(ctx, bson.M{})
		if err != nil {
			log.Error("loop get settings err", err)
		}

		for _, x := range settings {
			cl.c.Set(x.XId, x, cache.NoExpiration)
		}

	})

	cr.Start()

	return cl
}

func (t *PromptConfigCollection) Get(ctx context.Context, id string) (*aiagentpb.PromptConfig, error) {

	cachedSettings, ok := t.c.Get(id)
	if ok {
		return cachedSettings.(*aiagentpb.PromptConfig), nil
	}

	settings, err := t.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	t.c.Set(id, settings, cache.NoExpiration)

	return settings, nil
}
