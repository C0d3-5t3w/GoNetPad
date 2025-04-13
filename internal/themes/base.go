package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var BaseTheme = &baseTheme{}

type baseTheme struct{}

func (i *baseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.Black
	case theme.ColorNameInputBackground:
		return color.Black
	case theme.ColorNameButton, theme.ColorNameDisabledButton, theme.ColorNameHover, theme.ColorNameFocus:
		return theme.DefaultTheme().Color(name, variant)
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNamePrimary:
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (i *baseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (i *baseTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (i *baseTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
