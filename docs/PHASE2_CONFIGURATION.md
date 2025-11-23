# Phase 2: Component Configuration System

## Status: COMPLETED (Core Infrastructure)

Phase 2 adds comprehensive configuration capabilities to all infrastructure components, allowing players to select instance types, regions, availability zones, runtimes, and other deployment parameters.

## What Was Implemented

### 1. Instance Type Definitions (`internal/components/config/instance_types.go`)

**Compute Instance Types (15 types)**
- **Burstable:** t2.micro, t2.small, t2.medium, t3.micro, t3.small, t3.medium
  - Cost: $0.0104 - $0.0464/hour
  - Max RPS: 10 - 50 requests/sec
  - Use case: Development, low traffic applications

- **General Purpose:** m5.large, m5.xlarge, m5.2xlarge, m5.4xlarge
  - Cost: $0.096 - $0.768/hour
  - Max RPS: 100 - 800 requests/sec
  - Use case: Production workloads, balanced performance

- **Compute Optimized:** c5.large, c5.xlarge, c5.2xlarge
  - Cost: $0.085 - $0.34/hour
  - Max RPS: 120 - 480 requests/sec
  - Use case: CPU-intensive applications

- **Memory Optimized:** r5.large, r5.xlarge, r5.2xlarge
  - Cost: $0.126 - $0.504/hour
  - Max RPS: 90 - 360 requests/sec
  - Use case: Memory-intensive applications, in-memory databases

**Database Instance Types (9 types)**
- db.t2.micro - $0.017/hour, 50 connections
- db.t2.small - $0.034/hour, 100 connections
- db.t3.small - $0.034/hour, 150 connections
- db.t3.medium - $0.068/hour, 300 connections
- db.m5.large - $0.186/hour, 500 connections
- db.m5.xlarge - $0.372/hour, 1000 connections
- db.m5.2xlarge - $0.744/hour, 2000 connections
- db.r5.large - $0.24/hour, 700 connections
- db.r5.xlarge - $0.48/hour, 1500 connections

Each includes: vCPU count, memory, storage, IOPS, max connections, cost per hour + per GB storage

**Cache Instance Types (6 types)**
- cache.t2.micro - 0.5GB, $0.017/hour
- cache.t3.small - 1.5GB, $0.034/hour
- cache.m5.large - 6.4GB, $0.126/hour
- cache.m5.xlarge - 12.9GB, $0.252/hour
- cache.r5.large - 13.1GB, $0.188/hour
- cache.r5.xlarge - 26.3GB, $0.376/hour

### 2. Region & Availability Zone Definitions (`internal/components/config/regions.go`)

**5 Global Regions**
1. **US East (N. Virginia)** - us-east-1
   - AZs: us-east-1a, us-east-1b, us-east-1c, us-east-1d
   - Primary market: North America

2. **US West (N. California)** - us-west-1
   - AZs: us-west-1a, us-west-1b, us-west-1c
   - Low latency to US West Coast

3. **EU West (Ireland)** - eu-west-1
   - AZs: eu-west-1a, eu-west-1b, eu-west-1c
   - GDPR compliant, European market

4. **Asia Pacific (Singapore)** - ap-southeast-1
   - AZs: ap-southeast-1a, ap-southeast-1b, ap-southeast-1c
   - Asian market

5. **Asia Pacific (Sydney)** - ap-southeast-2
   - AZs: ap-southeast-2a, ap-southeast-2b, ap-southeast-2c
   - Australia/Oceania market

**Network Latency Matrix**
- Same region: 5ms
- US East ↔ US West: 70ms
- US East ↔ Europe: 90ms
- US East ↔ Singapore: 150ms
- US East ↔ Sydney: 180ms
- US West ↔ Europe: 140ms
- US West ↔ Singapore: 120ms
- US West ↔ Sydney: 130ms
- Europe ↔ Singapore: 120ms
- Europe ↔ Sydney: 220ms
- Singapore ↔ Sydney: 90ms

### 3. Runtime Configurations (`internal/components/config/runtimes.go`)

**13 Runtime Options Across 8 Languages**

**JavaScript/Node.js:**
- Node.js 18 LTS - 128MB overhead, 1000ms startup
- Node.js 20 LTS - 128MB overhead, 900ms startup

**Python:**
- Python 3.9 - 256MB overhead, 1500ms startup
- Python 3.11 - 256MB overhead, 1300ms startup (25% faster)

**Go:**
- Go 1.20 - 64MB overhead, 500ms startup
- Go 1.21 - 64MB overhead, 450ms startup

**Java:**
- Java 17 LTS - 512MB overhead, 3000ms startup
- Java 21 LTS - 512MB overhead, 2500ms startup (virtual threads)

**.NET/C#:**
- .NET 6 - 384MB overhead, 2000ms startup
- .NET 8 - 384MB overhead, 1800ms startup

**Ruby:**
- Ruby 3.2 - 256MB overhead, 2000ms startup

**PHP:**
- PHP 8.2 - 128MB overhead, 800ms startup

**Rust:**
- Rust 1.73 - 32MB overhead, 300ms startup (fastest)

**Additional Runtime Features:**
- Default ports (3000, 8000, 8080, 5000, etc.)
- Supported port configurations
- Memory overhead calculations
- Startup time simulation

**Deployment Strategies:**
- Rolling deployment
- Blue/Green deployment
- Canary deployment
- Recreate deployment

**Auto-Scaling Configuration:**
- Min/Max instances
- Target CPU %
- Target Memory %
- Scale up/down cooldown periods

**Health Check Configuration:**
- Health check path (/health)
- Interval, timeout, thresholds
- Healthy/unhealthy counts

### 4. Property Panel Widget (`internal/gui/widgets/property_panel.go`)

Comprehensive configuration UI for all component types:

**API Server Configuration:**
- Instance type dropdown (15 options)
- Region selection (5 regions)
- Availability zone selection
- Runtime selection (13 runtimes)
- Real-time cost display
- Performance metrics (RPS, latency)

**Database Configuration:**
- Database instance type (9 options)
- Database engine (PostgreSQL, MySQL, MongoDB, DynamoDB, Redis)
- Region selection
- Multi-AZ deployment toggle (doubles cost, 99.95% uptime)
- Storage size configuration
- Cost breakdown (instance + storage)

**Cache Configuration:**
- Cache instance type (6 options)
- Cache engine (Redis, Memcached)
- Region selection
- TTL configuration
- Eviction policy (LRU, LFU, TTL)
- Real-time cost display

**Load Balancer Configuration:**
- LB type (Application, Network, Classic)
- Routing algorithm (Round Robin, Least Connections, IP Hash, Weighted)
- Region selection
- Health check path & interval
- Cost display ($0.0225/hour base + usage)

**CDN Configuration:**
- CDN provider (CloudFront, Fastly, Akamai, Cloudflare)
- Edge location selection (regional or global)
- Cache TTL configuration
- Compression settings (gzip/brotli)
- Cost display ($0.085/GB for first 10TB)

**Features:**
- Right-click component → Properties
- Dropdown selectors for all options
- Real-time cost calculation
- Informational hints and descriptions
- Save & Apply button
- Validation and error handling

## Integration Points

### How Property Panel Works

1. **User Right-Clicks Component**
   ```go
   // In game_screen.go
   component.OnRightClick = func() {
       widgets.ShowPropertyPanel(component, window, func() {
           // Update callback - recalculate costs, refresh UI
           updateMetrics()
       })
   }
   ```

2. **Property Panel Displays**
   - Modal overlay with scrollable content
   - Context-specific fields based on component type
   - Dropdown selectors populated from config

3. **User Configures Settings**
   - Select instance type → Cost updates
   - Select region → Latency calculations update
   - Select runtime → Startup time/memory overhead updates

4. **Save & Apply**
   - Configuration persisted to component
   - Cost recalculated
   - Performance characteristics updated
   - UI refreshed

### Cost Calculation Example

```go
// API Server Configuration
instanceType := config.GetInstanceType("m5.large")  // $0.096/hour
region := "us-east-1"
runtime := config.GetRuntime("nodejs-18")

// Total monthly cost
monthlyCost := instanceType.CostPerHour * 730  // $70.08/month

// Database Configuration
dbInstanceType := config.GetDatabaseInstanceType("db.m5.large")  // $0.186/hour
storageGB := 100
multiAZ := true

instanceCost := dbInstanceType.CostPerHour * 730  // $135.78/month
storageCost := storageGB * dbInstanceType.CostPerGBStorage  // $10/month
multiAZMultiplier := 2.0 if multiAZ else 1.0

totalDBCost := (instanceCost + storageCost) * multiAZMultiplier  // $291.56/month
```

## Usage Examples

### Example 1: Level 1 - Minimal Cost Setup

**Requirement:** Handle 10 concurrent users, $10/month budget

**Optimal Configuration:**
- API Server: t2.micro ($0.0116/hour = $8.47/month)
- Database: db.t2.micro ($0.017/hour = $12.41/month) + 20GB storage ($2)
- **Problem:** Over budget by $2.88

**Solution:** Use smaller database storage (10GB) = $11.41 total
- Total: $8.47 + $11.41 = $9.88/month ✓

### Example 2: Level 4 - Multi-Region E-commerce

**Requirement:** 10K users, GDPR compliance, multi-region

**Configuration:**
- **US-East:**
  - 3x m5.xlarge API servers ($0.192/hour each × 3 = $421.44/month)
  - 1x db.m5.xlarge + Multi-AZ ($0.372 × 2 × 730 = $543.12/month)
  - 1x cache.m5.large ($0.126/hour = $91.98/month)

- **EU-West:**
  - 3x m5.xlarge API servers ($421.44/month)
  - 1x db.m5.xlarge + Multi-AZ ($543.12/month)
  - 1x cache.m5.large ($91.98/month)

- **CDN:** CloudFront global ($50/month estimated)

**Total:** ~$2,163/month

### Example 3: Runtime Selection Impact

**Scenario:** Same workload, different runtimes

**Node.js 18:**
- Startup: 1000ms
- Memory overhead: 128MB
- Instance: t3.medium (4GB RAM) - fits comfortably
- Cost: $0.0416/hour = $30.37/month

**Java 17:**
- Startup: 3000ms
- Memory overhead: 512MB
- Instance: t3.medium (4GB RAM) - tight fit, may need t3.large
- Cost: $0.0832/hour = $60.74/month (2x Node.js cost!)

**Go 1.21:**
- Startup: 450ms
- Memory overhead: 64MB
- Instance: t3.small (2GB RAM) - plenty of room
- Cost: $0.0208/hour = $15.18/month (50% Node.js cost!)

**Lesson:** Runtime choice significantly impacts cost and performance

## Benefits

### Educational Value
1. **Real-world decision-making** - Players learn cost/performance tradeoffs
2. **Regional considerations** - Understanding latency vs. cost
3. **Right-sizing** - Not over-provisioning or under-provisioning
4. **Multi-AZ tradeoffs** - High availability costs money

### Game Design
1. **Increased complexity** - More strategic decisions
2. **Replayability** - Try different configurations
3. **Cost optimization** - Budget becomes real constraint
4. **Performance tuning** - Balance speed vs. cost

### Technical Accuracy
1. **Realistic pricing** - Based on actual cloud provider costs
2. **Accurate latencies** - Real-world network performance
3. **Proper instance types** - Actual AWS/cloud offerings
4. **Runtime characteristics** - True memory/startup differences

## Files Created

1. **`internal/components/config/instance_types.go`** (540 lines)
   - 15 compute instance types
   - 9 database instance types
   - 6 cache instance types
   - Cost calculations, performance metrics

2. **`internal/components/config/regions.go`** (200 lines)
   - 5 global regions with coordinates
   - 20+ availability zones
   - Network latency matrix
   - Region selection utilities

3. **`internal/components/config/runtimes.go`** (300 lines)
   - 13 runtime configurations
   - Deployment strategies
   - Auto-scaling configs
   - Health check configs

4. **`internal/gui/widgets/property_panel.go`** (500 lines)
   - Comprehensive property editor
   - Context-specific fields
   - Real-time cost display
   - Modal UI with Save/Cancel

## Next Steps: Phase 3

With configuration infrastructure complete, Phase 3 will add:

### Regional Deployment & Networking
- VPC and subnet management
- Security groups and firewall rules
- DNS configuration
- Multi-region deployment UI
- Network topology visualization

See full roadmap in project documentation.

## Testing

Build successful:
```bash
go build ./...  # Compiles without errors
go test ./...   # All tests pass
```

## Conclusion

Phase 2 transforms the game from a simple component drag-and-drop into a realistic cloud architecture simulator. Players must now make informed decisions about:

- **Instance sizing** - Balance cost vs. capacity
- **Regional deployment** - Consider latency and compliance
- **Runtime selection** - Understand performance implications
- **High availability** - Multi-AZ costs 2x but provides failover

The property panel provides an intuitive UI for these complex decisions, while the configuration system ensures accurate cost and performance modeling.

This foundation enables Phase 3's networking features and Phase 4's code deployment system.
