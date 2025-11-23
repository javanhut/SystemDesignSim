package networking

import (
	"errors"
	"sync"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

// Gateway - Internet gateway or API gateway
type Gateway struct {
	ID           string
	Region       string
	Throughput   int
	Backend      engine.Component
	healthy      bool
	metrics      *engine.Metrics
	metricsMutex sync.RWMutex
	cost         float64
}

func NewGateway(id, region string) *Gateway {
	return &Gateway{
		ID:         id,
		Region:     region,
		Throughput: 10000, // requests per second
		healthy:    true,
		metrics:    &engine.Metrics{},
		cost:       0.03,
	}
}

func (g *Gateway) SetBackend(backend engine.Component) {
	g.Backend = backend
}

func (g *Gateway) GetID() string           { return g.ID }
func (g *Gateway) GetType() string         { return "gateway" }
func (g *Gateway) GetRegion() string       { return g.Region }
func (g *Gateway) IsHealthy() bool         { return g.healthy }
func (g *Gateway) SetHealthy(healthy bool) { g.healthy = healthy }
func (g *Gateway) GetCost() float64        { return g.cost }

func (g *Gateway) Process(req *engine.Request) (*engine.Response, error) {
	g.metricsMutex.Lock()
	g.metrics.RequestCount++
	g.metrics.Throughput++
	g.metricsMutex.Unlock()

	// Minimal latency for gateway (~1ms)
	time.Sleep(1 * time.Millisecond)

	if g.Backend != nil {
		return g.Backend.Process(req)
	}

	return &engine.Response{
		RequestID: req.ID,
		Success:   false,
		Latency:   1 * time.Millisecond,
		Error:     errors.New("no backend configured"),
	}, nil
}

func (g *Gateway) GetMetrics() *engine.Metrics {
	g.metricsMutex.RLock()
	defer g.metricsMutex.RUnlock()
	return g.metrics
}

// Firewall - Security filtering layer
type Firewall struct {
	ID             string
	Region         string
	Rules          []string
	BlockedCount   int
	Backend        engine.Component
	healthy        bool
	metrics        *engine.Metrics
	metricsMutex   sync.RWMutex
	cost           float64
}

func NewFirewall(id, region string) *Firewall {
	return &Firewall{
		ID:      id,
		Region:  region,
		Rules:   []string{"allow-http", "allow-https"},
		healthy: true,
		metrics: &engine.Metrics{},
		cost:    0.02,
	}
}

func (f *Firewall) SetBackend(backend engine.Component) {
	f.Backend = backend
}

func (f *Firewall) GetID() string           { return f.ID }
func (f *Firewall) GetType() string         { return "firewall" }
func (f *Firewall) GetRegion() string       { return f.Region }
func (f *Firewall) IsHealthy() bool         { return f.healthy }
func (f *Firewall) SetHealthy(healthy bool) { f.healthy = healthy }
func (f *Firewall) GetCost() float64        { return f.cost }

func (f *Firewall) Process(req *engine.Request) (*engine.Response, error) {
	f.metricsMutex.Lock()
	f.metrics.RequestCount++
	f.metrics.Throughput++
	f.metricsMutex.Unlock()

	// Firewall processing (~2ms)
	time.Sleep(2 * time.Millisecond)

	if f.Backend != nil {
		return f.Backend.Process(req)
	}

	return &engine.Response{
		RequestID: req.ID,
		Success:   false,
		Latency:   2 * time.Millisecond,
		Error:     errors.New("no backend configured"),
	}, nil
}

func (f *Firewall) GetMetrics() *engine.Metrics {
	f.metricsMutex.RLock()
	defer f.metricsMutex.RUnlock()
	return f.metrics
}

// NAT - Network Address Translation
type NAT struct {
	ID           string
	Region       string
	Backend      engine.Component
	healthy      bool
	metrics      *engine.Metrics
	metricsMutex sync.RWMutex
	cost         float64
}

func NewNAT(id, region string) *NAT {
	return &NAT{
		ID:      id,
		Region:  region,
		healthy: true,
		metrics: &engine.Metrics{},
		cost:    0.045,
	}
}

func (n *NAT) SetBackend(backend engine.Component) {
	n.Backend = backend
}

func (n *NAT) GetID() string           { return n.ID }
func (n *NAT) GetType() string         { return "nat" }
func (n *NAT) GetRegion() string       { return n.Region }
func (n *NAT) IsHealthy() bool         { return n.healthy }
func (n *NAT) SetHealthy(healthy bool) { n.healthy = healthy }
func (n *NAT) GetCost() float64        { return n.cost }

func (n *NAT) Process(req *engine.Request) (*engine.Response, error) {
	n.metricsMutex.Lock()
	n.metrics.RequestCount++
	n.metrics.Throughput++
	n.metricsMutex.Unlock()

	// NAT translation (~1ms)
	time.Sleep(1 * time.Millisecond)

	if n.Backend != nil {
		return n.Backend.Process(req)
	}

	return &engine.Response{
		RequestID: req.ID,
		Success:   false,
		Latency:   1 * time.Millisecond,
		Error:     errors.New("no backend configured"),
	}, nil
}

func (n *NAT) GetMetrics() *engine.Metrics {
	n.metricsMutex.RLock()
	defer n.metricsMutex.RUnlock()
	return n.metrics
}

// Router - Network routing layer
type Router struct {
	ID           string
	Region       string
	Routes       map[string]engine.Component
	healthy      bool
	metrics      *engine.Metrics
	metricsMutex sync.RWMutex
	cost         float64
}

func NewRouter(id, region string) *Router {
	return &Router{
		ID:      id,
		Region:  region,
		Routes:  make(map[string]engine.Component),
		healthy: true,
		metrics: &engine.Metrics{},
		cost:    0.015,
	}
}

func (r *Router) AddRoute(path string, component engine.Component) {
	r.Routes[path] = component
}

func (r *Router) GetID() string           { return r.ID }
func (r *Router) GetType() string         { return "router" }
func (r *Router) GetRegion() string       { return r.Region }
func (r *Router) IsHealthy() bool         { return r.healthy }
func (r *Router) SetHealthy(healthy bool) { r.healthy = healthy }
func (r *Router) GetCost() float64        { return r.cost }

func (r *Router) Process(req *engine.Request) (*engine.Response, error) {
	r.metricsMutex.Lock()
	r.metrics.RequestCount++
	r.metrics.Throughput++
	r.metricsMutex.Unlock()

	// Routing decision (~1ms)
	time.Sleep(1 * time.Millisecond)

	// Simple routing based on path
	if backend, ok := r.Routes[req.Path]; ok && backend != nil {
		return backend.Process(req)
	}

	// Default route
	for _, backend := range r.Routes {
		if backend != nil {
			return backend.Process(req)
		}
	}

	return &engine.Response{
		RequestID: req.ID,
		Success:   false,
		Latency:   1 * time.Millisecond,
		Error:     errors.New("no route found"),
	}, nil
}

func (r *Router) GetMetrics() *engine.Metrics {
	r.metricsMutex.RLock()
	defer r.metricsMutex.RUnlock()
	return r.metrics
}

// UserPool - Simulated user traffic source
type UserPool struct {
	ID             string
	Region         string
	UserCount      int
	RequestRate    int // requests per second per user
	healthy        bool
	metrics        *engine.Metrics
	metricsMutex   sync.RWMutex
	cost           float64
}

func NewUserPool(id, region string, userCount int) *UserPool {
	return &UserPool{
		ID:          id,
		Region:      region,
		UserCount:   userCount,
		RequestRate: 5, // 5 requests/sec per user
		healthy:     true,
		metrics:     &engine.Metrics{},
		cost:        0.0, // No cost for simulated users
	}
}

func (u *UserPool) GetID() string           { return u.ID }
func (u *UserPool) GetType() string         { return "user-pool" }
func (u *UserPool) GetRegion() string       { return u.Region }
func (u *UserPool) IsHealthy() bool         { return u.healthy }
func (u *UserPool) SetHealthy(healthy bool) { u.healthy = healthy }
func (u *UserPool) GetCost() float64        { return u.cost }

func (u *UserPool) Process(req *engine.Request) (*engine.Response, error) {
	// UserPool generates requests, doesn't process them
	return &engine.Response{
		RequestID: req.ID,
		Success:   true,
		Latency:   0,
	}, nil
}

func (u *UserPool) GetMetrics() *engine.Metrics {
	u.metricsMutex.RLock()
	defer u.metricsMutex.RUnlock()
	return u.metrics
}

func (u *UserPool) GetTotalRequestRate() int {
	return u.UserCount * u.RequestRate
}
