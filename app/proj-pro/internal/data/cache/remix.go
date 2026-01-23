package cache

import (
	"context"
	projpb "store/api/proj"
	"sync"
)

type RemixCache struct {
	remixCache *projpb.RemixOptions
	remixMu    sync.RWMutex
}

func NewRemixCache() *RemixCache {
	return &RemixCache{}
}

func (t *RemixCache) Set(ctx context.Context, options *projpb.RemixOptions) error {
	t.remixMu.Lock()
	defer t.remixMu.Unlock()
	t.remixCache = options
	return nil
}

func (t *RemixCache) Get(ctx context.Context) (*projpb.RemixOptions, error) {
	t.remixMu.RLock()
	defer t.remixMu.RUnlock()
	return t.remixCache, nil
}
