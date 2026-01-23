package repo

import (
	"fmt"
	"time"
)

type Limited struct {
}

func (t Limited) Error() string {
	return "limited"
}

type Usage map[string]int64

func (t Usage) FindTotalUsage(key string) int64 {
	return t[key]
}

func (t Usage) FindMonthUsage(key string) int64 {
	dd := time.Now().Format("2006-01")
	return t[fmt.Sprintf("%s:%s", key, dd)]
}

func (t Usage) FindDayUsage(key string) int64 {
	dd := time.Now().Format("2006-01-02")
	return t[fmt.Sprintf("%s:%s", key, dd)]
}
