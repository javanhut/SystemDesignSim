package gui

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/javanhut/systemdesignsim/internal/engine"
)

type ComponentType string

const (
	ComponentTypeAPIServer    ComponentType = "api-server"
	ComponentTypeDatabase     ComponentType = "database"
	ComponentTypeCache        ComponentType = "cache"
	ComponentTypeCDN          ComponentType = "cdn"
	ComponentTypeLoadBalancer ComponentType = "load-balancer"
	ComponentTypeDNS          ComponentType = "dns"
	ComponentTypeGateway      ComponentType = "gateway"
	ComponentTypeFirewall     ComponentType = "firewall"
	ComponentTypeNAT          ComponentType = "nat"
	ComponentTypeRouter       ComponentType = "router"
	ComponentTypeUserPool     ComponentType = "user-pool"
)

type VisualComponent struct {
	ID            string
	Type          ComponentType
	Position      fyne.Position
	Size          fyne.Size
	Component     engine.Component
	Connections   []*Connection
	Selected      bool
	Dragging      bool
	DragOffset    fyne.Position
	HealthStatus  HealthStatus
	Properties    map[string]interface{}
	mu            sync.RWMutex
}

type Connection struct {
	ID          string
	From        *VisualComponent
	To          *VisualComponent
	Color       color.Color
	Thickness   float32
	Animated    bool
	Particles   []*Particle
}

type Particle struct {
	Position  float32
	Speed     float32
	Color     color.Color
	Size      float32
}

type HealthStatus int

const (
	HealthStatusHealthy HealthStatus = iota
	HealthStatusWarning
	HealthStatusCritical
	HealthStatusDown
)

func NewVisualComponent(id string, compType ComponentType, pos fyne.Position) *VisualComponent {
	return &VisualComponent{
		ID:           id,
		Type:         compType,
		Position:     pos,
		Size:         fyne.NewSize(80, 80),
		Connections:  make([]*Connection, 0),
		Selected:     false,
		Dragging:     false,
		HealthStatus: HealthStatusHealthy,
		Properties:   make(map[string]interface{}),
	}
}

func (vc *VisualComponent) SetComponent(comp engine.Component) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.Component = comp
}

func (vc *VisualComponent) GetComponent() engine.Component {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.Component
}

func (vc *VisualComponent) AddConnection(conn *Connection) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.Connections = append(vc.Connections, conn)
}

func (vc *VisualComponent) RemoveConnection(connID string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	
	for i, conn := range vc.Connections {
		if conn.ID == connID {
			vc.Connections = append(vc.Connections[:i], vc.Connections[i+1:]...)
			return
		}
	}
}

func (vc *VisualComponent) GetCenter() fyne.Position {
	return fyne.NewPos(
		vc.Position.X+vc.Size.Width/2,
		vc.Position.Y+vc.Size.Height/2,
	)
}

func (vc *VisualComponent) Contains(pos fyne.Position) bool {
	return pos.X >= vc.Position.X && pos.X <= vc.Position.X+vc.Size.Width &&
		pos.Y >= vc.Position.Y && pos.Y <= vc.Position.Y+vc.Size.Height
}

func (vc *VisualComponent) UpdateHealthStatus() {
	if vc.Component == nil {
		return
	}

	if !vc.Component.IsHealthy() {
		vc.HealthStatus = HealthStatusDown
		return
	}

	metrics := vc.Component.GetMetrics()
	if metrics.ErrorRate > 0.1 {
		vc.HealthStatus = HealthStatusCritical
	} else if metrics.ErrorRate > 0.05 {
		vc.HealthStatus = HealthStatusWarning
	} else {
		vc.HealthStatus = HealthStatusHealthy
	}
}

func (vc *VisualComponent) GetColor() color.Color {
	switch vc.HealthStatus {
	case HealthStatusHealthy:
		return color.RGBA{R: 46, G: 204, B: 113, A: 255}
	case HealthStatusWarning:
		return color.RGBA{R: 241, G: 196, B: 15, A: 255}
	case HealthStatusCritical:
		return color.RGBA{R: 230, G: 126, B: 34, A: 255}
	case HealthStatusDown:
		return color.RGBA{R: 231, G: 76, B: 60, A: 255}
	default:
		return color.RGBA{R: 149, G: 165, B: 166, A: 255}
	}
}

func (vc *VisualComponent) GetTypeColor() color.Color {
	switch vc.Type {
	case ComponentTypeAPIServer:
		return color.RGBA{R: 52, G: 152, B: 219, A: 255} // Blue
	case ComponentTypeDatabase:
		return color.RGBA{R: 155, G: 89, B: 182, A: 255} // Purple
	case ComponentTypeCache:
		return color.RGBA{R: 26, G: 188, B: 156, A: 255} // Teal
	case ComponentTypeCDN:
		return color.RGBA{R: 52, G: 73, B: 94, A: 255} // Dark blue-gray
	case ComponentTypeLoadBalancer:
		return color.RGBA{R: 241, G: 196, B: 15, A: 255} // Yellow
	case ComponentTypeDNS:
		return color.RGBA{R: 142, G: 68, B: 173, A: 255} // Purple
	case ComponentTypeGateway:
		return color.RGBA{R: 46, G: 204, B: 113, A: 255} // Green
	case ComponentTypeFirewall:
		return color.RGBA{R: 231, G: 76, B: 60, A: 255} // Red
	case ComponentTypeNAT:
		return color.RGBA{R: 22, G: 160, B: 133, A: 255} // Sea green
	case ComponentTypeRouter:
		return color.RGBA{R: 52, G: 152, B: 219, A: 255} // Light blue
	case ComponentTypeUserPool:
		return color.RGBA{R: 149, G: 165, B: 166, A: 255} // Gray
	default:
		return color.RGBA{R: 127, G: 140, B: 141, A: 255}
	}
}

func NewConnection(id string, from, to *VisualComponent) *Connection {
	return &Connection{
		ID:        id,
		From:      from,
		To:        to,
		Color:     color.RGBA{R: 100, G: 100, B: 100, A: 200},
		Thickness: 2.0,
		Animated:  true,
		Particles: make([]*Particle, 0),
	}
}

func (c *Connection) AddParticle() {
	particle := &Particle{
		Position: 0.0,
		Speed:    0.01,
		Color:    color.RGBA{R: 52, G: 152, B: 219, A: 255},
		Size:     4.0,
	}
	c.Particles = append(c.Particles, particle)
}

func (c *Connection) UpdateParticles() {
	for i := len(c.Particles) - 1; i >= 0; i-- {
		c.Particles[i].Position += c.Particles[i].Speed
		
		if c.Particles[i].Position > 1.0 {
			c.Particles = append(c.Particles[:i], c.Particles[i+1:]...)
		}
	}
}
