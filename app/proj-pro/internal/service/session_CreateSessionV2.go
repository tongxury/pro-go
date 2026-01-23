package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
	"time"
)

func (t *SessionService) CreateSessionV2(ctx context.Context, req *projpb.CreateSessionV2Request) (*projpb.Session, error) {

	userId := krathelper.RequireUserId(ctx)

	commodity, err := t.data.Mongo.Commodity.GetById(ctx, req.CommodityId)
	if err != nil {
		return nil, err
	}

	chance := commodity.Chances[req.TargetChance]

	template, err := t.data.GrpcClients.ProjAdminClient.GetTemplate(ctx, &projpb.GetGetTemplateRequest{
		Id: req.TemplateId,
	})
	if err != nil {
		return nil, err
	}

	session, err := t.data.Mongo.Session.Insert(ctx, &projpb.Session{
		Commodity:    commodity,
		CreatedAt:    time.Now().Unix(),
		UserId:       userId,
		Status:       "created",
		TargetChance: chance,
		Segments:     nil,
		Template:     template,
	})
	if err != nil {
		return nil, err
	}

	var sessionSegments []*projpb.SessionSegment
	for _, x := range template.Segments {
		sessionSegments = append(sessionSegments, &projpb.SessionSegment{
			SessionId: session.XId,
			Status:    "created",
			Session:   session,
			Segment:   x,
		})
	}

	_, err = t.data.Mongo.SessionSegment.InsertMany(ctx, sessionSegments...)
	if err != nil {
		return nil, err
	}

	return &projpb.Session{
		XId: session.XId,
	}, nil
}
