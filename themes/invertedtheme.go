package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var InvertedTheme = &invertedTheme{}

type invertedTheme struct{}

func (i *invertedTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
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

func (i *invertedTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (i *invertedTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (i *invertedTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
