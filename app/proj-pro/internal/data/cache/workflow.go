package cache

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/rediz"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type WorkflowCache struct {
	redis  *rediz.RedisClient
	prefix string
}

func NewWorkflowCache(redis *rediz.RedisClient) *WorkflowCache {
	return &WorkflowCache{
		redis:  redis,
		prefix: "workflow_",
	}
}

func (t *WorkflowCache) ListIds(ctx context.Context) ([]string, error) {

	//t.redis.ZRange(ctx, t.prefix+":ids")
	return nil, nil
}

func (t *WorkflowCache) Set(ctx context.Context, workflow *projpb.Workflow, expiration time.Duration) error {

	err := t.redis.SAdd(ctx, t.prefix+":ids", workflow.XId).Err()
	if err != nil {
		return err
	}

	marshal, err := protojson.Marshal(workflow)
	if err != nil {
		return err
	}

	return t.redis.Set(ctx, t.prefix+workflow.XId, string(marshal), expiration).Err()
}

func (t *WorkflowCache) Get(ctx context.Context, id string) (*projpb.Workflow, error) {

	result, err := t.redis.Get(ctx, t.prefix+id).Bytes()
	if err != nil {
		return nil, err
	}

	var workflow projpb.Workflow
	err = protojson.Unmarshal(result, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil

}

func (t *WorkflowCache) Delete(ctx context.Context, id string) error {
	err := t.redis.ZRem(ctx, t.prefix+":ids", id).Err()
	if err != nil {
		return err
	}
	return t.redis.Del(ctx, t.prefix+id).Err()
}
