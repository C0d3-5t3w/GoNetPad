package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var OriginalTheme = &originalTheme{}

type originalTheme struct{}

func (o *originalTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.White
	case theme.ColorNameInputBackground:
		return color.White
	case theme.ColorNameButton, theme.ColorNameDisabledButton, theme.ColorNameHover, theme.ColorNameFocus:
		return theme.DefaultTheme().Color(name, variant)
	case theme.ColorNameForeground:
		return color.Black
	case theme.ColorNamePrimary:
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (o *originalTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (o *originalTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (o *originalTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
