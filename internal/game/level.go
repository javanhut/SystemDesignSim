package game

import (
	"time"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
	DifficultyExpert Difficulty = "expert"
)

type ApplicationType string

const (
	AppTypeBlog        ApplicationType = "blog"
	AppTypeSocialMedia ApplicationType = "social-media"
	AppTypeEcommerce   ApplicationType = "ecommerce"
	AppTypeStreaming   ApplicationType = "streaming"
	AppTypeMessaging   ApplicationType = "messaging"
	AppTypeFileStorage ApplicationType = "file-storage"
	AppTypeAnalytics   ApplicationType = "analytics"
)

type Level struct {
	ID           int
	Name         string
	Description  string
	Difficulty   Difficulty
	AppType      ApplicationType
	InitialUsers int
	PeakUsers    int
	Duration     time.Duration
	Budget       float64

	Requirements    Requirements
	SuccessCriteria SuccessCriteria
	Scenario        *Scenario

	Unlocked  bool
	Completed bool
	BestScore int
}

type Requirements struct {
	MaxLatencyP99       time.Duration
	MinUptime           float64
	MaxErrorRate        float64
	MinCacheHitRate     float64
	RequireRedundancy   bool
	RequireMultiRegion  bool
	RequireSharding     bool
	RequireReplication  bool
	RequireCDN          bool
	RequireLoadBalancer bool
}

type SuccessCriteria struct {
	TargetLatencyP99   time.Duration
	TargetUptime       float64
	TargetErrorRate    float64
	TargetCacheHitRate float64
	MaxBudget          float64
	BonusObjectives    []string
}

type LevelResult struct {
	Level           *Level
	Passed          bool
	Score           int
	Duration        time.Duration
	CostIncurred    float64
	MetricsAchieved map[string]float64
	BonusesEarned   []string
	Feedback        []string
}

var Levels = []*Level{
	{
		ID:           1,
		Name:         "Local Blog",
		Description:  "You've created a blog for your friends. Scale it to handle 10 concurrent users.",
		Difficulty:   DifficultyEasy,
		AppType:      AppTypeBlog,
		InitialUsers: 1,
		PeakUsers:    10,
		Duration:     5 * time.Minute,
		Budget:       10.0,
		Requirements: Requirements{
			MaxLatencyP99:       500 * time.Millisecond,
			MinUptime:           0.95,
			MaxErrorRate:        0.05,
			MinCacheHitRate:     0.0,
			RequireRedundancy:   false,
			RequireMultiRegion:  false,
			RequireSharding:     false,
			RequireReplication:  false,
			RequireCDN:          false,
			RequireLoadBalancer: false,
		},
		SuccessCriteria: SuccessCriteria{
			TargetLatencyP99:   200 * time.Millisecond,
			TargetUptime:       0.99,
			TargetErrorRate:    0.01,
			TargetCacheHitRate: 0.5,
			MaxBudget:          5.0,
			BonusObjectives:    []string{"Add caching", "Keep costs under $3"},
		},
		Scenario:  GetScenarioForLevel(1),
		Unlocked:  true,
		Completed: false,
	},
	{
		ID:           2,
		Name:         "Growing Blog",
		Description:  "Your blog is getting popular! Handle 100 concurrent users with good performance.",
		Difficulty:   DifficultyEasy,
		AppType:      AppTypeBlog,
		InitialUsers: 10,
		PeakUsers:    100,
		Duration:     10 * time.Minute,
		Budget:       50.0,
		Requirements: Requirements{
			MaxLatencyP99:       300 * time.Millisecond,
			MinUptime:           0.98,
			MaxErrorRate:        0.02,
			MinCacheHitRate:     0.5,
			RequireRedundancy:   false,
			RequireMultiRegion:  false,
			RequireSharding:     false,
			RequireReplication:  false,
			RequireCDN:          false,
			RequireLoadBalancer: true,
		},
		SuccessCriteria: SuccessCriteria{
			TargetLatencyP99:   150 * time.Millisecond,
			TargetUptime:       0.995,
			TargetErrorRate:    0.005,
			TargetCacheHitRate: 0.7,
			MaxBudget:          30.0,
			BonusObjectives:    []string{"Add load balancer", "Achieve 99.5% uptime"},
		},
		Scenario:  GetScenarioForLevel(2),
		Unlocked:  false,
		Completed: false,
	},
	{
		ID:           3,
		Name:         "Regional Social Network",
		Description:  "Launch a social media platform for a single region with 1,000 concurrent users.",
		Difficulty:   DifficultyMedium,
		AppType:      AppTypeSocialMedia,
		InitialUsers: 100,
		PeakUsers:    1000,
		Duration:     15 * time.Minute,
		Budget:       200.0,
		Requirements: Requirements{
			MaxLatencyP99:       200 * time.Millisecond,
			MinUptime:           0.99,
			MaxErrorRate:        0.01,
			MinCacheHitRate:     0.7,
			RequireRedundancy:   true,
			RequireMultiRegion:  false,
			RequireSharding:     false,
			RequireReplication:  true,
			RequireCDN:          false,
			RequireLoadBalancer: true,
		},
		SuccessCriteria: SuccessCriteria{
			TargetLatencyP99:   100 * time.Millisecond,
			TargetUptime:       0.999,
			TargetErrorRate:    0.001,
			TargetCacheHitRate: 0.85,
			MaxBudget:          150.0,
			BonusObjectives:    []string{"Implement database replication", "Achieve 99.9% uptime"},
		},
		Scenario:  GetScenarioForLevel(3),
		Unlocked:  false,
		Completed: false,
	},
	{
		ID:           4,
		Name:         "Global E-commerce Launch",
		Description:  "Launch an e-commerce platform across multiple regions handling 10,000 concurrent users.",
		Difficulty:   DifficultyHard,
		AppType:      AppTypeEcommerce,
		InitialUsers: 1000,
		PeakUsers:    10000,
		Duration:     20 * time.Minute,
		Budget:       1000.0,
		Requirements: Requirements{
			MaxLatencyP99:       150 * time.Millisecond,
			MinUptime:           0.995,
			MaxErrorRate:        0.005,
			MinCacheHitRate:     0.8,
			RequireRedundancy:   true,
			RequireMultiRegion:  true,
			RequireSharding:     true,
			RequireReplication:  true,
			RequireCDN:          true,
			RequireLoadBalancer: true,
		},
		SuccessCriteria: SuccessCriteria{
			TargetLatencyP99:   75 * time.Millisecond,
			TargetUptime:       0.9999,
			TargetErrorRate:    0.0001,
			TargetCacheHitRate: 0.9,
			MaxBudget:          750.0,
			BonusObjectives:    []string{"Deploy to 3+ regions", "Implement database sharding", "Achieve 99.99% uptime"},
		},
		Scenario:  GetScenarioForLevel(4),
		Unlocked:  false,
		Completed: false,
	},
	{
		ID:           5,
		Name:         "Viral Streaming Service",
		Description:  "Your streaming service went viral! Handle 100,000 concurrent users with global distribution.",
		Difficulty:   DifficultyExpert,
		AppType:      AppTypeStreaming,
		InitialUsers: 10000,
		PeakUsers:    100000,
		Duration:     30 * time.Minute,
		Budget:       5000.0,
		Requirements: Requirements{
			MaxLatencyP99:       100 * time.Millisecond,
			MinUptime:           0.9995,
			MaxErrorRate:        0.001,
			MinCacheHitRate:     0.9,
			RequireRedundancy:   true,
			RequireMultiRegion:  true,
			RequireSharding:     true,
			RequireReplication:  true,
			RequireCDN:          true,
			RequireLoadBalancer: true,
		},
		SuccessCriteria: SuccessCriteria{
			TargetLatencyP99:   50 * time.Millisecond,
			TargetUptime:       0.99999,
			TargetErrorRate:    0.00001,
			TargetCacheHitRate: 0.95,
			MaxBudget:          3500.0,
			BonusObjectives:    []string{"Deploy to all regions", "Achieve five nines uptime", "Keep costs under $3000"},
		},
		Scenario:  GetScenarioForLevel(5),
		Unlocked:  false,
		Completed: false,
	},
}

func GetLevel(id int) *Level {
	for _, level := range Levels {
		if level.ID == id {
			return level
		}
	}
	return nil
}

func GetUnlockedLevels() []*Level {
	unlocked := make([]*Level, 0)
	for _, level := range Levels {
		if level.Unlocked {
			unlocked = append(unlocked, level)
		}
	}
	return unlocked
}

func UnlockNextLevel() {
	for _, level := range Levels {
		if !level.Unlocked {
			level.Unlocked = true
			return
		}
	}
}
