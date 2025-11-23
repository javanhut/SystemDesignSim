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
│  │  [View Tutorial] [Design Patterns Tutorial]          │   │
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
└────┬───────────────────────┬────────────────────────────────┘
     │ Click [Play]          │ Click [Design Patterns]
     ▼                       ▼
                    ┌─────────────────────────────────┐
                    │  Design Pattern Selection       │
                    │  ┌────────────┐  ┌────────────┐ │
                    │  │Load Balance│  │Cache-Aside │ │
                    │  │★★☆☆☆      │  │★★☆☆☆      │ │
                    │  │[Watch Demo]│  │[Watch Demo]│ │
                    │  │[Practice]  │  │[Practice]  │ │
                    │  └────────────┘  └────────────┘ │
                    │  ┌────────────┐                 │
                    │  │Read Replica│                 │
                    │  │★★★☆☆      │                 │
                    │  │[Watch Demo]│                 │
                    │  │[Practice]  │                 │
                    │  └────────────┘                 │
                    └────────┬────────────────────────┘
                             │ Click [Watch Demo]
                             ▼
                    ┌─────────────────────────────────┐
                    │  Pattern Tutorial Screen        │
                    │  ┌──────┬──────────┬─────────┐ │
                    │  │ Info │ Canvas   │Controls │ │
                    │  │      │ (Animated│         │ │
                    │  │Problem│  Demo)  │▶ Play  │ │
                    │  │      │          │⏸ Pause │ │
                    │  │Solution│         │↻Restart│ │
                    │  │      │          │         │ │
                    │  │Benefits│         │Step 5/12│ │
                    │  │      │          │█████░░░│ │
                    │  │Tradeoffs│        │         │ │
                    │  └──────┴──────────┴─────────┘ │
                    └─────────────────────────────────┘
                             │
                             ▼
                    [Demo Complete - Try Practice]
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                       Game Screen                            │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │           Level 1: Local Blog                           │ │
│ │         Build • Connect • Simulate                      │ │
│ ├─────────────────────────────────────────────────────────┤ │
│ │  CLIENT: Personal Blog | BUSINESS: Blogging Platform   │ │
│ │  SITUATION: Growing audience, need scalable backend    │ │
│ │  USERS: 10 concurrent (peak: 50) | SESSION: 5min, 3pv  │ │
│ │  TRAFFIC: 80% reads | 15% writes | 5% static           │ │
│ │  CONSTRAINTS: Budget $10/hr | P99 < 500ms | Up > 95%   │ │
│ └─────────────────────────────────────────────────────────┘ │
│ ┌─────────┬─────────────────────────────┬─────────────────┐ │
│ │ Toolbox │        Canvas Area          │  Metrics Panel  │ │
│ │         │                             │                 │ │
│ │ [API]   │   ┌─────┐      ┌─────┐    │ Status: ✓ PASS  │ │
│ │ [DB]    │   │ API │─────▶│ DB  │    │                 │ │
│ │ [Cache] │   └─────┘      └─────┘    │ Objectives:     │ │
│ │ [LB]    │     (particles flowing)    │ ✓ Latency < 500ms│ │
│ │ [CDN]   │                            │ ✓ Uptime > 95%  │ │
│ │         │                            │ ✓ Cost < $10    │ │
│ │ [? Help]│                            │                 │ │
│ │         │                            │ Metrics:        │ │
│ │         │                            │ Requests: 1234  │ │
│ │         │                            │ RPS: 25         │ │
│ │         │                            │ P99: 45ms ✓    │ │
│ │         │                            │ Uptime: 99.5% ✓│ │
│ │         │                            │ Cost: $2.50 ✓  │ │
│ │         │                            │                 │ │
│ │         │                            │ Hints (scroll)  │ │
│ │         │                            │ Summary         │ │
│ │         │                            │ Test Plan       │ │
│ └─────────┴─────────────────────────────┴─────────────────┘ │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ [Start] [Stop] [Submit] [Show Hints]                   │ │
│ │ [Control Center] [System Plan] [Back to Levels]        │ │
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

## Architectural Hints System

### Hint Display Modes

```
1. Inline Hints (Right Panel)
   └─ Always visible, scrollable
   └─ Updates dynamically as architecture changes
   └─ Shows current architecture analysis

2. Popup Hints Dialog
   └─ Click "Show Hints" button
   └─ Modal popup (700x600)
   └─ Scrollable, focused view
   └─ Same content as inline hints
```

### Hint Content Structure

```
SYSTEM DESIGN PRINCIPLES:

1. REQUEST FLOW PATTERN
   ├─ Basic flow diagram
   ├─ Production flow diagram
   ├─ Component status (✓ present, ✗ missing)
   └─ WHY explanation for each component

2. HORIZONTAL SCALABILITY
   ├─ Scale out vs scale up principle
   ├─ Load balancer necessity
   ├─ Server count recommendations
   └─ Load distribution calculations

3. LATENCY OPTIMIZATION
   ├─ Target P99 latency from level requirements
   ├─ Cache impact (DB: 10-50ms vs Cache: 1-2ms)
   ├─ CDN benefits for static content
   └─ Real-world technologies (Redis, Memcached)

4. HIGH AVAILABILITY
   ├─ Uptime target with calculations
   │   (99.9% = 3 nines = 8.76hr downtime/year)
   ├─ SPOF (Single Point of Failure) detection
   ├─ Redundancy requirements (N+1, Active-Active)
   └─ Health check patterns

5. TRAFFIC HANDLING
   ├─ Peak load from level requirements
   ├─ Users per server calculations
   ├─ Database scaling strategies
   └─ Connection pooling recommendations

6. COST OPTIMIZATION
   ├─ Budget from level requirements
   ├─ Cost-saving strategies
   ├─ Right-sizing recommendations
   └─ Auto-scaling benefits

7. COMMON MISTAKES TO AVOID
   ├─ Single API server (SPOF)
   ├─ No caching (high latency + DB overload)
   ├─ No load balancer (can't scale horizontally)
   ├─ Direct DB access from internet
   └─ No monitoring
```

### Hint Update Flow

```
User adds/removes component
         ↓
Architecture analysis runs
         ↓
┌──────────────────────────────────┐
│ Analyze current components:      │
│ - Count by type                  │
│ - Check for critical components  │
│ - Calculate ratios               │
└──────────────────────────────────┘
         ↓
Generate dynamic hints
         ↓
┌──────────────────────────────────┐
│ For each principle:              │
│ - Show status (✓/✗)             │
│ - Calculate recommendations      │
│ - Include WHY explanations       │
│ - Show real-world examples       │
└──────────────────────────────────┘
         ↓
Update hint displays
         ↓
┌──────────────────────────────────┐
│ 1. Right panel label updates     │
│ 2. Popup (if open) updates       │
└──────────────────────────────────┘
```

## Educational Progression

```
Level 1: Basics
  Learn: Component placement, connections, scenario reading
  Architecture: Single API + DB
  Hints Focus: Basic flow, component roles

Level 2: Scaling
  Learn: Load balancing, horizontal scaling, hints usage
  Architecture: LB + Multiple APIs + DB
  Hints Focus: Redundancy, SPOF elimination

Level 3: Optimization
  Learn: Caching, replication, latency optimization
  Architecture: LB + APIs + Cache + Replicated DB
  Hints Focus: Cache strategies, performance tuning

Level 4: Distribution
  Learn: Multi-region, CDN, global distribution
  Architecture: Multi-region with CDN
  Hints Focus: Geographic distribution, edge caching

Level 5: Mastery
  Learn: Sharding, global distribution, cost optimization
  Architecture: Global multi-region with all features
  Hints Focus: Advanced patterns, trade-offs
```

## Design Patterns Tutorial System

### Pattern Learning Flow

```
User accesses Design Patterns
         ↓
Pattern Selection Screen
  ├─ Grid of pattern cards
  ├─ Category-based colors
  ├─ Difficulty indicators (★★★☆☆)
  └─ Two modes per pattern:
      ├─ Watch Demo
      └─ Try Practice
         ↓
User clicks "Watch Demo"
         ↓
Pattern Tutorial Screen loads
         ↓
┌──────────────────────────────────────┐
│ TutorialOrchestrator starts          │
│ Executes steps sequentially:         │
│  1. StepMessage → Show overlay       │
│  2. StepAddComponent → Fade-in       │
│  3. StepCreateConnection → Draw line │
│  4. StepShowTraffic → Spawn particles│
│  5-12. Continue pattern demo         │
└──────────────────────────────────────┘
         ↓
Progress bar updates (Step N/12)
         ↓
Demo completes
         ↓
"Try Practice Mode" button enabled
         ↓
User clicks "Try Practice"
         ↓
Practice Mode starts
  ├─ Show instruction: "Add a load balancer"
  ├─ User adds component
  ├─ Real-time validation
  ├─ Feedback: "✓ Perfect!" or "✗ Need 1 more..."
  └─ Next step unlocked when validated
         ↓
All practice steps completed
         ↓
"Practice Complete!" message
         ↓
Return to pattern selection or game
```

### Pattern Tutorial Architecture

```
PatternTutorialScreen (UI Layer)
  ├─ Left Panel: Pattern information
  │   ├─ Problem statement
  │   ├─ Solution overview
  │   ├─ Benefits list
  │   ├─ Trade-offs list
  │   └─ Real-world examples
  │
  ├─ Center Panel: Interactive canvas
  │   ├─ Component visualization
  │   ├─ Connection lines
  │   ├─ Particle animations
  │   └─ Message overlays
  │
  └─ Right Panel: Controls
      ├─ Demo controls (Play/Pause/Restart)
      ├─ Progress tracking
      ├─ Practice instructions
      └─ Validation feedback

TutorialOrchestrator (Logic Layer)
  ├─ Step execution engine
  ├─ Component spawning
  ├─ Connection automation
  ├─ Traffic simulation
  ├─ Practice validation
  └─ Animation timing

DesignPattern (Data Layer)
  ├─ Educational content
  ├─ Demo steps array
  ├─ Practice steps array
  └─ Validation criteria
```

### Available Patterns

1. **Load Balancing** (Difficulty: 2/5, Category: Scalability)
   - 12 demo steps showing LB + 3 API servers
   - 3 practice steps with validation
   - Real-world: Netflix, AWS ELB, NGINX, HAProxy

2. **Cache-Aside** (Difficulty: 2/5, Category: Performance)
   - 10 demo steps showing API → Cache → DB flow
   - 4 practice steps building cache chain
   - Real-world: Facebook, Twitter, Amazon, Stack Overflow

3. **Read Replicas** (Difficulty: 3/5, Category: Scalability)
   - 15 demo steps showing primary + 2 replicas
   - 3 practice steps with replication setup
   - Real-world: Instagram, GitHub, Shopify, Pinterest

### Pattern Integration Points

**Access from Level Select**:
- "Design Patterns Tutorial" button
- Standalone learning mode
- No prerequisites

**Access from Game Screen**:
- "Learn Patterns" button in controls
- Quick reference during gameplay
- Context-aware recommendations (future)

**Level-Specific Recommendations** (future):
- Level 2 → Load Balancing pattern
- Level 3 → Cache-Aside pattern
- Level 4 → Read Replicas pattern

## Success!

The game successfully teaches system design through:
1. Visual, interactive learning
2. Immediate feedback with real-time metrics
3. Progressive complexity across levels
4. Real-world constraints (budget, latency, uptime)
5. Hands-on experimentation
6. Comprehensive architectural guidance (hints system)
7. Context-aware scenario descriptions
8. System design principles with WHY explanations
9. **NEW**: Animated pattern tutorials with practice mode
10. **NEW**: Step-by-step guided learning with validation
