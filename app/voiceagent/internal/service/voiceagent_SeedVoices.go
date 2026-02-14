package service

import (
	"context"
	"store/api/voiceagent"
	"time"

	ucpb "store/api/usercenter"
	"store/pkg/sdk/third/cartesia"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) SeedSystemVoices(ctx context.Context) {

	// 1. Fetch voices from Cartesia API
	// Filter for Chinese voices directly via API
	voices, err := s.Data.Cartesia.ListVoices(ctx, &cartesia.ListVoicesRequest{
		Language: "zh",
	})
	if err != nil {
		log.Errorf("Failed to list voices from Cartesia: %v", err)
		return
	}

	// 2. Map and Upsert voices into MongoDB
	for _, v := range voices.Data {
		filter := bson.M{
			"voiceId":  v.Id,
			"user._id": "system",
		}

		one, err := s.Data.Mongo.Voice.FindOne(ctx, filter)
		if err != nil {
			continue
		}

		if one != nil {
			continue
		}

		log.Infof("Generating preview for system voice: %s (%s)", v.Name, v.Id)
		url, err := s.GeneratePreviewAudio(ctx, v.Id, "")
		if err != nil {
			log.Errorf("Failed to generate preview for voice %s: %v", v.Name, err)
			continue
		}

		_, err = s.Data.Mongo.Voice.Insert(ctx, &voiceagent.Voice{
			User:      &ucpb.User{XId: "system"},
			Name:      v.Name,
			Type:      "preset",
			VoiceId:   v.Id,
			Status:    "active",
			SampleUrl: url,
			CreatedAt: time.Now().Unix(),
		})
		if err != nil {
			continue
		}

	}
}
