package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowHelpPopup(window fyne.Window) {
	title := widget.NewLabel("Quick Reference Guide")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	componentsTitle := widget.NewLabel("Component Overview:")
	componentsTitle.TextStyle = fyne.TextStyle{Bold: true}

	components := widget.NewLabel(
		"API Server (Blue): Handles requests, has capacity limits\n" +
			"Database (Purple): Stores data, slower than cache\n" +
			"Cache (Green): Fast in-memory storage, aim for 70%+ hit rate\n" +
			"Load Balancer (Yellow): Distributes traffic across servers\n" +
			"CDN (Dark Blue): Edge caching for global users",
	)
	components.Wrapping = fyne.TextWrapWord

	controlsTitle := widget.NewLabel("Controls:")
	controlsTitle.TextStyle = fyne.TextStyle{Bold: true}

	controls := widget.NewLabel(
		"• Click buttons in toolbox to add components\n" +
			"• Right-click + drag to create connections\n" +
			"• Left-click to select components\n" +
			"• Watch color indicators for health status",
	)

	healthTitle := widget.NewLabel("Health Colors:")
	healthTitle.TextStyle = fyne.TextStyle{Bold: true}

	health := widget.NewLabel(
		"Green: Healthy (< 50% load)\n" +
			"Yellow: Warning (50-80% load)\n" +
			"Orange: Critical (> 80% load)\n" +
			"Red: Down/Failing",
	)

	connectionsTitle := widget.NewLabel("Valid Connections:")
	connectionsTitle.TextStyle = fyne.TextStyle{Bold: true}

	connections := widget.NewLabel(
		"Load Balancer → API Server\n" +
			"API Server → Database\n" +
			"API Server → Cache\n" +
			"Cache → Database\n" +
			"CDN → API Server",
	)

	tipsTitle := widget.NewLabel("Quick Tips:")
	tipsTitle.TextStyle = fyne.TextStyle{Bold: true}

	tips := widget.NewLabel(
		"• Use caching for read-heavy workloads\n" +
			"• Add load balancer for horizontal scaling\n" +
			"• Monitor metrics before submitting\n" +
			"• Balance cost vs performance",
	)
	tips.Wrapping = fyne.TextWrapWord

	closeButton := widget.NewButton("Close", func() {
		window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		componentsTitle,
		components,
		widget.NewSeparator(),
		controlsTitle,
		controls,
		widget.NewSeparator(),
		healthTitle,
		health,
		widget.NewSeparator(),
		connectionsTitle,
		connections,
		widget.NewSeparator(),
		tipsTitle,
		tips,
		widget.NewSeparator(),
		closeButton,
	)

	scrollContent := container.NewVScroll(content)
	scrollContent.SetMinSize(fyne.NewSize(500, 600))

	// This would be used as a dialog, but Fyne doesn't support overlays easily
	// For now, this is a placeholder for future implementation
}
