package capture

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"fyne.io/fyne/v2"
	"github.com/kbinani/screenshot"
)

func CaptureWindow(w fyne.Window) (string, error) {
	// Get window size
	bounds := w.Canvas().Size()

	// Create a rectangle for the capture region
	captureRect := image.Rect(0, 0, int(bounds.Width), int(bounds.Height))

	// Capture the screen region
	img, err := screenshot.CaptureRect(captureRect)
	if err != nil {
		return "", err
	}

	// Encode to base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
