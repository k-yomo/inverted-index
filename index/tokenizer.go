package index

import (
	"strings"
)

func tokenize(text string) []string {
	return strings.Split(text, " ")
}
