package database

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

type DatabaseType string

const (
	DatabaseTypeSQL      DatabaseType = "sql"
	DatabaseTypeNoSQL    DatabaseType = "nosql"
	DatabaseTypeKeyValue DatabaseType = "key-value"
	DatabaseTypeDocument DatabaseType = "document"
)

type Database struct {
	ID               string
	Type             DatabaseType
	Region           string
	Capacity         int64
	UsedCapacity     int64
	ReadLatency      time.Duration
	WriteLatency     time.Duration
	ReplicationLag   time.Duration
	Shards           []*Shard
	Replicas         []*Database
	IsPrimary        bool
	healthy          bool
	metrics          *engine.Metrics
	metricsMutex     sync.RWMutex
	costPerHour      float64
	data             map[string][]byte
	dataMutex        sync.RWMutex
}

type Shard struct {
	ID        string
	Database  *Database
	HashRange [2]uint32
}

func NewDatabase(id string, dbType DatabaseType, region string, capacity int64) *Database {
	return &Database{
		ID:           id,
		Type:         dbType,
		Region:       region,
		Capacity:     capacity,
		ReadLatency:  10 * time.Millisecond,
		WriteLatency: 15 * time.Millisecond,
		IsPrimary:    true,
		healthy:      true,
		metrics:      &engine.Metrics{},
		costPerHour:  0.05,
		data:         make(map[string][]byte),
		Shards:       make([]*Shard, 0),
		Replicas:     make([]*Database, 0),
	}
}

func (db *Database) GetID() string {
	return db.ID
}

func (db *Database) GetType() string {
	return fmt.Sprintf("database-%s", db.Type)
}

func (db *Database) Process(req *engine.Request) (*engine.Response, error) {
	start := time.Now()
	
	db.metricsMutex.Lock()
	db.metrics.RequestCount++
	db.metricsMutex.Unlock()

	if !db.healthy {
		db.metricsMutex.Lock()
		db.metrics.FailureCount++
		db.metricsMutex.Unlock()
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Latency:   time.Since(start),
			Error:     fmt.Errorf("database is unhealthy"),
		}, fmt.Errorf("database is unhealthy")
	}

	if len(db.Shards) > 0 {
		return db.processSharded(req)
	}

	var latency time.Duration
	var err error

	switch req.Type {
	case engine.RequestTypeRead:
		latency = db.ReadLatency
		err = db.read(req)
	case engine.RequestTypeWrite:
		if !db.IsPrimary {
			return &engine.Response{
				RequestID: req.ID,
				Success:   false,
				Latency:   time.Since(start),
				Error:     fmt.Errorf("cannot write to replica"),
			}, fmt.Errorf("cannot write to replica")
		}
		latency = db.WriteLatency
		err = db.write(req)
		db.replicateToReplicas(req)
	default:
		latency = db.ReadLatency
		err = db.read(req)
	}

	time.Sleep(latency)
	
	totalLatency := time.Since(start)
	
	db.metricsMutex.Lock()
	if err == nil {
		db.metrics.SuccessCount++
	} else {
		db.metrics.FailureCount++
	}
	db.metrics.TotalLatency += totalLatency
	db.metrics.AverageLatency = time.Duration(int64(db.metrics.TotalLatency) / db.metrics.RequestCount)
	db.metricsMutex.Unlock()

	return &engine.Response{
		RequestID: req.ID,
		Success:   err == nil,
		Latency:   totalLatency,
		DataSize:  req.DataSize,
		Error:     err,
		HopsTrace: []string{db.ID},
	}, err
}

func (db *Database) read(req *engine.Request) error {
	db.dataMutex.RLock()
	defer db.dataMutex.RUnlock()
	
	return nil
}

func (db *Database) write(req *engine.Request) error {
	db.dataMutex.Lock()
	defer db.dataMutex.Unlock()
	
	if db.UsedCapacity+req.DataSize > db.Capacity {
		return fmt.Errorf("database capacity exceeded")
	}
	
	db.data[req.Path] = make([]byte, req.DataSize)
	db.UsedCapacity += req.DataSize
	
	return nil
}

func (db *Database) replicateToReplicas(req *engine.Request) {
	for _, replica := range db.Replicas {
		go func(r *Database) {
			time.Sleep(r.ReplicationLag)
			r.write(req)
		}(replica)
	}
}

func (db *Database) processSharded(req *engine.Request) (*engine.Response, error) {
	shard := db.selectShard(req)
	if shard == nil {
		return &engine.Response{
			RequestID: req.ID,
			Success:   false,
			Error:     fmt.Errorf("no shard found for request"),
		}, fmt.Errorf("no shard found for request")
	}
	
	return shard.Database.Process(req)
}

func (db *Database) selectShard(req *engine.Request) *Shard {
	h := fnv.New32a()
	h.Write([]byte(req.UserID))
	hash := h.Sum32()
	
	for _, shard := range db.Shards {
		if hash >= shard.HashRange[0] && hash <= shard.HashRange[1] {
			return shard
		}
	}
	
	return nil
}

func (db *Database) AddShard(shard *Shard) {
	db.Shards = append(db.Shards, shard)
}

func (db *Database) AddReplica(replica *Database) {
	replica.IsPrimary = false
	db.Replicas = append(db.Replicas, replica)
}

func (db *Database) GetMetrics() *engine.Metrics {
	db.metricsMutex.RLock()
	defer db.metricsMutex.RUnlock()
	
	metricsCopy := *db.metrics
	if metricsCopy.RequestCount > 0 {
		metricsCopy.ErrorRate = float64(metricsCopy.FailureCount) / float64(metricsCopy.RequestCount)
	}
	return &metricsCopy
}

func (db *Database) GetCost() float64 {
	baseCost := db.costPerHour
	capacityCost := float64(db.Capacity) / (1024 * 1024 * 1024) * 0.01
	
	return baseCost + capacityCost
}

func (db *Database) IsHealthy() bool {
	return db.healthy
}

func (db *Database) SetHealthy(healthy bool) {
	db.healthy = healthy
}
