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

	e.SetMinRowsVisible(minVisibleEntryRows)
	e.OnChanged = func(text string) {
		rows := strings.Count(text, "\n") + 1

		if rows > maxVisibleEntryRows {
			rows = maxVisibleEntryRows
		}

		e.SetMinRowsVisible(rows)
	}
	//e.Scroll = fyne.ScrollVerticalOnly
	return e
}
