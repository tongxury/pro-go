package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/robfig/cron/v3"
	"store/app/aiagent/configs"
	"store/pkg/clients/grpcz"
	"store/pkg/sdk/helper/crond"
)

func main() {

	ctx := context.Background()

	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("failed get config: %v", err)
	}

	grpcClients, err := grpcz.NewClients(config.Component.Grpc)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugw("grpcClients", grpcClients)

	c := cron.New()

	// 执行提交的任务
	_, _ = c.AddJob("@every 1s", crond.NewJobWrapper(func() {
		_, _ = grpcClients.AiAgentClient.GenerateAnswerChunks(ctx, &empty.Empty{})
	}))

	_, _ = c.AddJob("@every 1s", crond.NewJobWrapper(func() {
		_, _ = grpcClients.AiAgentClient.FillItemRecord(context.Background(), &empty.Empty{})
	}))

	// 检查所有 generating状态的任务 5分钟内没结束就直接超时结束
	_, _ = c.AddJob("@every 1s", crond.NewJobWrapper(func() {
		_, _ = grpcClients.AiAgentClient.UpdateQuestionStatus(ctx, &empty.Empty{})
	}))
	c.Run()
}
