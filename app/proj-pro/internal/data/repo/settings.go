package repo

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SettingsCollection struct {
	*mgz.Core[*projpb.AppSettings]
}

func NewSettingsCollection(db *mongo.Database) *SettingsCollection {
	return &SettingsCollection{
		Core: mgz.NewCore[*projpb.AppSettings](db, "settings"),
	}
}

func (t *SettingsCollection) GetPrompt(ctx context.Context, key string) (*projpb.Prompt, error) {

	settings, err := t.FindOne(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, fmt.Errorf("setting not found")
	}

	p := settings.Prompts[key]

	if p == nil {
		return nil, fmt.Errorf("prompt not found")
	}

	log.Debugw("GetPrompt", key, "value", helper.SubString(p.GetContent(), 0, 100))

	return p, nil
}
