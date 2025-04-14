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
	// If we're not at the end of the history, truncate
	if h.current < len(h.entries)-1 {
		h.entries = h.entries[:h.current+1]
	}

	// Don't add if identical to current
	if len(h.entries) > 0 && h.entries[h.current] == text {
		return
	}

	h.entries = append(h.entries, text)
	h.current = len(h.entries) - 1
}

func (h *TextHistory) Undo() (string, bool) {
	if h.current > 0 {
		h.current--
		return h.entries[h.current], true
	}
	return h.entries[0], false
}

func (h *TextHistory) Redo() (string, bool) {
	if h.current < len(h.entries)-1 {
		h.current++
		return h.entries[h.current], true
	}
	return h.entries[h.current], false
}
