package engine

import (
	"time"
)

type RequestType string

const (
	RequestTypeRead  RequestType = "read"
	RequestTypeWrite RequestType = "write"
	RequestTypeAPI   RequestType = "api"
)

type Request struct {
	ID          string
	Type        RequestType
	Timestamp   time.Time
	UserID      string
	Region      string
	DataSize    int64
	Path        string
	Headers     map[string]string
	Metadata    map[string]interface{}
}

type Response struct {
	RequestID   string
	Success     bool
	Latency     time.Duration
	DataSize    int64
	Error       error
	CacheHit    bool
	HopsTrace   []string
	Metadata    map[string]interface{}
}

type Component interface {
	GetID() string
	GetType() string
	Process(req *Request) (*Response, error)
	GetMetrics() *Metrics
	GetCost() float64
	IsHealthy() bool
	SetHealthy(bool)
}

type Metrics struct {
	RequestCount    int64
	SuccessCount    int64
	FailureCount    int64
	TotalLatency    time.Duration
	AverageLatency  time.Duration
	P95Latency      time.Duration
	P99Latency      time.Duration
	Throughput      float64
	ErrorRate       float64
	CacheHitRate    float64
	DataTransferred int64
}

type Region string

const (
	RegionUSEast    Region = "us-east"
	RegionUSWest    Region = "us-west"
	RegionEurope    Region = "europe"
	RegionAsia      Region = "asia"
	RegionAustralia Region = "australia"
)

type ScalingType string

const (
	ScalingVertical   ScalingType = "vertical"
	ScalingHorizontal ScalingType = "horizontal"
)
