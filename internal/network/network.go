package network

import (
	"math/rand"
	"time"
)

type LatencyProfile struct {
	BaseLatency time.Duration
	Jitter      time.Duration
	PacketLoss  float64
}

var RegionalLatency = map[string]map[string]time.Duration{
	"us-east": {
		"us-east":    5 * time.Millisecond,
		"us-west":    70 * time.Millisecond,
		"europe":     80 * time.Millisecond,
		"asia":       180 * time.Millisecond,
		"australia":  200 * time.Millisecond,
	},
	"us-west": {
		"us-east":    70 * time.Millisecond,
		"us-west":    5 * time.Millisecond,
		"europe":     140 * time.Millisecond,
		"asia":       120 * time.Millisecond,
		"australia":  150 * time.Millisecond,
	},
	"europe": {
		"us-east":    80 * time.Millisecond,
		"us-west":    140 * time.Millisecond,
		"europe":     5 * time.Millisecond,
		"asia":       120 * time.Millisecond,
		"australia":  280 * time.Millisecond,
	},
	"asia": {
		"us-east":    180 * time.Millisecond,
		"us-west":    120 * time.Millisecond,
		"europe":     120 * time.Millisecond,
		"asia":       5 * time.Millisecond,
		"australia":  100 * time.Millisecond,
	},
	"australia": {
		"us-east":    200 * time.Millisecond,
		"us-west":    150 * time.Millisecond,
		"europe":     280 * time.Millisecond,
		"asia":       100 * time.Millisecond,
		"australia":  5 * time.Millisecond,
	},
}

func CalculateLatency(fromRegion, toRegion string, profile LatencyProfile) time.Duration {
	baseLatency := profile.BaseLatency
	
	if regionalLatency, ok := RegionalLatency[fromRegion]; ok {
		if latency, ok := regionalLatency[toRegion]; ok {
			baseLatency += latency
		}
	}

	jitter := time.Duration(rand.Int63n(int64(profile.Jitter)))
	return baseLatency + jitter
}

func SimulatePacketLoss(packetLossRate float64) bool {
	return rand.Float64() < packetLossRate
}

func CalculateBandwidth(dataSize int64, bandwidth int64) time.Duration {
	if bandwidth == 0 {
		return 0
	}
	
	bitsPerSecond := bandwidth * 1024 * 1024
	bits := dataSize * 8
	seconds := float64(bits) / float64(bitsPerSecond)
	
	return time.Duration(seconds * float64(time.Second))
}
