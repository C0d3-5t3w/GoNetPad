package components

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LineNumbersView struct {
	widget.BaseWidget
	container   *fyne.Container
	lineNumbers []string
	currentLine int
	background  *canvas.Rectangle
	labels      []*canvas.Text
	visible     bool
}

func NewLineNumbersView() *LineNumbersView {
	ln := &LineNumbersView{
		lineNumbers: []string{"1"},
		currentLine: 0,
		background:  canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255}),
		visible:     true,
		labels:      make([]*canvas.Text, 0),
	}

	ln.updateDisplay()
	ln.ExtendBaseWidget(ln)

	return ln
}

func (ln *LineNumbersView) SetLineCount(count int) {
	if count < 1 {
		count = 1
	}

	ln.lineNumbers = make([]string, count)
	for i := 0; i < count; i++ {
		ln.lineNumbers[i] = strconv.Itoa(i + 1)
	}

	ln.updateDisplay()
}

func (ln *LineNumbersView) SetCurrentLine(line int) {
	ln.currentLine = line
	ln.updateDisplay()
}

func (ln *LineNumbersView) Hide() {
	ln.visible = false
	ln.Refresh()
}

func (ln *LineNumbersView) Show() {
	ln.visible = true
	ln.Refresh()
}

func (ln *LineNumbersView) MinSize() fyne.Size {
	width := float32(30)

	maxDigits := len(fmt.Sprintf("%d", len(ln.lineNumbers)))
	if maxDigits > 2 {
		width = float32(maxDigits * 10)
	}

	return fyne.NewSize(width, 0)
}

func (ln *LineNumbersView) CreateRenderer() fyne.WidgetRenderer {
	ln.updateDisplay()
	return widget.NewSimpleRenderer(ln.container)
}

func (ln *LineNumbersView) updateDisplay() {
	if !ln.visible {
		if ln.container == nil {
			ln.container = container.NewStack(ln.background)
		}
		return
	}

	objects := []fyne.CanvasObject{ln.background}
	vbox := container.NewVBox()

	ln.labels = make([]*canvas.Text, len(ln.lineNumbers))

	for i, num := range ln.lineNumbers {
		label := canvas.NewText(num, color.White)
		label.Alignment = fyne.TextAlignTrailing
		label.TextStyle = fyne.TextStyle{Monospace: true}

		if i == ln.currentLine {
			label.Color = color.NRGBA{R: 255, G: 200, B: 0, A: 255}
			label.TextStyle.Bold = true
		}

		ln.labels[i] = label
		vbox.Add(label)
	}

	objects = append(objects, vbox)
	ln.container = container.NewStack(objects...) // Ensure container is initialized
	ln.Refresh()
}

func (ln *LineNumbersView) ExtendBaseWidget(w fyne.Widget) {
	if w == nil {
		ln.BaseWidget.ExtendBaseWidget(ln)
	} else {
		ln.BaseWidget.ExtendBaseWidget(w)
	}
}

func (ln *LineNumbersView) Move(position fyne.Position) {
	if ln.container != nil {
		ln.container.Move(position)
	}
}

func (ln *LineNumbersView) Position() fyne.Position {
	if ln.container != nil {
		return ln.container.Position()
	}
	return fyne.NewPos(0, 0)
}

func (ln *LineNumbersView) Size() fyne.Size {
	if ln.container != nil {
		return ln.container.Size()
	}
	return ln.MinSize()
}

func (ln *LineNumbersView) Resize(size fyne.Size) {
	if ln.container != nil {
		ln.container.Resize(size)
	}
}

func (ln *LineNumbersView) Visible() bool {
	return ln.visible
}
