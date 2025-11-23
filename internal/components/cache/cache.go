package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

type EvictionPolicy string

const (
	EvictionLRU   EvictionPolicy = "lru"
	EvictionLFU   EvictionPolicy = "lfu"
	EvictionFIFO  EvictionPolicy = "fifo"
	EvictionRandom EvictionPolicy = "random"
)

type CacheEntry struct {
	Key        string
	Data       []byte
	Size       int64
	Expiry     time.Time
	AccessTime time.Time
	AccessCount int64
}

type Cache struct {
	ID            string
	Type          string
	Region        string
	Capacity      int64
	UsedCapacity  int64
	Policy        EvictionPolicy
	TTL           time.Duration
	ReadLatency   time.Duration
	WriteLatency  time.Duration
	Backend       engine.Component
	healthy       bool
	metrics       *engine.Metrics
	metricsMutex  sync.RWMutex
	entries       map[string]*CacheEntry
	entriesMutex  sync.RWMutex
	costPerHour   float64
}

func NewCache(id, cacheType, region string, capacity int64, policy EvictionPolicy, ttl time.Duration) *Cache {
	return &Cache{
		ID:           id,
		Type:         cacheType,
		Region:       region,
		Capacity:     capacity,
		Policy:       policy,
		TTL:          ttl,
		ReadLatency:  time.Millisecond,
		WriteLatency: 2 * time.Millisecond,
		healthy:      true,
		metrics:      &engine.Metrics{},
		entries:      make(map[string]*CacheEntry),
		costPerHour:  0.02,
	}
}

func (c *Cache) SetBackend(backend engine.Component) {
	c.Backend = backend
}

func (c *Cache) GetID() string {
	return c.ID
}

func (c *Cache) GetType() string {
	return fmt.Sprintf("cache-%s", c.Type)
}

func (c *Cache) Process(req *engine.Request) (*engine.Response, error) {
	start := time.Now()
	
	c.metricsMutex.Lock()
	c.metrics.RequestCount++
	c.metricsMutex.Unlock()

	if !c.healthy {
		c.metricsMutex.Lock()
		c.metrics.FailureCount++
		c.metricsMutex.Unlock()
		
		if c.Backend != nil {
			return c.Backend.Process(req)
		}
		
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("cache is unhealthy"),
		}, fmt.Errorf("cache is unhealthy")
	}

	if req.Type == engine.RequestTypeRead {
		if entry := c.get(req.Path); entry != nil {
			time.Sleep(c.ReadLatency)
			
			totalLatency := time.Since(start)
			
			c.metricsMutex.Lock()
			c.metrics.SuccessCount++
			c.metrics.TotalLatency += totalLatency
			c.metrics.AverageLatency = time.Duration(int64(c.metrics.TotalLatency) / c.metrics.RequestCount)
			c.metricsMutex.Unlock()
			
			return &engine.Response{
				RequestID: req.ID,
				Success:   true,
				Latency:   totalLatency,
				DataSize:  entry.Size,
				CacheHit:  true,
				HopsTrace: []string{c.ID},
			}, nil
		}
	}

	if c.Backend != nil {
		resp, err := c.Backend.Process(req)
		
		if err == nil && resp.Success && req.Type == engine.RequestTypeRead {
			c.set(req.Path, resp.DataSize)
		}
		
		if resp != nil {
			resp.HopsTrace = append([]string{c.ID}, resp.HopsTrace...)
			resp.CacheHit = false
		}
		
		return resp, err
	}

	c.metricsMutex.Lock()
	c.metrics.FailureCount++
	c.metricsMutex.Unlock()
	
	return &engine.Response{
		RequestID: req.ID,
		Success:   false,
		Latency:   time.Since(start),
		Error:     fmt.Errorf("cache miss and no backend"),
	}, fmt.Errorf("cache miss and no backend")
}

func (c *Cache) get(key string) *CacheEntry {
	c.entriesMutex.RLock()
	defer c.entriesMutex.RUnlock()
	
	entry, exists := c.entries[key]
	if !exists {
		return nil
	}
	
	if time.Now().After(entry.Expiry) {
		go c.evict(key)
		return nil
	}
	
	entry.AccessTime = time.Now()
	entry.AccessCount++
	
	return entry
}

func (c *Cache) set(key string, size int64) {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
	
	for c.UsedCapacity+size > c.Capacity && len(c.entries) > 0 {
		c.evictOne()
	}
	
	entry := &CacheEntry{
		Key:         key,
		Data:        make([]byte, size),
		Size:        size,
		Expiry:      time.Now().Add(c.TTL),
		AccessTime:  time.Now(),
		AccessCount: 1,
	}
	
	c.entries[key] = entry
	c.UsedCapacity += size
}

func (c *Cache) evict(key string) {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
	
	if entry, exists := c.entries[key]; exists {
		c.UsedCapacity -= entry.Size
		delete(c.entries, key)
	}
}

func (c *Cache) evictOne() {
	var evictKey string
	var oldestTime time.Time = time.Now()
	var lowestCount int64 = int64(^uint64(0) >> 1)
	
	switch c.Policy {
	case EvictionLRU:
		for key, entry := range c.entries {
			if entry.AccessTime.Before(oldestTime) {
				oldestTime = entry.AccessTime
				evictKey = key
			}
		}
	case EvictionLFU:
		for key, entry := range c.entries {
			if entry.AccessCount < lowestCount {
				lowestCount = entry.AccessCount
				evictKey = key
			}
		}
	default:
		for key := range c.entries {
			evictKey = key
			break
		}
	}
	
	if evictKey != "" {
		if entry, exists := c.entries[evictKey]; exists {
			c.UsedCapacity -= entry.Size
			delete(c.entries, evictKey)
		}
	}
}

func (c *Cache) GetMetrics() *engine.Metrics {
	c.metricsMutex.RLock()
	defer c.metricsMutex.RUnlock()
	
	metricsCopy := *c.metrics
	if metricsCopy.RequestCount > 0 {
		metricsCopy.ErrorRate = float64(metricsCopy.FailureCount) / float64(metricsCopy.RequestCount)
		
		cacheHits := metricsCopy.SuccessCount
		metricsCopy.CacheHitRate = float64(cacheHits) / float64(metricsCopy.RequestCount)
	}
	return &metricsCopy
}

func (c *Cache) GetCost() float64 {
	baseCost := c.costPerHour
	capacityCost := float64(c.Capacity) / (1024 * 1024 * 1024) * 0.005
	
	return baseCost + capacityCost
}

func (c *Cache) IsHealthy() bool {
	return c.healthy
}

func (c *Cache) SetHealthy(healthy bool) {
	c.healthy = healthy
}
