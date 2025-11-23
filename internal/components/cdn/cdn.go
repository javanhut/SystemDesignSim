package cdn

import (
	"fmt"
	"sync"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

type EdgeLocation struct {
	Region      string
	Cache       map[string][]byte
	CacheMutex  sync.RWMutex
	HitCount    int64
	MissCount   int64
}

type CDN struct {
	ID            string
	EdgeLocations map[string]*EdgeLocation
	Origin        engine.Component
	TTL           time.Duration
	healthy       bool
	metrics       *engine.Metrics
	metricsMutex  sync.RWMutex
	costPerHour   float64
}

func NewCDN(id string, regions []string) *CDN {
	cdn := &CDN{
		ID:            id,
		EdgeLocations: make(map[string]*EdgeLocation),
		TTL:           time.Hour,
		healthy:       true,
		metrics:       &engine.Metrics{},
		costPerHour:   0.08,
	}
	
	for _, region := range regions {
		cdn.EdgeLocations[region] = &EdgeLocation{
			Region: region,
			Cache:  make(map[string][]byte),
		}
	}
	
	return cdn
}

func (cdn *CDN) SetOrigin(origin engine.Component) {
	cdn.Origin = origin
}

func (cdn *CDN) GetID() string {
	return cdn.ID
}

func (cdn *CDN) GetType() string {
	return "cdn"
}

func (cdn *CDN) Process(req *engine.Request) (*engine.Response, error) {
	start := time.Now()
	
	cdn.metricsMutex.Lock()
	cdn.metrics.RequestCount++
	cdn.metricsMutex.Unlock()

	if !cdn.healthy {
		cdn.metricsMutex.Lock()
		cdn.metrics.FailureCount++
		cdn.metricsMutex.Unlock()
		
		if cdn.Origin != nil {
			return cdn.Origin.Process(req)
		}
		
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("CDN is unhealthy"),
		}, fmt.Errorf("CDN is unhealthy")
	}

	edge, exists := cdn.EdgeLocations[req.Region]
	if !exists {
		for _, e := range cdn.EdgeLocations {
			edge = e
			break
		}
	}

	if edge != nil && req.Type == engine.RequestTypeRead {
		edge.CacheMutex.RLock()
		data, cached := edge.Cache[req.Path]
		edge.CacheMutex.RUnlock()
		
		if cached {
			edge.HitCount++
			
			time.Sleep(2 * time.Millisecond)
			
			totalLatency := time.Since(start)
			
			cdn.metricsMutex.Lock()
			cdn.metrics.SuccessCount++
			cdn.metrics.TotalLatency += totalLatency
			cdn.metrics.AverageLatency = time.Duration(int64(cdn.metrics.TotalLatency) / cdn.metrics.RequestCount)
			cdn.metricsMutex.Unlock()
			
			return &engine.Response{
				RequestID: req.ID,
				Success:   true,
				Latency:   totalLatency,
				DataSize:  int64(len(data)),
				CacheHit:  true,
				HopsTrace: []string{fmt.Sprintf("%s-edge-%s", cdn.ID, edge.Region)},
			}, nil
		}
		
		edge.MissCount++
	}

	if cdn.Origin != nil {
		resp, err := cdn.Origin.Process(req)
		
		if err == nil && resp.Success && req.Type == engine.RequestTypeRead && edge != nil {
			edge.CacheMutex.Lock()
			edge.Cache[req.Path] = make([]byte, resp.DataSize)
			edge.CacheMutex.Unlock()
		}
		
		if resp != nil {
			resp.HopsTrace = append([]string{cdn.ID}, resp.HopsTrace...)
			resp.CacheHit = false
		}
		
		return resp, err
	}

	cdn.metricsMutex.Lock()
	cdn.metrics.FailureCount++
	cdn.metricsMutex.Unlock()
	
	return &engine.Response{
		RequestID: req.ID,
		Success:   false,
		Latency:   time.Since(start),
		Error:     fmt.Errorf("CDN cache miss and no origin"),
	}, fmt.Errorf("CDN cache miss and no origin")
}

func (cdn *CDN) GetMetrics() *engine.Metrics {
	cdn.metricsMutex.RLock()
	defer cdn.metricsMutex.RUnlock()
	
	metricsCopy := *cdn.metrics
	if metricsCopy.RequestCount > 0 {
		metricsCopy.ErrorRate = float64(metricsCopy.FailureCount) / float64(metricsCopy.RequestCount)
		
		var totalHits, totalRequests int64
		for _, edge := range cdn.EdgeLocations {
			totalHits += edge.HitCount
			totalRequests += edge.HitCount + edge.MissCount
		}
		
		if totalRequests > 0 {
			metricsCopy.CacheHitRate = float64(totalHits) / float64(totalRequests)
		}
	}
	return &metricsCopy
}

func (cdn *CDN) GetCost() float64 {
	regionCost := float64(len(cdn.EdgeLocations)) * 0.01
	return cdn.costPerHour + regionCost
}

func (cdn *CDN) IsHealthy() bool {
	return cdn.healthy
}

func (cdn *CDN) SetHealthy(healthy bool) {
	cdn.healthy = healthy
}
