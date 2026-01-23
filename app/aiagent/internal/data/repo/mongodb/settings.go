package mongodb

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type SettingsCollection struct {
	*mgz.CoreCollection[aiagentpb.Settings]
	c *cache.Cache
}

func NewSettingsCollection(db *mongo.Database) *SettingsCollection {

	cl := &SettingsCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Settings](db, "settings"),
		c:              cache.New(cache.NoExpiration, cache.NoExpiration),
	}

	cr := cron.New(cron.WithSeconds())

	_, _ = cr.AddFunc("@every 2s", func() {
		ctx := context.Background()

		settings, err := cl.GetById(ctx, "1")
		if err != nil {
			log.Error("loop get settings err", err)
		}

		cl.c.Set("settings", settings, cache.NoExpiration)
	})

	cr.Start()

	return cl
}

func (t *SettingsCollection) Get(ctx context.Context) (*aiagentpb.Settings, error) {

	cachedSettings, ok := t.c.Get("settings")
	if ok {
		return cachedSettings.(*aiagentpb.Settings), nil
	}

	settings, err := t.GetById(ctx, "1")
	if err != nil {
		return nil, err
	}

	t.c.Set("settings", settings, cache.NoExpiration)

	return settings, nil
}
