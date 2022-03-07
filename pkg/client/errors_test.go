package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceLoadErrorError(t *testing.T) {
	type testCase struct {
		Name              string
		ResourceLoadError *ResourceLoadError
		ExpectedString    string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := tc.ResourceLoadError.Error()

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	validate(t, &testCase{
		Name: "Should return internal error",
		ResourceLoadError: &ResourceLoadError{
			Err:      assert.AnError,
			Resource: "test",
		},
		ExpectedString: "error loading test: assert.AnError general error for testing",
	})

	validate(t, &testCase{
		Name:              "Should return message when there is no internal error",
		ResourceLoadError: &ResourceLoadError{Resource: "test"},
		ExpectedString:    "error loading test",
	})
}

func TestResourceLoadErrorUnwrap(t *testing.T) {
	type testCase struct {
		Name              string
		ResourceLoadError *ResourceLoadError
		ExpectedError     error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.ResourceLoadError.Unwrap()

			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name: "Should return internal error",
		ResourceLoadError: &ResourceLoadError{
			Err:      assert.AnError,
			Resource: "test",
		},
		ExpectedError: assert.AnError,
	})
}

func TestNewResourceLoadError(t *testing.T) {
	type testCase struct {
		Name          string
		Resource      string
		Err           error
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := NewResourceLoadError(tc.Resource, tc.Err)

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name:     "Should wrap errors",
		Resource: "test",
		Err:      assert.AnError,
		ExpectedError: &ResourceLoadError{
			Err:      assert.AnError,
			Resource: "test",
		},
	})
}

func TestErrorWrappers(t *testing.T) {
	type testCase struct {
		Name          string
		Func          func(err error) error
		Err           error
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.Func(tc.Err)

			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name:          "NewSecretErr should wrap errors",
		Err:           assert.AnError,
		Func:          newSecretErr,
		ExpectedError: assert.AnError,
	})

	validate(t, &testCase{
		Name:          "NewConfigMapError should wrap errors",
		Err:           assert.AnError,
		Func:          newConfigMapError,
		ExpectedError: assert.AnError,
	})

	validate(t, &testCase{
		Name:          "NewWriteErr Should wrap errors",
		Err:           assert.AnError,
		Func:          newWriteError,
		ExpectedError: assert.AnError,
	})
}
