package loadbalancer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

type LoadBalancingStrategy string

const (
	StrategyRoundRobin     LoadBalancingStrategy = "round-robin"
	StrategyLeastConnected LoadBalancingStrategy = "least-connected"
	StrategyWeightedRandom LoadBalancingStrategy = "weighted-random"
	StrategyIPHash         LoadBalancingStrategy = "ip-hash"
)

type LoadBalancer struct {
	ID            string
	Region        string
	Strategy      LoadBalancingStrategy
	Backends      []engine.Component
	currentIndex  uint64
	healthy       bool
	metrics       *engine.Metrics
	metricsMutex  sync.RWMutex
	costPerHour   float64
	connections   map[string]int
	connMutex     sync.RWMutex
}

func NewLoadBalancer(id, region string, strategy LoadBalancingStrategy) *LoadBalancer {
	return &LoadBalancer{
		ID:          id,
		Region:      region,
		Strategy:    strategy,
		Backends:    make([]engine.Component, 0),
		healthy:     true,
		metrics:     &engine.Metrics{},
		costPerHour: 0.025,
		connections: make(map[string]int),
	}
}

func (lb *LoadBalancer) AddBackend(backend engine.Component) {
	lb.Backends = append(lb.Backends, backend)
}

func (lb *LoadBalancer) RemoveBackend(backendID string) {
	for i, backend := range lb.Backends {
		if backend.GetID() == backendID {
			lb.Backends = append(lb.Backends[:i], lb.Backends[i+1:]...)
			return
		}
	}
}

func (lb *LoadBalancer) GetID() string {
	return lb.ID
}

func (lb *LoadBalancer) GetType() string {
	return "load-balancer"
}

func (lb *LoadBalancer) Process(req *engine.Request) (*engine.Response, error) {
	start := time.Now()
	
	lb.metricsMutex.Lock()
	lb.metrics.RequestCount++
	lb.metricsMutex.Unlock()

	if !lb.healthy {
		lb.metricsMutex.Lock()
		lb.metrics.FailureCount++
		lb.metricsMutex.Unlock()
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("load balancer is unhealthy"),
		}, fmt.Errorf("load balancer is unhealthy")
	}

	backend := lb.selectBackend(req)
	if backend == nil {
		lb.metricsMutex.Lock()
		lb.metrics.FailureCount++
		lb.metricsMutex.Unlock()
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("no healthy backends available"),
		}, fmt.Errorf("no healthy backends available")
	}

	lbLatency := time.Millisecond * 2
	time.Sleep(lbLatency)

	resp, err := backend.Process(req)
	
	totalLatency := time.Since(start)
	
	lb.metricsMutex.Lock()
	if err == nil && resp.Success {
		lb.metrics.SuccessCount++
	} else {
		lb.metrics.FailureCount++
	}
	lb.metrics.TotalLatency += totalLatency
	lb.metrics.AverageLatency = time.Duration(int64(lb.metrics.TotalLatency) / lb.metrics.RequestCount)
	lb.metricsMutex.Unlock()

	if resp != nil {
		resp.Latency = totalLatency
		resp.HopsTrace = append([]string{lb.ID}, resp.HopsTrace...)
	}

	return resp, err
}

func (lb *LoadBalancer) selectBackend(req *engine.Request) engine.Component {
	healthyBackends := make([]engine.Component, 0)
	for _, backend := range lb.Backends {
		if backend.IsHealthy() {
			healthyBackends = append(healthyBackends, backend)
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	switch lb.Strategy {
	case StrategyRoundRobin:
		index := atomic.AddUint64(&lb.currentIndex, 1)
		return healthyBackends[index%uint64(len(healthyBackends))]
	
	case StrategyLeastConnected:
		lb.connMutex.RLock()
		defer lb.connMutex.RUnlock()
		
		minConn := int(^uint(0) >> 1)
		var selected engine.Component
		for _, backend := range healthyBackends {
			conn := lb.connections[backend.GetID()]
			if conn < minConn {
				minConn = conn
				selected = backend
			}
		}
		return selected
	
	default:
		index := atomic.AddUint64(&lb.currentIndex, 1)
		return healthyBackends[index%uint64(len(healthyBackends))]
	}
}

func (lb *LoadBalancer) GetMetrics() *engine.Metrics {
	lb.metricsMutex.RLock()
	defer lb.metricsMutex.RUnlock()
	
	metricsCopy := *lb.metrics
	if metricsCopy.RequestCount > 0 {
		metricsCopy.ErrorRate = float64(metricsCopy.FailureCount) / float64(metricsCopy.RequestCount)
	}
	return &metricsCopy
}

func (lb *LoadBalancer) GetCost() float64 {
	return lb.costPerHour
}

func (lb *LoadBalancer) IsHealthy() bool {
	return lb.healthy
}

func (lb *LoadBalancer) SetHealthy(healthy bool) {
	lb.healthy = healthy
}
