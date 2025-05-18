package cache

import (
	"gps-no-server/internal/core/models"
	"strings"
	"time"
)

type RangingCache struct {
	bySourceDestination *MemoryCache[string, *models.Ranging]
}

func NewRangingCache(ttl time.Duration, maxItems int) *RangingCache {
	return &RangingCache{
		bySourceDestination: NewMemoryCache[string, *models.Ranging](ttl, maxItems),
	}
}

func getCombinedKey(sourceMac, destinationMac string) string {
	return strings.Join([]string{sourceMac, destinationMac}, ":")
}

func (c *RangingCache) GetBySourceDestination(sourceMac, destinationMac string) (*models.Ranging, bool) {
	key := getCombinedKey(sourceMac, destinationMac)
	return c.bySourceDestination.Get(key)
}

func (c *RangingCache) Set(ranging *models.Ranging) {
	if ranging == nil || ranging.Source == nil || ranging.Destination == nil {
		return
	}

	rangingCopy := *ranging
	key := getCombinedKey(rangingCopy.Source.MacAddress, rangingCopy.Destination.MacAddress)

	c.bySourceDestination.Set(key, &rangingCopy)
}

func (c *RangingCache) Clear() {
	c.bySourceDestination.Clear()
}
