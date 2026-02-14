package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/sdk/third/cartesia"
)

func (s *VoiceAgentService) SeedSystemVoices(ctx context.Context) {
	if s.Data == nil || s.Data.Mongo == nil || s.Data.Mongo.Voice == nil {
		return
	}

	// 1. Fetch voices from Cartesia API
	// Filter for Chinese voices directly via API
	voices, err := s.Data.Cartesia.ListVoices(ctx, &cartesia.ListVoicesRequest{
		Language: "zh",
	})
	if err != nil {
		log.Errorf("Failed to list voices from Cartesia: %v", err)
		return
	}

	var systemVoices []*voiceagent.Voice

	// 2. Map to system voices
	for _, v := range voices {
		systemVoices = append(systemVoices, &voiceagent.Voice{
			User:      &ucpb.User{XId: "system"},
			Name:      v.Name,
			Type:      "preset",
			VoiceId:   v.Id,
			Status:    "active",
			SampleUrl: "",                // System voices usually don't need a sample URL for playback if they are preset
			CreatedAt: time.Now().Unix(), // This will be set on insert, ignored on update if not modified
		})
	}

	if len(systemVoices) == 0 {
		log.Warn("No Chinese voices found from Cartesia API")
	}

	// 3. Upsert voices into MongoDB
	for _, v := range systemVoices {
		filter := bson.M{
			"voiceId":  v.VoiceId,
			"user._id": "system",
		}

		update := bson.M{
			"$setOnInsert": bson.M{
				"_id":       primitive.NewObjectID().Hex(),
				"createdAt": time.Now().Unix(),
				"user":      v.User,
			},
			"$set": bson.M{
				"name":   v.Name,
				"type":   v.Type,
				"status": v.Status,
			},
		}

		// Use upsert to insert if not exists, or update mutable fields if exists
		opts := options.Update().SetUpsert(true)
		_, _ = s.Data.Mongo.Voice.C().UpdateOne(ctx, filter, update, opts)
	}
}
