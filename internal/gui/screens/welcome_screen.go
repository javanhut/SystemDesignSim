package screens

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/gui"
)

type WelcomeScreen struct {
	window      fyne.Window
	currentPage int
	totalPages  int
	content     *fyne.Container
	pageLabel   *widget.Label
	prevButton  *widget.Button
	nextButton  *widget.Button
}

func NewWelcomeScreen(window fyne.Window) *WelcomeScreen {
	ws := &WelcomeScreen{
		window:      window,
		currentPage: 0,
		totalPages:  7,
	}
	return ws
}

func (ws *WelcomeScreen) Build() fyne.CanvasObject {
	ws.content = container.NewVBox()
	ws.updateContent()

	ws.pageLabel = widget.NewLabel(fmt.Sprintf("Page %d of %d", ws.currentPage+1, ws.totalPages))
	ws.pageLabel.Alignment = fyne.TextAlignCenter

	ws.prevButton = widget.NewButton("Previous", func() {
		if ws.currentPage > 0 {
			ws.currentPage--
			ws.updateContent()
			ws.updateNavigation()
		}
	})

	ws.nextButton = widget.NewButton("Next", func() {
		if ws.currentPage < ws.totalPages-1 {
			ws.currentPage++
			ws.updateContent()
			ws.updateNavigation()
		} else {
			ws.finishTutorial()
		}
	})

	skipButton := widget.NewButton("Skip Tutorial", func() {
		ws.finishTutorial()
	})

	ws.updateNavigation()

	navigation := container.NewBorder(
		nil,
		nil,
		ws.prevButton,
		ws.nextButton,
		container.NewCenter(skipButton),
	)

	main := container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewSeparator(),
			ws.pageLabel,
			navigation,
		),
		nil,
		nil,
		container.NewVScroll(ws.content),
	)

	return main
}

func (ws *WelcomeScreen) updateNavigation() {
	ws.pageLabel.SetText(fmt.Sprintf("Page %d of %d", ws.currentPage+1, ws.totalPages))

	if ws.currentPage == 0 {
		ws.prevButton.Disable()
	} else {
		ws.prevButton.Enable()
	}

	if ws.currentPage == ws.totalPages-1 {
		ws.nextButton.SetText("Get Started!")
	} else {
		ws.nextButton.SetText("Next")
	}
}

func (ws *WelcomeScreen) updateContent() {
	ws.content.Objects = nil

	switch ws.currentPage {
	case 0:
		ws.content.Add(ws.createWelcomePage())
	case 1:
		ws.content.Add(ws.createScenarioPage())
	case 2:
		ws.content.Add(ws.createComponentsPage1())
	case 3:
		ws.content.Add(ws.createComponentsPage2())
	case 4:
		ws.content.Add(ws.createComponentsPage3())
	case 5:
		ws.content.Add(ws.createHowToPlayPage())
	case 6:
		ws.content.Add(ws.createScoringPage())
	}

	ws.content.Refresh()
}

func (ws *WelcomeScreen) createWelcomePage() fyne.CanvasObject {
	title := widget.NewLabel("Welcome to System Design Simulator!")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	intro := widget.NewLabel(
		"Learn distributed systems through hands-on simulation.\n\n" +
			"This game teaches you real-world system design by letting you build,\n" +
			"scale, and optimize infrastructure for different scenarios.",
	)
	intro.Wrapping = fyne.TextWrapWord

	whatYouLearn := widget.NewLabel("What You'll Learn:")
	whatYouLearn.TextStyle = fyne.TextStyle{Bold: true}

	topics := widget.NewLabel(
		"• Scaling strategies (horizontal vs vertical)\n" +
			"• Caching and performance optimization\n" +
			"• Load balancing and fault tolerance\n" +
			"• Global distribution with CDN\n" +
			"• Database sharding and replication\n" +
			"• Cost optimization in cloud infrastructure\n" +
			"• Real-world system design trade-offs",
	)

	howItWorks := widget.NewLabel(
		"\nHow It Works:\n" +
			"1. Choose a level with a specific scenario\n" +
			"2. Build your infrastructure by adding components\n" +
			"3. Connect components to create data flow\n" +
			"4. Run the simulation and watch real-time metrics\n" +
			"5. Optimize and submit your solution for scoring",
	)
	howItWorks.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		intro,
		widget.NewSeparator(),
		whatYouLearn,
		topics,
		howItWorks,
	)
}

func (ws *WelcomeScreen) createScenarioPage() fyne.CanvasObject {
	title := widget.NewLabel("Understanding the Scenarios")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	intro := widget.NewLabel(
		"Each level presents a real-world scenario where you need to build\n" +
			"infrastructure to handle specific user loads and requirements.\n",
	)
	intro.Wrapping = fyne.TextWrapWord

	scenarios := []struct {
		level       string
		scenario    string
		challenge   string
		learn       string
		requirement string
	}{
		{
			"Level 1 - Local Blog",
			"You've built a personal blog for friends",
			"Handle 10 concurrent readers",
			"Basic server + database setup",
			"Max latency: 500ms, Uptime: 95%, Budget: $10",
		},
		{
			"Level 2 - Growing Blog",
			"Your blog went viral on social media",
			"Scale from 10 to 100 concurrent users",
			"Load balancing and horizontal scaling",
			"Max latency: 300ms, Uptime: 98%, Budget: $50",
		},
		{
			"Level 3 - Regional Social Network",
			"Launch a Twitter-like app for your city",
			"1,000 users posting and reading",
			"Caching, database replication, redundancy",
			"Max latency: 200ms, Uptime: 99%, Budget: $200",
		},
		{
			"Level 4 - Global E-commerce",
			"Online store shipping worldwide",
			"Handle 10,000 concurrent shoppers",
			"Multi-region deployment, CDN, sharding",
			"Max latency: 150ms, Uptime: 99.5%, Budget: $1000",
		},
		{
			"Level 5 - Viral Streaming Service",
			"Your app is the next Netflix",
			"100,000 concurrent streamers globally",
			"Advanced optimization, five nines uptime",
			"Max latency: 100ms, Uptime: 99.95%, Budget: $5000",
		},
	}

	content := container.NewVBox(title, widget.NewSeparator(), intro)

	for _, s := range scenarios {
		levelTitle := widget.NewLabel(s.level)
		levelTitle.TextStyle = fyne.TextStyle{Bold: true}

		details := widget.NewLabel(
			fmt.Sprintf("Scenario: %s\nChallenge: %s\nLearn: %s\nRequirements: %s",
				s.scenario, s.challenge, s.learn, s.requirement),
		)
		details.Wrapping = fyne.TextWrapWord

		content.Add(levelTitle)
		content.Add(details)
		content.Add(widget.NewSeparator())
	}

	return content
}

func (ws *WelcomeScreen) createComponentsPage1() fyne.CanvasObject {
	title := widget.NewLabel("Infrastructure Components - Part 1")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	apiTitle := widget.NewLabel("API SERVER")
	apiTitle.TextStyle = fyne.TextStyle{Bold: true}

	apiDesc := widget.NewLabel(
		"Handles incoming user requests and processes business logic.\n\n" +
			"Details:\n" +
			"• Sizes: Small (10 concurrent), Medium (50), Large (200), XLarge (500)\n" +
			"• Processing Time: ~10-15ms per request\n" +
			"• Cost: $0.05 - $0.40/hour depending on size\n" +
			"• Connects To: Database (for data), Cache (for fast reads)\n" +
			"• Use When: You need to handle user requests\n\n" +
			"Capacity Matters: If you send more requests than the server can handle,\n" +
			"requests will fail! Scale horizontally with a load balancer.",
	)
	apiDesc.Wrapping = fyne.TextWrapWord

	dbTitle := widget.NewLabel("DATABASE")
	dbTitle.TextStyle = fyne.TextStyle{Bold: true}

	dbDesc := widget.NewLabel(
		"Stores persistent data permanently.\n\n" +
			"Details:\n" +
			"• Types: SQL, NoSQL, Key-Value, Document\n" +
			"• Latency: 10ms read, 15ms write (slower than cache!)\n" +
			"• Cost: $0.05/hour + storage costs\n" +
			"• Connects From: API Server (writes), Cache (fallback)\n" +
			"• Use When: You need permanent data storage\n\n" +
			"Advanced Features:\n" +
			"• Sharding: Split data across multiple databases for scale\n" +
			"• Replication: Copy data to multiple databases for redundancy",
	)
	dbDesc.Wrapping = fyne.TextWrapWord

	example := widget.NewLabel(
		"Typical Connection:\nAPI Server → Database (for data persistence)",
	)
	example.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		apiTitle,
		apiDesc,
		widget.NewSeparator(),
		dbTitle,
		dbDesc,
		widget.NewSeparator(),
		example,
	)
}

func (ws *WelcomeScreen) createComponentsPage2() fyne.CanvasObject {
	title := widget.NewLabel("Infrastructure Components - Part 2")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	cacheTitle := widget.NewLabel("CACHE")
	cacheTitle.TextStyle = fyne.TextStyle{Bold: true}

	cacheDesc := widget.NewLabel(
		"Speeds up reads with in-memory storage (10x faster than database!).\n\n" +
			"Details:\n" +
			"• Policies: LRU (Least Recently Used), LFU (Least Frequently), FIFO\n" +
			"• Latency: 1-2ms (vs 10ms for database)\n" +
			"• TTL: Data expires after timeout\n" +
			"• Cost: $0.02/hour + memory costs\n" +
			"• Connects To: Database (fallback on cache miss)\n" +
			"• Use When: High read traffic, repeated data access\n\n" +
			"How It Works:\n" +
			"1. API checks cache first (fast!)\n" +
			"2. If found (HIT): Return immediately\n" +
			"3. If not found (MISS): Get from database, store in cache\n" +
			"4. Aim for 70%+ cache hit rate for best performance",
	)
	cacheDesc.Wrapping = fyne.TextWrapWord

	lbTitle := widget.NewLabel("LOAD BALANCER")
	lbTitle.TextStyle = fyne.TextStyle{Bold: true}

	lbDesc := widget.NewLabel(
		"Distributes traffic across multiple servers for horizontal scaling.\n\n" +
			"Details:\n" +
			"• Strategies: Round-robin (equal distribution), Least-connected (send to least busy)\n" +
			"• Latency: 2ms overhead\n" +
			"• Cost: $0.025/hour\n" +
			"• Connects To: Multiple API Servers\n" +
			"• Use When: One server isn't enough\n\n" +
			"Benefits:\n" +
			"• High Availability: If one server fails, others continue\n" +
			"• Scalability: Add more servers to handle more traffic\n" +
			"• Even Distribution: No single server gets overloaded",
	)
	lbDesc.Wrapping = fyne.TextWrapWord

	example := widget.NewLabel(
		"Scaling Pattern:\n" +
			"Load Balancer → API Server 1\n" +
			"               → API Server 2 (distributes load evenly)\n" +
			"               → API Server 3\n\n" +
			"Caching Pattern:\n" +
			"API Server → Cache → Database (cache-aside pattern)",
	)
	example.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		cacheTitle,
		cacheDesc,
		widget.NewSeparator(),
		lbTitle,
		lbDesc,
		widget.NewSeparator(),
		example,
	)
}

func (ws *WelcomeScreen) createComponentsPage3() fyne.CanvasObject {
	title := widget.NewLabel("Infrastructure Components - Part 3")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	cdnTitle := widget.NewLabel("CDN (Content Delivery Network)")
	cdnTitle.TextStyle = fyne.TextStyle{Bold: true}

	cdnDesc := widget.NewLabel(
		"Caches content at edge locations close to users globally.\n\n" +
			"Details:\n" +
			"• Edge Locations: US-East, US-West, Europe, Asia, Australia\n" +
			"• Latency: 2ms for cache hit vs 200ms cross-region!\n" +
			"• Cost: $0.08/hour + per-region fees\n" +
			"• Connects To: Origin server (your API/backend)\n" +
			"• Use When: Global users, static content, high latency issues\n\n" +
			"Why CDN Matters - Regional Latency:\n" +
			"• US-East ↔ US-West: 70ms\n" +
			"• US-East ↔ Europe: 80ms\n" +
			"• US-East ↔ Asia: 180ms\n" +
			"• US-East ↔ Australia: 200ms\n\n" +
			"With CDN: Users get content from nearby edge (2ms)\n" +
			"Without CDN: Users wait for cross-region latency (200ms)\n\n" +
			"That's 100x faster!",
	)
	cdnDesc.Wrapping = fyne.TextWrapWord

	networkTitle := widget.NewLabel("NETWORK & LATENCY")
	networkTitle.TextStyle = fyne.TextStyle{Bold: true}

	networkDesc := widget.NewLabel(
		"Understanding network performance is key to system design.\n\n" +
			"Latency Sources:\n" +
			"• Component Processing: 10-15ms (API, Database)\n" +
			"• Caching: 1-2ms (Cache, CDN edge)\n" +
			"• Network Distance: 5-200ms (depends on regions)\n" +
			"• Load Balancer: 2ms overhead\n\n" +
			"Optimization Tips:\n" +
			"1. Use caching to avoid slow database reads\n" +
			"2. Deploy close to users (multi-region)\n" +
			"3. Use CDN for global content delivery\n" +
			"4. Minimize hops between components\n" +
			"5. Monitor P99 latency (99th percentile)",
	)
	networkDesc.Wrapping = fyne.TextWrapWord

	example := widget.NewLabel(
		"CDN Architecture:\n" +
			"          ┌─ CDN Edge (Europe) ─┐\n" +
			"User  →   ├─ CDN Edge (Asia)   ─┤  (cache hit = 2ms)\n" +
			"          └─ CDN Edge (US)     ─┘\n" +
			"                    ↓ (cache miss)\n" +
			"               Origin Server (your backend)",
	)
	example.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		cdnTitle,
		cdnDesc,
		widget.NewSeparator(),
		networkTitle,
		networkDesc,
		widget.NewSeparator(),
		example,
	)
}

func (ws *WelcomeScreen) createHowToPlayPage() fyne.CanvasObject {
	title := widget.NewLabel("How to Play - Controls & Interactions")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	addingTitle := widget.NewLabel("ADDING COMPONENTS:")
	addingTitle.TextStyle = fyne.TextStyle{Bold: true}

	addingSteps := widget.NewLabel(
		"1. Click component button in left toolbox\n" +
			"2. Component appears on canvas\n" +
			"3. Components automatically get unique IDs",
	)

	connectingTitle := widget.NewLabel("CONNECTING COMPONENTS:")
	connectingTitle.TextStyle = fyne.TextStyle{Bold: true}

	connectingSteps := widget.NewLabel(
		"1. Right-click on source component (e.g., API Server)\n" +
			"2. Drag to target component (e.g., Database)\n" +
			"3. Release to create connection\n" +
			"4. See animated particles flow along connection!\n\n" +
			"Valid Connections:\n" +
			"✓ Load Balancer → API Server (distribute traffic)\n" +
			"✓ API Server → Database (store data)\n" +
			"✓ API Server → Cache (speed up reads)\n" +
			"✓ Cache → Database (fallback on miss)\n" +
			"✓ CDN → API Server (edge caching)",
	)
	connectingSteps.Wrapping = fyne.TextWrapWord

	runningTitle := widget.NewLabel("RUNNING SIMULATION:")
	runningTitle.TextStyle = fyne.TextStyle{Bold: true}

	runningSteps := widget.NewLabel(
		"1. Build your architecture with components and connections\n" +
			"2. Click 'Start Simulation' button\n" +
			"3. Watch real-time metrics:\n" +
			"   • Request count and success rate\n" +
			"   • Latency (P99 - 99th percentile)\n" +
			"   • Cost accumulation\n" +
			"   • Component health (color indicators)\n" +
			"4. Monitor component colors for health status\n" +
			"5. Stop simulation when ready to submit",
	)
	runningSteps.Wrapping = fyne.TextWrapWord

	healthTitle := widget.NewLabel("COMPONENT HEALTH COLORS:")
	healthTitle.TextStyle = fyne.TextStyle{Bold: true}

	healthInfo := widget.NewLabel(
		"Green:  Healthy (< 50% capacity) - All good!\n" +
			"Yellow: Warning (50-80% capacity) - Getting busy\n" +
			"Orange: Critical (> 80% capacity) - Almost overloaded\n" +
			"Red:    Down (failing requests) - Add more capacity!",
	)

	submittingTitle := widget.NewLabel("SUBMITTING SOLUTION:")
	submittingTitle.TextStyle = fyne.TextStyle{Bold: true}

	submittingSteps := widget.NewLabel(
		"1. Run simulation until metrics stabilize\n" +
			"2. Click 'Stop Simulation'\n" +
			"3. Click 'Submit Solution'\n" +
			"4. See your score and detailed feedback!\n" +
			"5. Replay for higher scores or move to next level",
	)

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		addingTitle,
		addingSteps,
		widget.NewSeparator(),
		connectingTitle,
		connectingSteps,
		widget.NewSeparator(),
		runningTitle,
		runningSteps,
		widget.NewSeparator(),
		healthTitle,
		healthInfo,
		widget.NewSeparator(),
		submittingTitle,
		submittingSteps,
	)
}

func (ws *WelcomeScreen) createScoringPage() fyne.CanvasObject {
	title := widget.NewLabel("Level Progression & Scoring")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	reqTitle := widget.NewLabel("REQUIREMENTS (Must Pass):")
	reqTitle.TextStyle = fyne.TextStyle{Bold: true}

	reqDesc := widget.NewLabel(
		"Each level has mandatory requirements you must meet:\n\n" +
			"✓ Max Latency (P99): Keep response times below threshold\n" +
			"✓ Min Uptime: Percentage of successful requests\n" +
			"✓ Max Error Rate: Limit failed requests\n" +
			"✓ Budget: Stay within cost constraints\n" +
			"✓ Architecture: Some levels require specific components\n" +
			"   (e.g., Load Balancer required for Level 2+)\n\n" +
			"Fail ANY requirement = Level Failed",
	)
	reqDesc.Wrapping = fyne.TextWrapWord

	bonusTitle := widget.NewLabel("BONUS OBJECTIVES (Extra Points):")
	bonusTitle.TextStyle = fyne.TextStyle{Bold: true}

	bonusDesc := widget.NewLabel(
		"Exceed targets for bonus points:\n\n" +
			"• Excellent uptime (above target): +100 points\n" +
			"• Fast response time (below target): +100 points\n" +
			"• Low error rate (below target): +100 points\n" +
			"• Cost efficient (under target budget): +150 points\n" +
			"• Great cache hit rate (>70%): +50 points\n" +
			"• Cost savings (far under budget): up to +200 points",
	)
	bonusDesc.Wrapping = fyne.TextWrapWord

	scoringTitle := widget.NewLabel("SCORING FORMULA:")
	scoringTitle.TextStyle = fyne.TextStyle{Bold: true}

	scoringDesc := widget.NewLabel(
		"Base Score: 1000 points\n" +
			"Penalties: -200 per failed requirement\n" +
			"Bonuses: Up to +700 for excellence\n" +
			"Final Score: 0 to 1700 points\n\n" +
			"Example:\n" +
			"Base: 1000\n" +
			"Excellent uptime: +100\n" +
			"Fast latency: +100\n" +
			"Cost efficient: +150\n" +
			"Final Score: 1350 points!",
	)

	progressTitle := widget.NewLabel("PROGRESSION:")
	progressTitle.TextStyle = fyne.TextStyle{Bold: true}

	progressDesc := widget.NewLabel(
		"• Complete Level N to unlock Level N+1\n" +
			"• Replay levels for higher scores\n" +
			"• Each level teaches new concepts\n" +
			"• Difficulty increases progressively\n" +
			"• Best scores are saved",
	)

	tipsTitle := widget.NewLabel("TIPS FOR SUCCESS:")
	tipsTitle.TextStyle = fyne.TextStyle{Bold: true}

	tipsDesc := widget.NewLabel(
		"1. Start simple - don't over-engineer Level 1\n" +
			"2. Use caching for read-heavy workloads\n" +
			"3. Add load balancers before horizontal scaling\n" +
			"4. Monitor component health colors during simulation\n" +
			"5. Balance cost vs performance - cheapest isn't always best\n" +
			"6. Test your architecture before submitting!\n" +
			"7. Read the level description to understand the scenario",
	)
	tipsDesc.Wrapping = fyne.TextWrapWord

	ready := widget.NewLabel(
		"\n\nYou're ready to start!\n" +
			"Click 'Get Started' to begin with Level 1: Local Blog",
	)
	ready.Alignment = fyne.TextAlignCenter
	ready.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		reqTitle,
		reqDesc,
		widget.NewSeparator(),
		bonusTitle,
		bonusDesc,
		widget.NewSeparator(),
		scoringTitle,
		scoringDesc,
		widget.NewSeparator(),
		progressTitle,
		progressDesc,
		widget.NewSeparator(),
		tipsTitle,
		tipsDesc,
		ready,
	)
}

func (ws *WelcomeScreen) finishTutorial() {
	prefs, _ := gui.LoadPreferences()
	prefs.TutorialCompleted = true
	prefs.FirstLaunch = false
	gui.SavePreferences(prefs)

	ws.window.SetContent(NewLevelSelectScreen(ws.window).Build())
}
