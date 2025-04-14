package ui

import (
	"path/filepath"
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
	TextArea      *widget.Entry
	LineNumbers   *components.LineNumbersView
	Window        fyne.Window
	History       *TextHistory
	Language      string
	FilePath      string
	RichText      *widget.RichText
	StatusBar     *components.StatusBar
	CommandInput  *widget.Entry
	CurrentMode   EditorMode
	TabContainer  *container.DocTabs
	CurrentView   *container.Split
	MainContainer *fyne.Container
}

func NewEditor(window fyne.Window) *Editor {
	editor := &Editor{
		TextArea:    widget.NewMultiLineEntry(),
		LineNumbers: components.NewLineNumbersView(),
		Window:      window,
		History:     NewTextHistory(),
		Language:    "text",
		CurrentMode: NormalMode,
	}

	editor.TextArea.SetPlaceHolder("Enter Text Here...")
	editor.StatusBar = components.NewStatusBar()
	editor.CommandInput = widget.NewEntry()
	editor.CommandInput.SetPlaceHolder(":")
	editor.CommandInput.Hide()

	editor.TabContainer = container.NewDocTabs()

	editor.setupUI()
	editor.setupKeyBindings()
	return editor
}

func (e *Editor) setupUI() {
	formatBtn := widget.NewButtonWithIcon("Format", theme.DocumentSaveIcon(), e.handleFormatCode)

	languageOptions := []string{"text", "go", "javascript", "typescript", "html", "css", "python", "rust", "c", "c++", "java"}
	languageSelector := widget.NewSelect(languageOptions, func(selected string) {
		e.Language = selected
		e.updateSyntaxHighlighting(e.TextArea.Text)
		e.StatusBar.SetLanguage(selected)
	})
	languageSelector.SetSelected("text")

	e.TextArea.OnChanged = func(text string) {
		e.updateSyntaxHighlighting(text)
		e.updateLineNumbers(text)
		e.StatusBar.SetPosition(e.TextArea.CursorRow, e.TextArea.CursorColumn)
		e.History.Add(text) // Add to history on change
	}

	e.CommandInput.OnSubmitted = func(cmd string) {
		e.executeCommand(cmd)
		e.CommandInput.Hide()
		e.CurrentMode = NormalMode
		e.StatusBar.SetMode("NORMAL")
		e.TextArea.FocusGained()
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
		e.LineNumbers,
		e.TextArea,
	)
	editorWithLineNumbers.Offset = 0.05

	mainContent := container.NewBorder(
		toolbar,
		container.NewVBox(
			e.CommandInput,
			e.StatusBar,
		),
		nil, nil,
		editorWithLineNumbers,
	)

	firstTab := container.NewTabItem("Untitled", mainContent)
	e.TabContainer.Append(firstTab)

	e.MainContainer = container.NewBorder(
		nil, nil, nil, nil,
		e.TabContainer,
	)

	e.Window.SetContent(e.MainContainer)
}

func (e *Editor) setupKeyBindings() {
	e.Window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if e.CurrentMode == NormalMode {
			switch key.Name {
			case fyne.KeyI:
				e.CurrentMode = InsertMode
				e.StatusBar.SetMode("INSERT")
				return
			case fyne.KeyV:
				e.CurrentMode = VisualMode
				e.StatusBar.SetMode("VISUAL")
				return
			case fyne.KeyEscape:
				e.CurrentMode = NormalMode
				e.StatusBar.SetMode("NORMAL")
				return
			case fyne.KeySemicolon:
				e.CurrentMode = CommandMode
				e.StatusBar.SetMode("COMMAND")
				e.CommandInput.Show()
				e.CommandInput.FocusGained()
				return
			}
		} else if e.CurrentMode == InsertMode && key.Name == fyne.KeyEscape {
			e.CurrentMode = NormalMode
			e.StatusBar.SetMode("NORMAL")
			return
		} else if e.CurrentMode == VisualMode && key.Name == fyne.KeyEscape {
			e.CurrentMode = NormalMode
			e.StatusBar.SetMode("NORMAL")
			return
		}
	})
}

func (e *Editor) executeCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)

	switch {
	case cmd == "q" || cmd == "quit":
		e.Window.Close()
	case cmd == "w" || cmd == "write":
		e.saveFile()
	case cmd == "wq":
		e.saveFile()
		e.Window.Close()
	case strings.HasPrefix(cmd, "set"):
		parts := strings.Split(cmd, " ")
		if len(parts) >= 2 {
			option := parts[1]
			switch option {
			case "number":
				e.LineNumbers.Show()
			case "nonumber":
				e.LineNumbers.Hide()
			}
		}
	}
}

func (e *Editor) handleFormatCode() {
	formatted, err := tools.FormatCode(e.TextArea.Text)
	if err != nil {
		dialog.ShowError(err, e.Window)
		return
	}
	e.TextArea.SetText(formatted)
	e.StatusBar.ShowTemporaryMessage("Code formatted")
}

func (e *Editor) updateLineNumbers(text string) {
	lines := strings.Split(text, "\n")
	e.LineNumbers.SetLineCount(len(lines))
	e.LineNumbers.SetCurrentLine(e.TextArea.CursorRow)
}

func (e *Editor) SetFilePath(path string) {
	e.FilePath = path
	if path != "" {
		if e.TabContainer.Selected() != nil {
			filename := filepath.Base(path)
			e.TabContainer.Selected().Text = filename
			e.TabContainer.Refresh()
		}

		e.StatusBar.SetFilename(filepath.Base(path))

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
		case ".py":
			e.Language = "python"
		case ".rs":
			e.Language = "rust"
		case ".c":
			e.Language = "c"
		case ".cpp", ".cc", ".cxx":
			e.Language = "c++"
		case ".java":
			e.Language = "java"
		default:
			e.Language = tools.DetectLanguage(path, e.TextArea.Text)
		}

		e.StatusBar.SetLanguage(e.Language)
	}
}

func (e *Editor) AddNewTab(name string) {
	textArea := widget.NewMultiLineEntry()
	lineNumbers := components.NewLineNumbersView()

	editorWithLineNumbers := container.NewHSplit(
		lineNumbers,
		textArea,
	)
	editorWithLineNumbers.Offset = 0.05

	newTab := container.NewTabItem(name, editorWithLineNumbers)
	e.TabContainer.Append(newTab)
	e.TabContainer.Select(newTab)
}

func (e *Editor) saveFile() {
	e.StatusBar.ShowTemporaryMessage("File saved")
}

func (e *Editor) updateSyntaxHighlighting(text string) {
	if e.Language == "text" && e.FilePath != "" {
		e.Language = tools.DetectLanguage(e.FilePath, text)
	}

	highlighter := tools.NewSyntaxHighlighter(e.Language)
	coloredCodeView := highlighter.HighlightCode(text)

	size := e.TextArea.Size()

	scrollContainer := container.NewScroll(coloredCodeView)
	scrollContainer.Resize(size)

	if e.CurrentView == nil {
		e.CurrentView = container.NewHSplit(
			e.MainContainer,
			scrollContainer,
		)
		e.CurrentView.Offset = 0.6
		e.Window.SetContent(e.CurrentView)
	} else {
		e.CurrentView.Trailing = scrollContainer
		e.CurrentView.Refresh()
	}
}
