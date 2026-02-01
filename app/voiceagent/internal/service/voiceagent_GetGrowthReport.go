package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) GetGrowthReport(ctx context.Context, req *voiceagent.GetGrowthReportRequest) (*voiceagent.GrowthReport, error) {
	userId := krathelper.RequireUserId(ctx)

	// 确定时间范围
	now := time.Now()
	var startDate time.Time
	period := req.Period
	if period == "" {
		period = "week"
	}

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, 0, -7)
		period = "week"
	}

	report := &voiceagent.GrowthReport{
		Period:    period,
		StartDate: startDate.Unix(),
		EndDate:   now.Unix(),
	}

	// 统计对话次数
	convFilter := bson.M{
		"userId":    userId,
		"createdAt": bson.M{"$gte": startDate.Unix()},
	}
	conversations, total, err := s.Data.Mongo.Conversation.ListWithFilterAndSort(ctx, convFilter, bson.M{"createdAt": -1}, 1, 100)
	if err == nil {
		report.ConversationCount = int32(total)
		// 计算总时长（简单估算：每个对话平均5分钟）
		report.TotalDuration = int32(len(conversations) * 300)
	}

	// 获取新增记忆
	memFilter := bson.M{
		"userId":    userId,
		"createdAt": bson.M{"$gte": startDate.Unix()},
	}
	memories, _, err := s.Data.Mongo.Memory.ListWithFilterAndSort(ctx, memFilter, bson.M{"createdAt": -1}, 1, 10)
	if err == nil {
		report.NewMemories = memories
	}

	// 获取情绪统计
	emotionStats, err := s.getEmotionStatsForPeriod(ctx, userId, startDate, now)
	if err == nil {
		report.EmotionSummary = emotionStats
	}

	// 获取即将到来的事件（7天内）
	upcomingEvents, err := s.getUpcomingEventsForUser(ctx, userId, 7)
	if err == nil {
		report.UpcomingEvents = upcomingEvents
	}

	// 生成成长亮点（简单版本，后续可用 LLM 增强）
	report.Highlights = s.generateHighlights(report)
	report.Suggestions = s.generateSuggestions(report)

	return report, nil
}

// getEmotionStatsForPeriod 获取指定时间段的情绪统计
func (s *VoiceAgentService) getEmotionStatsForPeriod(ctx context.Context, userId string, start, end time.Time) (*voiceagent.EmotionStats, error) {
	filter := bson.M{
		"userId":    userId,
		"createdAt": bson.M{"$gte": start.Unix(), "$lte": end.Unix()},
	}

	logs, _, err := s.Data.Mongo.EmotionLog.ListWithFilterAndSort(ctx, filter, bson.M{"createdAt": 1}, 1, 1000)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return &voiceagent.EmotionStats{}, nil
	}

	// 统计各情绪出现次数和强度
	emotionCounts := make(map[string]int32)
	var totalIntensity float32
	for _, log := range logs {
		emotionCounts[log.Emotion]++
		totalIntensity += float32(log.Intensity)
	}

	// 找出主导情绪
	var dominantEmotion string
	var maxCount int32
	for emotion, count := range emotionCounts {
		if count > maxCount {
			maxCount = count
			dominantEmotion = emotion
		}
	}

	// 构建时间线
	timeline := make([]*voiceagent.EmotionDataPoint, 0, len(logs))
	for _, log := range logs {
		date := time.Unix(log.CreatedAt, 0).Format("2006-01-02")
		timeline = append(timeline, &voiceagent.EmotionDataPoint{
			Date:      date,
			Emotion:   log.Emotion,
			Intensity: log.Intensity,
		})
	}

	return &voiceagent.EmotionStats{
		DateRange:        start.Format("01/02") + " - " + end.Format("01/02"),
		DominantEmotion:  dominantEmotion,
		AverageIntensity: totalIntensity / float32(len(logs)),
		EmotionCounts:    emotionCounts,
		Timeline:         timeline,
	}, nil
}

// getUpcomingEventsForUser 获取即将到来的事件
func (s *VoiceAgentService) getUpcomingEventsForUser(ctx context.Context, userId string, days int) ([]*voiceagent.ImportantEvent, error) {
	// 复用已有的 GetUpcomingEvents 逻辑
	req := &voiceagent.GetUpcomingEventsRequest{Days: int32(days)}
	result, err := s.GetUpcomingEvents(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.List, nil
}

// generateHighlights 生成成长亮点
func (s *VoiceAgentService) generateHighlights(report *voiceagent.GrowthReport) []string {
	highlights := []string{}

	if report.ConversationCount > 0 {
		highlights = append(highlights, fmt.Sprintf("与 AI 进行了 %d 次深度对话", report.ConversationCount))
	}

	if len(report.NewMemories) > 0 {
		highlights = append(highlights, fmt.Sprintf("AI 记住了 %d 条关于你的新信息", len(report.NewMemories)))
	}

	if report.EmotionSummary != nil && report.EmotionSummary.DominantEmotion != "" {
		emotionLabels := map[string]string{
			"happy":   "开心",
			"calm":    "平静",
			"excited": "兴奋",
			"sad":     "难过",
			"anxious": "焦虑",
			"angry":   "生气",
			"neutral": "平常",
		}
		label := emotionLabels[report.EmotionSummary.DominantEmotion]
		if label == "" {
			label = report.EmotionSummary.DominantEmotion
		}
		highlights = append(highlights, "你这段时间的主要情绪是「"+label+"」")
	}

	if len(report.UpcomingEvents) > 0 {
		highlights = append(highlights, fmt.Sprintf("有 %d 个重要日子即将到来", len(report.UpcomingEvents)))
	}

	if len(highlights) == 0 {
		highlights = append(highlights, "这是一段新的开始，期待与你一起成长")
	}

	return highlights
}

// generateSuggestions 生成建议
func (s *VoiceAgentService) generateSuggestions(report *voiceagent.GrowthReport) []string {
	suggestions := []string{}

	if report.ConversationCount < 3 {
		suggestions = append(suggestions, "多和 AI 聊聊天，分享你的日常")
	}

	if report.EmotionSummary != nil {
		if report.EmotionSummary.DominantEmotion == "anxious" || report.EmotionSummary.DominantEmotion == "sad" {
			suggestions = append(suggestions, "最近似乎有些压力，试着放松一下")
		}
		if report.EmotionSummary.AverageIntensity > 7 {
			suggestions = append(suggestions, "情绪波动较大，记得照顾好自己")
		}
	}

	if len(report.NewMemories) == 0 {
		suggestions = append(suggestions, "分享更多关于你的故事，让 AI 更了解你")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "保持当前的状态，你做得很好！")
	}

	return suggestions
}
