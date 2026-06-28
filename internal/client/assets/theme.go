package assets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CompactTheme struct {
	fyne.Theme
}

func (t CompactTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameLineSpacing:
		return 0
	}
	return t.Theme.Size(name)
}
