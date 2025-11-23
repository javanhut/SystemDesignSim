package canvas

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/gui"
)

type GraphCanvas struct {
	widget.BaseWidget

	components      []*gui.VisualComponent
	connections     []*gui.Connection
	componentsMutex sync.RWMutex

	draggingComponent *gui.VisualComponent
	connectingFrom    *gui.VisualComponent
	connectingTo      fyne.Position

	onComponentClick func(*gui.VisualComponent)
	onComponentAdd   func(*gui.VisualComponent)
	onConnectionAdd  func(*gui.Connection)
}

func NewGraphCanvas() *GraphCanvas {
	gc := &GraphCanvas{
		components:  make([]*gui.VisualComponent, 0),
		connections: make([]*gui.Connection, 0),
	}
	gc.ExtendBaseWidget(gc)
	return gc
}

func (gc *GraphCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &graphCanvasRenderer{
		canvas:  gc,
		objects: []fyne.CanvasObject{},
	}
}

func (gc *GraphCanvas) AddComponent(comp *gui.VisualComponent) {
	gc.componentsMutex.Lock()
	gc.components = append(gc.components, comp)
	gc.componentsMutex.Unlock()

	if gc.onComponentAdd != nil {
		gc.onComponentAdd(comp)
	}
	gc.Refresh()
}

func (gc *GraphCanvas) RemoveComponent(id string) {
	gc.componentsMutex.Lock()

	for i, comp := range gc.components {
		if comp.ID == id {
			for j := len(gc.connections) - 1; j >= 0; j-- {
				conn := gc.connections[j]
				if conn.From.ID == id || conn.To.ID == id {
					gc.connections = append(gc.connections[:j], gc.connections[j+1:]...)
				}
			}

			gc.components = append(gc.components[:i], gc.components[i+1:]...)
			gc.componentsMutex.Unlock()
			gc.Refresh()
			return
		}
	}
	gc.componentsMutex.Unlock()
}

func (gc *GraphCanvas) AddConnection(from, to *gui.VisualComponent) {
	gc.componentsMutex.Lock()
	connID := fmt.Sprintf("conn-%s-%s", from.ID, to.ID)
	conn := gui.NewConnection(connID, from, to)
	gc.connections = append(gc.connections, conn)

	from.AddConnection(conn)
	to.AddConnection(conn)
	gc.componentsMutex.Unlock()

	if gc.onConnectionAdd != nil {
		gc.onConnectionAdd(conn)
	}
	gc.Refresh()
}

func (gc *GraphCanvas) GetComponents() []*gui.VisualComponent {
	gc.componentsMutex.RLock()
	defer gc.componentsMutex.RUnlock()
	
	comps := make([]*gui.VisualComponent, len(gc.components))
	copy(comps, gc.components)
	return comps
}

func (gc *GraphCanvas) GetComponentAt(pos fyne.Position) *gui.VisualComponent {
	gc.componentsMutex.RLock()
	defer gc.componentsMutex.RUnlock()

	for i := len(gc.components) - 1; i >= 0; i-- {
		if gc.components[i].Contains(pos) {
			return gc.components[i]
		}
	}
	return nil
}

func (gc *GraphCanvas) SetOnComponentClick(callback func(*gui.VisualComponent)) {
	gc.onComponentClick = callback
}

func (gc *GraphCanvas) SetOnComponentAdd(callback func(*gui.VisualComponent)) {
	gc.onComponentAdd = callback
}

func (gc *GraphCanvas) SetOnConnectionAdd(callback func(*gui.Connection)) {
	gc.onConnectionAdd = callback
}

// SpawnParticle adds a particle to a connection between two components
func (gc *GraphCanvas) SpawnParticle(fromID, toID string) {
	gc.componentsMutex.RLock()
	defer gc.componentsMutex.RUnlock()

	for _, conn := range gc.connections {
		if conn.From.ID == fromID && conn.To.ID == toID {
			conn.AddParticle()
			fyne.Do(func() {
				gc.Refresh()
			})
			return
		}
	}
}

// UpdateParticles updates all particles on all connections
func (gc *GraphCanvas) UpdateParticles() {
	gc.componentsMutex.RLock()
	defer gc.componentsMutex.RUnlock()

	for _, conn := range gc.connections {
		conn.UpdateParticles()
	}
	fyne.Do(func() {
		gc.Refresh()
	})
}

func (gc *GraphCanvas) Tapped(ev *fyne.PointEvent) {
	comp := gc.GetComponentAt(ev.Position)

	if gc.connectingFrom != nil {
		if comp != nil && comp != gc.connectingFrom {
			gc.AddConnection(gc.connectingFrom, comp)
		}
		gc.connectingFrom = nil
		gc.Refresh()
		return
	}

	if comp != nil && comp.Selected {
		// Second tap on the same node enters connection mode.
		gc.connectingFrom = comp
		gc.connectingTo = comp.GetCenter()
		gc.Refresh()
		return
	}

	gc.componentsMutex.Lock()
	for _, c := range gc.components {
		c.Selected = false
	}
	gc.componentsMutex.Unlock()

	if comp != nil {
		comp.Selected = true
		if gc.onComponentClick != nil {
			gc.onComponentClick(comp)
		}
	}

	gc.Refresh()
}

func (gc *GraphCanvas) TappedSecondary(ev *fyne.PointEvent) {
	comp := gc.GetComponentAt(ev.Position)
	if comp != nil {
		gc.connectingFrom = comp
		gc.connectingTo = ev.Position
		gc.Refresh()
	}
}

func (gc *GraphCanvas) Dragged(ev *fyne.DragEvent) {
	if gc.draggingComponent != nil {
		newPos := fyne.NewPos(
			ev.Position.X-gc.draggingComponent.DragOffset.X,
			ev.Position.Y-gc.draggingComponent.DragOffset.Y,
		)
		gc.draggingComponent.Position = newPos
		gc.Refresh()
	} else if gc.connectingFrom != nil {
		gc.connectingTo = ev.Position
		gc.Refresh()
	} else {
		// Start dragging if the gesture began over a component.
		if comp := gc.GetComponentAt(ev.Position); comp != nil {
			gc.draggingComponent = comp
			comp.Dragging = true
			comp.DragOffset = fyne.NewPos(ev.Position.X-comp.Position.X, ev.Position.Y-comp.Position.Y)
		}
	}
}

func (gc *GraphCanvas) DragEnd() {
	if gc.draggingComponent != nil {
		gc.draggingComponent.Dragging = false
		gc.draggingComponent = nil
	}
}

func (gc *GraphCanvas) MouseIn(*desktop.MouseEvent) {}
func (gc *GraphCanvas) MouseMoved(ev *desktop.MouseEvent) {
	if gc.connectingFrom != nil {
		gc.connectingTo = ev.Position
		gc.Refresh()
	}
}
func (gc *GraphCanvas) MouseOut() {}

type graphCanvasRenderer struct {
	canvas  *GraphCanvas
	objects []fyne.CanvasObject
}

func (r *graphCanvasRenderer) Layout(size fyne.Size) {
	// Background elements are sized in renderBackground; no extra layout needed here.
}

func (r *graphCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(800, 600)
}

func (r *graphCanvasRenderer) Refresh() {
	r.objects = []fyne.CanvasObject{}

	r.renderBackground()

	r.canvas.componentsMutex.RLock()
	defer r.canvas.componentsMutex.RUnlock()

	for _, conn := range r.canvas.connections {
		r.renderConnection(conn)
	}

	for _, comp := range r.canvas.components {
		r.renderComponent(comp)
	}

	if r.canvas.connectingFrom != nil {
		r.renderConnectingLine()
	}

	canvas.Refresh(r.canvas)
}

func (r *graphCanvasRenderer) renderConnection(conn *gui.Connection) {
	fromCenter := conn.From.GetCenter()
	toCenter := conn.To.GetCenter()

	line := canvas.NewLine(conn.Color)
	line.Position1 = fromCenter
	line.Position2 = toCenter
	line.StrokeWidth = conn.Thickness + 1
	r.objects = append(r.objects, line)

	for _, particle := range conn.Particles {
		x := fromCenter.X + (toCenter.X-fromCenter.X)*particle.Position
		y := fromCenter.Y + (toCenter.Y-fromCenter.Y)*particle.Position

		circle := canvas.NewCircle(particle.Color)
		circle.Position1 = fyne.NewPos(x-particle.Size/2, y-particle.Size/2)
		circle.Position2 = fyne.NewPos(x+particle.Size/2, y+particle.Size/2)
		circle.FillColor = particle.Color
		r.objects = append(r.objects, circle)
	}
}

func (r *graphCanvasRenderer) renderComponent(comp *gui.VisualComponent) {
	comp.UpdateHealthStatus()

	borderColor := comp.GetColor()
	fillColor := comp.GetTypeColor()

	if comp.Selected {
		borderColor = color.RGBA{R: 255, G: 225, B: 120, A: 255}
	}

	// Outer glow to make nodes feel alive.
	glow := canvas.NewRectangle(color.RGBA{R: fillColor.(color.RGBA).R, G: fillColor.(color.RGBA).G, B: fillColor.(color.RGBA).B, A: 60})
	glow.Move(fyne.NewPos(comp.Position.X-6, comp.Position.Y-6))
	glow.Resize(fyne.NewSize(comp.Size.Width+12, comp.Size.Height+12))
	glow.CornerRadius = 14
	r.objects = append(r.objects, glow)

	bgRect := canvas.NewRectangle(fillColor)
	bgRect.Move(comp.Position)
	bgRect.Resize(comp.Size)
	bgRect.CornerRadius = 10
	r.objects = append(r.objects, bgRect)

	border := canvas.NewRectangle(color.Transparent)
	border.Move(comp.Position)
	border.Resize(comp.Size)
	border.StrokeColor = borderColor
	border.StrokeWidth = 3
	border.CornerRadius = 10
	r.objects = append(r.objects, border)

	title := fmt.Sprintf("%s", strings.ToUpper(strings.ReplaceAll(string(comp.Type), "-", " ")))
	titleLabel := canvas.NewText(title, color.White)
	titleLabel.TextSize = 11
	titleLabel.Alignment = fyne.TextAlignCenter
	titlePos := fyne.NewPos(
		comp.Position.X+comp.Size.Width/2-titleLabel.MinSize().Width/2,
		comp.Position.Y+6,
	)
	titleLabel.Move(titlePos)
	r.objects = append(r.objects, titleLabel)

	if comp.Component != nil {
		metrics := comp.Component.GetMetrics()
		statusText := fmt.Sprintf("thr: %.0f rps", metrics.Throughput)
		statusLabel := canvas.NewText(statusText, color.White)
		statusLabel.TextSize = 9
		statusLabelPos := fyne.NewPos(
			comp.Position.X+8,
			comp.Position.Y+comp.Size.Height-18,
		)
		statusLabel.Move(statusLabelPos)
		r.objects = append(r.objects, statusLabel)

		health := canvas.NewRectangle(borderColor)
		health.Move(fyne.NewPos(comp.Position.X, comp.Position.Y+comp.Size.Height-6))
		health.Resize(fyne.NewSize(comp.Size.Width, 6))
		health.CornerRadius = 3
		r.objects = append(r.objects, health)
	}
}

func (r *graphCanvasRenderer) renderConnectingLine() {
	if r.canvas.connectingFrom == nil {
		return
	}

	fromCenter := r.canvas.connectingFrom.GetCenter()

	line := canvas.NewLine(color.RGBA{R: 100, G: 100, B: 100, A: 150})
	line.Position1 = fromCenter
	line.Position2 = r.canvas.connectingTo
	line.StrokeWidth = 2
	r.objects = append(r.objects, line)
}

func (r *graphCanvasRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *graphCanvasRenderer) Destroy() {}

func (r *graphCanvasRenderer) renderBackground() {
	size := r.canvas.Size()

	// Gradient backdrop for a more game-like board feel.
	bg := canvas.NewLinearGradient(
		color.RGBA{R: 16, G: 22, B: 34, A: 255},
		color.RGBA{R: 6, G: 10, B: 18, A: 255},
		90,
	)
	bg.Resize(size)
	r.objects = append(r.objects, bg)

	// Subtle grid to suggest placement and movement.
	gridColor := color.RGBA{R: 80, G: 120, B: 170, A: 40}
	spacing := float32(80)
	for x := float32(0); x < size.Width; x += spacing {
		line := canvas.NewLine(gridColor)
		line.Position1 = fyne.NewPos(x, 0)
		line.Position2 = fyne.NewPos(x, size.Height)
		line.StrokeWidth = 1
		r.objects = append(r.objects, line)
	}
	for y := float32(0); y < size.Height; y += spacing {
		line := canvas.NewLine(gridColor)
		line.Position1 = fyne.NewPos(0, y)
		line.Position2 = fyne.NewPos(size.Width, y)
		line.StrokeWidth = 1
		r.objects = append(r.objects, line)
	}

	// Empty-state hint when no components exist.
	if len(r.canvas.components) == 0 {
		hint := canvas.NewText("Drop components here and right-drag to link them", color.RGBA{R: 190, G: 208, B: 255, A: 200})
		hint.TextSize = 16
		hint.Alignment = fyne.TextAlignCenter
		center := fyne.NewPos(size.Width/2-hint.MinSize().Width/2, size.Height/2-hint.MinSize().Height/2)
		hint.Move(center)
		r.objects = append(r.objects, hint)
	}
}

func distance(p1, p2 fyne.Position) float32 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
