package ui

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/c0d3-5t3w/GoNetPad/internal/tools"
)

type Editor struct {
	TextArea    *widget.Entry
	LineNumbers *widget.Label
	Window      fyne.Window
	History     *TextHistory
	Language    string
	FilePath    string
	RichText    *widget.RichText
}

func NewEditor(window fyne.Window) *Editor {
	editor := &Editor{
		TextArea:    widget.NewMultiLineEntry(),
		LineNumbers: widget.NewLabel("1"),
		Window:      window,
		History:     NewTextHistory(),
		Language:    "text",
	}

	editor.TextArea.SetPlaceHolder("Enter Text Here...")
	editor.setupUI()
	return editor
}

func (e *Editor) setupUI() {
	formatBtn := widget.NewButton("Format Code", e.handleFormatCode)

	e.TextArea.OnChanged = func(text string) {
		e.updateSyntaxHighlighting(text)
	}

	languageOptions := []string{"text", "go", "javascript", "typescript", "html", "css"}
	languageSelector := widget.NewSelect(languageOptions, func(selected string) {
		e.Language = selected
		e.updateSyntaxHighlighting(e.TextArea.Text)
	})
	languageSelector.SetSelected("text")

	controls := container.NewHBox(formatBtn, widget.NewLabel("Language:"), languageSelector)
	content := container.NewBorder(nil, controls, e.LineNumbers, nil, e.TextArea)
	scroll := container.NewScroll(content)
	e.Window.SetContent(scroll)
}

func (e *Editor) handleFormatCode() {
	formatted, err := tools.FormatCode(e.TextArea.Text)
	if err != nil {
		dialog.ShowError(err, e.Window)
		return
	}
	e.TextArea.SetText(formatted)
}

func (e *Editor) SetFilePath(path string) {
	e.FilePath = path
	if path != "" {

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go":
			e.Language = "go"
		case ".js":
			e.Language = "javascript"
		case ".ts":
			e.Language = "typescript"
		case ".html":
			e.Language = "html"
		case ".css":
			e.Language = "css"
		default:

			e.Language = tools.DetectLanguage(path, e.TextArea.Text)
		}
	}
}

func (e *Editor) updateSyntaxHighlighting(text string) {

	if e.Language == "text" && e.FilePath == "" {
		e.Language = tools.DetectLanguage(e.FilePath, text)
	}

	if e.RichText == nil {
		e.RichText = tools.GenerateSyntaxHighlightedRichText(text, e.Language)

	} else {

		newRichText := tools.GenerateSyntaxHighlightedRichText(text, e.Language)
		e.RichText.Segments = newRichText.Segments
		e.RichText.Refresh()
	}
}
