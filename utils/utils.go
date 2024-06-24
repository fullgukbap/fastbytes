package utils

import (
	"strings"
)

func CleanUp(arg any) (r string) {
	var m string

	switch v := arg.(type) {
	case []byte:
		m = string(v)
	}

	r = strings.ReplaceAll(m, `"`, "")
	return r
}

func IsEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) != 0 {
			return false
		}
	}

	return true
}
