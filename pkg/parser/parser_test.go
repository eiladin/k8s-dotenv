package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type testCase struct {
		Name string

		ShouldExport bool
		Key          string
		Value        []byte

		ExpectedString string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := Parse(tc.ShouldExport, tc.Key, tc.Value)

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	validate(t, &testCase{Name: "Should not export", ShouldExport: false, Key: "key", Value: []byte("value"), ExpectedString: "key=\"value\"\n"})
	validate(t, &testCase{Name: "Should export", ShouldExport: true, Key: "key", Value: []byte("value"), ExpectedString: "export key=\"value\"\n"})
	validate(t, &testCase{Name: "Should remove dots from keys", ShouldExport: true, Key: "k.e.y", Value: []byte("value"), ExpectedString: "export key=\"value\"\n"})
	validate(t, &testCase{Name: "Should escape newlines in values", ShouldExport: true, Key: "key", Value: []byte("va\nlue"), ExpectedString: "export key=\"va\\nlue\"\n"})
}

func TestParseStr(t *testing.T) {
	type testCase struct {
		Name string

		ShouldExport bool
		Key          string
		Value        string

		ExpectedString string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := ParseStr(tc.ShouldExport, tc.Key, tc.Value)

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	validate(t, &testCase{Name: "Should not export", ShouldExport: false, Key: "key", Value: "value", ExpectedString: "key=\"value\"\n"})
	validate(t, &testCase{Name: "Should export", ShouldExport: true, Key: "key", Value: "value", ExpectedString: "export key=\"value\"\n"})
	validate(t, &testCase{Name: "Should remove dots from keys", ShouldExport: true, Key: "k.e.y", Value: "value", ExpectedString: "export key=\"value\"\n"})
}
