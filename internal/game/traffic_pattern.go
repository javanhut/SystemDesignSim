package game

import (
	"math"
	"time"
)

type TrafficGenerator struct {
	Pattern     *TrafficPattern
	UserProfile *UserProfile
	CurrentTime time.Time
	StartTime   time.Time
	BaselineRPS int
}

func NewTrafficGenerator(pattern *TrafficPattern, userProfile *UserProfile, baselineRPS int) *TrafficGenerator {
	now := time.Now()
	return &TrafficGenerator{
		Pattern:     pattern,
		UserProfile: userProfile,
		CurrentTime: now,
		StartTime:   now,
		BaselineRPS: baselineRPS,
	}
}

func (tg *TrafficGenerator) CalculateCurrentRPS(currentTime time.Time) int {
	tg.CurrentTime = currentTime

	hourOfDay := currentTime.Hour()
	dayMultiplier := tg.getDailyMultiplier(hourOfDay)

	baseRPS := float64(tg.BaselineRPS) * dayMultiplier

	return int(baseRPS)
}

func (tg *TrafficGenerator) getDailyMultiplier(hourOfDay int) float64 {
	switch tg.Pattern.DailyPattern {
	case "steady":
		return 1.0

	case "business-hours":
		if hourOfDay >= 9 && hourOfDay <= 17 {
			return tg.Pattern.PeakMultiplier
		} else if hourOfDay >= 7 && hourOfDay <= 19 {
			return (tg.Pattern.PeakMultiplier + 1.0) / 2.0
		}
		return 0.3

	case "evening-peak":
		if hourOfDay >= 19 && hourOfDay <= 23 {
			return tg.Pattern.PeakMultiplier
		} else if hourOfDay >= 17 || hourOfDay <= 1 {
			return (tg.Pattern.PeakMultiplier + 1.0) / 2.0
		} else if hourOfDay >= 12 && hourOfDay <= 14 {
			return (tg.Pattern.PeakMultiplier + 1.0) / 2.5
		}
		return 0.4

	default:
		return 1.0
	}
}

func (tg *TrafficGenerator) GetRequestType() string {
	readThreshold := tg.Pattern.ReadsPercentage
	writeThreshold := readThreshold + tg.Pattern.WritesPercentage

	rand := math.Mod(float64(time.Now().UnixNano()), 100.0)

	if rand < readThreshold {
		return "read"
	} else if rand < writeThreshold {
		return "write"
	}
	return "static"
}

func (tg *TrafficGenerator) IsStaticRequest() bool {
	return tg.GetRequestType() == "static"
}

func (tg *TrafficGenerator) GetExpectedDataSize() int64 {
	return tg.UserProfile.AverageSession.DataPerRequest
}

func (tg *TrafficGenerator) SimulateUserSession() *UserSession {
	return &UserSession{
		SessionID:         generateSessionID(),
		StartTime:         time.Now(),
		DurationMinutes:   tg.UserProfile.AverageSession.DurationMinutes,
		ExpectedPageViews: tg.UserProfile.AverageSession.PageViews,
		DataPerRequest:    tg.UserProfile.AverageSession.DataPerRequest,
		ActualPageViews:   0,
	}
}

type UserSession struct {
	SessionID         string
	StartTime         time.Time
	DurationMinutes   int
	ExpectedPageViews int
	DataPerRequest    int64
	ActualPageViews   int
}

func (us *UserSession) IsActive() bool {
	elapsed := time.Since(us.StartTime)
	return elapsed < time.Duration(us.DurationMinutes)*time.Minute
}

func (us *UserSession) ShouldMakeRequest() bool {
	if !us.IsActive() {
		return false
	}

	if us.ActualPageViews >= us.ExpectedPageViews {
		return false
	}

	elapsed := time.Since(us.StartTime).Minutes()
	expectedProgress := elapsed / float64(us.DurationMinutes)
	actualProgress := float64(us.ActualPageViews) / float64(us.ExpectedPageViews)

	return actualProgress < expectedProgress
}

func (us *UserSession) RecordRequest() {
	us.ActualPageViews++
}

func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" +
		string(rune('A'+time.Now().UnixNano()%26))
}

type LoadProjector struct {
	GrowthProjection *GrowthProjection
	InitialLoad      int
	StartDate        time.Time
}

func NewLoadProjector(growth *GrowthProjection, initialLoad int) *LoadProjector {
	return &LoadProjector{
		GrowthProjection: growth,
		InitialLoad:      initialLoad,
		StartDate:        time.Now(),
	}
}

func (lp *LoadProjector) ProjectLoad(targetDate time.Time) int {
	monthsElapsed := targetDate.Sub(lp.StartDate).Hours() / (24 * 30)

	if monthsElapsed <= 0 {
		return lp.InitialLoad
	}

	growth := math.Pow(lp.GrowthProjection.MonthlyGrowth, monthsElapsed)
	projected := float64(lp.InitialLoad) * growth

	if projected > float64(lp.GrowthProjection.ExpectedPeak) {
		return lp.GrowthProjection.ExpectedPeak
	}

	return int(projected)
}

func (lp *LoadProjector) GetTimeToReachPeak() time.Duration {
	return time.Duration(lp.GrowthProjection.TimeToReach) * 30 * 24 * time.Hour
}

func (lp *LoadProjector) GetMonthlyGrowthRate() float64 {
	return (lp.GrowthProjection.MonthlyGrowth - 1.0) * 100.0
}

type GeographicDistributor struct {
	Distribution map[string]float64
}

func NewGeographicDistributor(dist GeographicDistribution) *GeographicDistributor {
	return &GeographicDistributor{
		Distribution: dist.Distribution,
	}
}

func (gd *GeographicDistributor) SelectRegion() string {
	rand := math.Mod(float64(time.Now().UnixNano()), 1.0)

	cumulative := 0.0
	for region, percentage := range gd.Distribution {
		cumulative += percentage
		if rand < cumulative {
			return region
		}
	}

	for region := range gd.Distribution {
		return region
	}

	return "us-east"
}

func (gd *GeographicDistributor) GetDistribution() map[string]float64 {
	return gd.Distribution
}

func (gd *GeographicDistributor) GetPrimaryRegion() string {
	maxPercentage := 0.0
	primaryRegion := ""

	for region, percentage := range gd.Distribution {
		if percentage > maxPercentage {
			maxPercentage = percentage
			primaryRegion = region
		}
	}

	return primaryRegion
}

type DataGrowthCalculator struct {
	Requirements *DataRequirements
	StartDate    time.Time
}

func NewDataGrowthCalculator(requirements *DataRequirements) *DataGrowthCalculator {
	return &DataGrowthCalculator{
		Requirements: requirements,
		StartDate:    time.Now(),
	}
}

func (dgc *DataGrowthCalculator) CalculateCurrentSize(currentDate time.Time) int64 {
	monthsElapsed := int(currentDate.Sub(dgc.StartDate).Hours() / (24 * 30))

	if monthsElapsed <= 0 {
		return dgc.Requirements.InitialSize
	}

	growth := dgc.Requirements.GrowthPerMonth * int64(monthsElapsed)
	return dgc.Requirements.InitialSize + growth
}

func (dgc *DataGrowthCalculator) CalculateSizeAtDate(targetDate time.Time) int64 {
	return dgc.CalculateCurrentSize(targetDate)
}

func (dgc *DataGrowthCalculator) GetProjectedSize(months int) int64 {
	targetDate := dgc.StartDate.Add(time.Duration(months) * 30 * 24 * time.Hour)
	return dgc.CalculateCurrentSize(targetDate)
}

func (dgc *DataGrowthCalculator) GetStorageRequirements(months int) (int64, int64) {
	dataSize := dgc.GetProjectedSize(months)

	backupSize := int64(0)
	if dgc.Requirements.BackupRequired {
		backupMultiplier := dgc.Requirements.RetentionDays / 30
		if backupMultiplier < 1 {
			backupMultiplier = 1
		}
		backupSize = dataSize * int64(backupMultiplier)
	}

	return dataSize, backupSize
}

type LatencyCalculator struct {
	BaseLatency      time.Duration
	NetworkLatencies map[string]time.Duration
}

func NewLatencyCalculator(baseLatency time.Duration) *LatencyCalculator {
	return &LatencyCalculator{
		BaseLatency: baseLatency,
		NetworkLatencies: map[string]time.Duration{
			"us-east-us-east":     5 * time.Millisecond,
			"us-east-us-west":     70 * time.Millisecond,
			"us-east-europe":      90 * time.Millisecond,
			"us-east-asia":        150 * time.Millisecond,
			"us-east-australia":   180 * time.Millisecond,
			"us-west-us-west":     5 * time.Millisecond,
			"us-west-europe":      140 * time.Millisecond,
			"us-west-asia":        120 * time.Millisecond,
			"us-west-australia":   130 * time.Millisecond,
			"europe-europe":       5 * time.Millisecond,
			"europe-asia":         120 * time.Millisecond,
			"europe-australia":    220 * time.Millisecond,
			"asia-asia":           5 * time.Millisecond,
			"asia-australia":      90 * time.Millisecond,
			"australia-australia": 5 * time.Millisecond,
		},
	}
}

func (lc *LatencyCalculator) CalculateLatency(fromRegion, toRegion string) time.Duration {
	key1 := fromRegion + "-" + toRegion
	key2 := toRegion + "-" + fromRegion

	if latency, exists := lc.NetworkLatencies[key1]; exists {
		return lc.BaseLatency + latency
	}

	if latency, exists := lc.NetworkLatencies[key2]; exists {
		return lc.BaseLatency + latency
	}

	return lc.BaseLatency + 100*time.Millisecond
}

func (lc *LatencyCalculator) GetOptimalRegion(userRegion string, availableRegions []string) string {
	if len(availableRegions) == 0 {
		return userRegion
	}

	minLatency := time.Duration(math.MaxInt64)
	optimalRegion := availableRegions[0]

	for _, region := range availableRegions {
		latency := lc.CalculateLatency(userRegion, region)
		if latency < minLatency {
			minLatency = latency
			optimalRegion = region
		}
	}

	return optimalRegion
}
