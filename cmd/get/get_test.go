package get

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GetSuite struct {
	suite.Suite
}

func (suite GetSuite) TestNewCmd() {
	got := NewCmd(nil)
	suite.NotNil(got)
}

func TestGetSuite(t *testing.T) {
	suite.Run(t, new(GetSuite))
}
