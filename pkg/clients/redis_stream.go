package clients

import (
	"context"
	redisV9 "github.com/redis/go-redis/v9"
	"store/pkg/sdk/helper"
)

type RedisStream struct {
	redisClient *redisV9.ClusterClient
	group       string
	consumer    string
	stream      string
	handle      StreamHandleFunc
	ctx         context.Context
	cancel      context.CancelFunc
}

type StreamHandleFunc func(ctx context.Context, messages []redisV9.XMessage) error

func NewRedisStream(redisClient *redisV9.ClusterClient, group, consumer, stream string, handle StreamHandleFunc) *RedisStream {
	return &RedisStream{redisClient: redisClient, group: group, consumer: consumer, stream: stream, handle: handle}
}

func (t *RedisStream) run(ctx context.Context, handle StreamHandleFunc) error {
	// 读
	streams, err := t.redisClient.XReadGroup(ctx, &redisV9.XReadGroupArgs{
		Group:    t.group,
		Consumer: t.consumer,
		Streams:  []string{t.stream, ">"},
		Count:    1,
		Block:    0,
		NoAck:    false,
	}).Result()

	if err != nil {
		return err
	}

	var messages []redisV9.XMessage
	for _, x := range streams {
		messages = append(messages, x.Messages...)
	}

	// 处理
	if err = handle(ctx, messages); err != nil {
		return err
	}

	// ACK
	var ids []string
	for _, x := range messages {
		ids = append(ids, x.ID)
	}

	_, err = t.redisClient.XAck(ctx, t.stream, t.group, ids...).Result()
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisStream) Start() {

	t.ctx, t.cancel = context.WithCancel(context.Background())

	t.redisClient.XGroupCreateMkStream(t.ctx, t.stream, t.group, "$")

	go func() {
		defer helper.DeferFunc()

		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				if err := t.run(t.ctx, t.handle); err != nil {

				}
			}
		}
	}()
}

func (t *RedisStream) Stop() {
	t.cancel()
}
