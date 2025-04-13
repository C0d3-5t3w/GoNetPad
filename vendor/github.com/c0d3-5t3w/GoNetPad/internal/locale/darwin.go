package locale

import (
	"os"
	"strings"
)

func getSystemLocale() (string, error) {
	for _, env := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		if val := os.Getenv(env); val != "" {
			val = strings.Split(val, ".")[0]
			return strings.Replace(val, "_", "-", 1), nil
		}
	}

	return "en-US", nil
}
