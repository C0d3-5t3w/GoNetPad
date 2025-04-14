package locale

import (
	"os"
	"strings"

	"github.com/jeandeaual/go-locale"
)

func getSystemLocale() (string, error) {
	tag, err := locale.GetLocale()
	if err != nil {
		return fallbackGetSystemLocale()
	}

	return strings.Replace(tag, "_", "-", 1), nil
}

func fallbackGetSystemLocale() (string, error) {
	for _, env := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		if val := os.Getenv(env); val != "" {
			val = strings.Split(val, ".")[0]
			return strings.Replace(val, "_", "-", 1), nil
		}
	}

	return "en-US", nil
}
