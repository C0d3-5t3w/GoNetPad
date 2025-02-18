package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/c0d3-5t3w/GoNetPad/helpers"
)

type Editor struct {
	TextArea    *widget.Entry // Changed from TextEdit to Entry
	LineNumbers *widget.Label
	Window      fyne.Window
	History     *TextHistory
}

func NewEditor(window fyne.Window) *Editor {
	editor := &Editor{
		TextArea:    widget.NewMultiLineEntry(), // Using NewMultiLineEntry instead of NewTextEdit
		LineNumbers: widget.NewLabel("1"),
		Window:      window,
		History:     NewTextHistory(),
	}

	editor.TextArea.SetPlaceHolder("Enter Text Here...")
	editor.setupUI()
	return editor
}

func (e *Editor) setupUI() {
	formatBtn := widget.NewButton("Format Code", e.handleFormatCode)
	content := container.NewBorder(nil, container.NewHBox(formatBtn), e.LineNumbers, nil, e.TextArea)
	scroll := container.NewScroll(content)
	e.Window.SetContent(scroll)
}

func (e *Editor) handleFormatCode() {
	formatted, err := helpers.FormatCode(e.TextArea.Text) // Using Text property instead of Text() method
	if err != nil {
		dialog.ShowError(err, e.Window)
		return
	}
	e.TextArea.SetText(formatted)
}

// Additional methods would go here...
