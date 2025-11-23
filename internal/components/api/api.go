package api

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

type InstanceSize string

const (
	SizeSmall  InstanceSize = "small"
	SizeMedium InstanceSize = "medium"
	SizeLarge  InstanceSize = "large"
	SizeXLarge InstanceSize = "xlarge"
)

type APIServer struct {
	ID               string
	Region           string
	Size             InstanceSize
	MaxConcurrent    int
	CurrentLoad      int
	LoadMutex        sync.RWMutex
	ProcessingTime   time.Duration
	Database         engine.Component
	Cache            engine.Component
	healthy          bool
	metrics          *engine.Metrics
	metricsMutex     sync.RWMutex
	costPerHour      float64
	cpuCores         int
	memoryGB         int
}

func NewAPIServer(id, region string, size InstanceSize) *APIServer {
	api := &APIServer{
		ID:             id,
		Region:         region,
		Size:           size,
		ProcessingTime: 10 * time.Millisecond,
		healthy:        true,
		metrics:        &engine.Metrics{},
	}
	
	switch size {
	case SizeSmall:
		api.MaxConcurrent = 10
		api.costPerHour = 0.05
		api.cpuCores = 1
		api.memoryGB = 2
	case SizeMedium:
		api.MaxConcurrent = 50
		api.costPerHour = 0.10
		api.cpuCores = 2
		api.memoryGB = 4
	case SizeLarge:
		api.MaxConcurrent = 200
		api.costPerHour = 0.20
		api.cpuCores = 4
		api.memoryGB = 8
	case SizeXLarge:
		api.MaxConcurrent = 500
		api.costPerHour = 0.40
		api.cpuCores = 8
		api.memoryGB = 16
	}
	
	return api
}

func (api *APIServer) SetDatabase(db engine.Component) {
	api.Database = db
}

func (api *APIServer) SetCache(cache engine.Component) {
	api.Cache = cache
}

func (api *APIServer) GetID() string {
	return api.ID
}

func (api *APIServer) GetType() string {
	return "api-server"
}

func (api *APIServer) Process(req *engine.Request) (*engine.Response, error) {
	start := time.Now()
	
	api.metricsMutex.Lock()
	api.metrics.RequestCount++
	api.metricsMutex.Unlock()

	if !api.healthy {
		api.metricsMutex.Lock()
		api.metrics.FailureCount++
		api.metricsMutex.Unlock()
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("API server is unhealthy"),
		}, fmt.Errorf("API server is unhealthy")
	}

	api.LoadMutex.Lock()
	if api.CurrentLoad >= api.MaxConcurrent {
		api.LoadMutex.Unlock()
		
		api.metricsMutex.Lock()
		api.metrics.FailureCount++
		api.metricsMutex.Unlock()
		
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("server at capacity"),
		}, fmt.Errorf("server at capacity")
	}
	api.CurrentLoad++
	api.LoadMutex.Unlock()

	defer func() {
		api.LoadMutex.Lock()
		api.CurrentLoad--
		api.LoadMutex.Unlock()
	}()

	processingTime := api.ProcessingTime + time.Duration(rand.Int63n(int64(5*time.Millisecond)))
	time.Sleep(processingTime)

	var resp *engine.Response
	var err error

	if api.Cache != nil && req.Type == engine.RequestTypeRead {
		resp, err = api.Cache.Process(req)
	} else if api.Database != nil {
		resp, err = api.Database.Process(req)
	} else {
		resp = &engine.Response{
			RequestID: req.ID,
			Success:   true,
			Latency:   time.Since(start),
			DataSize:  1024,
			HopsTrace: []string{api.ID},
		}
	}

	totalLatency := time.Since(start)
	
	api.metricsMutex.Lock()
	if err == nil && (resp == nil || resp.Success) {
		api.metrics.SuccessCount++
	} else {
		api.metrics.FailureCount++
	}
	api.metrics.TotalLatency += totalLatency
	api.metrics.AverageLatency = time.Duration(int64(api.metrics.TotalLatency) / api.metrics.RequestCount)
	api.metricsMutex.Unlock()

	if resp != nil {
		resp.Latency = totalLatency
		resp.HopsTrace = append([]string{api.ID}, resp.HopsTrace...)
	}

	return resp, err
}

func (api *APIServer) GetMetrics() *engine.Metrics {
	api.metricsMutex.RLock()
	defer api.metricsMutex.RUnlock()
	
	metricsCopy := *api.metrics
	if metricsCopy.RequestCount > 0 {
		metricsCopy.ErrorRate = float64(metricsCopy.FailureCount) / float64(metricsCopy.RequestCount)
		metricsCopy.Throughput = float64(metricsCopy.SuccessCount) / time.Since(time.Now().Add(-time.Hour)).Seconds()
	}
	return &metricsCopy
}

func (api *APIServer) GetCost() float64 {
	return api.costPerHour
}

func (api *APIServer) IsHealthy() bool {
	return api.healthy
}

func (api *APIServer) SetHealthy(healthy bool) {
	api.healthy = healthy
}

func (api *APIServer) GetCurrentLoad() float64 {
	api.LoadMutex.RLock()
	defer api.LoadMutex.RUnlock()
	
	return float64(api.CurrentLoad) / float64(api.MaxConcurrent)
}
