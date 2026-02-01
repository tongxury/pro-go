package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) GetUpcomingEvents(ctx context.Context, req *voiceagent.GetUpcomingEventsRequest) (*voiceagent.EventList, error) {
	userId := krathelper.RequireUserId(ctx)

	days := req.Days
	if days <= 0 {
		days = 7
	}

	now := time.Now()
	endDate := now.AddDate(0, 0, int(days))

	// 获取用户所有事件
	filter := bson.M{"userId": userId}
	allEvents, _, err := s.Data.Mongo.ImportantEvent.ListWithFilterAndSort(ctx, filter, bson.M{"date": 1}, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 筛选即将到来的事件
	upcomingEvents := make([]*voiceagent.ImportantEvent, 0)
	currentYear := now.Year()
	todayStr := now.Format("01-02")   // MM-DD
	endStr := endDate.Format("01-02") // MM-DD

	for _, event := range allEvents {
		// 处理 YYYY-MM-DD 或 MM-DD 格式
		dateStr := event.Date
		var checkStr string

		if len(dateStr) == 10 { // YYYY-MM-DD
			if event.IsRecurring {
				// 循环事件，检查 MM-DD 部分
				checkStr = dateStr[5:] // 提取 MM-DD
			} else {
				// 非循环事件，检查完整日期
				eventDate, err := time.Parse("2006-01-02", dateStr)
				if err == nil && eventDate.After(now.AddDate(0, 0, -1)) && eventDate.Before(endDate.AddDate(0, 0, 1)) {
					upcomingEvents = append(upcomingEvents, event)
				}
				continue
			}
		} else if len(dateStr) == 5 { // MM-DD
			checkStr = dateStr
		} else {
			continue
		}

		// 检查 MM-DD 是否在范围内
		if isDateInRange(checkStr, todayStr, endStr, currentYear, now.Year() == endDate.Year()) {
			upcomingEvents = append(upcomingEvents, event)
		}
	}

	return &voiceagent.EventList{
		List:  upcomingEvents,
		Total: int64(len(upcomingEvents)),
	}, nil
}

// isDateInRange 检查 MM-DD 格式的日期是否在范围内
func isDateInRange(checkStr, startStr, endStr string, year int, sameYear bool) bool {
	check, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%s", year, checkStr))
	start, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%s", year, startStr))
	end, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%s", year, endStr))

	// 处理跨年的情况
	if !sameYear && end.Before(start) {
		end = end.AddDate(1, 0, 0)
		if check.Before(start) {
			check = check.AddDate(1, 0, 0)
		}
	}

	return !check.Before(start) && !check.After(end)
}
