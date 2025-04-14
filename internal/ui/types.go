package ui

// import (
// 	"fyne.io/fyne/v2"
// )

// type Editor struct {
// 	Window      fyne.Window
// 	CurrentView fyne.CanvasObject
// 	TextArea    *TextArea
// 	LineNumbers *LineNumbers
// 	StatusBar   *StatusBar
// 	History     *History
// 	FilePath    string
// }

// type TextArea struct {
// 	Text         string
// 	CursorRow    int
// 	CursorColumn int
// }

// type LineNumbers struct {
// 	Visible bool
// }

// type StatusBar struct {
// 	Message string
// }

// type History struct {
// 	Snapshots []string
// 	Position  int
// }

// func (h *History) Undo() (string, bool) {
// 	if h.Position <= 0 {
// 		return "", false
// 	}
// 	h.Position--
// 	return h.Snapshots[h.Position], true
// }

// func (h *History) Redo() (string, bool) {
// 	if h.Position >= len(h.Snapshots)-1 {
// 		return "", false
// 	}
// 	h.Position++
// 	return h.Snapshots[h.Position], true
// }

// func (ln *LineNumbers) Visible() bool {
// 	return ln.Visible
// }

// func (sb *StatusBar) Visible() bool {
// 	return sb.Message != ""
// }

// func (sb *StatusBar) Hide() {
// 	sb.Message = ""
// }

// func (sb *StatusBar) Show() {
// 	sb.Message = "Ready"
// }

// func (ta *TextArea) SelectedText() string {

// 	return ""
// }

// func (ta *TextArea) SetText(text string) {
// 	ta.Text = text
// }
