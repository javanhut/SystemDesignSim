package screens

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/game"
)

type LevelSelectScreen struct {
	window fyne.Window
}

func NewLevelSelectScreen(window fyne.Window) *LevelSelectScreen {
	return &LevelSelectScreen{
		window: window,
	}
}

func (ls *LevelSelectScreen) Build() fyne.CanvasObject {
	title := widget.NewLabel("System Design Simulator")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Select a Level")
	subtitle.Alignment = fyne.TextAlignCenter

	levelButtons := make([]fyne.CanvasObject, 0)

	for _, level := range game.Levels {
		lvl := level

		statusText := "Locked"
		if lvl.Unlocked {
			if lvl.Completed {
				statusText = fmt.Sprintf("Completed - Score: %d", lvl.BestScore)
			} else {
				statusText = "Not Completed"
			}
		}

		difficultyColor := map[game.Difficulty]string{
			game.DifficultyEasy:   "Easy",
			game.DifficultyMedium: "Medium",
			game.DifficultyHard:   "Hard",
			game.DifficultyExpert: "Expert",
		}

		levelInfo := fmt.Sprintf(
			"Level %d: %s\n%s\nDifficulty: %s\nUsers: %d | Budget: $%.0f\n%s",
			lvl.ID,
			lvl.Name,
			lvl.Description,
			difficultyColor[lvl.Difficulty],
			lvl.PeakUsers,
			lvl.Budget,
			statusText,
		)

		card := widget.NewCard(
			fmt.Sprintf("Level %d", lvl.ID),
			lvl.Name,
			widget.NewLabel(levelInfo),
		)

		playButton := widget.NewButton("Play", func() {
			ls.window.SetContent(NewGameScreen(ls.window, lvl).Build())
		})

		if !lvl.Unlocked {
			playButton.Disable()
		}

		levelContainer := container.NewVBox(
			card,
			playButton,
			widget.NewSeparator(),
		)

		levelButtons = append(levelButtons, levelContainer)
	}

	levelList := container.NewVBox(levelButtons...)

	scrollContainer := container.NewVScroll(levelList)
	scrollContainer.SetMinSize(fyne.NewSize(600, 500))

	tutorialButton := widget.NewButton("View Tutorial", func() {
		ls.window.SetContent(NewWelcomeScreen(ls.window).Build())
	})

	content := container.NewVBox(
		title,
		subtitle,
		tutorialButton,
		widget.NewSeparator(),
		scrollContainer,
	)

	return content
}
