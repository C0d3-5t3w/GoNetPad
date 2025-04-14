package ui

import (
	"fyne.io/fyne/v2/widget"
)

type TextEditHelper struct {
	editor      *Editor
	entryWidget *widget.Entry
}

func NewTextEditHelper(editor *Editor, entry *widget.Entry) *TextEditHelper {
	return &TextEditHelper{
		editor:      editor,
		entryWidget: entry,
	}
}

func (t *TextEditHelper) GetText() string {
	return t.entryWidget.Text
}

func (t *TextEditHelper) SetText(text string) {
	t.entryWidget.SetText(text)
	t.editor.TextArea.Text = text
}

func (t *TextEditHelper) GetSelectedText() string {
	return t.entryWidget.SelectedText()
}

func (t *TextEditHelper) CursorPosition() (row, col int) {
	row = t.entryWidget.CursorRow
	col = t.entryWidget.CursorColumn
	return row, col
}

func (t *TextEditHelper) SyncText() {
	t.editor.TextArea.Text = t.entryWidget.Text
	t.editor.TextArea.CursorRow = t.entryWidget.CursorRow
	t.editor.TextArea.CursorColumn = t.entryWidget.CursorColumn
}
