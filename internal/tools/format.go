package tools

import (
	"go/format"
)

func FormatCode(code string) (string, error) {
	if code == "" {
		return "", nil
	}

	formatted, err := format.Source([]byte(code))
	if err != nil {
		return code, err
	}

	return string(formatted), nil
}
