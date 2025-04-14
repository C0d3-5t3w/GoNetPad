package components

import (
	"fmt"
	"time"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	widget.BaseWidget
	container     *fyne.Container
	modeLabel     *widget.Label
	filenameLabel *widget.Label
	positionLabel *widget.Label
	messageLabel  *widget.Label
	languageLabel *widget.Label
	background    *canvas.Rectangle
}

func NewStatusBar() *StatusBar {
	sb := &StatusBar{
		modeLabel:     widget.NewLabel("NORMAL"),
		filenameLabel: widget.NewLabel("Untitled"),
		positionLabel: widget.NewLabel("Ln 1, Col 0"),
		messageLabel:  widget.NewLabel(""),
		languageLabel: widget.NewLabel("text"),
		background:    canvas.NewRectangle(color.NRGBA{R: 40, G: 40, B: 40, A: 255}),
	}

	sb.modeLabel.TextStyle = fyne.TextStyle{Bold: true}

	sb.container = container.NewStack(
		sb.background,
		container.NewHBox(
			container.NewPadded(sb.modeLabel),
			widget.NewSeparator(),
			container.NewPadded(sb.filenameLabel),
			widget.NewSeparator(),
			container.NewPadded(sb.positionLabel),
			widget.NewSeparator(),
			container.NewPadded(sb.languageLabel),
			widget.NewSeparator(),
			container.NewPadded(sb.messageLabel),
		),
	)

	return sb
}

func (sb *StatusBar) SetMode(mode string) {
	sb.modeLabel.SetText(mode)

	switch mode {
	case "NORMAL":
		sb.modeLabel.TextStyle.Bold = true
	case "INSERT":
		sb.modeLabel.TextStyle.Bold = true
	case "VISUAL":
		sb.modeLabel.TextStyle.Bold = true
	case "COMMAND":
		sb.modeLabel.TextStyle.Bold = true
	default:
		sb.modeLabel.TextStyle.Bold = true
	}

	sb.modeLabel.Refresh()
}

func (sb *StatusBar) SetFilename(filename string) {
	sb.filenameLabel.SetText(filename)
}

func (sb *StatusBar) SetPosition(line, column int) {
	sb.positionLabel.SetText(fmt.Sprintf("Ln %d, Col %d", line+1, column))
}

func (sb *StatusBar) SetLanguage(language string) {
	sb.languageLabel.SetText(language)
}

func (sb *StatusBar) ShowTemporaryMessage(message string) {
	sb.messageLabel.SetText(message)

	go func() {
		time.Sleep(3 * time.Second)
		sb.messageLabel.SetText("")
	}()
}

func (sb *StatusBar) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(sb.container)
}

func (sb *StatusBar) MinSize() fyne.Size {
	return sb.container.MinSize()
}

func (sb *StatusBar) Move(position fyne.Position) {
	sb.container.Move(position)
}

func (sb *StatusBar) Position() fyne.Position {
	return sb.container.Position()
}

func (sb *StatusBar) Size() fyne.Size {
	return sb.container.Size()
}

func (sb *StatusBar) Resize(size fyne.Size) {
	sb.container.Resize(size)
}

func (sb *StatusBar) Visible() bool {
	return sb.container.Visible()
}

func (sb *StatusBar) ExtendBaseWidget(w fyne.Widget) {
	if w == nil {
		sb.BaseWidget.ExtendBaseWidget(sb)
	} else {
		sb.BaseWidget.ExtendBaseWidget(w)
	}
}

func (sb *StatusBar) Hide() {
	sb.container.Hide()
}

func (sb *StatusBar) Show() {
	sb.container.Show()
}

func (sb *StatusBar) Refresh() {
	if sb.container == nil {
		sb.container = container.NewStack(sb.background) // Ensure container is initialized
	}
	sb.container.Refresh()
}
