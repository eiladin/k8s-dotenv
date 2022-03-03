package parser

import (
	"fmt"
	"strings"
)

// Parse builds an export statement given a k/v pair.
func Parse(shouldExport bool, key string, value []byte) string {
	export := ""
	if shouldExport {
		export = "export "
	}

	return fmt.Sprintf("%s%s=\"%s\"\n",
		export,
		strings.ReplaceAll(key, ".", ""),
		strings.ReplaceAll(string(value), "\n", "\\n"),
	)
}

// ParseStr builds an export statement given a k/v pair.
func ParseStr(shouldExport bool, key, value string) string {
	return Parse(shouldExport, key, []byte(value))
}
