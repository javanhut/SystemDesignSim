# System Design Simulator - Game Flow

## Visual Game Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Start                         │
│                         (main.go)                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Level Selection Screen                      │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  System Design Simulator                             │   │
│  │                                                       │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │ Level 1: Local Blog             [Unlocked]    │  │   │
│  │  │ Handle 10 users - Budget: $10                 │  │   │
│  │  │                                  [Play Button] │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  │                                                       │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │ Level 2: Growing Blog           [Locked]      │  │   │
│  │  │ Handle 100 users - Budget: $50                │  │   │
│  │  │                                  [Locked]      │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  │                                                       │   │
│  │  ... more levels ...                                 │   │
│  └─────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ Click [Play]
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                       Game Screen                            │
│ ┌─────────┬─────────────────────────────┬─────────────────┐ │
│ │ Toolbox │        Canvas Area          │  Metrics Panel  │ │
│ │         │                             │                 │ │
│ │ [API]   │   ┌─────┐      ┌─────┐    │ Level: Blog     │ │
│ │ [DB]    │   │ API │─────▶│ DB  │    │                 │ │
│ │ [Cache] │   └─────┘      └─────┘    │ Objectives:     │ │
│ │ [LB]    │                            │ ✓ Latency < 500ms│ │
│ │ [CDN]   │                            │ ✓ Uptime > 95%  │ │
│ │         │                            │ ○ Cost < $10    │ │
│ │         │                            │                 │ │
│ │         │                            │ Metrics:        │ │
│ │         │                            │ Requests: 1234  │ │
│ │         │                            │ Latency: 45ms   │ │
│ │         │                            │ Cost: $2.50/hr  │ │
│ └─────────┴─────────────────────────────┴─────────────────┘ │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ [Start] [Stop] [Submit Solution] [Back to Levels]      │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                         │ Click [Submit]
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Results Screen                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │             Level 1: Local Blog                      │   │
│  │                                                       │   │
│  │                    PASSED!                           │   │
│  │                                                       │   │
│  │                Score: 1,250                          │   │
│  │                                                       │   │
│  │  Metrics:                                            │   │
│  │  - Uptime: 99.5%       ✓                            │   │
│  │  - Latency: 45ms       ✓                            │   │
│  │  - Error Rate: 0.2%    ✓                            │   │
│  │  - Cost: $2.50         ✓                            │   │
│  │                                                       │   │
│  │  Bonuses Earned:                                     │   │
│  │  + Fast response time  (+100)                        │   │
│  │  + Cost efficient      (+150)                        │   │
│  │                                                       │   │
│  │                    [OK]                              │   │
│  └─────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ Click [OK]
                         ▼
                  Back to Level Selection
                  (Level 2 now unlocked)
```

## Interaction Flow

### Adding Components

```
User clicks [API Server] button
         ↓
Visual Component created at position
         ↓
Backend API Server component instantiated
         ↓
Component added to Simulator
         ↓
Visual component appears on canvas
```

### Connecting Components

```
User right-clicks API Server component
         ↓
User drags to Database component
         ↓
Connection created visually
         ↓
Backend components linked (API.SetDatabase(DB))
         ↓
Connection line drawn on canvas
         ↓
Animated particles flow along connection
```

### Running Simulation

```
User clicks [Start Simulation]
         ↓
Game state initialized
         ↓
Simulator starts processing
         ↓
┌─────────────────────────────────────┐
│  Three concurrent goroutines:       │
│  1. Traffic Generator (100ms tick)  │
│  2. Metrics Updater (500ms tick)    │
│  3. Visual Refresh (on demand)      │
└─────────────────────────────────────┘
         ↓
Requests flow through component graph
         ↓
Metrics collected and aggregated
         ↓
Visual feedback updated (colors, stats)
         ↓
User observes real-time behavior
```

## Request Flow Through Components

```
Traffic Generator
       ↓
   Request {
     ID, Type, UserID,
     Region, DataSize, Path
   }
       ↓
Simulator Queue
       ↓
Entry Component (e.g., Load Balancer)
       ↓
   [Round-robin selection]
       ↓
Backend API Server #1
       ↓
   [Check capacity]
   [Simulate processing time]
       ↓
Check for Cache?
  Yes ↓         No ↓
  Cache         Database
    ↓              ↓
  Hit?         [Simulate latency]
    ↓              ↓
  Return       Read/Write data
  fast           ↓
    ↓          Return data
    ↓              ↓
    └──────┬───────┘
           ↓
       Response {
         Success, Latency,
         DataSize, CacheHit,
         HopsTrace, Error
       }
           ↓
    Update Metrics
           ↓
    Visual Feedback
```

## Component State Machine

```
Component Created
       ↓
  [Healthy] ←──────┐
   (green)         │
       ↓           │
   Processing      │
   Requests        │
       │           │
       ├→ Load < 50%: Stay Healthy
       │           │
       ├→ Load 50-80%: Warning
       │    (yellow)│
       │           │
       └→ Load > 80%: Critical
            (orange)│
                   │
    Error Rate > 10%
                   ↓
              [Down]
               (red)
                   ↓
            No requests
            processed
```

## Metrics Calculation

```
Every 100ms Tick:
  ├─ Process queued requests
  ├─ Update component metrics
  │    ├─ Request count ++
  │    ├─ Success/Failure count
  │    ├─ Latency tracking
  │    └─ Cost accumulation
  └─ Aggregate to simulator metrics

Every 500ms:
  ├─ Calculate derived metrics
  │    ├─ Error rate
  │    ├─ Throughput
  │    ├─ Cache hit rate
  │    └─ P95/P99 latency
  └─ Update GUI displays
```

## Scoring Algorithm

```
Start: Base Score = 1000

Check Requirements:
  ├─ Uptime < required?     → -200
  ├─ Latency > max?         → -200
  ├─ Error rate > max?      → -200
  ├─ Cost > budget?         → -200
  ├─ Missing required LB?   → -100
  └─ Missing required CDN?  → -100

If PASSED:
  ├─ Uptime exceeds target?       → +100
  ├─ Latency under target?        → +100
  ├─ Error rate under target?     → +100
  ├─ Cost under target budget?    → +150
  ├─ Great cache hit rate?        → +50
  └─ Cost savings bonus           → up to +200

Final Score = max(0, calculated score)
```

## Data Structures

### Visual Component
```go
VisualComponent {
  ID: "api-1"
  Type: "api-server"
  Position: (200, 150)
  Size: (80, 80)
  Component: *APIServer
  Connections: []*Connection
  Selected: false
  HealthStatus: Healthy
}
```

### Connection
```go
Connection {
  ID: "conn-api-1-db-1"
  From: *VisualComponent
  To: *VisualComponent
  Particles: []*Particle {
    Position: 0.5,  // 50% along line
    Speed: 0.01,
    Color: blue
  }
}
```

### Request
```go
Request {
  ID: "req-12345"
  Type: Read
  UserID: "user-42"
  Region: "us-east"
  DataSize: 1024
  Path: "/api/data/42"
}
```

### Response
```go
Response {
  RequestID: "req-12345"
  Success: true
  Latency: 45ms
  CacheHit: true
  HopsTrace: ["cdn-1", "api-1", "cache-1"]
}
```

## Performance Characteristics

### Throughput
- Simulation: 1,000+ requests/second
- GUI: 60 FPS refresh rate
- Metrics: 2 updates/second

### Latency Simulation
- API: 10-15ms base
- Cache: 1-2ms
- Database: 10-15ms read, 15-20ms write
- Network: Varies by region (5ms-200ms)
- CDN: 2ms edge hit

### Resource Usage
- Memory: ~50MB for simulation
- CPU: Low (event-driven, not polling)
- GPU: Minimal (2D rendering only)

## Educational Progression

```
Level 1: Basics
  Learn: Component placement, connections
  Architecture: Single API + DB

Level 2: Scaling
  Learn: Load balancing, horizontal scaling
  Architecture: LB + Multiple APIs + DB

Level 3: Optimization
  Learn: Caching, replication
  Architecture: LB + APIs + Cache + Replicated DB

Level 4: Distribution
  Learn: Multi-region, CDN
  Architecture: Multi-region with CDN

Level 5: Mastery
  Learn: Sharding, global distribution, optimization
  Architecture: Global multi-region with all features
```

## Success!

The game successfully teaches system design through:
1. Visual, interactive learning
2. Immediate feedback
3. Progressive complexity
4. Real-world constraints
5. Hands-on experimentation
