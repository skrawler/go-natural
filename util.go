package natural

import (
	"fmt"
	"strings"
)

// upper case first letter in s
func ucFirst(s string) string {

	first := string([]rune(s)[0:1])
	rest := ""
	if len(s) > 1 {
		rest = string([]rune(s)[1:])
	}
	return strings.ToUpper(first) + strings.ToLower(rest)
}

func arrayIndex(s string, arr []string) (int, error) {
	for idx, x := range arr {
		if x == s {
			return idx, nil
		}
	}
	return 0, fmt.Errorf("not found")
}
