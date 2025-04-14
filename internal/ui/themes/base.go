package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type BaseTheme struct{}

var _ fyne.Theme = (*BaseTheme)(nil)

func NewBaseTheme() *BaseTheme {
	return &BaseTheme{}
}

func (t *BaseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 28, G: 28, B: 35, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 238, G: 238, B: 238, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 143, G: 188, B: 187, A: 255}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 94, G: 186, B: 235, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 50, G: 120, B: 150, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 36, G: 40, B: 59, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 45, G: 45, B: 55, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 120, G: 120, B: 120, A: 255}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 70, G: 70, B: 85, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 37, G: 38, B: 45, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 120, G: 120, B: 140, A: 255}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 60, G: 80, B: 120, A: 128}
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 33, G: 34, B: 44, A: 255}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *BaseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *BaseTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *BaseTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameScrollBar:
		return 8
	case theme.SizeNameScrollBarSmall:
		return 4
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 18
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInputBorder:
		return 1
	}

	return theme.DefaultTheme().Size(name)
}
