package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) GetEmotionStats(ctx context.Context, req *voiceagent.GetEmotionStatsRequest) (*voiceagent.EmotionStats, error) {
	userId := krathelper.RequireUserId(ctx)

	days := req.Days
	if days <= 0 {
		days = 30
	}

	startTime := time.Now().AddDate(0, 0, -int(days)).Unix()

	filter := bson.M{
		"user._id":  userId,
		"createdAt": bson.M{"$gte": startTime},
	}

	// 获取时间范围内的所有记录
	list, _, err := s.Data.Mongo.EmotionLog.ListWithFilterAndSort(ctx, filter, bson.M{"createdAt": 1}, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 统计各情绪出现次数
	emotionCounts := make(map[string]int32)
	var totalIntensity float32 = 0
	timeline := make([]*voiceagent.EmotionDataPoint, 0, len(list))

	for _, log := range list {
		emotionCounts[log.Emotion]++
		totalIntensity += float32(log.Intensity)

		timeline = append(timeline, &voiceagent.EmotionDataPoint{
			Date:      time.Unix(log.CreatedAt, 0).Format("2006-01-02"),
			Emotion:   log.Emotion,
			Intensity: log.Intensity,
		})
	}

	// 找出主要情绪
	dominantEmotion := ""
	maxCount := int32(0)
	for emotion, count := range emotionCounts {
		if count > maxCount {
			maxCount = count
			dominantEmotion = emotion
		}
	}

	// 计算平均强度
	var averageIntensity float32 = 0
	if len(list) > 0 {
		averageIntensity = totalIntensity / float32(len(list))
	}

	return &voiceagent.EmotionStats{
		DateRange:        time.Unix(startTime, 0).Format("2006-01-02") + " ~ " + time.Now().Format("2006-01-02"),
		DominantEmotion:  dominantEmotion,
		AverageIntensity: averageIntensity,
		EmotionCounts:    emotionCounts,
		Timeline:         timeline,
	}, nil
}
