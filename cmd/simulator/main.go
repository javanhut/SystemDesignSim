package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/javanhut/systemdesignsim/internal/gui"
	"github.com/javanhut/systemdesignsim/internal/gui/screens"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("System Design Simulator")

	window.Resize(fyne.NewSize(1200, 800))

	prefs, err := gui.LoadPreferences()
	if err != nil || prefs.FirstLaunch || !prefs.TutorialCompleted {
		welcome := screens.NewWelcomeScreen(window)
		window.SetContent(welcome.Build())
	} else {
		levelSelect := screens.NewLevelSelectScreen(window)
		window.SetContent(levelSelect.Build())
	}

	window.ShowAndRun()
}
