package game

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/javanhut/systemdesignsim/internal/engine"
	"github.com/javanhut/systemdesignsim/internal/gui"
	guicanvas "github.com/javanhut/systemdesignsim/internal/gui/canvas"
)

type TutorialMode string

const (
	ModeDemoWatch TutorialMode = "demo_watch"
	ModePractice  TutorialMode = "practice"
	ModeCompleted TutorialMode = "completed"
)

type TutorialOrchestrator struct {
	pattern      *DesignPattern
	canvas       *guicanvas.GraphCanvas
	currentStep  int
	mode         TutorialMode
	stepTimer    *time.Timer
	componentMap map[string]*gui.VisualComponent
	mu           sync.RWMutex

	animating      bool
	paused         bool
	stopChan       chan bool
	onStepComplete func(step int, total int)
	onTutorialEnd  func()
	onMessage      func(title, description string)

	componentCounter int
}

func NewTutorialOrchestrator(pattern *DesignPattern, canvas *guicanvas.GraphCanvas) *TutorialOrchestrator {
	return &TutorialOrchestrator{
		pattern:      pattern,
		canvas:       canvas,
		currentStep:  0,
		mode:         ModeDemoWatch,
		componentMap: make(map[string]*gui.VisualComponent),
		stopChan:     make(chan bool, 1),
	}
}

func (o *TutorialOrchestrator) SetOnStepComplete(callback func(step int, total int)) {
	o.onStepComplete = callback
}

func (o *TutorialOrchestrator) SetOnTutorialEnd(callback func()) {
	o.onTutorialEnd = callback
}

func (o *TutorialOrchestrator) SetOnMessage(callback func(title, description string)) {
	o.onMessage = callback
}

func (o *TutorialOrchestrator) StartDemo() {
	o.mu.Lock()
	o.mode = ModeDemoWatch
	o.currentStep = 0
	o.animating = true
	o.paused = false
	o.mu.Unlock()

	go o.executeDemo()
}

func (o *TutorialOrchestrator) Pause() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.paused = true
}

func (o *TutorialOrchestrator) Resume() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.paused = false
}

func (o *TutorialOrchestrator) Stop() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.animating = false
	o.paused = false
	select {
	case o.stopChan <- true:
	default:
	}
}

func (o *TutorialOrchestrator) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.currentStep = 0
	o.animating = false
	o.paused = false
	o.componentMap = make(map[string]*gui.VisualComponent)
	o.componentCounter = 0

	for _, comp := range o.canvas.GetComponents() {
		o.canvas.RemoveComponent(comp.ID)
	}
}

func (o *TutorialOrchestrator) executeDemo() {
	for o.currentStep < len(o.pattern.DemoSteps) {
		o.mu.RLock()
		if !o.animating {
			o.mu.RUnlock()
			return
		}
		if o.paused {
			o.mu.RUnlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}
		step := o.pattern.DemoSteps[o.currentStep]
		o.mu.RUnlock()

		select {
		case <-o.stopChan:
			return
		default:
		}

		o.executeStep(&step)

		o.mu.Lock()
		o.currentStep++
		if o.onStepComplete != nil {
			fyne.Do(func() {
				o.onStepComplete(o.currentStep, len(o.pattern.DemoSteps))
			})
		}
		o.mu.Unlock()
	}

	o.mu.Lock()
	o.animating = false
	o.mode = ModeCompleted
	o.mu.Unlock()

	if o.onTutorialEnd != nil {
		fyne.Do(func() {
			o.onTutorialEnd()
		})
	}
}

func (o *TutorialOrchestrator) executeStep(step *TutorialStep) {
	switch step.Type {
	case StepMessage:
		o.showMessage(step)
	case StepAddComponent:
		o.animateComponentAddition(step)
	case StepCreateConnection:
		o.animateConnection(step)
	case StepShowTraffic:
		o.animateTraffic(step)
	case StepHighlight:
		o.highlightArea(step)
	case StepWait:
		time.Sleep(step.Duration)
	}
}

func (o *TutorialOrchestrator) showMessage(step *TutorialStep) {
	if o.onMessage != nil {
		fyne.Do(func() {
			o.onMessage(step.Title, step.Description)
		})
	}

	if step.Duration > 0 {
		time.Sleep(step.Duration)
	}
}

func (o *TutorialOrchestrator) animateComponentAddition(step *TutorialStep) {
	o.componentCounter++
	id := fmt.Sprintf("%s-%d", step.ComponentID, o.componentCounter)

	var compType gui.ComponentType
	switch step.ComponentType {
	case "api-server":
		compType = gui.ComponentTypeAPIServer
	case "database":
		compType = gui.ComponentTypeDatabase
	case "cache":
		compType = gui.ComponentTypeCache
	case "load-balancer":
		compType = gui.ComponentTypeLoadBalancer
	case "cdn":
		compType = gui.ComponentTypeCDN
	case "gateway":
		compType = gui.ComponentTypeGateway
	case "firewall":
		compType = gui.ComponentTypeFirewall
	case "nat":
		compType = gui.ComponentTypeNAT
	case "router":
		compType = gui.ComponentTypeRouter
	default:
		compType = gui.ComponentTypeAPIServer
	}

	visualComp := gui.NewVisualComponent(id, compType, step.Position)

	var comp engine.Component
	switch compType {
	case gui.ComponentTypeAPIServer:
		comp = o.createAPIServer(id)
	case gui.ComponentTypeDatabase:
		comp = o.createDatabase(id)
	case gui.ComponentTypeCache:
		comp = o.createCache(id)
	case gui.ComponentTypeLoadBalancer:
		comp = o.createLoadBalancer(id)
	case gui.ComponentTypeCDN:
		comp = o.createCDN(id)
	case gui.ComponentTypeGateway:
		comp = o.createGateway(id)
	case gui.ComponentTypeFirewall:
		comp = o.createFirewall(id)
	case gui.ComponentTypeNAT:
		comp = o.createNAT(id)
	case gui.ComponentTypeRouter:
		comp = o.createRouter(id)
	}

	visualComp.SetComponent(comp)

	o.mu.Lock()
	o.componentMap[step.ComponentID] = visualComp
	o.mu.Unlock()

	fyne.Do(func() {
		o.canvas.AddComponent(visualComp)
	})

	if step.FadeIn && step.Duration > 0 {
		time.Sleep(step.Duration)
	}
}

func (o *TutorialOrchestrator) animateConnection(step *TutorialStep) {
	o.mu.RLock()
	fromComp, fromExists := o.componentMap[step.FromID]
	toComp, toExists := o.componentMap[step.ToID]
	o.mu.RUnlock()

	if !fromExists || !toExists {
		fmt.Printf("Warning: Cannot create connection from %s to %s - components not found\n", step.FromID, step.ToID)
		return
	}

	fyne.Do(func() {
		o.canvas.AddConnection(fromComp, toComp)
	})

	if step.ShowParticles {
		time.Sleep(500 * time.Millisecond)
		fyne.Do(func() {
			o.canvas.SpawnParticle(fromComp.ID, toComp.ID)
		})
	}

	if step.Duration > 0 {
		time.Sleep(step.Duration)
	}
}

func (o *TutorialOrchestrator) animateTraffic(step *TutorialStep) {
	particleCount := step.ParticleCount
	if particleCount == 0 {
		particleCount = 5
	}

	duration := step.Duration
	if duration == 0 {
		duration = 2 * time.Second
	}

	interval := duration / time.Duration(particleCount)

	for i := 0; i < particleCount; i++ {
		select {
		case <-o.stopChan:
			return
		default:
		}

		o.mu.RLock()
		if !o.animating || o.paused {
			o.mu.RUnlock()
			return
		}
		o.mu.RUnlock()

		components := o.getAllTutorialComponents()
		if len(components) > 0 {
			for _, comp := range components {
				if len(comp.Connections) > 0 {
					for _, conn := range comp.Connections {
						fyne.Do(func() {
							o.canvas.SpawnParticle(conn.From.ID, conn.To.ID)
						})
					}
				}
			}
		}

		time.Sleep(interval)
	}
}

func (o *TutorialOrchestrator) highlightArea(step *TutorialStep) {
	if step.Duration > 0 {
		time.Sleep(step.Duration)
	}
}

func (o *TutorialOrchestrator) getAllTutorialComponents() []*gui.VisualComponent {
	o.mu.RLock()
	defer o.mu.RUnlock()

	components := make([]*gui.VisualComponent, 0, len(o.componentMap))
	for _, comp := range o.componentMap {
		components = append(components, comp)
	}
	return components
}

func (o *TutorialOrchestrator) ValidatePracticeStep(stepIndex int, userComponents []*gui.VisualComponent) (bool, string) {
	if stepIndex < 0 || stepIndex >= len(o.pattern.PracticeSteps) {
		return false, "Invalid step index"
	}

	step := o.pattern.PracticeSteps[stepIndex]
	validation := step.Expected

	componentCounts := make(map[string]int)
	for _, comp := range userComponents {
		compTypeStr := string(comp.Type)
		componentCounts[compTypeStr]++
	}

	for requiredType, requiredCount := range validation.RequiredComponents {
		actualCount := componentCounts[requiredType]
		if actualCount < requiredCount {
			missing := requiredCount - actualCount
			return false, fmt.Sprintf("Need %d more %s component(s)", missing, requiredType)
		}
	}

	if validation.MinComponents > 0 && len(userComponents) < validation.MinComponents {
		return false, fmt.Sprintf("Need at least %d total components", validation.MinComponents)
	}

	if len(validation.RequiredConnections) > 0 {
		for _, reqConn := range validation.RequiredConnections {
			if !o.hasConnectionOfType(userComponents, reqConn.FromType, reqConn.ToType) {
				return false, fmt.Sprintf("Need connection: %s → %s", reqConn.FromType, reqConn.ToType)
			}
		}
	}

	return true, "Perfect! Step completed successfully."
}

func (o *TutorialOrchestrator) hasConnectionOfType(components []*gui.VisualComponent, fromType, toType string) bool {
	for _, comp := range components {
		if string(comp.Type) == fromType {
			for _, conn := range comp.Connections {
				if string(conn.To.Type) == toType {
					return true
				}
			}
		}
	}
	return false
}

func (o *TutorialOrchestrator) GetCurrentStep() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.currentStep
}

func (o *TutorialOrchestrator) GetTotalSteps() int {
	return len(o.pattern.DemoSteps)
}

func (o *TutorialOrchestrator) GetProgress() float64 {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if len(o.pattern.DemoSteps) == 0 {
		return 0
	}
	return float64(o.currentStep) / float64(len(o.pattern.DemoSteps))
}

func (o *TutorialOrchestrator) IsAnimating() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.animating
}

func (o *TutorialOrchestrator) IsPaused() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.paused
}

func (o *TutorialOrchestrator) GetMode() TutorialMode {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.mode
}

func (o *TutorialOrchestrator) SetMode(mode TutorialMode) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.mode = mode
}

func (o *TutorialOrchestrator) createAPIServer(id string) engine.Component {
	return &mockComponent{id: id, componentType: "api-server"}
}

func (o *TutorialOrchestrator) createDatabase(id string) engine.Component {
	return &mockComponent{id: id, componentType: "database"}
}

func (o *TutorialOrchestrator) createCache(id string) engine.Component {
	return &mockComponent{id: id, componentType: "cache"}
}

func (o *TutorialOrchestrator) createLoadBalancer(id string) engine.Component {
	return &mockComponent{id: id, componentType: "load-balancer"}
}

func (o *TutorialOrchestrator) createCDN(id string) engine.Component {
	return &mockComponent{id: id, componentType: "cdn"}
}

func (o *TutorialOrchestrator) createGateway(id string) engine.Component {
	return &mockComponent{id: id, componentType: "gateway"}
}

func (o *TutorialOrchestrator) createFirewall(id string) engine.Component {
	return &mockComponent{id: id, componentType: "firewall"}
}

func (o *TutorialOrchestrator) createNAT(id string) engine.Component {
	return &mockComponent{id: id, componentType: "nat"}
}

func (o *TutorialOrchestrator) createRouter(id string) engine.Component {
	return &mockComponent{id: id, componentType: "router"}
}

type mockComponent struct {
	id            string
	componentType string
}

func (m *mockComponent) GetID() string                    { return m.id }
func (m *mockComponent) GetType() string                  { return m.componentType }
func (m *mockComponent) GetRegion() string                { return "us-east" }
func (m *mockComponent) IsHealthy() bool                  { return true }
func (m *mockComponent) SetHealthy(healthy bool)          {}
func (m *mockComponent) GetCost() float64                 { return 0.0 }
func (m *mockComponent) Process(req *engine.Request) (*engine.Response, error) {
	return &engine.Response{RequestID: req.ID, Success: true}, nil
}
func (m *mockComponent) GetMetrics() *engine.Metrics {
	return &engine.Metrics{RequestCount: 0, Throughput: 0, ErrorRate: 0}
}
