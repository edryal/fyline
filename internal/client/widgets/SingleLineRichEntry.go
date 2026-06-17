package widgets

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SingleLineRichEntry struct {
	widget.Entry
}

const minVisibleEntryRows = 1
const maxVisibleEntryRows = 5

func NewSingleLineRichEntry() *SingleLineRichEntry {
	e := &SingleLineRichEntry{}
	e.ExtendBaseWidget(e)
	e.Wrapping = fyne.TextWrapWord
	e.MultiLine = true

	// Entry automatically becomes scrollable when there are more rows than maxVisibleEntryRows
	// e.Scroll = fyne.ScrollVerticalOnly

	e.SetMinRowsVisible(minVisibleEntryRows)
	e.OnChanged = func(text string) {
		rows := min(strings.Count(text, "\n") + 1, maxVisibleEntryRows)
		e.SetMinRowsVisible(rows)
	}
	return e
}
