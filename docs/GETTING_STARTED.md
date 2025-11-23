# Getting Started with System Design Simulator

## Quick Start

### Step 1: Install Prerequisites

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get install golang gcc libgl1-mesa-dev xorg-dev
```

#### Linux (Fedora/RHEL)
```bash
sudo dnf install golang gcc libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel
```

#### macOS
```bash
brew install go
```
(Xcode command line tools should be installed)

#### Windows
- Install Go from https://golang.org/dl/
- Install GCC via TDM-GCC or MinGW-w64

### Step 2: Clone and Build

```bash
git clone <repository-url>
cd SystemDesignSim
go mod download
go build -o systemdesignsim cmd/simulator/main.go
```

Note: First build may take 5-10 minutes as Fyne compiles with CGo.

### Step 3: Run the Game

```bash
./systemdesignsim
```

## Your First Level

### Level 1: Local Blog

**Objective**: Handle 10 concurrent users

**Requirements**:
- Max latency: 500ms
- Min uptime: 95%
- Budget: $10

**Suggested Architecture**:
1. Add an API Server (click "API Server" button)
2. Add a Database (click "Database" button)
3. Connect them: Right-click API Server, drag to Database
4. Click "Start Simulation"
5. Watch the metrics
6. Click "Submit Solution" when ready

### Understanding the UI

**Header - Level Title and Scenario**:
- Displays level name and number
- Shows detailed scenario information including:
  - Client name and business type
  - Current situation and context
  - User profile (concurrent users, session duration)
  - Traffic pattern (reads/writes/static content ratio)
  - Constraints (budget, latency, uptime requirements)

**Left Panel - Toolbox**:
- Click buttons to add components to the canvas
- Each component type has different capabilities and costs
- Includes component descriptions and latency information
- Help button provides detailed scenario guidance

**Center - Canvas**:
- Left-click: Select components
- Double-click then target: Create connections between components
- Drag components to rearrange them
- Visual particle animations show traffic flow
- Color-coded health indicators (Green=Healthy, Yellow=Warning, Orange=Critical, Red=Down)

**Right Panel - Metrics**:
- Real-time performance stats (RPS, latency, uptime)
- Pass/fail indicators for each objective
- Cost tracking with budget status
- Level objectives checklist
- Architecture hints (scrollable)
- Architecture summary
- Test plan information

**Bottom - Controls**:
- Start/Stop simulation
- Submit your solution
- Show Hints (popup with detailed system design principles)
- Control Center (advanced configuration)
- System Plan (architecture overview)
- Return to level select

## Component Types Explained

### API Server
- Handles incoming requests
- Different sizes: Small (10 concurrent), Medium (50), Large (200), XLarge (500)
- Cost: $0.05 - $0.40/hour
- Can connect to: Database, Cache

### Database
- Stores persistent data
- Types: SQL, NoSQL, Key-Value, Document
- Supports sharding and replication
- Cost: $0.05/hour + capacity cost
- Latency: Read 10ms, Write 15ms

### Cache
- Speeds up read operations
- Eviction policies: LRU, LFU, FIFO
- High hit rates reduce database load
- Cost: $0.02/hour + capacity cost
- Latency: 1-2ms

### Load Balancer
- Distributes traffic across multiple servers
- Strategies: Round-robin, least-connected
- Enables horizontal scaling
- Cost: $0.025/hour
- Latency: 2ms overhead

### CDN
- Caches static content at edge locations
- Reduces latency for global users
- Requires origin server
- Cost: $0.08/hour + per-region
- Latency: 2ms (edge hit)

## Using Architectural Hints

The game provides comprehensive architectural hints to help you design better systems:

**Show Hints Button**:
- Click "Show Hints" in the controls panel to open the hints dialog
- Hints are also displayed in the right panel (scrollable)
- Hints dynamically update based on your current architecture

**What Hints Cover**:
1. Request Flow Pattern - Recommended architecture and component flow
2. Horizontal Scalability - How to scale out with load balancers and multiple servers
3. Latency Optimization - Cache strategies, CDN usage, performance targets
4. High Availability - Eliminating single points of failure, redundancy patterns
5. Traffic Handling - Load distribution, database scaling, connection pooling
6. Cost Optimization - Budget management, right-sizing, auto-scaling
7. Common Mistakes - Anti-patterns to avoid

**Real System Design Principles**:
- Hints explain WHY each component is needed, not just WHAT to add
- Includes real-world technologies (Redis, AWS ELB, CloudFront, etc.)
- Shows specific performance numbers (DB: 10-50ms, Cache: 1-2ms)
- Explains algorithms (Round-robin, Least-connections, LRU caching)
- Provides uptime calculations (99.9% = 3 nines = 8.76hr downtime/year)

**How to Use Hints Effectively**:
1. Read the scenario block at the top to understand the requirements
2. Click "Show Hints" to see recommended architecture
3. Add components based on the hints
4. Watch for checkmarks (✓) or crosses (✗) indicating missing components
5. Start simulation and monitor real-time pass/fail indicators
6. Iterate based on performance metrics

## Tips and Strategies

### Level 1-2: Start Simple
- Single API server + Database is enough
- Focus on meeting basic requirements
- Don't over-engineer

### Level 3: Add Redundancy
- Use load balancer with multiple API servers
- Add database replication for high availability
- Consider adding cache for read performance

### Level 4-5: Go Global
- Deploy CDN for edge caching
- Multi-region database setup
- Implement sharding for scale
- Optimize costs while meeting SLAs

### Cost Optimization
- Right-size your components
- Use caching aggressively
- Only use CDN when needed
- Balance performance vs cost

### Performance Optimization
- Add cache layers for read-heavy workloads
- Use load balancers to distribute load
- Replicate databases for read scalability
- Consider sharding for very large datasets

## Common Issues

### Simulation Won't Start
- Ensure you have at least one API server
- Check that components are connected properly
- Database must be connected to API server

### High Error Rates
- API servers may be overloaded (check capacity)
- Add more servers behind load balancer
- Check for failed components (red borders)

### High Latency
- Add caching layers
- Use CDN for static content
- Reduce cross-region hops
- Increase API server sizes

### Over Budget
- Use smaller instance sizes
- Remove unnecessary components
- Optimize cache hit rates to reduce database calls
- Consider cost vs performance tradeoffs

## Advanced Techniques

### Sharding Strategy
- Use consistent hashing for even distribution
- Balance shard sizes
- Monitor per-shard metrics

### Caching Strategy
- Set appropriate TTLs
- Use LRU for general purpose
- Monitor cache hit rates (target >70%)
- Layer caches (CDN → Application → Database)

### Load Balancing
- Use least-connected for variable workloads
- Round-robin for uniform requests
- Monitor backend health

## Keyboard Shortcuts

- `Ctrl+C`: Copy selected component
- `Delete`: Remove selected component
- `Space`: Play/Pause simulation
- `Esc`: Deselect all

## Next Steps

- Complete Level 1 with a perfect score
- Try different architectures for the same level
- Read the Architecture documentation in docs/
- Explore the source code to understand simulations
- Contribute new components or levels!

## Troubleshooting

### Build Issues

**"Cannot find package"**
```bash
go mod download
go mod tidy
```

**CGo/GCC errors**
- Install C compiler and development headers
- See prerequisites section above

**Fyne compilation slow**
- First compilation is slow (10+ minutes)
- Subsequent builds are faster
- Consider using `go build -tags debug` for faster debug builds

### Runtime Issues

**Window doesn't appear**
- Check graphics drivers are installed
- Ensure OpenGL 2.1+ support
- Try setting `FYNE_SCALE=1.0` environment variable

**Poor performance**
- Reduce simulation speed
- Limit number of components
- Close other applications

## Getting Help

- Check the README.md for overview
- Read docs/ARCHITECTURE.md for technical details
- Open an issue on GitHub
- Check Fyne documentation at https://fyne.io

## Contributing

We welcome contributions! Areas where you can help:
- New infrastructure components
- Additional game levels
- Performance optimizations
- Documentation improvements
- Bug fixes
- Visual enhancements

See CONTRIBUTING.md for guidelines.
