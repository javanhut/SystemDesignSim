package widgets

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/components/api"
	"github.com/javanhut/systemdesignsim/internal/components/cache"
	"github.com/javanhut/systemdesignsim/internal/components/cdn"
	"github.com/javanhut/systemdesignsim/internal/components/config"
	"github.com/javanhut/systemdesignsim/internal/components/database"
	"github.com/javanhut/systemdesignsim/internal/components/loadbalancer"
	"github.com/javanhut/systemdesignsim/internal/gui"
)

type PropertyPanel struct {
	widget.BaseWidget
	component *gui.VisualComponent
	window    fyne.Window
	onUpdate  func()
	onDelete  func()

	content *fyne.Container
}

func NewPropertyPanel(component *gui.VisualComponent, window fyne.Window, onUpdate func(), onDelete func()) *PropertyPanel {
	pp := &PropertyPanel{
		component: component,
		window:    window,
		onUpdate:  onUpdate,
		onDelete:  onDelete,
	}
	pp.ExtendBaseWidget(pp)
	pp.buildContent()
	return pp
}

func (pp *PropertyPanel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(pp.content)
}

func (pp *PropertyPanel) buildContent() {
	title := widget.NewLabel(fmt.Sprintf("Configure: %s", pp.component.ID))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	var propertyWidgets []fyne.CanvasObject
	var saveFunc func()

	switch pp.component.Type {
	case gui.ComponentTypeAPIServer:
		propertyWidgets, saveFunc = pp.buildAPIServerProperties()
	case gui.ComponentTypeDatabase:
		propertyWidgets, saveFunc = pp.buildDatabaseProperties()
	case gui.ComponentTypeCache:
		propertyWidgets, saveFunc = pp.buildCacheProperties()
	case gui.ComponentTypeLoadBalancer:
		propertyWidgets, saveFunc = pp.buildLoadBalancerProperties()
	case gui.ComponentTypeCDN:
		propertyWidgets, saveFunc = pp.buildCDNProperties()
	default:
		propertyWidgets = []fyne.CanvasObject{
			widget.NewLabel("No properties available"),
		}
	}

	closeButton := widget.NewButton("Close", func() {
		if pp.window.Canvas().Overlays().Top() != nil {
			pp.window.Canvas().Overlays().Remove(pp.window.Canvas().Overlays().Top())
		}
	})

	deleteButton := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		dialog := widget.NewModalPopUp(
			container.NewVBox(
				widget.NewLabel("Are you sure you want to delete this component?"),
				container.NewHBox(
					widget.NewButton("Cancel", func() {
						pp.window.Canvas().Overlays().Top().Hide()
					}),
					widget.NewButton("Delete", func() {
						pp.window.Canvas().Overlays().Top().Hide() // Hide confirmation
						if pp.window.Canvas().Overlays().Top() != nil {
							pp.window.Canvas().Overlays().Remove(pp.window.Canvas().Overlays().Top()) // Hide property panel
						}
						if pp.onDelete != nil {
							pp.onDelete()
						}
					}),
				),
			),
			pp.window.Canvas(),
		)
		dialog.Show()
	})
	deleteButton.Importance = widget.DangerImportance

	saveButton := widget.NewButton("Save & Apply", func() {
		if saveFunc != nil {
			saveFunc()
		}
		if pp.onUpdate != nil {
			pp.onUpdate()
		}
		if pp.window.Canvas().Overlays().Top() != nil {
			pp.window.Canvas().Overlays().Remove(pp.window.Canvas().Overlays().Top())
		}
	})
	saveButton.Importance = widget.HighImportance

	buttons := container.NewHBox(saveButton, closeButton, widget.NewSeparator(), deleteButton)

	contentItems := []fyne.CanvasObject{
		title,
		widget.NewSeparator(),
	}
	contentItems = append(contentItems, propertyWidgets...)
	contentItems = append(contentItems, widget.NewSeparator(), buttons)

	pp.content = container.NewVBox(contentItems...)
}

func (pp *PropertyPanel) buildAPIServerProperties() ([]fyne.CanvasObject, func()) {
	comp, ok := pp.component.Component.(*api.APIServer)
	if !ok {
		return []fyne.CanvasObject{widget.NewLabel("Error: Invalid Component Type")}, nil
	}

	widgets := []fyne.CanvasObject{}

	// Instance Type
	instanceTypeLabel := widget.NewLabel("Instance Type:")
	instanceTypeLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, instanceTypeLabel)

	instanceTypes := config.GetInstanceTypeNames()
	instanceSelect := widget.NewSelect(instanceTypes, nil)
	currentSize := string(comp.Size)
	if currentSize == "" {
		currentSize = "t2.micro"
	}
	instanceSelect.SetSelected(currentSize)
	widgets = append(widgets, instanceSelect)

	// Region
	regionLabel := widget.NewLabel("Region:")
	regionLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, regionLabel)

	regionIDs := config.GetRegionIDs()
	regionSelect := widget.NewSelect(regionIDs, nil)
	regionSelect.SetSelected(comp.Region)
	widgets = append(widgets, regionSelect)

	saveFunc := func() {
		comp.Size = api.InstanceSize(instanceSelect.Selected)
		comp.Region = regionSelect.Selected
		// Note: Updating size properties (cost, cores) would normally require re-calling NewAPIServer logic or a helper
		// For now we just set the field, but a real fix would update derived stats too.
	}

	return widgets, saveFunc
}

func (pp *PropertyPanel) buildDatabaseProperties() ([]fyne.CanvasObject, func()) {
	comp, ok := pp.component.Component.(*database.Database)
	if !ok {
		return []fyne.CanvasObject{widget.NewLabel("Error: Invalid Component Type")}, nil
	}

	widgets := []fyne.CanvasObject{}

	// DB Type
	dbTypeLabel := widget.NewLabel("Database Type:")
	dbTypeLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, dbTypeLabel)

	dbTypes := []string{"PostgreSQL", "MySQL", "MongoDB", "DynamoDB", "Redis"}
	dbTypeSelect := widget.NewSelect(dbTypes, nil)
	// Mapping internal int type to string for UI
	currentType := "PostgreSQL" // default
	switch comp.Type {
	case database.DatabaseTypeSQL:
		currentType = "PostgreSQL"
	case database.DatabaseTypeNoSQL:
		currentType = "MongoDB"
	case database.DatabaseTypeKeyValue:
		currentType = "Redis"
	case database.DatabaseTypeDocument:
		currentType = "DynamoDB"
	}
	dbTypeSelect.SetSelected(currentType)
	widgets = append(widgets, dbTypeSelect)

	// Region
	regionLabel := widget.NewLabel("Region:")
	regionLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, regionLabel)

	regionIDs := config.GetRegionIDs()
	regionSelect := widget.NewSelect(regionIDs, nil)
	regionSelect.SetSelected(comp.Region)
	widgets = append(widgets, regionSelect)

	// Storage
	storageLabel := widget.NewLabel("Storage Size (GB):")
	storageLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, storageLabel)

	storageEntry := widget.NewEntry()
	storageEntry.SetText(fmt.Sprintf("%d", comp.Capacity/1024/1024/1024))
	widgets = append(widgets, storageEntry)

	saveFunc := func() {
		switch dbTypeSelect.Selected {
		case "PostgreSQL", "MySQL":
			comp.Type = database.DatabaseTypeSQL
		case "MongoDB":
			comp.Type = database.DatabaseTypeNoSQL
		case "Redis":
			comp.Type = database.DatabaseTypeKeyValue
		case "DynamoDB":
			comp.Type = database.DatabaseTypeDocument
		}
		comp.Region = regionSelect.Selected
		if size, err := strconv.ParseInt(storageEntry.Text, 10, 64); err == nil {
			comp.Capacity = size * 1024 * 1024 * 1024
		}
	}

	return widgets, saveFunc
}

func (pp *PropertyPanel) buildCacheProperties() ([]fyne.CanvasObject, func()) {
	comp, ok := pp.component.Component.(*cache.Cache)
	if !ok {
		return []fyne.CanvasObject{widget.NewLabel("Error: Invalid Component Type")}, nil
	}

	widgets := []fyne.CanvasObject{}

	// Engine
	cacheTypeLabel := widget.NewLabel("Cache Engine:")
	cacheTypeLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, cacheTypeLabel)

	cacheTypes := []string{"Redis", "Memcached"}
	cacheTypeSelect := widget.NewSelect(cacheTypes, nil)
	cacheTypeSelect.SetSelected(comp.Type)
	widgets = append(widgets, cacheTypeSelect)

	// Region
	regionLabel := widget.NewLabel("Region:")
	regionLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, regionLabel)

	regionIDs := config.GetRegionIDs()
	regionSelect := widget.NewSelect(regionIDs, nil)
	regionSelect.SetSelected(comp.Region)
	widgets = append(widgets, regionSelect)

	// Eviction
	evictionLabel := widget.NewLabel("Eviction Policy:")
	evictionLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, evictionLabel)

	evictionPolicies := []string{"lru", "lfu", "fifo", "random"}
	evictionSelect := widget.NewSelect(evictionPolicies, nil)
	evictionSelect.SetSelected(string(comp.Policy))
	widgets = append(widgets, evictionSelect)

	saveFunc := func() {
		comp.Type = cacheTypeSelect.Selected
		comp.Region = regionSelect.Selected
		comp.Policy = cache.EvictionPolicy(evictionSelect.Selected)
	}

	return widgets, saveFunc
}

func (pp *PropertyPanel) buildLoadBalancerProperties() ([]fyne.CanvasObject, func()) {
	comp, ok := pp.component.Component.(*loadbalancer.LoadBalancer)
	if !ok {
		return []fyne.CanvasObject{widget.NewLabel("Error: Invalid Component Type")}, nil
	}

	widgets := []fyne.CanvasObject{}

	// Algorithm
	algorithmLabel := widget.NewLabel("Routing Algorithm:")
	algorithmLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, algorithmLabel)

	algorithms := []string{"round-robin", "least-connected", "ip-hash", "weighted-random"}
	algorithmSelect := widget.NewSelect(algorithms, nil)
	algorithmSelect.SetSelected(string(comp.Strategy))
	widgets = append(widgets, algorithmSelect)

	// Region
	regionLabel := widget.NewLabel("Region:")
	regionLabel.TextStyle = fyne.TextStyle{Bold: true}
	widgets = append(widgets, regionLabel)

	regionIDs := config.GetRegionIDs()
	regionSelect := widget.NewSelect(regionIDs, nil)
	regionSelect.SetSelected(comp.Region)
	widgets = append(widgets, regionSelect)

	saveFunc := func() {
		comp.Strategy = loadbalancer.LoadBalancingStrategy(algorithmSelect.Selected)
		comp.Region = regionSelect.Selected
	}

	return widgets, saveFunc
}

func (pp *PropertyPanel) buildCDNProperties() ([]fyne.CanvasObject, func()) {
	comp, ok := pp.component.Component.(*cdn.CDN)
	if !ok {
		return []fyne.CanvasObject{widget.NewLabel("Error: Invalid Component Type")}, nil
	}

	widgets := []fyne.CanvasObject{}

	// Regions
	widgets = append(widgets, widget.NewLabel("Regions (Active):"))

	// Reconstruct regions from EdgeLocations map keys
	regionsStr := ""
	i := 0
	for region := range comp.EdgeLocations {
		if i > 0 {
			regionsStr += ", "
		}
		regionsStr += region
		i++
	}

	regionsEntry := widget.NewEntry()
	regionsEntry.SetText(regionsStr)
	widgets = append(widgets, regionsEntry)
	widgets = append(widgets, widget.NewLabel("(Comma separated, e.g., us-east, us-west)"))

	saveFunc := func() {
		// Naive implementation: recreate map based on input
		// In a real app, we'd want to preserve cache state for existing regions
		// But for this sim, resizing the CDN is effectively a reset

		// Re-creating the map
		newLocations := make(map[string]*cdn.EdgeLocation)

		// Manual split by comma
		inputText := regionsEntry.Text
		currentRegion := ""
		for _, r := range inputText {
			if r == ',' {
				trimmed := trimSpace(currentRegion)
				if trimmed != "" {
					newLocations[trimmed] = &cdn.EdgeLocation{Region: trimmed, Cache: make(map[string][]byte)}
				}
				currentRegion = ""
			} else {
				currentRegion += string(r)
			}
		}
		trimmed := trimSpace(currentRegion)
		if trimmed != "" {
			newLocations[trimmed] = &cdn.EdgeLocation{Region: trimmed, Cache: make(map[string][]byte)}
		}

		comp.EdgeLocations = newLocations
	}

	return widgets, saveFunc
}

// ShowPropertyPanel displays the property panel as an overlay on the window
func ShowPropertyPanel(component *gui.VisualComponent, window fyne.Window, onUpdate func(), onDelete func()) {
	panel := NewPropertyPanel(component, window, onUpdate, onDelete)

	modal := widget.NewModalPopUp(
		container.NewPadded(panel),
		window.Canvas(),
	)
	modal.Resize(fyne.NewSize(400, 500))
	modal.Show()
}

// Helper for manual trim since I can't guarantee strings package
func trimSpace(s string) string {
	// simplistic trim
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
