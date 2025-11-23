package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Simulator struct {
	components     map[string]Component
	componentMutex sync.RWMutex
	eventQueue     chan *Request
	responseQueue  chan *Response
	running        bool
	ctx            context.Context
	cancel         context.CancelFunc
	tickRate       time.Duration
	currentTime    time.Time
	metrics        *AggregateMetrics
}

type AggregateMetrics struct {
	TotalRequests     int64
	TotalSuccesses    int64
	TotalFailures     int64
	TotalCost         float64
	TotalLatency      time.Duration
	RecentLatencies   []time.Duration // Track last 1000 latencies for P99
	ComponentMetrics  map[string]*Metrics
	mu                sync.RWMutex
}

func NewSimulator(tickRate time.Duration) *Simulator {
	ctx, cancel := context.WithCancel(context.Background())
	return &Simulator{
		components:    make(map[string]Component),
		eventQueue:    make(chan *Request, 10000),
		responseQueue: make(chan *Response, 10000),
		running:       false,
		ctx:           ctx,
		cancel:        cancel,
		tickRate:      tickRate,
		currentTime:   time.Now(),
		metrics: &AggregateMetrics{
			ComponentMetrics:  make(map[string]*Metrics),
			RecentLatencies: make([]time.Duration, 0, 1000),
		},
	}
}

func (s *Simulator) RegisterComponent(component Component) error {
	s.componentMutex.Lock()
	defer s.componentMutex.Unlock()

	id := component.GetID()
	if _, exists := s.components[id]; exists {
		return fmt.Errorf("component with ID %s already exists", id)
	}

	s.components[id] = component
	return nil
}

func (s *Simulator) UnregisterComponent(id string) error {
	s.componentMutex.Lock()
	defer s.componentMutex.Unlock()

	if _, exists := s.components[id]; !exists {
		return fmt.Errorf("component with ID %s not found", id)
	}

	delete(s.components, id)
	return nil
}

func (s *Simulator) GetComponent(id string) (Component, error) {
	s.componentMutex.RLock()
	defer s.componentMutex.RUnlock()

	component, exists := s.components[id]
	if !exists {
		return nil, fmt.Errorf("component with ID %s not found", id)
	}

	return component, nil
}

func (s *Simulator) SubmitRequest(req *Request) {
	if s.running {
		s.eventQueue <- req
	}
}

func (s *Simulator) Start() {
	s.running = true
	go s.processEvents()
	go s.tick()
}

func (s *Simulator) Stop() {
	s.running = false
	s.cancel()
	close(s.eventQueue)
	close(s.responseQueue)
}

func (s *Simulator) processEvents() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case req := <-s.eventQueue:
			go s.handleRequest(req)
		}
	}
}

func (s *Simulator) handleRequest(req *Request) {
	s.componentMutex.RLock()
	var entryPoint Component
	
	// Strategy: Try to find CDN first, then LB, then API
	// In a real sim, we'd check Region match too, but keeping it simple for now.
	for _, comp := range s.components {
		if comp.GetType() == "cdn" {
			entryPoint = comp
			break
		}
	}
	if entryPoint == nil {
		for _, comp := range s.components {
			if comp.GetType() == "load-balancer" {
				entryPoint = comp
				break
			}
		}
	}
	if entryPoint == nil {
		for _, comp := range s.components {
			if comp.GetType() == "api-server" {
				entryPoint = comp
				break
			}
		}
	}
	s.componentMutex.RUnlock()

	s.metrics.mu.Lock()
	s.metrics.TotalRequests++
	s.metrics.mu.Unlock()

	if entryPoint == nil {
		s.metrics.mu.Lock()
		s.metrics.TotalFailures++
		s.metrics.mu.Unlock()
		return
	}

	// Process request
	resp, err := entryPoint.Process(req)

	s.metrics.mu.Lock()
	if err == nil && (resp == nil || resp.Success) {
		s.metrics.TotalSuccesses++
	} else {
		s.metrics.TotalFailures++
	}
	if resp != nil {
		s.metrics.TotalLatency += resp.Latency
		// Track individual latencies for P99 calculation (keep last 1000)
		s.metrics.RecentLatencies = append(s.metrics.RecentLatencies, resp.Latency)
		if len(s.metrics.RecentLatencies) > 1000 {
			s.metrics.RecentLatencies = s.metrics.RecentLatencies[1:]
		}
	}
	s.metrics.mu.Unlock()
}

func (s *Simulator) tick() {
	ticker := time.NewTicker(s.tickRate)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.currentTime = s.currentTime.Add(s.tickRate)
			s.updateMetrics()
		}
	}
}

func (s *Simulator) updateMetrics() {
	s.componentMutex.RLock()
	defer s.componentMutex.RUnlock()

	s.metrics.mu.Lock()
	defer s.metrics.mu.Unlock()

	var totalCost float64
	for id, component := range s.components {
		metrics := component.GetMetrics()
		s.metrics.ComponentMetrics[id] = metrics
		totalCost += component.GetCost()
	}
	s.metrics.TotalCost = totalCost
}

func (s *Simulator) GetMetrics() *AggregateMetrics {
	s.metrics.mu.RLock()
	defer s.metrics.mu.RUnlock()
	return s.metrics
}

// GetP99Latency calculates the 99th percentile latency from recent latencies
func (s *Simulator) GetP99Latency() time.Duration {
	s.metrics.mu.RLock()
	defer s.metrics.mu.RUnlock()

	if len(s.metrics.RecentLatencies) == 0 {
		return 0
	}

	// Sort latencies to find P99
	latencies := make([]time.Duration, len(s.metrics.RecentLatencies))
	copy(latencies, s.metrics.RecentLatencies)

	// Simple bubble sort for small datasets
	for i := 0; i < len(latencies); i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}

	// Calculate 99th percentile index
	p99Index := int(float64(len(latencies)) * 0.99)
	if p99Index >= len(latencies) {
		p99Index = len(latencies) - 1
	}

	return latencies[p99Index]
}

func (s *Simulator) GetCurrentTime() time.Time {
	return s.currentTime
}
