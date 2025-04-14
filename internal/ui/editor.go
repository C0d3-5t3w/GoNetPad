package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/C0d3-5t3w/GoNetPad/internal/ui/components"
	"github.com/C0d3-5t3w/GoNetPad/internal/ui/tools"
)

type EditorMode int

const (
	NormalMode EditorMode = iota
	InsertMode
	VisualMode
	CommandMode
)

type Editor struct {
	Window      fyne.Window
	TextArea    *TextArea
	History     *History
	FilePath    string
	LineNumbers *LineNumbers
	StatusBar   *StatusBar
	CurrentMode EditorMode
	CurrentView fyne.CanvasObject
}

func (e *Editor) SetText(text string) {
	e.TextArea.Text = text
}

func (e *Editor) GetText() string {
	return e.TextArea.Text
}

func (e *Editor) GetSelectedText() string {

	return ""
}

type TextArea struct {
	Text         string
	CursorRow    int
	CursorColumn int
}

func (t *TextArea) SetText(text string) {
	t.Text = text
}

type History struct {
	Snapshots []string
	Position  int
}

func (h *History) Add(text string) {

	if len(h.Snapshots) > 0 && h.Snapshots[len(h.Snapshots)-1] == text {
		return
	}
	h.Snapshots = append(h.Snapshots, text)
	h.Position = len(h.Snapshots) - 1
}

func (h *History) Undo() (string, bool) {
	if h.Position <= 0 {
		return "", false
	}
	h.Position--
	return h.Snapshots[h.Position], true
}

func (h *History) Redo() (string, bool) {
	if h.Position >= len(h.Snapshots)-1 {
		return "", false
	}
	h.Position++
	return h.Snapshots[h.Position], true
}

type LineNumbers struct {
	Visible bool
}

func (ln *LineNumbers) Show() {
	ln.Visible = true
}

func (ln *LineNumbers) Hide() {
	ln.Visible = false
}

type StatusBar struct {
	Message string
}

func (sb *StatusBar) ShowTemporaryMessage(msg string) {
	sb.Message = msg

}

func (sb *StatusBar) Visible() bool {
	return true
}

func (sb *StatusBar) Show() {

}

func (sb *StatusBar) Hide() {

}

func NewEditor(window fyne.Window) *Editor {
	editor := &Editor{
		Window:   window,
		TextArea: &TextArea{},
		History:  &History{},
		FilePath: "",
	}

	fyneTextArea := widget.NewMultiLineEntry()
	fyneTextArea.SetPlaceHolder("Enter Text Here...")

	editor.LineNumbers = &LineNumbers{}
	editor.StatusBar = &StatusBar{}

	lineNumbersView := components.NewLineNumbersView()
	statusBar := components.NewStatusBar()
	commandInput := widget.NewEntry()
	commandInput.SetPlaceHolder(":")
	commandInput.Hide()

	fyneTextArea.OnChanged = func(text string) {
		editor.TextArea.Text = text
		editor.updateLineNumbers(text)
		lineNumbersView.UpdateLineNumbers(text)
		lineNumbersView.Refresh()
		statusBar.SetPosition(fyneTextArea.CursorRow, fyneTextArea.CursorColumn)
		editor.History.Add(text)
	}

	tabContainer := container.NewDocTabs()

	setupUI(editor, fyneTextArea, lineNumbersView, statusBar, commandInput, tabContainer)
	setupKeyBindings(editor, fyneTextArea, commandInput, statusBar)

	return editor
}

func setupUI(e *Editor,
	textArea *widget.Entry,
	lineNumbers *components.LineNumbersView,
	statusBar *components.StatusBar,
	commandInput *widget.Entry,
	tabContainer *container.DocTabs) {

	formatBtn := widget.NewButtonWithIcon("Format", theme.DocumentSaveIcon(), func() {
		handleFormatCode(e, textArea, statusBar)
	})

	languageOptions := []string{"text", "go", "javascript", "typescript", "html", "css", "python", "rust", "c", "c++", "java"}
	languageSelector := widget.NewSelect(languageOptions, func(selected string) {
		updateSyntaxHighlighting(e, textArea.Text, selected, tabContainer)
		statusBar.SetLanguage(selected)
	})
	languageSelector.SetSelected("text")

	commandInput.OnSubmitted = func(cmd string) {
		executeCommand(e, cmd, lineNumbers)
		commandInput.Hide()
		e.CurrentMode = NormalMode
		statusBar.SetMode("NORMAL")
		textArea.FocusGained()
	}

	toolbar := container.NewHBox(
		formatBtn,
		widget.NewLabel("Language:"),
		languageSelector,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.ViewFullScreenIcon(), func() {
			e.Window.SetFullScreen(!e.Window.FullScreen())
		}),
	)

	editorWithLineNumbers := container.NewHSplit(
		lineNumbers,
		textArea,
	)
	editorWithLineNumbers.Offset = 0.05

	mainContent := container.NewBorder(
		toolbar,
		container.NewVBox(
			commandInput,
			statusBar,
		),
		nil, nil,
		editorWithLineNumbers,
	)

	firstTab := container.NewTabItem("Untitled", mainContent)
	tabContainer.Append(firstTab)

	mainContainer := container.NewBorder(
		nil, nil, nil, nil,
		tabContainer,
	)

	e.CurrentView = mainContainer
	e.Window.SetContent(mainContainer)
}

func setupKeyBindings(e *Editor, textArea *widget.Entry, commandInput *widget.Entry, statusBar *components.StatusBar) {
	e.Window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if e.CurrentMode == NormalMode {
			switch key.Name {
			case fyne.KeyI:
				e.CurrentMode = InsertMode
				statusBar.SetMode("INSERT")
				return
			case fyne.KeyV:
				e.CurrentMode = VisualMode
				statusBar.SetMode("VISUAL")
				return
			case fyne.KeyEscape:
				e.CurrentMode = NormalMode
				statusBar.SetMode("NORMAL")
				return
			case fyne.KeySemicolon:
				e.CurrentMode = CommandMode
				statusBar.SetMode("COMMAND")
				commandInput.Show()
				commandInput.FocusGained()
				return
			}
		} else if e.CurrentMode == InsertMode && key.Name == fyne.KeyEscape {
			e.CurrentMode = NormalMode
			statusBar.SetMode("NORMAL")
			return
		} else if e.CurrentMode == VisualMode && key.Name == fyne.KeyEscape {
			e.CurrentMode = NormalMode
			statusBar.SetMode("NORMAL")
			return
		}
	})
}

func executeCommand(e *Editor, cmd string, lineNumbers *components.LineNumbersView) {
	cmd = strings.TrimSpace(cmd)

	switch {
	case cmd == "q" || cmd == "quit":
		e.Window.Close()
	case cmd == "w" || cmd == "write":
		saveFile(e)
	case cmd == "wq":
		saveFile(e)
		e.Window.Close()
	case strings.HasPrefix(cmd, "set"):
		parts := strings.Split(cmd, " ")
		if len(parts) >= 2 {
			option := parts[1]
			switch option {
			case "number":
				lineNumbers.Show()
				e.LineNumbers.Show()
			case "nonumber":
				lineNumbers.Hide()
				e.LineNumbers.Hide()
			}
		}
	}
}

func handleFormatCode(e *Editor, textArea *widget.Entry, statusBar *components.StatusBar) {
	formatted, err := tools.FormatCode(textArea.Text)
	if err != nil {
		dialog.ShowError(err, e.Window)
		return
	}
	textArea.SetText(formatted)
	e.TextArea.Text = formatted
	statusBar.ShowTemporaryMessage("Code formatted")
}

func (e *Editor) updateLineNumbers(text string) {
	lines := strings.Split(text, "\n")
	_ = lines
	e.TextArea.CursorRow = 0
	e.TextArea.CursorColumn = 0
}

func (e *Editor) SetFilePath(path string) {
	e.FilePath = path
}

func (e *Editor) AddNewTab(name string) {

}

func saveFile(e *Editor) {
	e.StatusBar.ShowTemporaryMessage("File saved")
}

func updateSyntaxHighlighting(e *Editor, text string, language string, tabContainer *container.DocTabs) {
	if language == "text" && e.FilePath != "" {
		language = tools.DetectLanguage(e.FilePath, text)
	}

	highlighter := tools.NewSyntaxHighlighter(language)
	coloredCodeView := highlighter.HighlightCode(text)

	scrollContainer := container.NewScroll(coloredCodeView)

	if e.CurrentView == nil {

		e.CurrentView = container.NewBorder(nil, nil, nil, nil, tabContainer)
		e.Window.SetContent(e.CurrentView)
	} else if split, ok := e.CurrentView.(*container.Split); ok {
		split.Trailing = scrollContainer
		split.Refresh()
	} else {

		e.CurrentView = container.NewBorder(nil, nil, nil, nil, tabContainer)
		e.Window.SetContent(e.CurrentView)
	}
}
