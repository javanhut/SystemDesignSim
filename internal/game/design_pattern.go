package game

import (
	"time"

	"fyne.io/fyne/v2"
)

type DesignPattern struct {
	ID          string
	Name        string
	Category    string
	Description string
	Difficulty  int

	Problem   string
	Solution  string
	Benefits  []string
	Tradeoffs []string
	RealWorld []string

	DemoSteps     []TutorialStep
	PracticeSteps []PracticeStep

	Requirements PatternRequirements
}

type StepType string

const (
	StepAddComponent      StepType = "add_component"
	StepCreateConnection  StepType = "create_connection"
	StepShowTraffic       StepType = "show_traffic"
	StepHighlight         StepType = "highlight"
	StepMessage           StepType = "message"
	StepWait              StepType = "wait"
)

type TutorialStep struct {
	Order       int
	Type        StepType
	Title       string
	Description string
	Duration    time.Duration

	ComponentType string
	ComponentID   string
	Position      fyne.Position

	FromID string
	ToID   string

	FadeIn        bool
	ShowParticles bool

	ParticleCount int
}

type PracticeStep struct {
	Order       int
	Instruction string
	Hint        string
	Expected    StepValidation
}

type StepValidation struct {
	RequiredComponents  map[string]int
	RequiredConnections []ConnectionPair
	MinComponents       int
	CustomValidator     func(components map[string]interface{}) (bool, string)
}

type ConnectionPair struct {
	FromType string
	ToType   string
}

type PatternRequirements struct {
	MinComponents  int
	RequiredTypes  []string
	MustHaveConnections bool
}

var DesignPatterns = map[string]*DesignPattern{
	"load-balancing": LoadBalancingPattern(),
	"cache-aside": CacheAsidePattern(),
	"read-replicas": ReadReplicasPattern(),
}

func LoadBalancingPattern() *DesignPattern {
	return &DesignPattern{
		ID:          "load-balancing",
		Name:        "Load Balancing",
		Category:    "Scalability",
		Description: "Distribute traffic across multiple servers for horizontal scaling and high availability",
		Difficulty:  2,

		Problem: "A single server becomes a bottleneck as traffic increases. If it fails, the entire system goes down.",

		Solution: "Place a load balancer in front of multiple servers. The load balancer distributes incoming requests across all healthy servers using algorithms like round-robin or least-connections.",

		Benefits: []string{
			"Horizontal scalability - add more servers to handle more traffic",
			"High availability - system survives individual server failures",
			"Better resource utilization - distribute load evenly",
			"Zero-downtime deployments - update servers one at a time",
		},

		Tradeoffs: []string{
			"Additional component cost and complexity",
			"Load balancer itself can become a bottleneck (use multiple)",
			"Session affinity may require sticky sessions",
			"Increased latency (small, ~2ms overhead)",
		},

		RealWorld: []string{
			"Netflix - Zuul load balancer routes requests to microservices",
			"AWS ELB - Elastic Load Balancing for millions of requests",
			"NGINX - Used by 400M+ websites for load balancing",
			"HAProxy - Powers GitHub, Stack Overflow traffic distribution",
		},

		DemoSteps: []TutorialStep{
			{
				Order:       1,
				Type:        StepMessage,
				Title:       "Welcome to Load Balancing",
				Description: "Load balancing is the foundation of horizontal scalability.\n\nWatch as we build a load balanced architecture that can handle increased traffic and survive server failures.",
				Duration:    3 * time.Second,
			},
			{
				Order:         2,
				Type:          StepMessage,
				Title:         "The Problem",
				Description:   "A single API server handling all traffic:\n\n- Limited capacity (max ~1000 req/sec)\n- Single point of failure\n- Can't handle traffic spikes\n- No way to update without downtime",
				Duration:      3 * time.Second,
			},
			{
				Order:         3,
				Type:          StepAddComponent,
				Title:         "Adding Load Balancer",
				Description:   "The load balancer will distribute traffic across multiple servers",
				ComponentType: "load-balancer",
				ComponentID:   "lb-1",
				Position:      fyne.NewPos(300, 200),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         4,
				Type:          StepAddComponent,
				Title:         "Adding API Server #1",
				Description:   "First backend server to handle requests",
				ComponentType: "api-server",
				ComponentID:   "api-1",
				Position:      fyne.NewPos(500, 150),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         5,
				Type:          StepCreateConnection,
				Title:         "Connecting LB to API #1",
				Description:   "Load balancer routes traffic to the first server",
				FromID:        "lb-1",
				ToID:          "api-1",
				ShowParticles: true,
				Duration:      1 * time.Second,
			},
			{
				Order:         6,
				Type:          StepAddComponent,
				Title:         "Adding API Server #2",
				Description:   "Second server for redundancy and capacity",
				ComponentType: "api-server",
				ComponentID:   "api-2",
				Position:      fyne.NewPos(500, 250),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         7,
				Type:          StepCreateConnection,
				Title:         "Connecting LB to API #2",
				Description:   "Now traffic can be distributed across both servers",
				FromID:        "lb-1",
				ToID:          "api-2",
				ShowParticles: true,
				Duration:      1 * time.Second,
			},
			{
				Order:         8,
				Type:          StepShowTraffic,
				Title:         "Round-Robin Distribution",
				Description:   "Watch traffic alternate between servers:\nRequest 1 → API #1\nRequest 2 → API #2\nRequest 3 → API #1\n\nEven distribution = balanced load",
				ParticleCount: 10,
				Duration:      4 * time.Second,
			},
			{
				Order:         9,
				Type:          StepAddComponent,
				Title:         "Scaling Up: API Server #3",
				Description:   "Need more capacity? Just add another server!",
				ComponentType: "api-server",
				ComponentID:   "api-3",
				Position:      fyne.NewPos(500, 350),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         10,
				Type:          StepCreateConnection,
				Title:         "Connecting API #3",
				Description:   "Load balancer automatically includes new server",
				FromID:        "lb-1",
				ToID:          "api-3",
				ShowParticles: true,
				Duration:      1 * time.Second,
			},
			{
				Order:         11,
				Type:          StepShowTraffic,
				Title:         "3-Way Distribution",
				Description:   "Now traffic is split across 3 servers:\n\n- 3x capacity\n- Can survive 2 server failures\n- Each server handles ~33% of load\n\nThis is horizontal scaling!",
				ParticleCount: 15,
				Duration:      4 * time.Second,
			},
			{
				Order:       12,
				Type:        StepMessage,
				Title:       "Benefits Recap",
				Description: "Load Balancing enables:\n\n✓ Horizontal scaling (add more servers)\n✓ High availability (survive failures)\n✓ Better resource utilization\n✓ Zero-downtime deployments\n\nUsed by: Netflix, AWS, GitHub, Stack Overflow",
				Duration:    4 * time.Second,
			},
		},

		PracticeSteps: []PracticeStep{
			{
				Order:       1,
				Instruction: "Add a Load Balancer to distribute traffic",
				Hint:        "Click the 'Load Balancer' button in the toolbox",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"load-balancer": 1,
					},
				},
			},
			{
				Order:       2,
				Instruction: "Add at least 2 API Servers",
				Hint:        "Click the 'API Server' button twice to add two servers",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"load-balancer": 1,
						"api-server":    2,
					},
				},
			},
			{
				Order:       3,
				Instruction: "Connect the Load Balancer to both API Servers",
				Hint:        "Double-click the Load Balancer, then click each API Server to create connections",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"load-balancer": 1,
						"api-server":    2,
					},
					RequiredConnections: []ConnectionPair{
						{FromType: "load-balancer", ToType: "api-server"},
					},
					MinComponents: 3,
				},
			},
		},

		Requirements: PatternRequirements{
			MinComponents:       3,
			RequiredTypes:       []string{"load-balancer", "api-server"},
			MustHaveConnections: true,
		},
	}
}

func CacheAsidePattern() *DesignPattern {
	return &DesignPattern{
		ID:          "cache-aside",
		Name:        "Cache-Aside (Lazy Loading)",
		Category:    "Performance",
		Description: "Cache frequently accessed data in memory to reduce database load and latency",
		Difficulty:  2,

		Problem: "Database queries are slow (10-50ms) and expensive. Repeated queries for the same data waste resources and increase latency.",

		Solution: "Add a cache layer between the application and database. On reads, check the cache first. On cache miss, query the database and populate the cache. Writes go directly to the database and invalidate the cache.",

		Benefits: []string{
			"10-50x faster reads (1-2ms vs 10-50ms)",
			"Reduced database load (80-90% fewer queries)",
			"Lower costs (cache cheaper than database)",
			"Better scalability under read-heavy workloads",
		},

		Tradeoffs: []string{
			"Cache invalidation complexity",
			"Stale data risk (cache and DB out of sync)",
			"Additional infrastructure cost",
			"Cold cache problem on startup",
		},

		RealWorld: []string{
			"Facebook - Memcached caching layer (billions of reads/sec)",
			"Twitter - Redis for timeline caching",
			"Amazon - ElastiCache for product catalog",
			"Stack Overflow - Redis for user sessions and view counts",
		},

		DemoSteps: []TutorialStep{
			{
				Order:       1,
				Type:        StepMessage,
				Title:       "Cache-Aside Pattern",
				Description: "Learn how caching can reduce database load by 80-90% and speed up response times by 10-50x.\n\nWatch as we build a cache-aside architecture.",
				Duration:    3 * time.Second,
			},
			{
				Order:       2,
				Type:        StepMessage,
				Title:       "The Problem",
				Description: "Without caching:\n\n- Every request hits the database (10-50ms)\n- Same data queried repeatedly\n- Database becomes bottleneck\n- High latency, high costs",
				Duration:    3 * time.Second,
			},
			{
				Order:         3,
				Type:          StepAddComponent,
				Title:         "Adding API Server",
				Description:   "Application server that handles requests",
				ComponentType: "api-server",
				ComponentID:   "api-1",
				Position:      fyne.NewPos(200, 200),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         4,
				Type:          StepAddComponent,
				Title:         "Adding Cache (Redis)",
				Description:   "In-memory cache for fast data access",
				ComponentType: "cache",
				ComponentID:   "cache-1",
				Position:      fyne.NewPos(400, 200),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         5,
				Type:          StepAddComponent,
				Title:         "Adding Database",
				Description:   "Persistent storage (slower but durable)",
				ComponentType: "database",
				ComponentID:   "db-1",
				Position:      fyne.NewPos(600, 200),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         6,
				Type:          StepCreateConnection,
				Title:         "Connecting API to Cache",
				Description:   "API checks cache first",
				FromID:        "api-1",
				ToID:          "cache-1",
				ShowParticles: false,
				Duration:      1 * time.Second,
			},
			{
				Order:         7,
				Type:          StepCreateConnection,
				Title:         "Connecting Cache to Database",
				Description:   "Cache queries database on miss",
				FromID:        "cache-1",
				ToID:          "db-1",
				ShowParticles: false,
				Duration:      1 * time.Second,
			},
			{
				Order:         8,
				Type:          StepShowTraffic,
				Title:         "Cache Miss Flow",
				Description:   "First request (cache empty):\n\n1. API → Cache (check)\n2. Cache → Database (miss, query DB)\n3. Database → Cache (store result)\n4. Cache → API (return data)\n\nLatency: ~15ms",
				ParticleCount: 3,
				Duration:      4 * time.Second,
			},
			{
				Order:         9,
				Type:          StepShowTraffic,
				Title:         "Cache Hit Flow",
				Description:   "Subsequent requests (data cached):\n\n1. API → Cache (check)\n2. Cache → API (return data)\n\nNo database query!\nLatency: ~2ms (7x faster)\n\n90% cache hit rate = 90% faster",
				ParticleCount: 8,
				Duration:      4 * time.Second,
			},
			{
				Order:       10,
				Type:        StepMessage,
				Title:       "Benefits Recap",
				Description: "Cache-Aside Pattern:\n\n✓ 10-50x faster reads\n✓ 80-90% fewer database queries\n✓ Lower costs and better scalability\n✓ Simple to implement\n\nUsed by: Facebook, Twitter, Amazon, Reddit",
				Duration:    4 * time.Second,
			},
		},

		PracticeSteps: []PracticeStep{
			{
				Order:       1,
				Instruction: "Add an API Server",
				Hint:        "Click the 'API Server' button in the toolbox",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
					},
				},
			},
			{
				Order:       2,
				Instruction: "Add a Cache component",
				Hint:        "Click the 'Cache' button to add in-memory caching",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
						"cache":      1,
					},
				},
			},
			{
				Order:       3,
				Instruction: "Add a Database",
				Hint:        "Click the 'Database' button for persistent storage",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
						"cache":      1,
						"database":   1,
					},
				},
			},
			{
				Order:       4,
				Instruction: "Create the cache-aside flow: API → Cache → Database",
				Hint:        "Double-click API Server, then click Cache. Then double-click Cache and click Database",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
						"cache":      1,
						"database":   1,
					},
					RequiredConnections: []ConnectionPair{
						{FromType: "api-server", ToType: "cache"},
						{FromType: "cache", ToType: "database"},
					},
				},
			},
		},

		Requirements: PatternRequirements{
			MinComponents:       3,
			RequiredTypes:       []string{"api-server", "cache", "database"},
			MustHaveConnections: true,
		},
	}
}

func ReadReplicasPattern() *DesignPattern {
	return &DesignPattern{
		ID:          "read-replicas",
		Name:        "Database Read Replicas",
		Category:    "Scalability",
		Description: "Scale read-heavy workloads by replicating data to multiple read-only database instances",
		Difficulty:  3,

		Problem: "A single database can't handle high read volumes. Most applications are read-heavy (80-95% reads), so the database becomes a bottleneck.",

		Solution: "Create read-only replicas of the primary database. All writes go to the primary, which asynchronously replicates to replicas. Reads can be distributed across all replicas, multiplying read capacity.",

		Benefits: []string{
			"Multiply read capacity (N replicas = Nx reads)",
			"Geographic distribution (place replicas near users)",
			"Improved availability (reads survive primary failure)",
			"Offload reporting queries to replicas",
		},

		Tradeoffs: []string{
			"Replication lag (replicas slightly behind primary)",
			"Eventual consistency (not strong consistency)",
			"Increased infrastructure costs",
			"Application must route reads vs writes correctly",
		},

		RealWorld: []string{
			"Instagram - Postgres read replicas for photo metadata",
			"GitHub - MySQL read replicas for repositories and issues",
			"Shopify - Read replicas in every region",
			"Pinterest - Hundreds of replicas for recommendation engine",
		},

		DemoSteps: []TutorialStep{
			{
				Order:       1,
				Type:        StepMessage,
				Title:       "Database Read Replicas",
				Description: "Learn how read replicas multiply your database read capacity.\n\nWatch as we build a scalable database architecture for read-heavy workloads.",
				Duration:    3 * time.Second,
			},
			{
				Order:       2,
				Type:        StepMessage,
				Title:       "The Problem",
				Description: "Read-heavy workload (90% reads, 10% writes):\n\n- Single database bottleneck\n- Can't handle traffic spikes\n- High latency for distant users\n- Limited to ~5000 reads/sec",
				Duration:    3 * time.Second,
			},
			{
				Order:         3,
				Type:          StepAddComponent,
				Title:         "Adding API Server",
				Description:   "Application server handling read/write requests",
				ComponentType: "api-server",
				ComponentID:   "api-1",
				Position:      fyne.NewPos(200, 250),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         4,
				Type:          StepAddComponent,
				Title:         "Adding Primary Database",
				Description:   "Primary database handles all writes and some reads",
				ComponentType: "database",
				ComponentID:   "db-primary",
				Position:      fyne.NewPos(450, 150),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         5,
				Type:          StepCreateConnection,
				Title:         "Connecting API to Primary",
				Description:   "All writes must go to the primary database",
				FromID:        "api-1",
				ToID:          "db-primary",
				ShowParticles: true,
				Duration:      1 * time.Second,
			},
			{
				Order:         6,
				Type:          StepAddComponent,
				Title:         "Adding Read Replica #1",
				Description:   "Read-only copy of the primary database",
				ComponentType: "database",
				ComponentID:   "db-replica-1",
				Position:      fyne.NewPos(450, 280),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         7,
				Type:          StepCreateConnection,
				Title:         "Replication Connection",
				Description:   "Primary asynchronously replicates data to replica",
				FromID:        "db-primary",
				ToID:          "db-replica-1",
				ShowParticles: false,
				Duration:      1 * time.Second,
			},
			{
				Order:         8,
				Type:          StepShowTraffic,
				Title:         "Write Traffic Flow",
				Description:   "Writes always go to primary:\n\nAPI → Primary DB\n\nPrimary replicates to replica in background\n\n(Replication lag: 10-100ms typical)",
				ParticleCount: 2,
				Duration:      4 * time.Second,
			},
			{
				Order:         9,
				Type:          StepCreateConnection,
				Title:         "Connecting API to Replica",
				Description:   "API can now read from replica",
				FromID:        "api-1",
				ToID:          "db-replica-1",
				ShowParticles: false,
				Duration:      1 * time.Second,
			},
			{
				Order:         10,
				Type:          StepShowTraffic,
				Title:         "Read Distribution",
				Description:   "Reads can go to either database:\n\n- Read 1 → Primary\n- Read 2 → Replica\n- Read 3 → Primary\n\n2x read capacity!",
				ParticleCount: 8,
				Duration:      4 * time.Second,
			},
			{
				Order:         11,
				Type:          StepAddComponent,
				Title:         "Adding Read Replica #2",
				Description:   "Scale reads further with another replica",
				ComponentType: "database",
				ComponentID:   "db-replica-2",
				Position:      fyne.NewPos(450, 410),
				FadeIn:        true,
				Duration:      800 * time.Millisecond,
			},
			{
				Order:         12,
				Type:          StepCreateConnection,
				Title:         "Replicating to Replica #2",
				Description:   "Primary replicates to both replicas",
				FromID:        "db-primary",
				ToID:          "db-replica-2",
				ShowParticles: false,
				Duration:      1 * time.Second,
			},
			{
				Order:         13,
				Type:          StepCreateConnection,
				Title:         "API to Replica #2",
				Description:   "Now 3 databases available for reads",
				FromID:        "api-1",
				ToID:          "db-replica-2",
				ShowParticles: false,
				Duration:      1 * time.Second,
			},
			{
				Order:         14,
				Type:          StepShowTraffic,
				Title:         "3-Way Read Distribution",
				Description:   "Reads distributed across 3 databases:\n\n- Primary: 33% of reads\n- Replica 1: 33% of reads\n- Replica 2: 33% of reads\n\n3x read capacity!\nCan handle 15,000 reads/sec",
				ParticleCount: 12,
				Duration:      5 * time.Second,
			},
			{
				Order:       15,
				Type:        StepMessage,
				Title:       "Benefits Recap",
				Description: "Read Replicas Pattern:\n\n✓ Nx read capacity (N replicas)\n✓ Geographic distribution\n✓ Improved availability\n✓ Offload analytics to replicas\n\nTrade-off: Eventual consistency\n\nUsed by: Instagram, GitHub, Shopify",
				Duration:    4 * time.Second,
			},
		},

		PracticeSteps: []PracticeStep{
			{
				Order:       1,
				Instruction: "Add an API Server and a Primary Database",
				Hint:        "Add both API Server and Database components",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
						"database":   1,
					},
				},
			},
			{
				Order:       2,
				Instruction: "Add at least 2 Read Replica databases",
				Hint:        "Click 'Database' button 2 more times to create replicas",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
						"database":   3,
					},
				},
			},
			{
				Order:       3,
				Instruction: "Connect: API → Primary, Primary → Replicas, API → Replicas",
				Hint:        "Create replication connections from primary to replicas, and read connections from API to all databases",
				Expected: StepValidation{
					RequiredComponents: map[string]int{
						"api-server": 1,
						"database":   3,
					},
					MinComponents: 4,
				},
			},
		},

		Requirements: PatternRequirements{
			MinComponents:       4,
			RequiredTypes:       []string{"api-server", "database"},
			MustHaveConnections: true,
		},
	}
}

func GetAllPatterns() []*DesignPattern {
	return []*DesignPattern{
		LoadBalancingPattern(),
		CacheAsidePattern(),
		ReadReplicasPattern(),
	}
}

func GetPatternByID(id string) *DesignPattern {
	return DesignPatterns[id]
}

func GetPatternsByCategory(category string) []*DesignPattern {
	var patterns []*DesignPattern
	for _, pattern := range DesignPatterns {
		if pattern.Category == category {
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}
