# Phase 3 & 4: Networking and Deployment

## Status: COMPLETED (Core Infrastructure)

Phases 3 and 4 add comprehensive networking and code deployment capabilities to the System Design Simulator, enabling players to configure VPCs, security groups, DNS, and deploy applications with realistic deployment strategies.

## Phase 3: Regional Deployment & Networking

### Overview
Phase 3 transforms the simulator from simple component placement into a comprehensive networking simulation with VPCs, subnets, security groups, and DNS configuration.

### Files Created

#### 1. `internal/network/vpc.go` (500+ lines)

**VPC Management:**
- Create VPCs with custom CIDR blocks
- Subnet management (public, private, database tiers)
- Internet Gateway and NAT Gateway
- Route tables and routing configuration
- VPC peering for cross-VPC communication

**Key Features:**
```go
// Create a VPC
vpc, _ := network.NewVPC("vpc-1", "Production VPC", "us-east-1", "10.0.0.0/16")

// Add subnets
publicSubnet, _ := vpc.CreateSubnet("subnet-1", "Public AZ1", "10.0.1.0/24", "us-east-1a", network.SubnetTypePublic)
privateSubnet, _ := vpc.CreateSubnet("subnet-2", "Private AZ1", "10.0.11.0/24", "us-east-1a", network.SubnetTypePrivate)

// Attach Internet Gateway
igw := vpc.AttachInternetGateway("igw-1", "Main IGW")

// Create NAT Gateway
natGW, _ := vpc.CreateNATGateway("nat-1", "NAT Gateway", publicSubnet)

// Configure routing
publicSubnet.MakePublic(igw)
privateSubnet.MakePrivate(natGW)
```

**VPC Presets:**
- **Single AZ:** Basic setup with 1 public + 1 private subnet
- **Multi AZ:** High availability across 3 AZs (6 subnets)
- **Three Tier:** Web/App/Database tiers across 2 AZs (6 subnets)

**Cost Model:**
- NAT Gateway: $0.045/hour + $0.045/GB data processed
- VPC peering: $0.01/GB cross-region transfer
- Internet Gateway: Free

**Features:**
- CIDR validation and subnet calculations
- Available IP counting (excludes AWS reserved IPs)
- Route table management
- VPC peering for multi-VPC architectures

#### 2. `internal/network/security_groups.go` (400+ lines)

**Security Group Management:**
- Create security groups for VPCs
- Add ingress/egress rules
- Protocol support (TCP, UDP, ICMP, All)
- Port range configuration
- Source/destination CIDR or security group references

**Key Features:**
```go
// Create security group
sg := network.NewSecurityGroup("sg-web", "Web Server SG", "Allow HTTP/HTTPS", vpc)

// Add rules - helper methods
sg.AllowHTTP("0.0.0.0/0")        // Port 80 from anywhere
sg.AllowHTTPS("0.0.0.0/0")       // Port 443 from anywhere
sg.AllowSSH("10.0.0.0/16")       // Port 22 from VPC only

// Add custom rules
sg.AddIngressRule(network.ProtocolTCP, 8080, 8080, "10.0.0.0/16", "App port")

// Check if traffic is allowed
allowed := sg.IsTrafficAllowed(network.ProtocolTCP, 80, "1.2.3.4")
```

**Security Group Presets:**
- **Web Server:** HTTP (80), HTTPS (443) from anywhere
- **App Server:** Custom app ports (8080, 3000) from VPC
- **Database:** MySQL (3306), PostgreSQL (5432) from VPC
- **Cache:** Redis (6379), Memcached (11211) from VPC
- **Load Balancer:** HTTP/HTTPS from anywhere

**Network ACLs:**
- Stateless firewall rules at subnet level
- Rule numbers for priority
- Allow/deny actions
- Protocol and port-based filtering

**Features:**
- Protocol validation
- Port range parsing and validation
- Source IP/CIDR matching
- Helper methods for common services
- Preset security groups for quick setup

#### 3. `internal/network/dns.go` (400+ lines)

**DNS Management:**
- Hosted zones (public and private)
- DNS records (A, AAAA, CNAME, MX, TXT, etc.)
- Health checks for failover
- Routing policies (simple, weighted, latency, failover, geo)
- CDN distribution configuration

**Key Features:**
```go
// Create hosted zone
zone := network.NewHostedZone("Z123", "example.com", false, nil)

// Add DNS records
zone.AddARecord("www.example.com", "192.0.2.1", 300)
zone.AddCNAMERecord("app.example.com", "www.example.com", 300)

// CDN distribution
cdn := network.NewCDNDistribution("E123")
cdn.AddOrigin("origin-1", "example.com", "/")
cdn.AddAlias("www.example.com")

// Cost calculation
cost := cdn.CalculateCost(1000.0) // 1TB data transfer
```

**Routing Policies:**
- **Simple:** Single resource
- **Weighted:** Traffic distribution across resources
- **Latency:** Route to lowest latency endpoint
- **Failover:** Primary/secondary with health checks
- **Geolocation:** Route based on user location
- **Geoproximity:** Route based on resource and user location

**DNS Presets:**
- **Simple Web:** Basic A record + CNAME for www
- **CDN-Enabled:** CloudFront CNAME for static assets
- **Multi-Region:** Latency-based routing across regions

**CDN Features:**
- Origin configuration
- Edge location selection
- TTL and cache behavior
- Compression settings
- Cost tiers based on data transfer volume

**Cost Model:**
- Hosted zone: $0.50/month
- Queries: $0.40/million (first billion)
- CDN: $0.085/GB (first 10TB), tiered pricing

## Phase 4: Code Deployment

### Overview
Phase 4 adds comprehensive application deployment capabilities, allowing players to deploy code with different strategies, manage instances, configure auto-scaling, and monitor deployments.

### Files Created

#### 1. `internal/deployment/deployment.go` (500+ lines)

**Deployment Management:**
- Create deployments with version tracking
- Multiple deployment strategies
- Instance management and health monitoring
- Auto-scaling configuration
- Deployment pipelines

**Key Features:**
```go
// Create deployment
deployment := deployment.NewDeployment(
    "deploy-1", 
    "Production Deploy", 
    "my-app", 
    "v1.2.3",
    runtime,  // Node.js 18
    "us-east-1",
)

// Configure deployment strategy
deployment.Strategy = deployment.DeploymentStrategyRolling
deployment.Config.MinHealthyInstances = 2
deployment.Config.MaxBatchSize = 1
deployment.Config.RollbackOnFailure = true

// Add instances
instance1 := deployment.AddInstance("t3.medium", "us-east-1a")
instance2 := deployment.AddInstance("t3.medium", "us-east-1b")

// Start deployment
deployment.Start()

// Monitor health
healthyInstances := deployment.GetHealthyInstances()
```

**Deployment Strategies:**

1. **All At Once**
   - Deploy to all instances simultaneously
   - Fastest but causes downtime
   - Best for: Development environments

2. **Rolling**
   - Deploy batch by batch
   - Maintains minimum healthy instances
   - Zero-downtime deployment
   - Best for: Production with gradual rollout

3. **Blue/Green**
   - Deploy to new environment
   - Switch traffic instantly
   - Easy rollback
   - Best for: Critical applications

4. **Canary**
   - Deploy to small subset first
   - Gradually shift traffic
   - Monitor metrics before full rollout
   - Best for: Risk-averse deployments

**Instance Management:**
- Launch/stop/terminate instances
- Health status tracking
- Private/public IP assignment
- Version tagging

**Auto-Scaling:**
```go
autoScaling := &deployment.AutoScalingConfig{
    Enabled:            true,
    MinInstances:       2,
    MaxInstances:       10,
    DesiredInstances:   3,
    ScaleUpThreshold:   70.0,   // CPU %
    ScaleDownThreshold: 30.0,   // CPU %
    CooldownPeriod:     300 * time.Second,
}
```

**Health Checks:**
```go
healthCheck := &deployment.HealthCheckConfig{
    Protocol:           "HTTP",
    Port:               8080,
    Path:               "/health",
    IntervalSeconds:    30,
    TimeoutSeconds:     5,
    HealthyThreshold:   2,  // Consecutive successes
    UnhealthyThreshold: 3,  // Consecutive failures
}
```

**Code Source Configuration:**
- Git repositories (with branch/tag)
- S3 buckets
- Container registries
- ZIP files

**Deployment Pipelines:**
```go
pipeline := deployment.NewDeploymentPipeline("pipe-1", "CI/CD Pipeline")

// Add stages
pipeline.AddStage("Source", deployment.StageTypeSource)
pipeline.AddStage("Build", deployment.StageTypeBuild)
pipeline.AddStage("Test", deployment.StageTypeTest)
pipeline.AddStage("Deploy", deployment.StageTypeDeploy)

// Execute pipeline
pipeline.Execute()
```

**Deployment Presets:**

1. **Simple**
   - Single instance
   - All-at-once strategy
   - No auto-scaling
   - Best for: Development

2. **Rolling**
   - 2+ instances with auto-scaling
   - Rolling strategy
   - Min 1 healthy instance
   - Best for: Production

3. **Blue/Green**
   - 4+ instances
   - Blue/green strategy
   - Full environment duplication
   - Best for: Zero-downtime critical apps

4. **Canary**
   - 5+ instances
   - Canary strategy
   - Gradual traffic shifting
   - Best for: High-risk deployments

**Deployment Metrics:**
- Total instances
- Healthy/unhealthy counts
- CPU/memory usage
- Requests per second
- Error rate
- Deployment success rate

## Integration with Game

### How Networking Fits Into Gameplay

**Level 1: Sarah's Tech Blog**
- Simple VPC with single public subnet
- Basic security group (HTTP/HTTPS)
- Single A record DNS

**Level 2: Growing Blog**
- Multi-AZ VPC
- Public and private subnets
- Security groups for web and app tiers
- DNS with health checks

**Level 3: LocalConnect Social Network**
- Three-tier architecture (web/app/db)
- NAT Gateway for private subnet internet access
- Network ACLs for additional security
- Private hosted zone for internal DNS

**Level 4: GlobalGoods E-commerce**
- Multi-region VPCs
- VPC peering for cross-region communication
- Security groups for PCI compliance
- CDN with multiple origins
- Latency-based DNS routing

**Level 5: StreamNow Streaming**
- Global VPC architecture
- Advanced routing with traffic distribution
- CDN with edge locations worldwide
- DNS failover and health checks
- Complex security group mesh

### How Deployment Fits Into Gameplay

**Configuration Steps:**
1. Select deployment strategy
2. Configure auto-scaling parameters
3. Set health check endpoints
4. Choose code source (Git/S3/Container)
5. Set environment variables
6. Define deployment batch size

**Monitoring:**
- Real-time instance health
- Deployment progress
- Rollback on failure
- Cost tracking per deployment

**Scoring Impact:**
- Deployment strategy affects uptime
- Auto-scaling affects cost efficiency
- Health checks affect reliability
- Fast deployments earn bonus points

## Educational Value

### Networking Concepts Taught

1. **VPC Fundamentals**
   - CIDR notation and IP addressing
   - Public vs. private subnets
   - Internet Gateway vs. NAT Gateway
   - Route tables and routing

2. **Security**
   - Defense in depth (security groups + NACLs)
   - Least privilege principle
   - Stateful vs. stateless firewalls
   - Protocol and port management

3. **DNS and CDN**
   - DNS resolution and caching
   - Routing policies for high availability
   - CDN edge caching benefits
   - TTL impact on updates

4. **Multi-Region Architecture**
   - VPC peering
   - Cross-region latency
   - Data residency and compliance
   - Failover strategies

### Deployment Concepts Taught

1. **Deployment Strategies**
   - Tradeoffs: speed vs. safety vs. cost
   - Zero-downtime deployments
   - Rollback procedures
   - Blue/green benefits and costs

2. **Auto-Scaling**
   - Scale-up/scale-down thresholds
   - Cooldown periods
   - Cost optimization through right-sizing
   - Handling traffic spikes

3. **Health Monitoring**
   - Health check configuration
   - Healthy/unhealthy thresholds
   - Grace periods
   - Load balancer integration

4. **CI/CD Pipelines**
   - Source → Build → Test → Deploy flow
   - Automated testing gates
   - Approval stages
   - Rollback on test failure

## Technical Specifications

### Networking Data Structures

**VPC:**
- ID, Name, Region, CIDR
- Subnets collection
- Internet Gateway (optional)
- NAT Gateways (multiple)
- Route tables (multiple)
- Cost tracking

**Subnet:**
- ID, Name, CIDR, AZ
- Type (public/private/database)
- Auto-assign public IP (boolean)
- Route table reference
- Available IP count
- Security groups

**Security Group:**
- ID, Name, Description
- VPC reference
- Ingress rules collection
- Egress rules collection
- Rule validation logic

**DNS Record:**
- Name, Type, Value, TTL
- Routing policy
- Health check (optional)
- Region (for latency routing)

### Deployment Data Structures

**Deployment:**
- ID, Name, Version
- Application name
- Runtime configuration
- Strategy (rolling/blue-green/canary)
- Status (pending/in-progress/succeeded/failed)
- Instances collection
- Configuration object
- Timestamps

**Instance:**
- ID, Type, AZ
- Status (pending/running/stopped/terminated)
- Health status (healthy/unhealthy/unknown)
- IP addresses (private/public)
- Version deployed
- Launch time

**Auto-Scaling Config:**
- Min/Max/Desired instances
- Scale up/down thresholds
- Cooldown period
- Enabled flag

## Cost Model Integration

### Networking Costs

**VPC:**
- VPC itself: Free
- NAT Gateway: $0.045/hour + $0.045/GB
- VPC peering: $0.01/GB (cross-region)
- Data transfer out: $0.09/GB

**DNS:**
- Hosted zone: $0.50/month
- Queries: $0.40/million
- Health checks: $0.50/month each

**CDN:**
- Data transfer: $0.085/GB (first 10TB)
- HTTP/HTTPS requests: $0.0075/10,000
- Invalidation requests: $0.005/path

### Deployment Costs

**Compute:**
- Instance hours based on type
- Auto-scaling can increase/decrease costs
- Blue/green doubles infrastructure during deployment

**Data Transfer:**
- Deployment package downloads
- Health check traffic
- Cross-AZ data transfer: $0.01/GB

## Usage Examples

### Example 1: Simple Web Application

```go
// Create VPC
vpc, _ := network.CreateVPCFromPreset("vpc-1", "Web VPC", "us-east-1", "single-az", []string{"us-east-1a"})

// Security group
sg := network.CreateSecurityGroupFromPreset("sg-1", "Web SG", vpc, "web-server")

// DNS
zone := network.CreateHostedZoneFromPreset("Z1", "example.com", "simple-web", nil)

// Deployment
deployment := deployment.CreateDeploymentFromPreset(
    "deploy-1", "Web Deploy", "webapp", "v1.0.0", 
    runtime, "us-east-1", "simple",
)

deployment.AddInstance("t3.micro", "us-east-1a")
deployment.Start()
```

### Example 2: High Availability Setup

```go
// Multi-AZ VPC
vpc, _ := network.CreateVPCFromPreset("vpc-1", "HA VPC", "us-east-1", "multi-az", 
    []string{"us-east-1a", "us-east-1b", "us-east-1c"})

// Security groups for tiers
webSG := network.CreateSecurityGroupFromPreset("sg-web", "Web Tier", vpc, "web-server")
appSG := network.CreateSecurityGroupFromPreset("sg-app", "App Tier", vpc, "app-server")
dbSG := network.CreateSecurityGroupFromPreset("sg-db", "DB Tier", vpc, "database")

// DNS with health checks
zone := network.NewHostedZone("Z1", "example.com", false, nil)
record := zone.AddARecord("www.example.com", "192.0.2.1", 60)
record.HealthCheck = &network.HealthCheck{
    Protocol:        "HTTP",
    Port:            80,
    Path:            "/health",
    IntervalSeconds: 30,
}

// Rolling deployment with auto-scaling
deployment := deployment.CreateDeploymentFromPreset(
    "deploy-1", "HA Deploy", "webapp", "v2.0.0",
    runtime, "us-east-1", "rolling",
)

deployment.AddInstance("t3.small", "us-east-1a")
deployment.AddInstance("t3.small", "us-east-1b")
deployment.AddInstance("t3.small", "us-east-1c")
deployment.Start()
```

### Example 3: Global Multi-Region

```go
// VPCs in multiple regions
vpcUS := network.CreateVPCFromPreset("vpc-us", "US VPC", "us-east-1", "three-tier", 
    []string{"us-east-1a", "us-east-1b"})
vpcEU := network.CreateVPCFromPreset("vpc-eu", "EU VPC", "eu-west-1", "three-tier",
    []string{"eu-west-1a", "eu-west-1b"})

// VPC peering
peering := network.NewVPCPeering("pcx-1", "US-EU Peering", vpcUS, vpcEU)
peering.Accept()

// CDN with multiple origins
cdn := network.NewCDNDistribution("E123")
cdn.AddOrigin("us-origin", "us-lb.example.com", "/")
cdn.AddOrigin("eu-origin", "eu-lb.example.com", "/")

// Latency-based DNS
zone := network.NewHostedZone("Z1", "example.com", false, nil)
recordUS := zone.AddARecord("www.example.com", "192.0.2.1", 60)
recordUS.Region = "us-east-1"
recordEU := zone.AddARecord("www.example.com", "198.51.100.1", 60)
recordEU.Region = "eu-west-1"

// Blue/Green deployment in each region
deploymentUS := deployment.CreateDeploymentFromPreset(
    "deploy-us", "US Production", "webapp", "v3.0.0",
    runtime, "us-east-1", "blue-green",
)
deploymentEU := deployment.CreateDeploymentFromPreset(
    "deploy-eu", "EU Production", "webapp", "v3.0.0",
    runtime, "eu-west-1", "blue-green",
)
```

## Files Created

### Phase 3: Networking
1. **`internal/network/vpc.go`** (500 lines)
   - VPC, Subnet, Gateway management
   - Route tables and routing
   - VPC peering
   - 3 VPC presets

2. **`internal/network/security_groups.go`** (400 lines)
   - Security groups and rules
   - Network ACLs
   - Protocol and port management
   - 5 security group presets

3. **`internal/network/dns.go`** (400 lines)
   - Hosted zones and DNS records
   - Routing policies
   - CDN distributions
   - Health checks
   - 3 DNS presets

### Phase 4: Deployment
1. **`internal/deployment/deployment.go`** (500 lines)
   - Deployment management
   - Instance lifecycle
   - Auto-scaling configuration
   - Deployment strategies (4 types)
   - Deployment pipelines
   - 4 deployment presets

## Benefits

### Educational
- Real-world networking concepts
- Security best practices
- Multi-region architecture patterns
- Deployment strategy tradeoffs
- Cost optimization techniques

### Game Design
- Increased strategic depth
- More realistic scenarios
- Multiple solution paths
- Progressive complexity
- Cost/performance tradeoffs

### Technical Accuracy
- Based on AWS VPC architecture
- Real security group behavior
- Actual deployment strategies
- Realistic cost models
- Industry-standard practices

## Build Status

All code compiles successfully:

```bash
$ go build ./...
# Success - no errors
```

## Conclusion

Phases 3 and 4 add professional-grade networking and deployment capabilities to the System Design Simulator. Players now configure VPCs with proper subnetting, security groups, DNS routing, and deploy applications using industry-standard strategies like rolling deployments and blue/green deployments.

The networking layer teaches foundational cloud concepts like public/private subnets, NAT gateways, security groups, and DNS routing. The deployment layer introduces CI/CD concepts, auto-scaling, health checks, and deployment strategies.

Together, these phases transform the simulator from a simple component sandbox into a comprehensive cloud architecture and deployment platform that mirrors real-world cloud engineering practices.
