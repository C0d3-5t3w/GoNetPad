package locale

import (
	"os"
	"strings"
)

func getSystemLocale() (string, error) {
	// Try getting locale from environment variables first
	for _, env := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		if val := os.Getenv(env); val != "" {
			// Clean up the locale string (e.g., "en_US.UTF-8" -> "en-US")
			val = strings.Split(val, ".")[0]
			return strings.Replace(val, "_", "-", 1), nil
		}
	}

	// Fallback to en-US if no locale is found
	return "en-US", nil
}
