# System Design Simulator - Architecture Documentation

## Overview

The System Design Simulator is built using a layered architecture that separates concerns between simulation logic, infrastructure components, game mechanics, and the graphical user interface.

## Core Layers

### 1. Simulation Engine Layer (`internal/engine`)

The simulation engine is the heart of the application, responsible for:
- Processing incoming requests asynchronously
- Managing component lifecycle
- Collecting and aggregating metrics
- Simulating time progression

#### Key Components:

**Simulator**
- Event-driven request processing
- Component registry
- Metrics aggregation
- Clock simulation with configurable tick rate

**Request/Response Flow**
```
User → Request Queue → Component Graph → Response Queue → Metrics
```

**Component Interface**
All infrastructure components implement this interface:
```go
type Component interface {
    GetID() string
    GetType() string
    Process(req *Request) (*Response, error)
    GetMetrics() *Metrics
    GetCost() float64
    IsHealthy() bool
    SetHealthy(bool)
}
```

### 2. Infrastructure Components Layer (`internal/components`)

Each infrastructure component simulates a specific type of system:

#### API Server (`api/`)
- Simulates concurrent request handling
- Configurable instance sizes (small, medium, large, xlarge)
- Capacity limits and load tracking
- Processing time simulation
- Backend integration (database, cache)

#### Database (`database/`)
- Multiple database types (SQL, NoSQL, key-value, document)
- Read/write latency simulation
- Sharding support with consistent hashing
- Replication with configurable lag
- Capacity management

#### Cache (`cache/`)
- Multiple eviction policies (LRU, LFU, FIFO)
- TTL-based expiration
- Cache hit/miss tracking
- Backend fallback on miss
- Memory capacity management

#### Load Balancer (`loadbalancer/`)
- Multiple balancing strategies:
  - Round-robin
  - Least-connected
  - Weighted-random
  - IP hash
- Health checking of backends
- Connection tracking

#### CDN (`cdn/`)
- Edge locations in multiple regions
- Origin server fallback
- Regional cache management
- Hit rate optimization

### 3. Network Simulation Layer (`internal/network`)

Simulates realistic network conditions:
- Regional latency matrices
- Bandwidth calculations
- Packet loss simulation
- Jitter modeling

### 4. Game Logic Layer (`internal/game`)

Manages game progression and scoring:

#### Level System
- Progressive difficulty
- Objective tracking
- Success criteria validation
- Bonus objectives

#### Scoring Algorithm
```
Base Score: 1000
- Penalties for requirement failures: -200 each
- Bonuses for exceeding targets: +50 to +150 each
- Cost savings bonus: up to +200
- Minimum score: 0
```

#### Level Progression
- Linear unlocking (complete level N to unlock N+1)
- Best score tracking
- Replayability for higher scores

### 5. GUI Layer (`internal/gui`)

Provides interactive visual interface using Fyne:

#### Visual Component System
- Component positioning on canvas
- Connection graph representation
- Health status visualization
- Drag-and-drop interaction

#### Graph Canvas
- Custom Fyne widget
- Real-time rendering
- Event handling (tap, drag, etc.)
- Animated particle effects for traffic

#### Screen System
- Level selection screen
- Game screen with panels:
  - Left: Component toolbox
  - Center: Interactive canvas
  - Right: Metrics and objectives
  - Bottom: Controls

## Data Flow

### Request Processing Flow
```
1. Traffic Generator creates Request
2. Request enters Simulator queue
3. Simulator dispatches to entry Component
4. Component processes and may forward to connected components
5. Component returns Response
6. Metrics updated
7. Visual feedback updated
```

### Component Linking Flow
```
1. User creates visual connection in GUI
2. Connection callback triggered
3. Backend components linked based on types
4. Connection validated
5. Visual connection rendered
```

### Metrics Update Flow
```
1. Components track local metrics
2. Simulator aggregates every tick
3. GUI polls metrics periodically
4. Visual indicators updated (colors, labels)
5. Metrics panel refreshed
```

## Design Patterns

### Observer Pattern
- GUI observes simulation state
- Callbacks for component addition/connection
- Metrics polling for updates

### Strategy Pattern
- Load balancing strategies
- Cache eviction policies
- Routing algorithms

### Composite Pattern
- Components can contain other components
- Sharded databases
- Load balancers with backends

### Factory Pattern
- Component creation based on type
- Visual component instantiation

## Concurrency Model

### Thread Safety
- Components use mutexes for metric updates
- Canvas uses RWMutex for component access
- Simulator processes requests concurrently

### Goroutine Usage
- Request processing (per request)
- Metrics updates (periodic ticker)
- Traffic generation (periodic ticker)
- GUI refresh (Fyne event loop)

## Performance Considerations

### Optimization Techniques
- Component-level metric caching
- Lazy visual updates
- Request batching
- Efficient connection lookup

### Scalability Limits
- Max components: ~100 (GUI performance)
- Max requests/sec: ~1000 (simulation accuracy)
- Tick rate: 100ms (balance between accuracy and performance)

## Extension Points

### Adding New Components
1. Implement `Component` interface
2. Add to `gui.ComponentType`
3. Update factory in game screen
4. Add visual representation
5. Update documentation

### Adding New Levels
1. Create `Level` struct in `game/level.go`
2. Define requirements and criteria
3. Add to `Levels` array
4. Balance difficulty

### Adding New Metrics
1. Extend `Metrics` struct
2. Update component metric collection
3. Add to aggregation logic
4. Update GUI display

## Testing Strategy

### Unit Tests
- Component behavior
- Metric calculations
- Scoring algorithm
- Network simulation

### Integration Tests
- Component interactions
- Request flow through graph
- Metric aggregation

### Manual Testing
- GUI interactions
- Level completion
- Performance under load

## Dependencies

### External Libraries
- **Fyne**: GUI framework
- **Go standard library**: Core functionality

### Why Fyne?
- Cross-platform (Windows, Mac, Linux)
- Pure Go (no C dependencies)
- Modern Material Design
- Active development
- Good documentation

## Future Architecture Improvements

1. **Plugin System**: Load components dynamically
2. **State Persistence**: Save/load architectures
3. **Network Layer**: Actual network simulation vs shortcuts
4. **Auto-scaling**: Automatic component scaling based on load
5. **Replay System**: Record and replay scenarios
6. **Multi-user**: Collaborative architecture design
