package ui

type TextHistory struct {
	entries []string
	current int
}

func NewTextHistory() *TextHistory {
	return &TextHistory{
		entries: []string{""},
		current: 0,
	}
}

func (h *TextHistory) Add(text string) {
	h.entries = append(h.entries[:h.current+1], text)
	h.current++
}

func (h *TextHistory) Undo() (string, bool) {
	if h.current > 0 {
		h.current--
		return h.entries[h.current], true
	}
	return h.entries[h.current], false
}

func (h *TextHistory) Redo() (string, bool) {
	if h.current < len(h.entries)-1 {
		h.current++
		return h.entries[h.current], true
	}
	return h.entries[h.current], false
}
