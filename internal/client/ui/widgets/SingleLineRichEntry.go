package widgets

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type SingleLineRichEntry struct {
	widget.Entry
	isShiftHeld bool
}

const (
	minVisibleEntryRows = 1
	maxVisibleEntryRows = 5
)

func NewSingleLineRichEntry() *SingleLineRichEntry {
	e := &SingleLineRichEntry{}
	e.ExtendBaseWidget(e)
	e.Wrapping = fyne.TextWrapWord
	e.MultiLine = true

	// Entry automatically becomes scrollable when there are more rows than maxVisibleEntryRows
	// e.Scroll = fyne.ScrollVerticalOnly

	e.SetMinRowsVisible(minVisibleEntryRows)
	e.OnChanged = func(text string) {
		rows := min(strings.Count(text, "\n")+1, maxVisibleEntryRows)
		e.SetMinRowsVisible(rows)
	}
	return e
}

func (e *SingleLineRichEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		if e.isShiftHeld == true {
			fmt.Println("Shift is being help: ", e.isShiftHeld)
			fmt.Println("key: ", key)
			// Yes shift -> new line
			e.Entry.TypedKey(key)
		} else if e.OnSubmitted != nil {
			// No shift → submit
			e.OnSubmitted(e.Text)
		}
		return
	}
	e.Entry.TypedKey(key)
}

func (e *SingleLineRichEntry) KeyDown(key *fyne.KeyEvent) {
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.isShiftHeld = true
	}
}

func (e *SingleLineRichEntry) KeyUp(key *fyne.KeyEvent) {
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		e.isShiftHeld = false
	}
}
