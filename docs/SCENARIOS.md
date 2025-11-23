# Enhanced Scenarios System

## Overview

The enhanced scenarios system provides detailed customer briefs with realistic business requirements for each level. This transforms the game from abstract technical challenges into real-world consulting scenarios.

## Phase 1: Enhanced Scenarios (COMPLETED)

### What We Built

#### 1. Scenario Data Structure (`internal/game/scenario.go`)

Each scenario includes:

**Customer Context:**
- Customer name and business type
- Current situation narrative
- Business requirements

**User Profile:**
- Concurrent users (initial and peak)
- Average session profile
- Peak times
- Geographic distribution

**Traffic Pattern:**
- Read/write/static request percentages
- Peak multiplier
- Daily traffic pattern (steady, business-hours, evening-peak)

**Data Requirements:**
- Initial data size
- Monthly growth projection
- Data types and counts
- Backup and retention requirements

**Technical Specifications:**
- Read/write ratio
- Growth projections
- Geographic spread
- Compliance needs (GDPR, PCI-DSS, etc.)

**Tasks:**
- Step-by-step implementation tasks
- Task types (infrastructure, networking, deployment, configuration, optimization)
- Mandatory vs. optional tasks
- Hints for each task

**Bonus Objectives:**
- Optional challenges for higher scores

#### 2. Requirement Tracking System (`internal/game/requirements.go`)

**RequirementTracker:**
- Tracks progress on all scenario tasks
- Records start/completion times
- Marks tasks as not started, in progress, completed, or skipped
- Generates detailed progress reports
- Validates mandatory task completion

**RequirementValidator:**
- Validates task completion against actual architecture
- Checks infrastructure requirements
- Validates network configuration
- Validates deployment setup
- Validates cost optimization

#### 3. Traffic Pattern Utilities (`internal/game/traffic_pattern.go`)

**TrafficGenerator:**
- Generates realistic traffic based on scenario patterns
- Applies daily traffic patterns (business hours, evening peak, steady)
- Calculates request types (read/write/static) based on percentages
- Simulates user sessions with realistic behavior

**LoadProjector:**
- Projects future load based on growth patterns
- Calculates time to reach peak capacity
- Provides monthly growth rate calculations

**GeographicDistributor:**
- Selects regions based on user distribution
- Identifies primary regions
- Calculates optimal deployment locations

**DataGrowthCalculator:**
- Calculates data storage requirements over time
- Projects backup storage needs
- Accounts for retention policies

**LatencyCalculator:**
- Calculates realistic network latency between regions
- Determines optimal regions for deployment
- Simulates cross-region communication delays

#### 4. Level Integration (`internal/game/level.go`)

- Added `Scenario` field to Level struct
- Each level now loads its corresponding scenario
- Scenarios provide rich context for technical requirements

#### 5. GUI Integration (`internal/gui/screens/game_screen.go`)

**Enhanced Help Screen:**
- Displays customer brief with business context
- Shows user profile and traffic patterns
- Lists constraints and compliance requirements
- Provides step-by-step task checklist with hints
- Shows bonus objectives

**Quick Tasks Panel:**
- Shows first 3 mandatory tasks in side panel
- Reminds player to check full help for details
- Provides at-a-glance task overview

## Scenario Levels

### Level 1: Sarah's Tech Blog

**Customer:** Personal technology blog
**Users:** 10 concurrent (peak)
**Budget:** $10/month

**Situation:**
Sarah is a software engineer who writes technical tutorials. She wants to move from a free platform to her own infrastructure to learn system design.

**Key Requirements:**
- 70% of users in US-East
- Evening peak traffic
- 85% reads, 15% writes
- 6MB initial data, 500KB/month growth

**Tasks:**
1. Select deployment region (US-East)
2. Deploy application server (t2.micro)
3. Configure database (PostgreSQL/MySQL)
4. Configure networking (HTTP/HTTPS)
5. Deploy application code (Node.js)
6. Test with traffic
7. Optimize for budget (<$10/month)

**Bonus:**
- Add caching
- Keep cost under $5
- Achieve 99% uptime
- P99 latency <200ms

### Level 2: Growing Blog

**Customer:** Growing technology blog
**Users:** 100 concurrent (peak)
**Budget:** $50/month

**Situation:**
Sarah's blog went viral after a Reddit/HN post. Traffic increased 10x overnight. Needs scaling without downtime.

**Key Requirements:**
- International audience (50% US-East, 20% US-West, 20% Europe)
- Business hours peak
- 88% reads, 12% writes
- 50MB data, 10MB/month growth

**Tasks:**
1. Add load balancer
2. Deploy multiple API servers (2-3)
3. Add caching layer (Redis)
4. Configure auto-scaling

**Bonus:**
- 70%+ cache hit rate
- 99.5% uptime
- Cost under $30
- P99 latency <150ms

### Level 3: LocalConnect Social Network

**Customer:** Regional social network
**Users:** 1,000 concurrent (peak)
**Budget:** $200/month

**Situation:**
Twitter-like platform for a metropolitan area. Real-time updates, high write volume, needs high availability.

**Key Requirements:**
- Single region (US-East)
- Lunch and evening peak times
- 65% reads, 30% writes
- 2GB data, 500MB/month growth

**Tasks:**
1. Implement database replication (master-replica)
2. Deploy multi-AZ infrastructure
3. Implement aggressive caching (80%+ hit rate)

**Bonus:**
- 99.9% uptime
- Database replication lag <5s
- 85%+ cache hit rate
- Cost under $150

### Level 4: GlobalGoods E-commerce

**Customer:** International e-commerce platform
**Users:** 10,000 concurrent (peak)
**Budget:** $1,000/month

**Situation:**
Launching online marketplace for artisan products. Worldwide users, GDPR compliance required, peak traffic during holidays (10x multiplier).

**Key Requirements:**
- Multi-region (US-East 35%, Europe 35%, Asia 10%, etc.)
- GDPR compliance (EU data stays in EU)
- 70% reads, 25% writes
- Holiday traffic spikes

**Tasks:**
1. Deploy multi-region infrastructure (US, EU, Asia)
2. Implement CDN for product images
3. Configure database sharding by region

**Bonus:**
- Deploy to 3+ regions
- 99.99% uptime
- P99 latency <75ms globally
- Cost under $750

### Level 5: StreamNow Streaming Service

**Customer:** Video streaming platform
**Users:** 100,000 concurrent (peak)
**Budget:** $5,000/month

**Situation:**
Competing with Netflix/YouTube. Viral marketing campaign successful. Need five nines uptime, minimal buffering, global distribution.

**Key Requirements:**
- Global distribution (5 regions)
- Prime time peak (7-11 PM all timezones)
- 95% reads (video streaming), 2% writes
- 2MB per request (video chunks)
- Five nines uptime required

**Tasks:**
1. Deploy global CDN (all 5 regions)
2. Implement database sharding (user_id as key)
3. Configure multi-layer caching (CDN + app + DB)

**Bonus:**
- 99.999% uptime (five nines)
- P99 latency <50ms globally
- 95%+ CDN cache hit rate
- Cost under $3,500

## How to Use Scenarios in Game

### For Players

1. **Click "? Help" button** on game screen
2. **Read customer brief** to understand business context
3. **Review task checklist** for step-by-step guidance
4. **Follow hints** for each task
5. **Complete bonus objectives** for higher scores

### For Developers

#### Access Scenario Data

```go
level := game.GetLevel(1)
scenario := level.Scenario

// Customer context
fmt.Println(scenario.CustomerName)
fmt.Println(scenario.CurrentSituation)

// User profile
users := scenario.UserProfile.PeakConcurrent
sessionDuration := scenario.UserProfile.AverageSession.DurationMinutes

// Traffic pattern
readPercent := scenario.TrafficPattern.ReadsPercentage
peakMultiplier := scenario.TrafficPattern.PeakMultiplier

// Tasks
for _, task := range scenario.Tasks {
    fmt.Printf("%d. %s\n", task.Step, task.Title)
    if task.Mandatory {
        fmt.Println("   [REQUIRED]")
    }
}
```

#### Track Requirements

```go
tracker := game.NewRequirementTracker(scenario)

// Start a task
tracker.StartTask(1)

// Complete a task
tracker.CompleteTask(1, "Deployed to us-east-1")

// Mark bonus completed
tracker.MarkBonusCompleted("Add caching")

// Check progress
completed, total := tracker.GetProgress()
fmt.Printf("%d/%d tasks completed\n", completed, total)

// Validate completion
validator := game.NewRequirementValidator(tracker)
valid, errors := validator.ValidateCompletion()
```

#### Generate Traffic

```go
generator := game.NewTrafficGenerator(
    scenario.TrafficPattern,
    scenario.UserProfile,
    baselineRPS,
)

// Calculate current RPS based on time of day
currentTime := time.Now()
rps := generator.CalculateCurrentRPS(currentTime)

// Get request type
requestType := generator.GetRequestType() // "read", "write", or "static"

// Simulate user session
session := generator.SimulateUserSession()
if session.ShouldMakeRequest() {
    session.RecordRequest()
}
```

#### Project Growth

```go
projector := game.NewLoadProjector(
    scenario.GrowthProjection,
    scenario.UserProfile.InitialConcurrent,
)

// Project load 6 months from now
futureDate := time.Now().Add(6 * 30 * 24 * time.Hour)
projectedLoad := projector.ProjectLoad(futureDate)

// Get time to reach peak
timeToReach := projector.GetTimeToReachPeak()
```

## Next Steps: Phase 2

The foundation is complete. Next phase will add:

### Component Configuration
- Instance size selection (t2.micro, t2.small, m5.large, etc.)
- Region and availability zone placement
- Runtime configuration (Node.js, Python, Go, etc.)
- Resource limits (CPU, memory, storage)

### Property Panels
- Right-click component → Edit properties
- Dropdown selectors for all options
- Cost calculation updates in real-time
- Visual feedback for configuration changes

See `IMPLEMENTATION_PLAN.md` for detailed Phase 2 roadmap.

## Benefits

### Educational Value
1. **Real business context** - Players understand WHY they're building systems, not just HOW
2. **Realistic constraints** - Budget, compliance, geography match real-world challenges
3. **Progressive complexity** - Scenarios grow from simple blog to global streaming platform
4. **Step-by-step guidance** - Tasks with hints teach best practices

### Game Design
1. **Clear objectives** - Task checklist removes ambiguity
2. **Replayability** - Bonus objectives encourage optimization
3. **Engagement** - Stories make technical challenges relatable
4. **Progression** - Each level builds on previous knowledge

### Technical Design
1. **Modular architecture** - Scenarios separate from game engine
2. **Extensible system** - Easy to add new scenarios
3. **Validation framework** - Automated checking of requirements
4. **Traffic simulation** - Realistic load patterns based on scenario

## Files Created/Modified

### Created
- `internal/game/scenario.go` (800+ lines) - Scenario definitions
- `internal/game/requirements.go` (400+ lines) - Requirement tracking
- `internal/game/traffic_pattern.go` (400+ lines) - Traffic utilities
- `docs/SCENARIOS.md` (this file)

### Modified
- `internal/game/level.go` - Added Scenario field to Level
- `internal/gui/screens/game_screen.go` - Enhanced help and task display

## Testing

The code compiles and builds successfully:

```bash
go test ./...     # All tests pass
go build ./...    # Builds successfully
```

## Conclusion

Phase 1 transforms System Design Simulator from a technical sandbox into an immersive consulting experience. Players now receive detailed customer briefs with realistic business requirements, making the learning experience more engaging and practical.

The requirement tracking and validation systems provide structure and guidance, while the traffic pattern utilities enable realistic simulation of real-world load patterns.

This foundation sets the stage for Phase 2's component configuration system, which will allow players to configure every aspect of their infrastructure.
