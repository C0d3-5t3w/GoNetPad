package locale

import (
	"fyne.io/fyne/v2"
)

func Configure(app fyne.App) error {
	_, err := getSystemLocale()
	if err != nil {
		return err
	}

	return nil
}
