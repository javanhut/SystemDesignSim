# System Design Simulator

An interactive, GUI-based game that teaches system design concepts through hands-on simulation of distributed systems. Build, scale, and optimize real-world infrastructure while learning about load balancing, caching, databases, CDNs, and more.

## Features

### Interactive Visual Builder
- Drag-and-drop components onto a canvas
- Draw connections between components to create data flow
- Visual feedback showing component health (green, yellow, red)
- Real-time animated traffic flow visualization

### Comprehensive Infrastructure Components
- **API Servers**: Different sizes (small, medium, large, xlarge) with varying capacity
- **Databases**: SQL, NoSQL, key-value, and document stores with sharding and replication
- **Caching**: Redis/Memcached with LRU, LFU, FIFO eviction policies
- **Load Balancers**: Round-robin, least-connected, weighted-random strategies
- **CDN**: Multi-region edge caching with origin fallback
- **DNS**: Regional routing (coming soon)

### Real Simulation Engine
- Event-driven architecture processing thousands of requests
- Network latency simulation based on geographic regions
- Realistic performance metrics and bottleneck detection
- Cost tracking for infrastructure spending
- Fault injection and chaos engineering

### Progressive Game Levels
1. **Local Blog** - Handle 10 users with a simple setup
2. **Growing Blog** - Scale to 100 users with load balancing
3. **Regional Social Network** - 1,000 users with redundancy
4. **Global E-commerce** - 10,000 users across multiple regions
5. **Viral Streaming Service** - 100,000 users with five nines uptime

### Learning Objectives
- Horizontal vs vertical scaling
- Database sharding and replication
- Caching strategies and hit rates
- Regional distribution and CDN usage
- Load balancing algorithms
- Cost optimization
- Fault tolerance and redundancy
- Performance metrics (latency, throughput, error rates)

## Installation

### Prerequisites
- Go 1.21 or higher
- Graphics drivers supporting OpenGL 2.1+ (for GUI)

### Build from Source

```bash
git clone https://github.com/javanhut/systemdesignsim.git
cd systemdesignsim
go mod download
go build -o systemdesignsim cmd/simulator/main.go
./systemdesignsim
```

## How to Play

### Basic Controls
1. **Add Components**: Click buttons in the toolbox to add infrastructure components
2. **Connect Components**: Right-click on a component and drag to another to create connections
3. **Configure**: Click components to view and modify their properties
4. **Simulate**: Click "Start Simulation" to begin traffic simulation
5. **Monitor**: Watch real-time metrics in the right panel
6. **Submit**: When ready, click "Submit Solution" to see your score

### Connection Rules
- **Load Balancer → API Servers**: Distribute traffic across multiple backends
- **CDN → Origin Server**: Cache static content at the edge
- **API Server → Database**: Store and retrieve data
- **API Server → Cache**: Speed up reads with caching layer
- **Cache → Database**: Cache misses fall back to database

### Winning Strategy
Each level has specific requirements:
- **Latency**: Keep P99 latency below the threshold
- **Uptime**: Maintain minimum availability percentage
- **Error Rate**: Keep errors below maximum allowed
- **Budget**: Stay within the cost constraints
- **Architecture**: Use required components (load balancer, CDN, etc.)

Bonus points for exceeding targets and optimizing costs!

## Architecture Overview

```
SystemDesignSim/
├── cmd/
│   └── simulator/           # Main application entry
├── internal/
│   ├── engine/             # Core simulation engine
│   │   ├── types.go        # Request/Response types
│   │   └── simulator.go    # Event processing
│   ├── components/         # Infrastructure components
│   │   ├── api/           # API server implementation
│   │   ├── database/      # Database with sharding
│   │   ├── cache/         # Cache with eviction policies
│   │   ├── cdn/           # CDN with edge locations
│   │   └── loadbalancer/  # Load balancing strategies
│   ├── network/           # Network simulation (latency, bandwidth)
│   ├── game/              # Game logic and levels
│   │   ├── level.go       # Level definitions
│   │   └── game.go        # Scoring and validation
│   └── gui/               # GUI implementation
│       ├── visual_component.go  # Visual node representation
│       ├── canvas/              # Graph rendering
│       └── screens/             # UI screens
├── docs/                  # Documentation
└── scenarios/            # Level scenario definitions
```

## Technical Details

### Simulation Engine
The simulation uses an event-driven architecture:
- Requests are queued and processed asynchronously
- Each component implements the `Component` interface
- Metrics are collected and aggregated in real-time
- Time-based ticker simulates clock progression

### Component Interface
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

### Network Simulation
Regional latencies based on real-world measurements:
- us-east ↔ us-west: 70ms
- us-east ↔ europe: 80ms
- us-east ↔ asia: 180ms
- us-east ↔ australia: 200ms

### Metrics Tracked
- Request count, success/failure rates
- Average, P95, P99 latency
- Throughput (requests/second)
- Cache hit rate
- Cost per hour and total cost
- Data transferred

## Development

### Running Tests
```bash
go test ./...
```

### Contributing
This is an educational project. Contributions welcome!

Ideas for contributions:
- More component types (message queues, object storage, etc.)
- Additional levels and scenarios
- Better visualizations (graphs, charts)
- Tutorial system
- Save/load architectures
- Multiplayer challenges

## License

MIT License - See LICENSE file for details

## Learning Resources

This simulator teaches concepts from:
- Designing Data-Intensive Applications (Martin Kleppmann)
- System Design Interview (Alex Xu)
- AWS Well-Architected Framework
- Google SRE Books

## Roadmap

- [ ] DNS component with geolocation routing
- [ ] Message queue components (Kafka, RabbitMQ)
- [ ] Object storage (S3-like)
- [ ] Auto-scaling capabilities
- [ ] More sophisticated chaos engineering
- [ ] Tutorial mode with guided walkthroughs
- [ ] Architecture templates and patterns
- [ ] Multiplayer/competitive mode
- [ ] Mobile support

## Screenshots

(Screenshots will be added after first successful run)

## Credits

Created as an educational tool for learning distributed systems and system design principles.

Built with:
- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- Go - Backend simulation engine

## Support

For questions, issues, or suggestions:
- Open an issue on GitHub
- Check the docs/ folder for detailed documentation
