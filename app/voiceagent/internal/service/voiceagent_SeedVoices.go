package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"

	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/sdk/third/cartesia"
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

	var systemVoices []*voiceagent.Voice

	// 2. Map to system voices
	for _, v := range voices.Data {
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

		// 改成实时查找 (Real-time lookup)
		count, err := s.Data.Mongo.Voice.C().CountDocuments(ctx, filter)
		if err != nil {
			log.Errorf("Failed to count voice %s: %v", v.VoiceId, err)
			continue
		}

		if count > 0 {
			// 如果数据库有 就跳过 (If exists, skip)
			continue
		}

		// 没有就插入 (If not exists, insert)
		// 不要生成 _id (Do not generate _id manually, let Mongo handle it)
		doc := bson.M{
			"voiceId":   v.VoiceId,
			"user":      v.User,
			"name":      v.Name,
			"type":      v.Type,
			"status":    v.Status,
			"sampleUrl": v.SampleUrl,
			"createdAt": time.Now().Unix(),
		}

		_, err = s.Data.Mongo.Voice.C().InsertOne(ctx, doc)
		if err != nil {
			log.Errorf("Failed to insert voice %s: %v", v.Name, err)
		}
	}
}
