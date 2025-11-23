package screens

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/components/api"
	"github.com/javanhut/systemdesignsim/internal/components/cache"
	"github.com/javanhut/systemdesignsim/internal/components/cdn"
	"github.com/javanhut/systemdesignsim/internal/components/database"
	"github.com/javanhut/systemdesignsim/internal/components/loadbalancer"
	"github.com/javanhut/systemdesignsim/internal/components/networking"
	"github.com/javanhut/systemdesignsim/internal/engine"
	"github.com/javanhut/systemdesignsim/internal/game"
	"github.com/javanhut/systemdesignsim/internal/gui"
	guicanvas "github.com/javanhut/systemdesignsim/internal/gui/canvas"
	"github.com/javanhut/systemdesignsim/internal/gui/widgets"
	"github.com/javanhut/systemdesignsim/internal/network"
)

type GameScreen struct {
	window    fyne.Window
	level     *game.Level
	gameState *game.GameState
	canvas    *guicanvas.GraphCanvas

	metricsLabel   *widget.Label
	costLabel      *widget.Label
	statusLabel    *widget.Label
	summaryLabel   *widget.Label
	testPlanLabel  *widget.Label
	userCountLabel *widget.Label
	hintsLabel     *widget.Label

	playButton   *widget.Button
	stopButton   *widget.Button
	submitButton *widget.Button

	componentCounter int

	running          bool
	stopChan         chan bool
	trafficGenerator *game.TrafficGenerator

	networkSettings    networkConfig
	securitySettings   securityConfig
	dnsSettings        dnsConfig
	deploymentSettings deploymentConfig
	monitoringSettings monitoringConfig
}

type networkConfig struct {
	vpcPreset string
	region    string
	nat       bool
	sgPreset  string
}

type securityConfig struct {
	webSG bool
	appSG bool
	dbSG  bool
	waf   bool
}

type dnsConfig struct {
	provider      string
	routingPolicy string
	edgeScope     string
}

type deploymentConfig struct {
	strategy        string
	batchSize       int
	healthGraceSecs int
	autoScaling     bool
	minInstances    int
	maxInstances    int
}

type monitoringConfig struct {
	metrics   bool
	alerts    bool
	synthetic bool
	backups   bool
	drRegion  string
}

func NewGameScreen(window fyne.Window, level *game.Level) *GameScreen {
	gs := &GameScreen{
		window:    window,
		level:     level,
		gameState: game.NewGame(),
		canvas:    guicanvas.NewGraphCanvas(),
		stopChan:  make(chan bool),

		networkSettings: networkConfig{
			vpcPreset: "single-az",
			region:    "us-east-1",
			nat:       true,
			sgPreset:  "web",
		},
		securitySettings: securityConfig{
			webSG: true, appSG: true, dbSG: true, waf: false,
		},
		dnsSettings: dnsConfig{
			provider:      "CloudFront",
			routingPolicy: "Latency",
			edgeScope:     "Global",
		},
		deploymentSettings: deploymentConfig{
			strategy:        "Rolling",
			batchSize:       2,
			healthGraceSecs: 30,
			autoScaling:     true,
			minInstances:    2,
			maxInstances:    6,
		},
		monitoringSettings: monitoringConfig{
			metrics: true, alerts: true, synthetic: true, backups: true, drRegion: "us-west-1",
		},
	}

	gs.setupCallbacks()

	return gs
}

func (gs *GameScreen) setupCallbacks() {
	gs.canvas.SetOnComponentAdd(func(vc *gui.VisualComponent) {
		if vc.Component != nil {
			gs.gameState.AddComponent(vc.Component)
		}
	})

	gs.canvas.SetOnConnectionAdd(func(conn *gui.Connection) {
		gs.linkComponents(conn.From, conn.To)
	})

	gs.canvas.SetOnComponentClick(func(vc *gui.VisualComponent) {
		widgets.ShowPropertyPanel(vc, gs.window, func() {
			gs.statusLabel.SetText(fmt.Sprintf("Updated %s", vc.ID))
		}, func() {
			gs.canvas.RemoveComponent(vc.ID)
			// We don't check gs.running here because RemoveComponent in GameState handles nil simulator gracefully now
			gs.gameState.RemoveComponent(vc.ID)
			gs.statusLabel.SetText(fmt.Sprintf("Deleted %s", vc.ID))
		})
	})
}

func (gs *GameScreen) linkComponents(from, to *gui.VisualComponent) {
	fromComp := from.GetComponent()
	toComp := to.GetComponent()

	if fromComp == nil || toComp == nil {
		return
	}

	switch fromComp.GetType() {
	case "load-balancer":
		if lb, ok := fromComp.(*loadbalancer.LoadBalancer); ok {
			lb.AddBackend(toComp)
		}
	case "cache-redis", "cache-memcached":
		if c, ok := fromComp.(*cache.Cache); ok {
			c.SetBackend(toComp)
		}
	case "cdn":
		if c, ok := fromComp.(*cdn.CDN); ok {
			c.SetOrigin(toComp)
		}
	case "api-server":
		if apiServer, ok := fromComp.(*api.APIServer); ok {
			switch toComp.GetType() {
			case "database-sql", "database-nosql", "database-key-value", "database-document":
				apiServer.SetDatabase(toComp)
			case "cache-redis", "cache-memcached":
				apiServer.SetCache(toComp)
			}
		}
	}
}

func (gs *GameScreen) Build() fyne.CanvasObject {
	toolbox := gs.createToolbox()

	canvasContainer := container.NewMax(gs.canvas)

	metricsPanel := gs.createMetricsPanel()

	controlsPanel := gs.createControlsPanel()

	// Make toolbox scrollable with max height
	toolboxScroll := container.NewVScroll(container.NewVBox(toolbox))
	toolboxScroll.SetMinSize(fyne.NewSize(250, 600))
	leftPanel := gs.wrapPanel(toolboxScroll)

	// Make metrics panel scrollable with max height
	metricsScroll := container.NewVScroll(metricsPanel)
	metricsScroll.SetMinSize(fyne.NewSize(300, 600))
	rightPanel := gs.wrapPanel(metricsScroll)

	header := gs.createHeader()
	scenarioBlock := gs.createScenarioBlock()

	// Combine header and scenario block
	topSection := container.NewVBox(header, scenarioBlock)

	centerPanel := container.NewBorder(topSection, controlsPanel, nil, nil, canvasContainer)

	mainContent := container.NewBorder(
		nil,
		nil,
		leftPanel,
		rightPanel,
		centerPanel,
	)

	return mainContent
}

func (gs *GameScreen) createToolbox() *widget.Card {
	apiBtn := widget.NewButton("API Server", func() {
		gs.addComponent(gui.ComponentTypeAPIServer)
	})
	apiDesc := widget.NewLabel("Processes requests, business logic. ~10ms latency")
	apiDesc.Wrapping = fyne.TextWrapWord

	dbBtn := widget.NewButton("Database", func() {
		gs.addComponent(gui.ComponentTypeDatabase)
	})
	dbDesc := widget.NewLabel("Persistent storage. SQL/NoSQL. ~10ms reads")
	dbDesc.Wrapping = fyne.TextWrapWord

	cacheBtn := widget.NewButton("Cache", func() {
		gs.addComponent(gui.ComponentTypeCache)
	})
	cacheDesc := widget.NewLabel("In-memory fast reads. Redis/Memcached. ~1-2ms")
	cacheDesc.Wrapping = fyne.TextWrapWord

	lbBtn := widget.NewButton("Load Balancer", func() {
		gs.addComponent(gui.ComponentTypeLoadBalancer)
	})
	lbDesc := widget.NewLabel("Distributes traffic across servers. High availability")
	lbDesc.Wrapping = fyne.TextWrapWord

	cdnBtn := widget.NewButton("CDN", func() {
		gs.addComponent(gui.ComponentTypeCDN)
	})
	cdnDesc := widget.NewLabel("Edge caching, global distribution. Static content")
	cdnDesc.Wrapping = fyne.TextWrapWord

	gatewayBtn := widget.NewButton("Gateway", func() {
		gs.addComponent(gui.ComponentTypeGateway)
	})
	gatewayDesc := widget.NewLabel("Internet/API gateway. Entry point. ~1ms")
	gatewayDesc.Wrapping = fyne.TextWrapWord

	firewallBtn := widget.NewButton("Firewall", func() {
		gs.addComponent(gui.ComponentTypeFirewall)
	})
	firewallDesc := widget.NewLabel("Security filtering layer. WAF rules. ~2ms")
	firewallDesc.Wrapping = fyne.TextWrapWord

	natBtn := widget.NewButton("NAT", func() {
		gs.addComponent(gui.ComponentTypeNAT)
	})
	natDesc := widget.NewLabel("Network address translation. Private subnets")
	natDesc.Wrapping = fyne.TextWrapWord

	routerBtn := widget.NewButton("Router", func() {
		gs.addComponent(gui.ComponentTypeRouter)
	})
	routerDesc := widget.NewLabel("Network routing layer. Path-based routing")
	routerDesc.Wrapping = fyne.TextWrapWord

	userPoolBtn := widget.NewButton("User Pool", func() {
		gs.addComponent(gui.ComponentTypeUserPool)
	})
	userPoolDesc := widget.NewLabel("Simulated users. Traffic source. Configurable")
	userPoolDesc.Wrapping = fyne.TextWrapWord

	helpBtn := widget.NewButton("? Help", func() {
		gs.showScenarioHelp()
	})

	toolboxContent := container.NewVBox(
		widget.NewLabel("App Components"),
		widget.NewSeparator(),
		apiBtn,
		apiDesc,
		widget.NewSeparator(),
		dbBtn,
		dbDesc,
		widget.NewSeparator(),
		cacheBtn,
		cacheDesc,
		widget.NewSeparator(),
		lbBtn,
		lbDesc,
		widget.NewSeparator(),
		cdnBtn,
		cdnDesc,
		widget.NewSeparator(),
		widget.NewLabel("Network & Users"),
		widget.NewSeparator(),
		gatewayBtn,
		gatewayDesc,
		widget.NewSeparator(),
		firewallBtn,
		firewallDesc,
		widget.NewSeparator(),
		natBtn,
		natDesc,
		widget.NewSeparator(),
		routerBtn,
		routerDesc,
		widget.NewSeparator(),
		userPoolBtn,
		userPoolDesc,
		widget.NewSeparator(),
		widget.NewLabel("Quick Guide:"),
		widget.NewLabel("• Click to add"),
		widget.NewLabel("• 2x click then target to connect"),
		widget.NewLabel("• Green=healthy, Red=overload"),
		widget.NewSeparator(),
		helpBtn,
	)

	return widget.NewCard("Toolbox", "", toolboxContent)
}

func (gs *GameScreen) showScenarioHelp() {
	helpTitle := widget.NewLabel("Level Scenario & Help")
	helpTitle.TextStyle = fyne.TextStyle{Bold: true}
	helpTitle.Alignment = fyne.TextAlignCenter

	var scenarioText string
	if gs.level.Scenario != nil {
		s := gs.level.Scenario

		scenarioText = fmt.Sprintf("CUSTOMER BRIEF\n\n")
		scenarioText += fmt.Sprintf("Client: %s\n", s.CustomerName)
		scenarioText += fmt.Sprintf("Business: %s\n\n", s.BusinessType)
		scenarioText += fmt.Sprintf("Situation:\n%s\n\n", s.CurrentSituation)

		scenarioText += fmt.Sprintf("USER PROFILE\n")
		scenarioText += fmt.Sprintf("• Concurrent Users: %d (peak: %d)\n",
			s.UserProfile.InitialConcurrent, s.UserProfile.PeakConcurrent)
		scenarioText += fmt.Sprintf("• Session: %d min, %d page views\n",
			s.UserProfile.AverageSession.DurationMinutes, s.UserProfile.AverageSession.PageViews)
		scenarioText += fmt.Sprintf("• Peak Times: %s\n\n", s.UserProfile.PeakTimes[0])

		scenarioText += fmt.Sprintf("TRAFFIC PATTERN\n")
		scenarioText += fmt.Sprintf("• Reads: %.0f%% | Writes: %.0f%% | Static: %.0f%%\n",
			s.TrafficPattern.ReadsPercentage*100,
			s.TrafficPattern.WritesPercentage*100,
			s.TrafficPattern.StaticPercentage*100)
		scenarioText += fmt.Sprintf("• Peak Multiplier: %.1fx\n", s.TrafficPattern.PeakMultiplier)
		scenarioText += fmt.Sprintf("• Pattern: %s\n\n", s.TrafficPattern.DailyPattern)

		scenarioText += fmt.Sprintf("CONSTRAINTS\n")
		scenarioText += fmt.Sprintf("• Budget: $%.2f/month\n", gs.level.Budget)
		scenarioText += fmt.Sprintf("• Max Latency: %dms (P99)\n", gs.level.Requirements.MaxLatencyP99.Milliseconds())
		scenarioText += fmt.Sprintf("• Min Uptime: %.1f%%\n", gs.level.Requirements.MinUptime*100)
		if len(s.ComplianceNeeds) > 0 {
			scenarioText += fmt.Sprintf("• Compliance: %v\n", s.ComplianceNeeds)
		}
	} else {
		scenarioText = fmt.Sprintf("Scenario: %s\n\n%s\n\nObjective:\nHandle %d concurrent users within budget of $%.2f\n\nRequirements:\n• Max Latency: %dms\n• Min Uptime: %.1f%%\n• Max Error Rate: %.1f%%",
			gs.level.Name,
			gs.level.Description,
			gs.level.PeakUsers,
			gs.level.Budget,
			gs.level.Requirements.MaxLatencyP99.Milliseconds(),
			gs.level.Requirements.MinUptime*100,
			gs.level.Requirements.MaxErrorRate*100,
		)
	}

	scenarioInfo := widget.NewLabel(scenarioText)
	scenarioInfo.Wrapping = fyne.TextWrapWord

	var tasksText string
	if gs.level.Scenario != nil && len(gs.level.Scenario.Tasks) > 0 {
		tasksText = "\nTASKS TO COMPLETE\n\n"
		for _, task := range gs.level.Scenario.Tasks {
			mandatoryStr := ""
			if task.Mandatory {
				mandatoryStr = " [REQUIRED]"
			}
			tasksText += fmt.Sprintf("%d. %s%s\n", task.Step, task.Title, mandatoryStr)
			tasksText += fmt.Sprintf("   %s\n", task.Description)
			if task.Hint != "" {
				tasksText += fmt.Sprintf("   💡 Hint: %s\n", task.Hint)
			}
			tasksText += "\n"
		}

		if len(gs.level.Scenario.BonusObjectives) > 0 {
			tasksText += "BONUS OBJECTIVES\n"
			for _, bonus := range gs.level.Scenario.BonusObjectives {
				tasksText += fmt.Sprintf("⭐ %s\n", bonus)
			}
		}
	}

	tasksInfo := widget.NewLabel(tasksText)
	tasksInfo.Wrapping = fyne.TextWrapWord

	componentGuide := widget.NewLabel(
		"\nCOMPONENT GUIDE\n" +
			"• API Server (Blue): Handles requests - watch capacity!\n" +
			"• Database (Purple): Stores data - 10ms latency\n" +
			"• Cache (Green): Fast reads - 1-2ms latency\n" +
			"• Load Balancer (Yellow): Distribute traffic\n" +
			"• CDN (Dark): Global edge caching\n\n" +
			"Health Colors:\n" +
			"Green = Healthy | Yellow = Busy | Orange = Critical | Red = Failing\n\n" +
			"Valid Connections:\n" +
			"Load Balancer → API Server\n" +
			"API Server → Database/Cache\n" +
			"Cache → Database\n" +
			"CDN → API Server",
	)
	componentGuide.Wrapping = fyne.TextWrapWord

	closeBtn := widget.NewButton("Close", func() {
		gs.window.SetContent(gs.Build())
	})

	content := container.NewVBox(
		helpTitle,
		widget.NewSeparator(),
		scenarioInfo,
		widget.NewSeparator(),
		tasksInfo,
		widget.NewSeparator(),
		componentGuide,
		widget.NewSeparator(),
		closeBtn,
	)

	scrollContent := container.NewVScroll(content)
	gs.window.SetContent(scrollContent)
}

func (gs *GameScreen) createMetricsPanel() *widget.Card {
	gs.metricsLabel = widget.NewLabel("Traffic Metrics:\nStatus: ⏳ Ready\nTotal Requests: 0\nRPS (current): 0\nSuccess Rate: 0.0%\nError Rate: 0.0%\nAvg Latency: 0ms\nP99 Latency: 0ms\nUptime: 0.0%")
	gs.costLabel = widget.NewLabel("Total Cost: $0.00/hr")
	gs.statusLabel = widget.NewLabel("Status: Ready")
	gs.userCountLabel = widget.NewLabel(fmt.Sprintf("Users: 0 / %d (Peak)", gs.level.PeakUsers))
	gs.userCountLabel.TextStyle = fyne.TextStyle{Bold: true}
	gs.summaryLabel = widget.NewLabel(gs.buildSummaryText())
	gs.summaryLabel.Wrapping = fyne.TextWrapWord
	gs.testPlanLabel = widget.NewLabel(gs.buildTestPlanText())
	gs.testPlanLabel.Wrapping = fyne.TextWrapWord

	// Add hints label
	gs.hintsLabel = widget.NewLabel(gs.getArchitecturalHints())
	gs.hintsLabel.Wrapping = fyne.TextWrapWord

	objectivesText := fmt.Sprintf(
		"Level: %s\n\nObjectives:\n- Max Latency: %dms\n- Min Uptime: %.1f%%\n- Budget: $%.2f\n- Users: %d",
		gs.level.Name,
		gs.level.Requirements.MaxLatencyP99.Milliseconds(),
		gs.level.Requirements.MinUptime*100,
		gs.level.Budget,
		gs.level.PeakUsers,
	)

	objectives := widget.NewLabel(objectivesText)
	objectives.TextStyle = fyne.TextStyle{Bold: true}

	var taskChecklistText string
	if gs.level.Scenario != nil && len(gs.level.Scenario.Tasks) > 0 {
		taskChecklistText = "\nQuick Tasks:\n"
		mandatoryCount := 0
		for _, task := range gs.level.Scenario.Tasks {
			if task.Mandatory {
				mandatoryCount++
				if mandatoryCount <= 3 {
					taskChecklistText += fmt.Sprintf("☐ %s\n", task.Title)
				}
			}
		}
		if mandatoryCount > 3 {
			taskChecklistText += fmt.Sprintf("... and %d more tasks\n", mandatoryCount-3)
		}
		taskChecklistText += "\nClick '? Help' for details"
	}

	tasksLabel := widget.NewLabel(taskChecklistText)
	tasksLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		objectives,
		widget.NewSeparator(),
		tasksLabel,
		widget.NewSeparator(),
		gs.statusLabel,
		gs.userCountLabel,
		gs.metricsLabel,
		gs.costLabel,
		widget.NewSeparator(),
		widget.NewLabel("Architecture Hints"),
		gs.hintsLabel,
		widget.NewSeparator(),
		widget.NewLabel("Architecture Summary"),
		gs.summaryLabel,
		widget.NewSeparator(),
		widget.NewLabel("Test Plan"),
		gs.testPlanLabel,
	)

	return widget.NewCard("Level Info", "", content)
}

func (gs *GameScreen) createHeader() fyne.CanvasObject {
	title := canvas.NewText(fmt.Sprintf("Level %d: %s", gs.level.ID, gs.level.Name), color.RGBA{R: 180, G: 214, B: 255, A: 255})
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText("Build • Connect • Simulate", color.RGBA{R: 150, G: 190, B: 255, A: 220})
	subtitle.TextSize = 12
	subtitle.Alignment = fyne.TextAlignCenter

	bg := canvas.NewLinearGradient(
		color.RGBA{R: 18, G: 29, B: 48, A: 255},
		color.RGBA{R: 10, G: 18, B: 32, A: 255},
		0,
	)

	return container.NewMax(
		bg,
		container.NewVBox(
			container.NewCenter(title),
			container.NewCenter(subtitle),
		),
	)
}

func (gs *GameScreen) createScenarioBlock() fyne.CanvasObject {
	if gs.level.Scenario == nil {
		// Show basic scenario for levels without detailed scenario
		scenarioText := fmt.Sprintf(
			"Objective: %s\n\n"+
				"Challenge: Handle %d concurrent users with max %dms latency and %.1f%% uptime.\n"+
				"Budget: $%.2f/hr | Max Error Rate: %.1f%%",
			gs.level.Description,
			gs.level.PeakUsers,
			gs.level.Requirements.MaxLatencyP99.Milliseconds(),
			gs.level.Requirements.MinUptime*100,
			gs.level.Budget,
			gs.level.Requirements.MaxErrorRate*100,
		)

		label := widget.NewLabel(scenarioText)
		label.Wrapping = fyne.TextWrapWord
		label.Alignment = fyne.TextAlignCenter

		bg := canvas.NewRectangle(color.RGBA{R: 24, G: 36, B: 52, A: 255})
		bg.StrokeColor = color.RGBA{R: 70, G: 130, B: 255, A: 100}
		bg.StrokeWidth = 1
		bg.CornerRadius = 8

		return container.NewMax(
			bg,
			container.NewPadded(label),
		)
	}

	// Detailed scenario block for levels with scenario
	s := gs.level.Scenario

	scenarioText := fmt.Sprintf(
		"CLIENT: %s | BUSINESS: %s\n\n"+
			"SITUATION: %s\n\n"+
			"USERS: %d concurrent (peak: %d) | SESSION: %d min, %d page views\n"+
			"TRAFFIC: %.0f%% reads | %.0f%% writes | %.0f%% static | Peak: %.1fx | Pattern: %s\n\n"+
			"CONSTRAINTS: Budget $%.2f/hr | P99 Latency < %dms | Uptime > %.1f%%",
		s.CustomerName,
		s.BusinessType,
		s.CurrentSituation,
		s.UserProfile.InitialConcurrent,
		s.UserProfile.PeakConcurrent,
		s.UserProfile.AverageSession.DurationMinutes,
		s.UserProfile.AverageSession.PageViews,
		s.TrafficPattern.ReadsPercentage*100,
		s.TrafficPattern.WritesPercentage*100,
		s.TrafficPattern.StaticPercentage*100,
		s.TrafficPattern.PeakMultiplier,
		s.TrafficPattern.DailyPattern,
		gs.level.Budget,
		gs.level.Requirements.MaxLatencyP99.Milliseconds(),
		gs.level.Requirements.MinUptime*100,
	)

	label := widget.NewLabel(scenarioText)
	label.Wrapping = fyne.TextWrapWord

	bg := canvas.NewRectangle(color.RGBA{R: 24, G: 36, B: 52, A: 255})
	bg.StrokeColor = color.RGBA{R: 70, G: 130, B: 255, A: 100}
	bg.StrokeWidth = 1
	bg.CornerRadius = 8

	return container.NewMax(
		bg,
		container.NewPadded(label),
	)
}

func (gs *GameScreen) wrapPanel(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 18, G: 24, B: 38, A: 255})
	bg.StrokeColor = color.RGBA{R: 70, G: 130, B: 255, A: 120}
	bg.StrokeWidth = 1.5
	bg.CornerRadius = 12

	return container.NewMax(
		bg,
		container.NewPadded(content),
	)
}

func (gs *GameScreen) createControlsPanel() fyne.CanvasObject {
	gs.playButton = widget.NewButton("Start Simulation", func() {
		gs.startSimulation()
	})

	gs.stopButton = widget.NewButton("Stop Simulation", func() {
		gs.stopSimulation()
	})
	gs.stopButton.Disable()

	gs.submitButton = widget.NewButton("Submit Solution", func() {
		gs.submitSolution()
	})
	gs.submitButton.Disable()

	backButton := widget.NewButton("Back to Levels", func() {
		gs.window.SetContent(NewLevelSelectScreen(gs.window).Build())
	})

	controlCenterBtn := widget.NewButton("Control Center", func() {
		gs.showControlCenter()
	})

	planBtn := widget.NewButton("System Plan", func() {
		gs.showSystemPlan()
	})

	hintsBtn := widget.NewButton("Show Hints", func() {
		gs.showArchitecturalHintsDialog()
	})

	return container.NewHBox(
		gs.playButton,
		gs.stopButton,
		gs.submitButton,
		hintsBtn,
		controlCenterBtn,
		planBtn,
		backButton,
	)
}

func (gs *GameScreen) showControlCenter() {
	tabs := container.NewAppTabs(
		container.NewTabItem("VPC & Network", gs.networkTab()),
		container.NewTabItem("Security", gs.securityTab()),
		container.NewTabItem("DNS / CDN", gs.dnsTab()),
		container.NewTabItem("Deployment", gs.deploymentTab()),
		container.NewTabItem("Monitoring/DR", gs.monitoringTab()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	var modal *widget.PopUp
	closeBtn := widget.NewButton("Close", func() {
		if modal != nil {
			modal.Hide()
		}
	})

	content := container.NewBorder(nil, closeBtn, nil, nil, container.NewMax(tabs))
	modal = widget.NewModalPopUp(content, gs.window.Canvas())
	modal.Resize(fyne.NewSize(640, 520))
	modal.Show()
}

func (gs *GameScreen) networkTab() fyne.CanvasObject {
	vpcSelect := widget.NewSelect(network.GetVPCPresetNames(), func(selected string) {
		gs.networkSettings.vpcPreset = selected
	})
	vpcSelect.SetSelected(gs.networkSettings.vpcPreset)

	regions := []string{"us-east-1", "us-west-1", "eu-west-1", "ap-southeast-1", "ap-southeast-2"}
	regionSelect := widget.NewSelect(regions, func(selected string) {
		gs.networkSettings.region = selected
	})
	regionSelect.SetSelected(gs.networkSettings.region)

	natCheck := widget.NewCheck("Include NAT Gateway for private subnets", func(c bool) {
		gs.networkSettings.nat = c
	})
	natCheck.SetChecked(gs.networkSettings.nat)

	sgPreset := widget.NewSelect([]string{"web", "app", "db", "cache", "lb"}, func(s string) {
		gs.networkSettings.sgPreset = s
	})
	sgPreset.SetSelected(gs.networkSettings.sgPreset)

	vpcInfo := widget.NewLabel("Presets: single-az, multi-az, three-tier. NAT enables private subnets to reach the internet. SG presets mirror typical web/app/db tiers.")
	vpcInfo.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		widget.NewLabel("VPC Preset"),
		vpcSelect,
		widget.NewLabel("Region"),
		regionSelect,
		natCheck,
		widget.NewSeparator(),
		widget.NewLabel("Security Group Preset"),
		sgPreset,
		vpcInfo,
	)
}

func (gs *GameScreen) securityTab() fyne.CanvasObject {
	web := widget.NewCheck("Web SG (80/443 open)", func(c bool) { gs.securitySettings.webSG = c })
	web.SetChecked(gs.securitySettings.webSG)
	app := widget.NewCheck("App SG (only LB allowed)", func(c bool) { gs.securitySettings.appSG = c })
	app.SetChecked(gs.securitySettings.appSG)
	db := widget.NewCheck("DB SG (only App allowed)", func(c bool) { gs.securitySettings.dbSG = c })
	db.SetChecked(gs.securitySettings.dbSG)
	waf := widget.NewCheck("WAF Enabled", func(c bool) { gs.securitySettings.waf = c })
	waf.SetChecked(gs.securitySettings.waf)

	info := widget.NewLabel("Compose defense-in-depth: web SG for ingress, app SG for east-west, db SG for storage tier. WAF blocks common exploits.")
	info.Wrapping = fyne.TextWrapWord

	return container.NewVBox(web, app, db, waf, widget.NewSeparator(), info)
}

func (gs *GameScreen) dnsTab() fyne.CanvasObject {
	providers := []string{"CloudFront", "Fastly", "Akamai", "Cloudflare"}
	providerSelect := widget.NewSelect(providers, func(s string) { gs.dnsSettings.provider = s })
	providerSelect.SetSelected(gs.dnsSettings.provider)

	routing := []string{"Simple", "Weighted", "Latency", "Failover", "Geo"}
	routingSelect := widget.NewSelect(routing, func(s string) { gs.dnsSettings.routingPolicy = s })
	routingSelect.SetSelected(gs.dnsSettings.routingPolicy)

	edge := []string{"Global", "NA+EU", "APAC", "Custom"}
	edgeSelect := widget.NewSelect(edge, func(s string) { gs.dnsSettings.edgeScope = s })
	edgeSelect.SetSelected(gs.dnsSettings.edgeScope)

	info := widget.NewLabel("Pick CDN provider and DNS policy. Latency/Geo best for performance; Failover for DR. Edge scope limits where caches deploy.")
	info.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		widget.NewLabel("CDN / DNS Provider"),
		providerSelect,
		widget.NewLabel("Routing Policy"),
		routingSelect,
		widget.NewLabel("Edge Scope"),
		edgeSelect,
		widget.NewSeparator(),
		info,
	)
}

func (gs *GameScreen) deploymentTab() fyne.CanvasObject {
	strategies := []string{"All-at-once", "Rolling", "Blue/Green", "Canary"}
	strategySelect := widget.NewSelect(strategies, func(s string) { gs.deploymentSettings.strategy = s })
	strategySelect.SetSelected(gs.deploymentSettings.strategy)

	batchEntry := widget.NewEntry()
	batchEntry.SetText(fmt.Sprintf("%d", gs.deploymentSettings.batchSize))
	batchEntry.OnChanged = func(val string) {
		if v, err := strconv.Atoi(val); err == nil {
			gs.deploymentSettings.batchSize = v
		}
	}

	healthEntry := widget.NewEntry()
	healthEntry.SetText(fmt.Sprintf("%d", gs.deploymentSettings.healthGraceSecs))
	healthEntry.OnChanged = func(val string) {
		if v, err := strconv.Atoi(val); err == nil {
			gs.deploymentSettings.healthGraceSecs = v
		}
	}

	autoScale := widget.NewCheck("Enable Auto-scaling", func(c bool) { gs.deploymentSettings.autoScaling = c })
	autoScale.SetChecked(gs.deploymentSettings.autoScaling)

	minEntry := widget.NewEntry()
	minEntry.SetText(fmt.Sprintf("%d", gs.deploymentSettings.minInstances))
	minEntry.OnChanged = func(val string) {
		if v, err := strconv.Atoi(val); err == nil {
			gs.deploymentSettings.minInstances = v
		}
	}
	maxEntry := widget.NewEntry()
	maxEntry.SetText(fmt.Sprintf("%d", gs.deploymentSettings.maxInstances))
	maxEntry.OnChanged = func(val string) {
		if v, err := strconv.Atoi(val); err == nil {
			gs.deploymentSettings.maxInstances = v
		}
	}

	info := widget.NewLabel("Choose rollout style and health grace. Auto-scaling bounds control elasticity. Canary/Blue-Green reduce blast radius.")
	info.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		widget.NewLabel("Strategy"),
		strategySelect,
		widget.NewLabel("Batch Size"),
		batchEntry,
		widget.NewLabel("Health Grace (sec)"),
		healthEntry,
		autoScale,
		widget.NewLabel("Min Instances"),
		minEntry,
		widget.NewLabel("Max Instances"),
		maxEntry,
		widget.NewSeparator(),
		info,
	)
}

func (gs *GameScreen) monitoringTab() fyne.CanvasObject {
	metrics := widget.NewCheck("Metrics Dashboard", func(c bool) { gs.monitoringSettings.metrics = c })
	metrics.SetChecked(gs.monitoringSettings.metrics)
	alerts := widget.NewCheck("Alerting (pager/email)", func(c bool) { gs.monitoringSettings.alerts = c })
	alerts.SetChecked(gs.monitoringSettings.alerts)
	synth := widget.NewCheck("Synthetic probes", func(c bool) { gs.monitoringSettings.synthetic = c })
	synth.SetChecked(gs.monitoringSettings.synthetic)
	backups := widget.NewCheck("Automated backups", func(c bool) { gs.monitoringSettings.backups = c })
	backups.SetChecked(gs.monitoringSettings.backups)

	drEntry := widget.NewEntry()
	drEntry.SetText(gs.monitoringSettings.drRegion)
	drEntry.OnChanged = func(val string) { gs.monitoringSettings.drRegion = val }

	info := widget.NewLabel("Monitoring/DR: enable dashboards, paging, synthetic checks, backups, and designate a DR region for failover rehearsals.")
	info.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		metrics,
		alerts,
		synth,
		backups,
		widget.NewLabel("DR Region"),
		drEntry,
		widget.NewSeparator(),
		info,
	)
}

func (gs *GameScreen) showSystemPlan() {
	plan := widget.NewLabel(fmt.Sprintf(
		"Architecture Plan:\n%s\n\nComponent Tips:\n- LB → API → Cache → DB\n- CDN fronts read-heavy/static paths\n- SG: web->app->db chain\n\nTest Harness:\n%s\n\nNext Moves:\n- Add APIs, DB, Cache, LB, CDN.\n- Wire connections, start sim, tune configs.",
		gs.buildSummaryText(),
		gs.buildTestPlanText(),
	))
	plan.Wrapping = fyne.TextWrapWord

	var modal *widget.PopUp
	closeBtn := widget.NewButton("Close", func() {
		if modal != nil {
			modal.Hide()
		}
	})

	content := container.NewBorder(nil, closeBtn, nil, nil, container.NewVScroll(plan))
	modal = widget.NewModalPopUp(content, gs.window.Canvas())
	modal.Resize(fyne.NewSize(520, 480))
	modal.Show()
}

func (gs *GameScreen) showArchitecturalHintsDialog() {
	title := widget.NewLabel("System Design Hints")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	hintsText := gs.getArchitecturalHints()
	hintsLabel := widget.NewLabel(hintsText)
	hintsLabel.Wrapping = fyne.TextWrapWord

	var modal *widget.PopUp
	closeBtn := widget.NewButton("Close", func() {
		if modal != nil {
			modal.Hide()
		}
	})

	scrollContent := container.NewVScroll(hintsLabel)
	content := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()),
		closeBtn,
		nil,
		nil,
		scrollContent,
	)

	modal = widget.NewModalPopUp(content, gs.window.Canvas())
	modal.Resize(fyne.NewSize(700, 600))
	modal.Show()
}

func (gs *GameScreen) addComponent(compType gui.ComponentType) {
	gs.componentCounter++
	id := fmt.Sprintf("%s-%d", compType, gs.componentCounter)

	pos := fyne.NewPos(200+float32(gs.componentCounter*20), 200+float32(gs.componentCounter*20))
	visualComp := gui.NewVisualComponent(id, compType, pos)

	var comp engine.Component
	switch compType {
	case gui.ComponentTypeAPIServer:
		comp = api.NewAPIServer(id, "us-east", api.SizeMedium)
	case gui.ComponentTypeDatabase:
		comp = database.NewDatabase(id, database.DatabaseTypeSQL, "us-east", 10*1024*1024*1024)
	case gui.ComponentTypeCache:
		comp = cache.NewCache(id, "redis", "us-east", 1024*1024*1024, cache.EvictionLRU, time.Hour)
	case gui.ComponentTypeLoadBalancer:
		comp = loadbalancer.NewLoadBalancer(id, "us-east", loadbalancer.StrategyRoundRobin)
	case gui.ComponentTypeCDN:
		comp = cdn.NewCDN(id, []string{"us-east", "us-west", "europe"})
	case gui.ComponentTypeGateway:
		comp = networking.NewGateway(id, "us-east")
	case gui.ComponentTypeFirewall:
		comp = networking.NewFirewall(id, "us-east")
	case gui.ComponentTypeNAT:
		comp = networking.NewNAT(id, "us-east")
	case gui.ComponentTypeRouter:
		comp = networking.NewRouter(id, "us-east")
	case gui.ComponentTypeUserPool:
		comp = networking.NewUserPool(id, "us-east", gs.level.PeakUsers)
	}

	visualComp.SetComponent(comp)
	gs.canvas.AddComponent(visualComp)
}

func (gs *GameScreen) startSimulation() {
	err := gs.gameState.StartLevel(gs.level)
	if err != nil {
		gs.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	// Register all components with the new simulator
	for _, vc := range gs.canvas.GetComponents() {
		if vc.Component != nil {
			if err := gs.gameState.AddComponent(vc.Component); err != nil {
				fmt.Printf("Error adding component %s: %v\n", vc.ID, err)
			}
		}
	}

	// Initialize traffic generator if scenario has traffic pattern
	if gs.level.Scenario != nil {
		baselineRPS := 50 // Base 50 requests/sec, will be modulated by pattern
		gs.trafficGenerator = game.NewTrafficGenerator(
			&gs.level.Scenario.TrafficPattern,
			&gs.level.Scenario.UserProfile,
			baselineRPS,
		)
	}

	gs.running = true
	gs.playButton.Disable()
	gs.stopButton.Enable()
	gs.submitButton.Disable()
	gs.statusLabel.SetText("Status: Running")

	go gs.updateMetrics()
	go gs.simulateTraffic()
	go gs.animateParticles()
}

func (gs *GameScreen) stopSimulation() {
	gs.running = false
	gs.stopChan <- true

	gs.playButton.Enable()
	gs.stopButton.Disable()
	gs.submitButton.Enable()
	gs.statusLabel.SetText("Status: Stopped")
}

func (gs *GameScreen) submitSolution() {
	result := gs.gameState.StopLevel()
	if result == nil {
		return
	}

	gs.running = false

	resultText := fmt.Sprintf(
		"Level %s\n\n%s\n\nScore: %d\n\nMetrics:\n- Uptime: %.2f%%\n- Avg Latency: %.0fms\n- Error Rate: %.2f%%\n- Cost: $%.2f\n\nFeedback:\n",
		result.Level.Name,
		map[bool]string{true: "PASSED", false: "FAILED"}[result.Passed],
		result.Score,
		result.MetricsAchieved["uptime"]*100,
		result.MetricsAchieved["avg_latency_ms"],
		result.MetricsAchieved["error_rate"]*100,
		result.CostIncurred,
	)

	for _, feedback := range result.Feedback {
		resultText += "- " + feedback + "\n"
	}

	if len(result.BonusesEarned) > 0 {
		resultText += "\nBonuses:\n"
		for _, bonus := range result.BonusesEarned {
			resultText += "- " + bonus + "\n"
		}
	}

	dialog := widget.NewLabel(resultText)
	okButton := widget.NewButton("OK", func() {
		gs.window.SetContent(NewLevelSelectScreen(gs.window).Build())
	})

	content := container.NewVBox(dialog, okButton)
	gs.window.SetContent(content)
}

func (gs *GameScreen) updateMetrics() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-gs.stopChan:
			return
		case <-ticker.C:
			if !gs.running {
				return
			}

			metrics := gs.gameState.Simulator.GetMetrics()
			p99Latency := gs.gameState.Simulator.GetP99Latency()

			// Calculate traffic metrics
			successRate := 0.0
			errorRate := 0.0
			uptime := 0.0
			if metrics.TotalRequests > 0 {
				successRate = (float64(metrics.TotalSuccesses) / float64(metrics.TotalRequests)) * 100
				errorRate = (float64(metrics.TotalFailures) / float64(metrics.TotalRequests)) * 100
				uptime = successRate
			}

			avgLatency := int64(0)
			if metrics.TotalRequests > 0 {
				avgLatency = metrics.TotalLatency.Milliseconds() / metrics.TotalRequests
			}

			// Estimate current RPS
			currentRPS := "0"
			if gs.trafficGenerator != nil {
				rps := gs.trafficGenerator.CalculateCurrentRPS(time.Now())
				currentRPS = fmt.Sprintf("%d", rps)
			}

			// Check victory conditions
			passedLatency := p99Latency <= gs.level.Requirements.MaxLatencyP99
			passedUptime := uptime >= gs.level.Requirements.MinUptime*100
			passedBudget := metrics.TotalCost <= gs.level.Budget

			statusIcon := "⏳"
			statusText := "Running"
			if metrics.TotalRequests > 100 { // Only check after enough samples
				if passedLatency && passedUptime && passedBudget {
					statusIcon = "✓"
					statusText = "PASSING"
				} else {
					statusIcon = "✗"
					statusText = "FAILING"
				}
			}

			metricsText := fmt.Sprintf(
				"Traffic Metrics:\n"+
					"Status: %s %s\n"+
					"Total Requests: %d\n"+
					"RPS (current): %s\n"+
					"Success Rate: %.1f%% %s\n"+
					"Error Rate: %.1f%%\n"+
					"Avg Latency: %dms\n"+
					"P99 Latency: %dms %s\n"+
					"Uptime: %.1f%% %s",
				statusIcon,
				statusText,
				metrics.TotalRequests,
				currentRPS,
				successRate,
				gs.getCheckmark(passedUptime),
				errorRate,
				avgLatency,
				p99Latency.Milliseconds(),
				gs.getCheckmark(passedLatency),
				uptime,
				gs.getCheckmark(passedUptime),
			)
			costText := fmt.Sprintf("Cost: $%.2f/hr %s", metrics.TotalCost, gs.getCheckmark(passedBudget))

			// Calculate simulated user count based on request volume
			// Assume ~50 requests per user session, so users ≈ total requests / 50
			estimatedUsers := int(metrics.TotalRequests / 50)
			if estimatedUsers > gs.level.PeakUsers {
				estimatedUsers = gs.level.PeakUsers
			}
			if estimatedUsers < 0 {
				estimatedUsers = 0
			}
			userCountText := fmt.Sprintf("Users: %d / %d (Peak)", estimatedUsers, gs.level.PeakUsers)

			fyne.Do(func() {
				gs.metricsLabel.SetText(metricsText)
				gs.costLabel.SetText(costText)
				gs.userCountLabel.SetText(userCountText)
				gs.summaryLabel.SetText(gs.buildSummaryText())
				gs.testPlanLabel.SetText(gs.buildTestPlanText())
				gs.hintsLabel.SetText(gs.getArchitecturalHints())
			})
		}
	}
}

func (gs *GameScreen) simulateTraffic() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	requestCounter := 0

	for {
		select {
		case <-gs.stopChan:
			return
		case <-ticker.C:
			if !gs.running {
				return
			}

			// Calculate realistic request count based on traffic pattern
			requestCount := 5 // Default fallback
			if gs.trafficGenerator != nil {
				currentRPS := gs.trafficGenerator.CalculateCurrentRPS(time.Now())
				// Convert RPS to requests per 100ms tick
				requestCount = currentRPS / 10
				if requestCount < 1 {
					requestCount = 1
				}
			}

			for i := 0; i < requestCount; i++ {
				requestCounter++

				// Get realistic request type from traffic pattern
				reqType := engine.RequestTypeRead
				if gs.trafficGenerator != nil {
					reqTypeStr := gs.trafficGenerator.GetRequestType()
					switch reqTypeStr {
					case "read":
						reqType = engine.RequestTypeRead
					case "write":
						reqType = engine.RequestTypeWrite
					case "static":
						reqType = engine.RequestTypeAPI
					}
				}

				req := &engine.Request{
					ID:        fmt.Sprintf("req-%d", requestCounter),
					Type:      reqType,
					Timestamp: time.Now(),
					UserID:    fmt.Sprintf("user-%d", requestCounter%1000),
					Region:    "us-east",
					DataSize:  1024,
					Path:      fmt.Sprintf("/data/%d", requestCounter%100),
				}

				gs.gameState.Simulator.SubmitRequest(req)

				// Spawn particles on connections to visualize traffic
				components := gs.canvas.GetComponents()
				for _, comp := range components {
					for _, conn := range comp.Connections {
						// Spawn particle with some randomness to avoid overwhelming
						if requestCounter%3 == 0 {
							gs.canvas.SpawnParticle(conn.From.ID, conn.To.ID)
						}
					}
				}
			}
		}
	}
}

func (gs *GameScreen) animateParticles() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-gs.stopChan:
			return
		case <-ticker.C:
			if !gs.running {
				return
			}

			gs.canvas.UpdateParticles()
		}
	}
}

func (gs *GameScreen) getCheckmark(passed bool) string {
	if passed {
		return "✓"
	}
	return "✗"
}

func (gs *GameScreen) getArchitecturalHints() string {
	hints := "SYSTEM DESIGN PRINCIPLES:\n\n"

	// Analyze current architecture
	hasCache := false
	hasLoadBalancer := false
	hasCDN := false
	hasGateway := false
	hasFirewall := false
	hasDatabase := false
	apiServerCount := 0
	totalComponents := 0

	for _, comp := range gs.canvas.GetComponents() {
		totalComponents++
		switch comp.Type {
		case gui.ComponentTypeCache:
			hasCache = true
		case gui.ComponentTypeLoadBalancer:
			hasLoadBalancer = true
		case gui.ComponentTypeCDN:
			hasCDN = true
		case gui.ComponentTypeGateway:
			hasGateway = true
		case gui.ComponentTypeFirewall:
			hasFirewall = true
		case gui.ComponentTypeDatabase:
			hasDatabase = true
		case gui.ComponentTypeAPIServer:
			apiServerCount++
		}
	}

	// Always show core architecture pattern
	hints += "1. REQUEST FLOW PATTERN\n"
	hints += "   Basic: User → API Server → Database\n"
	hints += "   Production: User → Gateway → Firewall → LB → API → Cache → DB\n\n"

	hints += "   WHY: Defense in Depth + Performance + Scalability\n"
	if totalComponents == 0 {
		hints += "   ⚠ Start by adding components from the toolbox!\n"
		hints += "   Quick start: Add API Server + Database\n"
	} else {
		if !hasGateway {
			hints += "   • Gateway (MISSING ✗): Single entry point.\n"
			hints += "     Handles: Rate limiting, request validation, TLS\n"
		} else {
			hints += "   • Gateway ✓: Entry point secured\n"
		}

		if !hasFirewall {
			hints += "   • Firewall (MISSING ✗): Network security layer.\n"
			hints += "     Handles: DDoS, IP filtering, WAF rules\n"
		} else {
			hints += "   • Firewall ✓: Network security active\n"
		}

		if !hasDatabase {
			hints += "   • Database (MISSING ✗): Persistent storage needed!\n"
			hints += "     Handles: Data persistence, queries, transactions\n"
		} else {
			hints += "   • Database ✓: Data persistence ready\n"
		}
	}

	// Always show scalability principles
	hints += "\n2. HORIZONTAL SCALABILITY\n"
	hints += "   Principle: Scale out (add servers), not up (bigger servers)\n"
	hints += "   WHY: Cost-effective, fault-tolerant, elastic\n\n"

	if !hasLoadBalancer {
		hints += "   • Load Balancer (MISSING ✗):\n"
		hints += "     WHY: Distributes requests across N servers\n"
		hints += "     Enables: Zero-downtime, auto-failover, health checks\n"
		hints += "     Algorithms: Round-robin, least-conn, IP-hash\n"
		hints += "     Real-world: AWS ELB, NGINX, HAProxy\n"
	} else {
		hints += "   • Load Balancer ✓: Ready for horizontal scaling\n"
		if apiServerCount < 2 {
			hints += "   ⚠ Add 2+ API servers for true redundancy!\n"
			hints += "     Pattern: Active-Active (N servers sharing load)\n"
		} else {
			hints += "   • " + fmt.Sprintf("%d", apiServerCount) + " API servers ✓: Redundancy achieved\n"
			hints += "     Load distributed: ~" + fmt.Sprintf("%.0f", 100.0/float64(apiServerCount)) + "% per server\n"
		}
	}

	// Latency optimization - show for all levels with specific targets
	hints += "\n3. LATENCY OPTIMIZATION\n"
	targetLatency := gs.level.Requirements.MaxLatencyP99.Milliseconds()
	hints += fmt.Sprintf("   Target: P99 < %dms\n", targetLatency)
	hints += "   Principle: Every millisecond matters at scale\n\n"

	if !hasCache {
		hints += "   • Cache (MISSING ✗): Critical for low latency!\n"
		hints += "     WHY: Database queries = 10-50ms\n"
		hints += "          Cache reads = 1-2ms (10-50x faster)\n"
		hints += "     Patterns: Cache-aside, write-through, write-behind\n"
		hints += "     Use for: User sessions, hot data, computed results\n"
		hints += "     Tech: Redis, Memcached, ElastiCache\n"
	} else {
		hints += "   • Cache ✓: Latency killer deployed\n"
		hints += "     Expected: 80-95% cache hit rate\n"
		hints += "     Connect: API → Cache (check) → DB (on miss)\n"
	}

	if !hasCDN && targetLatency < 200 {
		hints += "   • CDN (MISSING ✗): Geographic latency matters!\n"
		hints += "     WHY: Edge caching + geographic distribution\n"
		hints += "     Benefit: 50-90% faster for static assets\n"
		hints += "     Handles: Images, CSS, JS, videos\n"
		hints += "     Tech: CloudFront, Cloudflare, Fastly\n"
	} else if hasCDN {
		hints += "   • CDN ✓: Edge caching active\n"
		hints += "     Global PoPs serve content from nearest location\n"
	}

	// High availability - always important
	minUptime := gs.level.Requirements.MinUptime * 100
	hints += "\n4. HIGH AVAILABILITY\n"
	hints += fmt.Sprintf("   Target: %.2f%% uptime ", minUptime)
	if minUptime >= 99.9 {
		hints += "(3 nines = 8.76hr downtime/year)\n"
	} else if minUptime >= 99.0 {
		hints += "(2 nines = 3.65 days downtime/year)\n"
	} else {
		hints += "\n"
	}
	hints += "   Principle: Eliminate Single Points of Failure (SPOFs)\n\n"

	hints += "   Required components:\n"
	if !hasLoadBalancer {
		hints += "   • Load Balancer ✗: Needed for health checks & failover\n"
		hints += "     Detects failures: TCP checks every 5-30s\n"
		hints += "     Auto-failover: Routes traffic to healthy servers\n"
	} else {
		hints += "   • Load Balancer ✓: Health monitoring active\n"
	}

	if apiServerCount < 2 {
		hints += "   • 2+ API Servers ✗: SPOF exists!\n"
		hints += "     WHY: If 1 server fails with no backup = 100% downtime\n"
		hints += "     N servers: Failure of 1 = only 1/N capacity loss\n"
	} else {
		hints += "   • " + fmt.Sprintf("%d", apiServerCount) + " API Servers ✓: No SPOF\n"
		hints += "     Can survive " + fmt.Sprintf("%d", apiServerCount-1) + " server failures\n"
	}

	hints += "\n   Best practices:\n"
	hints += "   • Multi-AZ: Deploy across 2+ availability zones\n"
	hints += "   • Health checks: Heartbeat + deep health endpoint\n"
	hints += "   • Graceful degradation: Serve cached data on DB failure\n"

	// Traffic handling - show for all levels
	hints += "\n5. TRAFFIC HANDLING\n"
	hints += fmt.Sprintf("   Peak load: %d concurrent users\n", gs.level.PeakUsers)
	hints += "   Principle: Design for peak, not average\n\n"

	hints += "   Load distribution:\n"
	hints += "   • Load Balancer: Distribute requests evenly\n"
	hints += "   • Cache: Reduce DB queries by 80-90%\n"
	if apiServerCount > 1 {
		hints += fmt.Sprintf("   • %d servers handle: ~%d users each\n", apiServerCount, gs.level.PeakUsers/apiServerCount)
	}

	hints += "\n   Database scaling:\n"
	hints += "   • Read replicas: Scale read-heavy workloads\n"
	hints += "   • Write sharding: Partition data across DBs\n"
	hints += "   • Connection pooling: Reuse connections (max ~100/server)\n"
	hints += "   • Indexing: B-tree indexes on frequent queries\n"

	// Cost optimization - always relevant
	hints += "\n6. COST OPTIMIZATION\n"
	hints += "   Budget: $" + fmt.Sprintf("%.2f", gs.level.Budget) + "/hr\n"
	hints += "   Principle: Optimize for cost/performance ratio\n\n"

	hints += "   Cost savers:\n"
	hints += "   • Cache: 90% fewer DB queries → smaller DB instance\n"
	hints += "   • CDN: Reduced origin bandwidth costs\n"
	hints += "   • Right-sizing: Don't over-provision\n"
	hints += "   • Auto-scaling: Scale down during low traffic\n"
	hints += "   • Reserved instances: 30-60% discount vs on-demand\n"

	// Common anti-patterns
	hints += "\n7. COMMON MISTAKES TO AVOID\n"
	hints += "   ✗ Single API server (SPOF)\n"
	hints += "   ✗ No caching (high latency + DB overload)\n"
	hints += "   ✗ No load balancer (can't scale horizontally)\n"
	hints += "   ✗ Direct DB access from internet (security risk)\n"
	hints += "   ✗ No monitoring (can't debug issues)\n"

	return hints
}

func (gs *GameScreen) buildSummaryText() string {
	return fmt.Sprintf(
		"VPC: %s (%s) NAT:%v | SG: %s WAF:%v\nDNS/CDN: %s [%s/%s]\nDeploy: %s batch %d HC %ds AS:%v (%d-%d)\nMon: metrics:%v alerts:%v synthetic:%v backups:%v DR:%s",
		gs.networkSettings.vpcPreset,
		gs.networkSettings.region,
		gs.networkSettings.nat,
		gs.networkSettings.sgPreset,
		gs.securitySettings.waf,
		gs.dnsSettings.provider,
		gs.dnsSettings.routingPolicy,
		gs.dnsSettings.edgeScope,
		gs.deploymentSettings.strategy,
		gs.deploymentSettings.batchSize,
		gs.deploymentSettings.healthGraceSecs,
		gs.deploymentSettings.autoScaling,
		gs.deploymentSettings.minInstances,
		gs.deploymentSettings.maxInstances,
		gs.monitoringSettings.metrics,
		gs.monitoringSettings.alerts,
		gs.monitoringSettings.synthetic,
		gs.monitoringSettings.backups,
		gs.monitoringSettings.drRegion,
	)
}

func (gs *GameScreen) buildTestPlanText() string {
	return "Smoke: health endpoints • Synthetic: ping/login • Load: 5x for 60s • Failover: AZ down + SG checks • DNS: flip routing/edge hit • Deploy: canary 10% then full • Backups: hourly snapshot • DR: warm region cutover."
}
