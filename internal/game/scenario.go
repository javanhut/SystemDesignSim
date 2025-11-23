package game

type Scenario struct {
	// Customer Context
	CustomerName     string
	BusinessType     string
	CurrentSituation string

	// Business Requirements
	UserProfile      UserProfile
	TrafficPattern   TrafficPattern
	DataRequirements DataRequirements
	ComplianceNeeds  []string

	// Technical Specifications
	ReadWriteRatio   ReadWriteRatio
	GrowthProjection GrowthProjection

	// Constraints
	DataResidency    []string
	GeographicSpread GeographicDistribution

	// Tasks
	Tasks           []Task
	BonusObjectives []string
}

type UserProfile struct {
	InitialConcurrent int
	PeakConcurrent    int
	AverageSession    SessionProfile
	PeakTimes         []string
	GeographicMix     map[string]float64 // Region -> percentage
}

type SessionProfile struct {
	DurationMinutes int
	PageViews       int
	DataPerRequest  int64 // bytes
}

type TrafficPattern struct {
	ReadsPercentage  float64
	WritesPercentage float64
	StaticPercentage float64
	PeakMultiplier   float64
	DailyPattern     string // "steady", "business-hours", "evening-peak"
}

type DataRequirements struct {
	InitialSize    int64 // bytes
	GrowthPerMonth int64
	Types          []DataType
	BackupRequired bool
	RetentionDays  int
}

type DataType struct {
	Name        string
	Count       int
	SizePerItem int64
	GrowthRate  string
}

type ReadWriteRatio struct {
	ReadPercent  float64
	WritePercent float64
}

type GrowthProjection struct {
	MonthlyGrowth float64 // e.g., 2.0 = 2x per month
	ExpectedPeak  int
	TimeToReach   int // months
}

type GeographicDistribution struct {
	Primary      string             // Primary region
	Distribution map[string]float64 // Region -> percentage
}

type Task struct {
	Step        int
	Title       string
	Description string
	Type        TaskType
	Mandatory   bool
	Hint        string
}

type TaskType string

const (
	TaskTypeInfrastructure TaskType = "infrastructure"
	TaskTypeNetworking     TaskType = "networking"
	TaskTypeDeployment     TaskType = "deployment"
	TaskTypeConfiguration  TaskType = "configuration"
	TaskTypeOptimization   TaskType = "optimization"
)

// Scenario templates for each level
func GetScenarioForLevel(levelID int) *Scenario {
	scenarios := map[int]*Scenario{
		1: createLocalBlogScenario(),
		2: createGrowingBlogScenario(),
		3: createSocialNetworkScenario(),
		4: createEcommerceScenario(),
		5: createStreamingScenario(),
	}

	return scenarios[levelID]
}

func createLocalBlogScenario() *Scenario {
	return &Scenario{
		CustomerName:     "Sarah's Tech Blog",
		BusinessType:     "Personal Technology Blog",
		CurrentSituation: "Sarah is a software engineer who writes technical tutorials. She currently hosts her blog on a free platform but wants to move to her own infrastructure for better control and to learn system design. Her blog focuses on web development tutorials and attracts primarily US-based developers.",

		UserProfile: UserProfile{
			InitialConcurrent: 5,
			PeakConcurrent:    10,
			AverageSession: SessionProfile{
				DurationMinutes: 8,
				PageViews:       4,
				DataPerRequest:  50 * 1024, // 50KB per page
			},
			PeakTimes: []string{"Tuesday evenings (8-10 PM EST)", "Weekend mornings"},
			GeographicMix: map[string]float64{
				"us-east": 0.70,
				"us-west": 0.20,
				"europe":  0.08,
				"asia":    0.02,
			},
		},

		TrafficPattern: TrafficPattern{
			ReadsPercentage:  0.75,
			WritesPercentage: 0.15,
			StaticPercentage: 0.10,
			PeakMultiplier:   2.0,
			DailyPattern:     "evening-peak",
		},

		DataRequirements: DataRequirements{
			InitialSize:    6 * 1024 * 1024, // 6MB
			GrowthPerMonth: 500 * 1024,      // 500KB/month
			Types: []DataType{
				{Name: "Blog Posts", Count: 50, SizePerItem: 5 * 1024, GrowthRate: "5 posts/month"},
				{Name: "Images", Count: 100, SizePerItem: 50 * 1024, GrowthRate: "10 images/month"},
				{Name: "Comments", Count: 100, SizePerItem: 1024, GrowthRate: "20 comments/month"},
			},
			BackupRequired: true,
			RetentionDays:  30,
		},

		ReadWriteRatio: ReadWriteRatio{
			ReadPercent:  85.0,
			WritePercent: 15.0,
		},

		GrowthProjection: GrowthProjection{
			MonthlyGrowth: 1.5, // 50% growth if content goes viral
			ExpectedPeak:  50,
			TimeToReach:   6,
		},

		DataResidency: []string{},

		GeographicSpread: GeographicDistribution{
			Primary: "us-east",
			Distribution: map[string]float64{
				"us-east": 0.90,
				"us-west": 0.10,
			},
		},

		Tasks: []Task{
			{
				Step:        1,
				Title:       "Select Deployment Region",
				Description: "Choose the AWS region closest to your users (primarily US East Coast)",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "70% of users are in US-East. Deploying there minimizes latency.",
			},
			{
				Step:        2,
				Title:       "Deploy Application Server",
				Description: "Set up a server to run the Node.js application (blog engine)",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Start with a small instance (t2.micro) to keep costs low. You can scale up later.",
			},
			{
				Step:        3,
				Title:       "Configure Database",
				Description: "Set up a database to store blog posts and comments permanently",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "PostgreSQL or MySQL work well for blogs. Start with 10GB storage.",
			},
			{
				Step:        4,
				Title:       "Configure Networking",
				Description: "Set up how users will reach your blog (DNS, security rules)",
				Type:        TaskTypeNetworking,
				Mandatory:   true,
				Hint:        "You'll need to allow HTTP/HTTPS traffic (ports 80/443) from the internet.",
			},
			{
				Step:        5,
				Title:       "Deploy Application Code",
				Description: "Deploy the blog application code to your server",
				Type:        TaskTypeDeployment,
				Mandatory:   true,
				Hint:        "Configure Node.js 18 runtime and set environment variables for database connection.",
			},
			{
				Step:        6,
				Title:       "Test with Traffic",
				Description: "Run the simulation to see if your architecture handles the load",
				Type:        TaskTypeConfiguration,
				Mandatory:   true,
				Hint:        "Watch for high latency or errors. Green components = healthy, Red = overloaded.",
			},
			{
				Step:        7,
				Title:       "Optimize for Budget",
				Description: "Ensure total cost stays under $10/month while meeting performance requirements",
				Type:        TaskTypeOptimization,
				Mandatory:   true,
				Hint:        "t2.micro + small database should cost ~$8-9/month. Avoid over-provisioning.",
			},
		},

		BonusObjectives: []string{
			"Add caching to improve read performance",
			"Keep total monthly cost under $5",
			"Achieve 99% uptime (instead of required 95%)",
			"Keep P99 latency under 200ms (instead of 500ms)",
		},
	}
}

func createGrowingBlogScenario() *Scenario {
	return &Scenario{
		CustomerName:     "Sarah's Tech Blog (Expanded)",
		BusinessType:     "Growing Technology Blog",
		CurrentSituation: "Sarah's blog went viral after a popular tutorial was shared on Reddit and Hacker News. Traffic has increased 10x overnight. She needs to scale her infrastructure to handle 100 concurrent users without downtime. The blog now attracts international readers and needs better performance globally.",

		UserProfile: UserProfile{
			InitialConcurrent: 50,
			PeakConcurrent:    100,
			AverageSession: SessionProfile{
				DurationMinutes: 10,
				PageViews:       6,
				DataPerRequest:  50 * 1024,
			},
			PeakTimes: []string{"Business hours US/EU (9 AM - 5 PM)", "Evening US (7-11 PM)"},
			GeographicMix: map[string]float64{
				"us-east":   0.50,
				"us-west":   0.20,
				"europe":    0.20,
				"asia":      0.08,
				"australia": 0.02,
			},
		},

		TrafficPattern: TrafficPattern{
			ReadsPercentage:  0.80,
			WritesPercentage: 0.12,
			StaticPercentage: 0.08,
			PeakMultiplier:   3.0,
			DailyPattern:     "business-hours",
		},

		DataRequirements: DataRequirements{
			InitialSize:    50 * 1024 * 1024, // 50MB
			GrowthPerMonth: 10 * 1024 * 1024, // 10MB/month
			Types: []DataType{
				{Name: "Blog Posts", Count: 200, SizePerItem: 8 * 1024, GrowthRate: "20 posts/month"},
				{Name: "Images", Count: 500, SizePerItem: 75 * 1024, GrowthRate: "50 images/month"},
				{Name: "Comments", Count: 2000, SizePerItem: 1024, GrowthRate: "200 comments/month"},
			},
			BackupRequired: true,
			RetentionDays:  90,
		},

		ReadWriteRatio: ReadWriteRatio{
			ReadPercent:  88.0,
			WritePercent: 12.0,
		},

		GrowthProjection: GrowthProjection{
			MonthlyGrowth: 1.3,
			ExpectedPeak:  300,
			TimeToReach:   12,
		},

		Tasks: []Task{
			{
				Step:        1,
				Title:       "Add Load Balancer",
				Description: "Deploy a load balancer to distribute traffic across multiple servers",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Load balancer enables horizontal scaling by routing requests to multiple API servers.",
			},
			{
				Step:        2,
				Title:       "Deploy Multiple API Servers",
				Description: "Add 2-3 API servers behind the load balancer for redundancy",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Multiple servers provide high availability. If one fails, others continue serving.",
			},
			{
				Step:        3,
				Title:       "Add Caching Layer",
				Description: "Implement Redis cache to reduce database load for popular posts",
				Type:        TaskTypeInfrastructure,
				Mandatory:   false,
				Hint:        "Cache hit rate of 70%+ significantly improves performance and reduces costs.",
			},
			{
				Step:        4,
				Title:       "Configure Auto-Scaling",
				Description: "Set up auto-scaling rules to handle traffic spikes automatically",
				Type:        TaskTypeConfiguration,
				Mandatory:   false,
				Hint:        "Scale up when CPU > 70%, scale down when CPU < 30% to optimize costs.",
			},
		},

		BonusObjectives: []string{
			"Implement caching with 70%+ hit rate",
			"Achieve 99.5% uptime",
			"Keep costs under $30/month",
			"P99 latency under 150ms",
		},
	}
}

func createSocialNetworkScenario() *Scenario {
	return &Scenario{
		CustomerName:     "LocalConnect",
		BusinessType:     "Regional Social Network",
		CurrentSituation: "LocalConnect is a Twitter-like social network for a metropolitan area with 1 million residents. Users post updates, photos, and interact with local businesses. The platform needs to handle 1,000 concurrent users during peak hours with real-time updates and high availability. Data consistency is important for user posts and interactions.",

		UserProfile: UserProfile{
			InitialConcurrent: 500,
			PeakConcurrent:    1000,
			AverageSession: SessionProfile{
				DurationMinutes: 25,
				PageViews:       30,
				DataPerRequest:  20 * 1024,
			},
			PeakTimes: []string{"Lunch (12-1 PM)", "Evening (6-10 PM)", "Weekend afternoons"},
			GeographicMix: map[string]float64{
				"us-east": 1.0, // Regional, single metro area
			},
		},

		TrafficPattern: TrafficPattern{
			ReadsPercentage:  0.65,
			WritesPercentage: 0.30,
			StaticPercentage: 0.05,
			PeakMultiplier:   4.0,
			DailyPattern:     "business-hours",
		},

		DataRequirements: DataRequirements{
			InitialSize:    2 * 1024 * 1024 * 1024, // 2GB
			GrowthPerMonth: 500 * 1024 * 1024,      // 500MB/month
			Types: []DataType{
				{Name: "User Posts", Count: 100000, SizePerItem: 2 * 1024, GrowthRate: "50K posts/month"},
				{Name: "User Profiles", Count: 50000, SizePerItem: 10 * 1024, GrowthRate: "5K users/month"},
				{Name: "Photos", Count: 20000, SizePerItem: 200 * 1024, GrowthRate: "10K photos/month"},
			},
			BackupRequired: true,
			RetentionDays:  365,
		},

		Tasks: []Task{
			{
				Step:        1,
				Title:       "Implement Database Replication",
				Description: "Set up master-replica database configuration for read scalability",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Write to master, read from replicas. This scales read-heavy workloads.",
			},
			{
				Step:        2,
				Title:       "Deploy Multi-AZ Infrastructure",
				Description: "Deploy across multiple availability zones for high availability",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Multi-AZ provides 99.9%+ uptime. If one zone fails, others continue.",
			},
			{
				Step:        3,
				Title:       "Implement Aggressive Caching",
				Description: "Cache user feeds, profiles, and popular content",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Social networks are read-heavy. 80%+ cache hit rate is achievable.",
			},
		},

		BonusObjectives: []string{
			"Achieve 99.9% uptime",
			"Implement database replication with <5s lag",
			"Cache hit rate >85%",
			"Keep costs under $150/month",
		},
	}
}

func createEcommerceScenario() *Scenario {
	return &Scenario{
		CustomerName:     "GlobalGoods",
		BusinessType:     "International E-commerce Platform",
		CurrentSituation: "GlobalGoods is launching an online marketplace selling artisan products worldwide. They expect 10,000 concurrent users during peak shopping seasons (holidays). The platform must handle product catalogs, user accounts, shopping carts, and payment processing. Data residency requirements mandate that EU customer data stays in EU regions (GDPR compliance).",

		UserProfile: UserProfile{
			InitialConcurrent: 5000,
			PeakConcurrent:    10000,
			AverageSession: SessionProfile{
				DurationMinutes: 15,
				PageViews:       12,
				DataPerRequest:  100 * 1024,
			},
			PeakTimes: []string{"Black Friday", "Cyber Monday", "Christmas Season (Nov-Dec)"},
			GeographicMix: map[string]float64{
				"us-east":   0.35,
				"us-west":   0.15,
				"europe":    0.35,
				"asia":      0.10,
				"australia": 0.05,
			},
		},

		TrafficPattern: TrafficPattern{
			ReadsPercentage:  0.70,
			WritesPercentage: 0.25,
			StaticPercentage: 0.05,
			PeakMultiplier:   10.0, // Holiday spikes
			DailyPattern:     "business-hours",
		},

		ComplianceNeeds: []string{"GDPR", "PCI-DSS (payment data)"},
		DataResidency:   []string{"EU data must stay in EU region"},

		Tasks: []Task{
			{
				Step:        1,
				Title:       "Deploy Multi-Region Infrastructure",
				Description: "Deploy in US-East, Europe, and Asia regions for global coverage",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Multi-region reduces latency for global users from 200ms to <50ms.",
			},
			{
				Step:        2,
				Title:       "Implement CDN",
				Description: "Deploy CDN for product images and static assets",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "CDN edge caching is 100x faster than cross-region data transfer.",
			},
			{
				Step:        3,
				Title:       "Configure Database Sharding",
				Description: "Shard database by user region for scalability and compliance",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Shard by region: EU users → EU DB, US users → US DB (GDPR compliance).",
			},
		},

		BonusObjectives: []string{
			"Deploy to 3+ regions",
			"Achieve 99.99% uptime",
			"P99 latency <75ms globally",
			"Keep costs under $750/month",
		},
	}
}

func createStreamingScenario() *Scenario {
	return &Scenario{
		CustomerName:     "StreamNow",
		BusinessType:     "Video Streaming Platform",
		CurrentSituation: "StreamNow is a video streaming service competing with Netflix and YouTube. After a viral marketing campaign, they need to handle 100,000 concurrent streams globally. The platform must deliver video content with minimal buffering, support multiple quality levels, and maintain 99.999% uptime (five nines). Content must be distributed globally with <100ms latency to edge servers.",

		UserProfile: UserProfile{
			InitialConcurrent: 50000,
			PeakConcurrent:    100000,
			AverageSession: SessionProfile{
				DurationMinutes: 45,
				PageViews:       5,
				DataPerRequest:  2 * 1024 * 1024, // 2MB/request (video chunks)
			},
			PeakTimes: []string{"Prime time (7-11 PM all timezones)", "Weekends all day"},
			GeographicMix: map[string]float64{
				"us-east":   0.25,
				"us-west":   0.15,
				"europe":    0.30,
				"asia":      0.20,
				"australia": 0.10,
			},
		},

		TrafficPattern: TrafficPattern{
			ReadsPercentage:  0.95,
			WritesPercentage: 0.02,
			StaticPercentage: 0.03,
			PeakMultiplier:   5.0,
			DailyPattern:     "evening-peak",
		},

		Tasks: []Task{
			{
				Step:        1,
				Title:       "Deploy Global CDN",
				Description: "Deploy CDN in all 5 regions with edge caching for video content",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "CDN is essential for video streaming. 95%+ of traffic should hit edge cache.",
			},
			{
				Step:        2,
				Title:       "Implement Database Sharding",
				Description: "Shard user data and metadata across regions for global scale",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "100K+ concurrent users require sharding. Use user_id as shard key.",
			},
			{
				Step:        3,
				Title:       "Configure Advanced Caching",
				Description: "Multi-layer caching: CDN edge + application cache + database cache",
				Type:        TaskTypeInfrastructure,
				Mandatory:   true,
				Hint:        "Layer 1: CDN (video), Layer 2: App cache (metadata), Layer 3: DB cache (queries)",
			},
		},

		BonusObjectives: []string{
			"Achieve five nines uptime (99.999%)",
			"P99 latency <50ms globally",
			"CDN cache hit rate >95%",
			"Keep costs under $3500/month",
		},
	}
}
