package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjAdminService) UpdateSettings(ctx context.Context, req *projpb.UpdateSettingsRequest) (*projpb.AppSettings, error) {

	op := mgz.Op()

	for k, v := range req.Prompts {
		op = op.Set("prompts."+k, v)
	}

	for _, k := range req.DeletePrompts {
		op = op.Unset("prompts." + k)
	}

	//if request.GetVideoHighlight() != nil {
	//	request.GetVideoHighlight().UpdatedAt = time.Now().Unix()
	//	op = op.Set("videoHighlight", request.GetVideoHighlight())
	//} else if request.GetVideoTemplate() != nil {
	//	request.GetVideoTemplate().UpdatedAt = time.Now().Unix()
	//	op = op.Set("videoTemplate", request.GetVideoTemplate())
	//} else if request.GetVideoScript() != nil {
	//	request.GetVideoScript().UpdatedAt = time.Now().Unix()
	//	op = op.Set("videoScript", request.GetVideoScript())
	//} else if request.GetVideoGenerate() != nil {
	//	request.GetVideoGenerate().UpdatedAt = time.Now().Unix()
	//	op = op.Set("videoGenerate", request.GetVideoGenerate())
	//}

	one, err := t.data.Mongo.Settings.FindOne(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if one == nil {
		return &projpb.AppSettings{}, nil
	}

	_, err = t.data.Mongo.Settings.UpdateByIDIfExists(ctx, one.XId, op)
	if err != nil {
		log.Errorw("UpdateByIDIfExists err", err, "id", one.XId)
		return nil, err
	}

	//id, err := t.data.Mongo.Settings.FindByID(ctx, req.Id)
	//if err != nil {
	//	return nil, err
	//}

	return &projpb.AppSettings{
		Prompts: req.Prompts,
	}, nil
}
