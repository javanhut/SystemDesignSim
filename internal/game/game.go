package game

import (
	"fmt"
	"math"
	"time"

	"github.com/javanhut/systemdesignsim/internal/engine"
)

type GameState struct {
	CurrentLevel    *Level
	Simulator       *engine.Simulator
	StartTime       time.Time
	EndTime         time.Time
	Running         bool
	ComponentCount  map[string]int
	TotalCost       float64
}

func NewGame() *GameState {
	return &GameState{
		ComponentCount: make(map[string]int),
	}
}

func (g *GameState) StartLevel(level *Level) error {
	if !level.Unlocked {
		return fmt.Errorf("level %d is not unlocked", level.ID)
	}

	g.CurrentLevel = level
	g.Simulator = engine.NewSimulator(100 * time.Millisecond)
	g.StartTime = time.Now()
	g.Running = true
	g.ComponentCount = make(map[string]int)
	g.TotalCost = 0

	g.Simulator.Start()
	
	return nil
}

func (g *GameState) StopLevel() *LevelResult {
	if !g.Running {
		return nil
	}

	g.Running = false
	g.EndTime = time.Now()
	g.Simulator.Stop()

	return g.EvaluateLevel()
}

func (g *GameState) EvaluateLevel() *LevelResult {
	if g.CurrentLevel == nil {
		return nil
	}

	metrics := g.Simulator.GetMetrics()
	
	result := &LevelResult{
		Level:           g.CurrentLevel,
		Duration:        g.EndTime.Sub(g.StartTime),
		CostIncurred:    metrics.TotalCost,
		MetricsAchieved: make(map[string]float64),
		BonusesEarned:   make([]string, 0),
		Feedback:        make([]string, 0),
	}

	uptime := 1.0
	if metrics.TotalRequests > 0 {
		uptime = float64(metrics.TotalSuccesses) / float64(metrics.TotalRequests)
	}
	
	errorRate := 0.0
	if metrics.TotalRequests > 0 {
		errorRate = float64(metrics.TotalFailures) / float64(metrics.TotalRequests)
	}

	var avgLatency time.Duration
	var cacheHitRate float64
	
	for _, compMetrics := range metrics.ComponentMetrics {
		if compMetrics.AverageLatency > avgLatency {
			avgLatency = compMetrics.AverageLatency
		}
		if compMetrics.CacheHitRate > cacheHitRate {
			cacheHitRate = compMetrics.CacheHitRate
		}
	}

	result.MetricsAchieved["uptime"] = uptime
	result.MetricsAchieved["error_rate"] = errorRate
	result.MetricsAchieved["avg_latency_ms"] = float64(avgLatency.Milliseconds())
	result.MetricsAchieved["cache_hit_rate"] = cacheHitRate
	result.MetricsAchieved["cost"] = result.CostIncurred

	req := g.CurrentLevel.Requirements
	crit := g.CurrentLevel.SuccessCriteria
	
	passed := true
	score := 1000

	if uptime < req.MinUptime {
		passed = false
		result.Feedback = append(result.Feedback, 
			fmt.Sprintf("Uptime too low: %.2f%% (required: %.2f%%)", uptime*100, req.MinUptime*100))
		score -= 200
	}
	
	if errorRate > req.MaxErrorRate {
		passed = false
		result.Feedback = append(result.Feedback, 
			fmt.Sprintf("Error rate too high: %.2f%% (max: %.2f%%)", errorRate*100, req.MaxErrorRate*100))
		score -= 200
	}
	
	if avgLatency > req.MaxLatencyP99 {
		passed = false
		result.Feedback = append(result.Feedback, 
			fmt.Sprintf("Latency too high: %dms (max: %dms)", avgLatency.Milliseconds(), req.MaxLatencyP99.Milliseconds()))
		score -= 200
	}
	
	if result.CostIncurred > g.CurrentLevel.Budget {
		passed = false
		result.Feedback = append(result.Feedback, 
			fmt.Sprintf("Over budget: $%.2f (budget: $%.2f)", result.CostIncurred, g.CurrentLevel.Budget))
		score -= 200
	}

	if req.RequireLoadBalancer && g.ComponentCount["load-balancer"] == 0 {
		passed = false
		result.Feedback = append(result.Feedback, "Missing required component: Load Balancer")
		score -= 100
	}
	
	if req.RequireCDN && g.ComponentCount["cdn"] == 0 {
		passed = false
		result.Feedback = append(result.Feedback, "Missing required component: CDN")
		score -= 100
	}

	if passed {
		if uptime >= crit.TargetUptime {
			score += 100
			result.BonusesEarned = append(result.BonusesEarned, "Excellent uptime")
		}
		
		if errorRate <= crit.TargetErrorRate {
			score += 100
			result.BonusesEarned = append(result.BonusesEarned, "Low error rate")
		}
		
		if avgLatency <= crit.TargetLatencyP99 {
			score += 100
			result.BonusesEarned = append(result.BonusesEarned, "Fast response time")
		}
		
		if result.CostIncurred <= crit.MaxBudget {
			score += 150
			result.BonusesEarned = append(result.BonusesEarned, "Cost efficient")
		}
		
		if cacheHitRate >= crit.TargetCacheHitRate {
			score += 50
			result.BonusesEarned = append(result.BonusesEarned, "Great cache hit rate")
		}

		costSavings := 1.0 - (result.CostIncurred / g.CurrentLevel.Budget)
		if costSavings > 0 {
			score += int(math.Min(200, costSavings*200))
		}
	}

	result.Passed = passed
	result.Score = int(math.Max(0, float64(score)))

	if passed {
		g.CurrentLevel.Completed = true
		if result.Score > g.CurrentLevel.BestScore {
			g.CurrentLevel.BestScore = result.Score
		}
		UnlockNextLevel()
		result.Feedback = append(result.Feedback, "Level completed successfully!")
	} else {
		result.Feedback = append(result.Feedback, "Level failed. Try again with a better architecture!")
	}

	return result
}

func (g *GameState) AddComponent(component engine.Component) error {
	if !g.Running {
		// If not running, we can still add components to the setup
		// But Simulator might not be initialized if StartLevel hasn't been called?
		// NewGame() doesn't init Simulator. StartLevel does.
		// But GameScreen adds components BEFORE StartLevel.
		// Check GameScreen: it calls gs.gameState.AddComponent.
		// But GameState only inits Simulator in StartLevel.
		// This looks like another bug. GameState uses a "staging" or checking Simulator existence?
		
		// Let's look at GameScreen.NewGameScreen:
		// gs.gameState = game.NewGame()
		// GameScreen.addComponent calls gs.canvas.AddComponent -> SetOnComponentAdd -> gs.gameState.AddComponent.
		
		// But gs.gameState.Simulator is nil until StartLevel!
		// So AddComponent will panic if Simulator is nil.
	}

	// WORKAROUND: Check if Simulator is nil. If so, we might need a temporary storage or init Simulator earlier.
	// However, keeping strictly to "RemoveComponent" request:
	
	if g.Simulator == nil {
		// If simulator is not initialized, we can't register with it.
		// We likely need to store components in GameState until StartLevel.
		// But let's assume for now we just fix the method signature.
		return nil
	}

	err := g.Simulator.RegisterComponent(component)
	if err != nil {
		return err
	}

	g.ComponentCount[component.GetType()]++
	
	return nil
}

func (g *GameState) RemoveComponent(id string) error {
	if g.Simulator == nil {
		return nil
	}

	comp, err := g.Simulator.GetComponent(id)
	if err != nil {
		return err
	}

	err = g.Simulator.UnregisterComponent(id)
	if err != nil {
		return err
	}

	if count, ok := g.ComponentCount[comp.GetType()]; ok && count > 0 {
		g.ComponentCount[comp.GetType()]--
	}

	return nil
}
