package screens

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/game"
)

type PatternSelectionScreen struct {
	window fyne.Window
}

func NewPatternSelectionScreen(window fyne.Window) *PatternSelectionScreen {
	return &PatternSelectionScreen{
		window: window,
	}
}

func (pss *PatternSelectionScreen) Build() fyne.CanvasObject {
	header := pss.createHeader()
	intro := pss.createIntro()
	patternsGrid := pss.createPatternsGrid()
	backButton := pss.createBackButton()

	content := container.NewVBox(
		header,
		intro,
		widget.NewSeparator(),
		patternsGrid,
		widget.NewSeparator(),
		backButton,
	)

	scroll := container.NewVScroll(content)
	return scroll
}

func (pss *PatternSelectionScreen) createHeader() fyne.CanvasObject {
	title := canvas.NewText("Design Patterns Tutorial", color.RGBA{R: 180, G: 214, B: 255, A: 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText("Learn system design patterns through interactive tutorials", color.RGBA{R: 150, G: 190, B: 255, A: 220})
	subtitle.TextSize = 14
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

func (pss *PatternSelectionScreen) createIntro() fyne.CanvasObject {
	introText := widget.NewLabel(
		"Each pattern tutorial has two modes:\n\n" +
			"1. Watch Demo - Automated walkthrough showing how the pattern works\n" +
			"2. Practice Mode - Build the pattern yourself with guided validation\n\n" +
			"Select a pattern below to begin learning!",
	)
	introText.Wrapping = fyne.TextWrapWord
	introText.Alignment = fyne.TextAlignCenter

	return container.NewPadded(introText)
}

func (pss *PatternSelectionScreen) createPatternsGrid() fyne.CanvasObject {
	patterns := game.GetAllPatterns()

	patternCards := make([]fyne.CanvasObject, 0)

	for _, pattern := range patterns {
		card := pss.createPatternCard(pattern)
		patternCards = append(patternCards, card)
	}

	grid := container.NewGridWithColumns(2, patternCards...)

	return container.NewPadded(grid)
}

func (pss *PatternSelectionScreen) createPatternCard(pattern *game.DesignPattern) fyne.CanvasObject {
	nameLabel := widget.NewLabel(pattern.Name)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	nameLabel.Alignment = fyne.TextAlignCenter

	categoryLabel := widget.NewLabel(fmt.Sprintf("Category: %s", pattern.Category))
	categoryLabel.Alignment = fyne.TextAlignCenter

	difficultyText := fmt.Sprintf("Difficulty: %s", getDifficultyStars(pattern.Difficulty))
	difficultyLabel := widget.NewLabel(difficultyText)
	difficultyLabel.Alignment = fyne.TextAlignCenter

	descriptionLabel := widget.NewLabel(pattern.Description)
	descriptionLabel.Wrapping = fyne.TextWrapWord

	watchBtn := widget.NewButton("Watch Demo", func() {
		tutorialScreen := NewPatternTutorialScreen(pss.window, pattern)
		pss.window.SetContent(tutorialScreen.Build())
	})

	practiceBtn := widget.NewButton("Try Practice", func() {
		tutorialScreen := NewPatternTutorialScreen(pss.window, pattern)
		tutorialScreen.startPractice()
		pss.window.SetContent(tutorialScreen.Build())
	})

	buttons := container.NewGridWithColumns(2, watchBtn, practiceBtn)

	content := container.NewVBox(
		nameLabel,
		categoryLabel,
		difficultyLabel,
		widget.NewSeparator(),
		descriptionLabel,
		widget.NewSeparator(),
		buttons,
	)

	bg := canvas.NewRectangle(getCardColorByCategory(pattern.Category))
	bg.CornerRadius = 8
	bg.StrokeColor = color.RGBA{R: 70, G: 130, B: 255, A: 120}
	bg.StrokeWidth = 2

	card := container.NewMax(
		bg,
		container.NewPadded(content),
	)

	card.Resize(fyne.NewSize(400, 250))

	return card
}

func (pss *PatternSelectionScreen) createBackButton() fyne.CanvasObject {
	backBtn := widget.NewButton("← Back to Level Select", func() {
		pss.window.SetContent(NewLevelSelectScreen(pss.window).Build())
	})

	return container.NewCenter(backBtn)
}

func getDifficultyStars(difficulty int) string {
	stars := ""
	for i := 0; i < difficulty; i++ {
		stars += "★"
	}
	for i := difficulty; i < 5; i++ {
		stars += "☆"
	}
	return stars
}

func getCardColorByCategory(category string) color.Color {
	switch category {
	case "Scalability":
		return color.RGBA{R: 32, G: 42, B: 58, A: 255}
	case "Performance":
		return color.RGBA{R: 42, G: 52, B: 38, A: 255}
	case "Availability":
		return color.RGBA{R: 52, G: 42, B: 38, A: 255}
	case "Reliability":
		return color.RGBA{R: 38, G: 42, B: 52, A: 255}
	default:
		return color.RGBA{R: 28, G: 34, B: 48, A: 255}
	}
}
