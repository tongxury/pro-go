package volcengine

import (
	"context"
	"testing"
)

func TestClient_QueryMixCutTaskResult(t *testing.T) {

	c := NewClient()

	c.QueryMixCutTaskResult(context.Background(), QueryMixCutTaskResultParams{
		//TaskKey: "142673e2-9ae2-11f0-9a17-3436ac120025",
		TaskKey: "5addab02-9ae7-11f0-a21d-3436ac120709",
	})
}

type T4 struct {
	Code             int    `json:"Code"`
	Message          string `json:"Message"`
	ResponseMetadata struct {
		Action    string `json:"Action"`
		Region    string `json:"Region"`
		RequestId string `json:"RequestId"`
		Service   string `json:"Service"`
		Version   string `json:"Version"`
	} `json:"ResponseMetadata"`
	Result struct {
		Data struct {
			Task struct {
				CreatedAt string `json:"CreatedAt"`
				Status    int    `json:"Status"`
				TaskKey   string `json:"TaskKey"`
				UpdatedAt string `json:"UpdatedAt"`
				VideoKey  string `json:"VideoKey"`
			} `json:"Task"`
		} `json:"Data"`
	} `json:"Result"`
}
