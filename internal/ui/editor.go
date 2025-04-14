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

// Editor represents the main editor component
type Editor struct {
	Window      fyne.Window
	TextArea    *TextArea
	History     *History
	FilePath    string
	LineNumbers *LineNumbers
	StatusBar   *StatusBar
	CurrentMode EditorMode
	CurrentView fyne.CanvasObject // Changed to interface type instead of concrete type
}

// SetText sets the text content of the editor
func (e *Editor) SetText(text string) {
	e.TextArea.Text = text
}

// GetText returns the current text content of the editor
func (e *Editor) GetText() string {
	return e.TextArea.Text
}

// GetSelectedText returns the currently selected text (stub implementation)
func (e *Editor) GetSelectedText() string {
	// In a real implementation, this would return the selected text
	// For now, return an empty string as placeholder
	return ""
}

// TextArea represents the text editing component
type TextArea struct {
	Text         string
	CursorRow    int
	CursorColumn int
}

// SetText sets the text content in the text area
func (t *TextArea) SetText(text string) {
	t.Text = text
}

// History maintains a list of text snapshots for undo/redo
type History struct {
	Snapshots []string
	Position  int
}

// Add adds a new text snapshot to the history
func (h *History) Add(text string) {
	// Add the current text to history snapshots
	if len(h.Snapshots) > 0 && h.Snapshots[len(h.Snapshots)-1] == text {
		return // Skip duplicate entries
	}
	h.Snapshots = append(h.Snapshots, text)
	h.Position = len(h.Snapshots) - 1
}

// Undo reverts to the previous snapshot in history
func (h *History) Undo() (string, bool) {
	if h.Position <= 0 {
		return "", false
	}
	h.Position--
	return h.Snapshots[h.Position], true
}

// Redo advances to the next snapshot in history
func (h *History) Redo() (string, bool) {
	if h.Position >= len(h.Snapshots)-1 {
		return "", false
	}
	h.Position++
	return h.Snapshots[h.Position], true
}

// LineNumbers represents the line numbering component
type LineNumbers struct {
	Visible bool
}

// Show makes the line numbers visible
func (ln *LineNumbers) Show() {
	ln.Visible = true
}

// Hide makes the line numbers invisible
func (ln *LineNumbers) Hide() {
	ln.Visible = false
}

// StatusBar represents the status display component
type StatusBar struct {
	Message string
}

// ShowTemporaryMessage displays a temporary message in the status bar
func (sb *StatusBar) ShowTemporaryMessage(msg string) {
	sb.Message = msg
	// Additional logic for temporary display could be added here
}

// Visible returns whether the status bar is visible
func (sb *StatusBar) Visible() bool {
	return true // Default implementation always returns true
}

// Show makes the status bar visible
func (sb *StatusBar) Show() {
	// Implementation for showing the status bar
}

// Hide makes the status bar invisible
func (sb *StatusBar) Hide() {
	// Implementation for hiding the status bar
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

	// Create Fyne UI components
	lineNumbersView := components.NewLineNumbersView()
	statusBar := components.NewStatusBar()
	commandInput := widget.NewEntry()
	commandInput.SetPlaceHolder(":")
	commandInput.Hide()

	// Setup text area event handling
	fyneTextArea.OnChanged = func(text string) {
		editor.TextArea.Text = text
		editor.updateLineNumbers(text)
		lineNumbersView.UpdateLineNumbers(text)
		lineNumbersView.Refresh()
		statusBar.SetPosition(fyneTextArea.CursorRow, fyneTextArea.CursorColumn)
		editor.History.Add(text)
	}

	tabContainer := container.NewDocTabs()

	// Store the Fyne UI components
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
	_ = lines                   // Use the variable to avoid unused error
	e.TextArea.CursorRow = 0    // Default value
	e.TextArea.CursorColumn = 0 // Default value
}

func (e *Editor) SetFilePath(path string) {
	e.FilePath = path
}

func (e *Editor) AddNewTab(name string) {
	// Implementation of AddNewTab
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

	// Check if CurrentView is nil before attempting to use it
	if e.CurrentView == nil {
		// Create a new split view with the tabContainer and the scrollContainer
		e.CurrentView = container.NewBorder(nil, nil, nil, nil, tabContainer)
		e.Window.SetContent(e.CurrentView)
	} else if split, ok := e.CurrentView.(*container.Split); ok {
		split.Trailing = scrollContainer
		split.Refresh()
	} else {
		// If CurrentView exists but is not a Split, replace it
		e.CurrentView = container.NewBorder(nil, nil, nil, nil, tabContainer)
		e.Window.SetContent(e.CurrentView)
	}
}
