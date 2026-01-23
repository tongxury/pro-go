package cache

import (
	"store/pkg/rediz"
)

type Caches struct {
	Workflow *WorkflowCache
	Remix    *RemixCache
}

func NewCaches(redis *rediz.RedisClient) *Caches {
	return &Caches{
		Workflow: NewWorkflowCache(redis),
		Remix:    NewRemixCache(),
	}
}
