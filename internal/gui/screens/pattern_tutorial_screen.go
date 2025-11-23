package screens

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/javanhut/systemdesignsim/internal/game"
	"github.com/javanhut/systemdesignsim/internal/gui"
	guicanvas "github.com/javanhut/systemdesignsim/internal/gui/canvas"
)

type PatternTutorialScreen struct {
	window       fyne.Window
	pattern      *game.DesignPattern
	orchestrator *game.TutorialOrchestrator
	canvas       *guicanvas.GraphCanvas

	titleLabel       *widget.Label
	stepLabel        *widget.Label
	progressBar      *widget.ProgressBar
	messageTitle     *widget.Label
	messageBody      *widget.Label
	messageBox       *fyne.Container
	instructionLabel *widget.Label
	feedbackLabel    *widget.Label

	playButton     *widget.Button
	pauseButton    *widget.Button
	restartButton  *widget.Button
	practiceButton *widget.Button
	nextButton     *widget.Button
	checkButton    *widget.Button
	backButton     *widget.Button

	mode              game.TutorialMode
	currentPracticeStep int
	practiceStarted   bool

	animationTicker *time.Ticker
	stopAnimChan    chan bool
}

func NewPatternTutorialScreen(window fyne.Window, pattern *game.DesignPattern) *PatternTutorialScreen {
	pts := &PatternTutorialScreen{
		window:             window,
		pattern:            pattern,
		canvas:             guicanvas.NewGraphCanvas(),
		mode:               game.ModeDemoWatch,
		currentPracticeStep: 0,
		stopAnimChan:       make(chan bool),
	}

	pts.orchestrator = game.NewTutorialOrchestrator(pattern, pts.canvas)
	pts.setupCallbacks()

	return pts
}

func (pts *PatternTutorialScreen) setupCallbacks() {
	pts.orchestrator.SetOnStepComplete(func(step int, total int) {
		pts.updateProgress()
	})

	pts.orchestrator.SetOnTutorialEnd(func() {
		pts.onDemoComplete()
	})

	pts.orchestrator.SetOnMessage(func(title, description string) {
		pts.showStepMessage(title, description)
	})
}

func (pts *PatternTutorialScreen) Build() fyne.CanvasObject {
	leftPanel := pts.createInfoPanel()
	centerPanel := pts.createCanvasPanel()
	rightPanel := pts.createControlsPanel()

	header := pts.createHeader()

	mainContent := container.NewBorder(
		header,
		nil,
		leftPanel,
		rightPanel,
		centerPanel,
	)

	pts.startParticleAnimation()

	return mainContent
}

func (pts *PatternTutorialScreen) createHeader() fyne.CanvasObject {
	pts.titleLabel = widget.NewLabel(fmt.Sprintf("Design Pattern: %s", pts.pattern.Name))
	pts.titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	pts.titleLabel.Alignment = fyne.TextAlignCenter

	titleText := canvas.NewText(pts.titleLabel.Text, color.RGBA{R: 180, G: 214, B: 255, A: 255})
	titleText.TextSize = 18
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter

	subtitle := canvas.NewText(fmt.Sprintf("Category: %s | Difficulty: %d/5", pts.pattern.Category, pts.pattern.Difficulty),
		color.RGBA{R: 150, G: 190, B: 255, A: 220})
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
			container.NewCenter(titleText),
			container.NewCenter(subtitle),
		),
	)
}

func (pts *PatternTutorialScreen) createInfoPanel() fyne.CanvasObject {
	problemLabel := widget.NewLabel("THE PROBLEM")
	problemLabel.TextStyle = fyne.TextStyle{Bold: true}

	problemText := widget.NewLabel(pts.pattern.Problem)
	problemText.Wrapping = fyne.TextWrapWord

	solutionLabel := widget.NewLabel("THE SOLUTION")
	solutionLabel.TextStyle = fyne.TextStyle{Bold: true}

	solutionText := widget.NewLabel(pts.pattern.Solution)
	solutionText.Wrapping = fyne.TextWrapWord

	benefitsLabel := widget.NewLabel("BENEFITS")
	benefitsLabel.TextStyle = fyne.TextStyle{Bold: true}

	benefitsText := ""
	for _, benefit := range pts.pattern.Benefits {
		benefitsText += "✓ " + benefit + "\n"
	}
	benefitsList := widget.NewLabel(benefitsText)
	benefitsList.Wrapping = fyne.TextWrapWord

	tradeoffsLabel := widget.NewLabel("TRADE-OFFS")
	tradeoffsLabel.TextStyle = fyne.TextStyle{Bold: true}

	tradeoffsText := ""
	for _, tradeoff := range pts.pattern.Tradeoffs {
		tradeoffsText += "⚠ " + tradeoff + "\n"
	}
	tradeoffsList := widget.NewLabel(tradeoffsText)
	tradeoffsList.Wrapping = fyne.TextWrapWord

	realWorldLabel := widget.NewLabel("REAL-WORLD EXAMPLES")
	realWorldLabel.TextStyle = fyne.TextStyle{Bold: true}

	realWorldText := ""
	for _, example := range pts.pattern.RealWorld {
		realWorldText += "• " + example + "\n"
	}
	realWorldList := widget.NewLabel(realWorldText)
	realWorldList.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		problemLabel,
		problemText,
		widget.NewSeparator(),
		solutionLabel,
		solutionText,
		widget.NewSeparator(),
		benefitsLabel,
		benefitsList,
		widget.NewSeparator(),
		tradeoffsLabel,
		tradeoffsList,
		widget.NewSeparator(),
		realWorldLabel,
		realWorldList,
	)

	scroll := container.NewVScroll(content)
	scroll.SetMinSize(fyne.NewSize(300, 600))

	return pts.wrapPanel(scroll)
}

func (pts *PatternTutorialScreen) createCanvasPanel() fyne.CanvasObject {
	pts.messageTitle = widget.NewLabel("")
	pts.messageTitle.TextStyle = fyne.TextStyle{Bold: true}
	pts.messageTitle.Alignment = fyne.TextAlignCenter

	pts.messageBody = widget.NewLabel("")
	pts.messageBody.Wrapping = fyne.TextWrapWord
	pts.messageBody.Alignment = fyne.TextAlignCenter

	messageBg := canvas.NewRectangle(color.RGBA{R: 42, G: 52, B: 72, A: 240})
	messageBg.CornerRadius = 8

	pts.messageBox = container.NewMax(
		messageBg,
		container.NewPadded(
			container.NewVBox(
				pts.messageTitle,
				pts.messageBody,
			),
		),
	)
	pts.messageBox.Hide()

	canvasContainer := container.NewMax(pts.canvas)

	overlay := container.NewStack(
		canvasContainer,
		container.NewCenter(pts.messageBox),
	)

	return overlay
}

func (pts *PatternTutorialScreen) createControlsPanel() fyne.CanvasObject {
	pts.stepLabel = widget.NewLabel("Step 0 of 0")
	pts.stepLabel.Alignment = fyne.TextAlignCenter

	pts.progressBar = widget.NewProgressBar()
	pts.progressBar.Min = 0
	pts.progressBar.Max = 1

	pts.instructionLabel = widget.NewLabel("")
	pts.instructionLabel.Wrapping = fyne.TextWrapWord
	pts.instructionLabel.Hide()

	pts.feedbackLabel = widget.NewLabel("")
	pts.feedbackLabel.Wrapping = fyne.TextWrapWord
	pts.feedbackLabel.Hide()

	pts.playButton = widget.NewButton("▶ Watch Demo", func() {
		pts.startDemo()
	})

	pts.pauseButton = widget.NewButton("⏸ Pause", func() {
		pts.pauseDemo()
	})
	pts.pauseButton.Disable()

	pts.restartButton = widget.NewButton("↻ Restart", func() {
		pts.restartDemo()
	})

	pts.practiceButton = widget.NewButton("✋ Try Practice Mode", func() {
		pts.startPractice()
	})
	pts.practiceButton.Disable()

	pts.nextButton = widget.NewButton("Next Step →", func() {
		pts.nextPracticeStep()
	})
	pts.nextButton.Hide()

	pts.checkButton = widget.NewButton("Check Progress", func() {
		pts.checkPracticeStep()
	})
	pts.checkButton.Hide()

	pts.backButton = widget.NewButton("← Back to Patterns", func() {
		pts.orchestrator.Stop()
		pts.window.SetContent(NewPatternSelectionScreen(pts.window).Build())
	})

	demoControls := container.NewVBox(
		widget.NewLabel("Demo Controls"),
		widget.NewSeparator(),
		pts.playButton,
		pts.pauseButton,
		pts.restartButton,
		widget.NewSeparator(),
		pts.stepLabel,
		pts.progressBar,
	)

	practiceControls := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("Practice Mode"),
		widget.NewSeparator(),
		pts.practiceButton,
		pts.instructionLabel,
		pts.checkButton,
		pts.nextButton,
		pts.feedbackLabel,
	)

	controls := container.NewVBox(
		demoControls,
		practiceControls,
		widget.NewSeparator(),
		pts.backButton,
	)

	scroll := container.NewVScroll(controls)
	scroll.SetMinSize(fyne.NewSize(250, 600))

	return pts.wrapPanel(scroll)
}

func (pts *PatternTutorialScreen) wrapPanel(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 18, G: 24, B: 38, A: 255})
	bg.StrokeColor = color.RGBA{R: 70, G: 130, B: 255, A: 120}
	bg.StrokeWidth = 1.5
	bg.CornerRadius = 12

	return container.NewMax(
		bg,
		container.NewPadded(content),
	)
}

func (pts *PatternTutorialScreen) startDemo() {
	pts.playButton.Disable()
	pts.pauseButton.Enable()
	pts.practiceButton.Disable()

	pts.orchestrator.Reset()
	pts.orchestrator.StartDemo()
}

func (pts *PatternTutorialScreen) pauseDemo() {
	if pts.orchestrator.IsPaused() {
		pts.orchestrator.Resume()
		pts.pauseButton.SetText("⏸ Pause")
	} else {
		pts.orchestrator.Pause()
		pts.pauseButton.SetText("▶ Resume")
	}
}

func (pts *PatternTutorialScreen) restartDemo() {
	pts.orchestrator.Stop()
	time.Sleep(100 * time.Millisecond)
	pts.orchestrator.Reset()

	pts.playButton.Enable()
	pts.pauseButton.Disable()
	pts.pauseButton.SetText("⏸ Pause")
	pts.practiceButton.Disable()

	pts.updateProgress()
	pts.hideStepMessage()
}

func (pts *PatternTutorialScreen) onDemoComplete() {
	pts.playButton.Enable()
	pts.pauseButton.Disable()
	pts.practiceButton.Enable()

	pts.showStepMessage("Demo Complete!",
		"Great! You've seen how the "+pts.pattern.Name+" pattern works.\n\nNow try building it yourself in Practice Mode!")
}

func (pts *PatternTutorialScreen) startPractice() {
	pts.mode = game.ModePractice
	pts.currentPracticeStep = 0
	pts.practiceStarted = true

	pts.playButton.Disable()
	pts.pauseButton.Disable()
	pts.practiceButton.Disable()

	pts.checkButton.Show()
	pts.instructionLabel.Show()
	pts.feedbackLabel.Show()

	pts.orchestrator.Stop()
	pts.orchestrator.Reset()
	pts.orchestrator.SetMode(game.ModePractice)

	pts.canvas.SetOnComponentAdd(func(vc *gui.VisualComponent) {
		pts.checkPracticeStep()
	})

	pts.showPracticeInstruction()
	pts.hideStepMessage()
}

func (pts *PatternTutorialScreen) showPracticeInstruction() {
	if pts.currentPracticeStep >= len(pts.pattern.PracticeSteps) {
		pts.onPracticeComplete()
		return
	}

	step := pts.pattern.PracticeSteps[pts.currentPracticeStep]
	pts.instructionLabel.SetText(fmt.Sprintf("Step %d/%d:\n%s\n\nHint: %s",
		pts.currentPracticeStep+1,
		len(pts.pattern.PracticeSteps),
		step.Instruction,
		step.Hint))

	pts.feedbackLabel.SetText("")
	pts.feedbackLabel.Importance = widget.MediumImportance
}

func (pts *PatternTutorialScreen) checkPracticeStep() {
	if pts.currentPracticeStep >= len(pts.pattern.PracticeSteps) {
		return
	}

	components := pts.canvas.GetComponents()
	passed, feedback := pts.orchestrator.ValidatePracticeStep(pts.currentPracticeStep, components)

	if passed {
		pts.feedbackLabel.SetText("✓ " + feedback)
		pts.feedbackLabel.Importance = widget.SuccessImportance
		pts.nextButton.Show()
	} else {
		pts.feedbackLabel.SetText("✗ " + feedback)
		pts.feedbackLabel.Importance = widget.DangerImportance
		pts.nextButton.Hide()
	}

	pts.feedbackLabel.Refresh()
}

func (pts *PatternTutorialScreen) nextPracticeStep() {
	pts.currentPracticeStep++
	pts.nextButton.Hide()
	pts.showPracticeInstruction()
}

func (pts *PatternTutorialScreen) onPracticeComplete() {
	pts.checkButton.Hide()
	pts.nextButton.Hide()
	pts.instructionLabel.Hide()
	pts.feedbackLabel.Hide()

	pts.showStepMessage("Practice Complete!",
		"Excellent work! You've successfully implemented the "+pts.pattern.Name+" pattern.\n\n"+
		"You now understand:\n"+
		"• When to use this pattern\n"+
		"• How to implement it\n"+
		"• The benefits and trade-offs\n\n"+
		"Try another pattern or return to the game!")

	pts.backButton.SetText("← Back to Patterns")
}

func (pts *PatternTutorialScreen) showStepMessage(title, description string) {
	pts.messageTitle.SetText(title)
	pts.messageBody.SetText(description)
	pts.messageBox.Show()
}

func (pts *PatternTutorialScreen) hideStepMessage() {
	pts.messageBox.Hide()
}

func (pts *PatternTutorialScreen) updateProgress() {
	step := pts.orchestrator.GetCurrentStep()
	total := pts.orchestrator.GetTotalSteps()
	progress := pts.orchestrator.GetProgress()

	pts.stepLabel.SetText(fmt.Sprintf("Step %d of %d", step, total))
	pts.progressBar.SetValue(progress)
}

func (pts *PatternTutorialScreen) startParticleAnimation() {
	pts.animationTicker = time.NewTicker(50 * time.Millisecond)

	go func() {
		for {
			select {
			case <-pts.stopAnimChan:
				return
			case <-pts.animationTicker.C:
				if pts.orchestrator.IsAnimating() || pts.practiceStarted {
					pts.canvas.UpdateParticles()
				}
			}
		}
	}()
}

func (pts *PatternTutorialScreen) StopAnimations() {
	pts.stopAnimChan <- true
	if pts.animationTicker != nil {
		pts.animationTicker.Stop()
	}
	pts.orchestrator.Stop()
}
