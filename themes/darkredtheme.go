package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var DarkRedTheme = &darkRedTheme{}

type darkRedTheme struct{}

func (d *darkRedTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{80, 0, 0, 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{110, 0, 0, 255}
	case theme.ColorNameButton, theme.ColorNameDisabledButton, theme.ColorNameHover, theme.ColorNameFocus:
		return theme.DefaultTheme().Color(name, variant)
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNamePrimary:
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (d *darkRedTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (d *darkRedTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (d *darkRedTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
