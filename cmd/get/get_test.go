package get

import (
	"testing"
)

func TestNewCmd(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		got := NewCmd(nil)
		if got == nil {
			t.Errorf("NewCmd() = nil, want not nil")
		}
	})
}
