package ui

import (
	"fyne.io/fyne/v2/widget"
)

// TextEditHelper manages the interactions between the Editor and Fyne's Entry widget
type TextEditHelper struct {
	editor      *Editor
	entryWidget *widget.Entry
}

// NewTextEditHelper creates a new helper for text editing operations
func NewTextEditHelper(editor *Editor, entry *widget.Entry) *TextEditHelper {
	return &TextEditHelper{
		editor:      editor,
		entryWidget: entry,
	}
}

// GetText returns the current text content
func (t *TextEditHelper) GetText() string {
	return t.entryWidget.Text
}

// SetText sets the text content and updates the entry widget
func (t *TextEditHelper) SetText(text string) {
	t.entryWidget.SetText(text)
	t.editor.TextArea.Text = text
}

// GetSelectedText returns the currently selected text from the entry widget
func (t *TextEditHelper) GetSelectedText() string {
	return t.entryWidget.SelectedText()
}

// CursorPosition returns the current cursor position in the text
func (t *TextEditHelper) CursorPosition() (row, col int) {
	row = t.entryWidget.CursorRow
	col = t.entryWidget.CursorColumn
	return row, col
}

// SyncText ensures the Editor's state matches the Entry widget
func (t *TextEditHelper) SyncText() {
	t.editor.TextArea.Text = t.entryWidget.Text
	t.editor.TextArea.CursorRow = t.entryWidget.CursorRow
	t.editor.TextArea.CursorColumn = t.entryWidget.CursorColumn
}
