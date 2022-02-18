package parser

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite
}

var cases = []struct {
	shouldExport bool
	key          string
	value        string
	expected     string
}{
	{shouldExport: false, key: "key", value: "value", expected: "key=\"value\"\n"},
	{shouldExport: true, key: "key", value: "value", expected: "export key=\"value\"\n"},
	{shouldExport: false, key: "k.e.y", value: "value", expected: "key=\"value\"\n"},
	{shouldExport: false, key: "key", value: "val\nue", expected: "key=\"val\\nue\"\n"},
}

func (suite ParserSuite) TestParse() {
	for _, c := range cases {
		got := Parse(c.shouldExport, c.key, []byte(c.value))
		suite.Equal(c.expected, got)
	}
}

func (suite ParserSuite) TestParseStr() {
	for _, c := range cases {
		got := ParseStr(c.shouldExport, c.key, c.value)
		suite.Equal(c.expected, got)
	}
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
