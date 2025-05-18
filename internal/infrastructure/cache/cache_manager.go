package cache

import "time"

type CacheManager struct {
	RangingCache *RangingCache
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		RangingCache: NewRangingCache(1*time.Minute, 5000),
	}
}

func (cm *CacheManager) ClearAll() {
	cm.RangingCache.Clear()
}
