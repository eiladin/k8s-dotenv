package parser

import (
	"fmt"
	"strings"
)

func Parse(key string, value []byte) string {
	return fmt.Sprintf("export %s=\"%s\"\n", strings.ReplaceAll(key, ".", ""), strings.ReplaceAll(string(value), "\n", "\\n"))
}

func ParseStr(key, value string) string {
	return Parse(key, []byte(value))
}
