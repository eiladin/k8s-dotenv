package deployment

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	cl := mock.NewFakeClient(mock.Deployment("test", "test", nil, nil, nil))

	got := NewCmd(&options.Options{Client: cl, Namespace: "test"})
	assert.NotNil(t, got)

	objs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"test"}, objs)

	actualError := got.RunE(got, []string{})
	assert.Equal(t, ErrResourceNameRequired, actualError)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name        string
		Opt         *options.Options
		Args        []string
		ExpectError bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := run(tc.Opt, tc.Args)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	validate(t, &testCase{
		Name:        "Should error with no args",
		ExpectError: true,
	})

	cl := mock.NewFakeClient(mock.Deployment("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))

	validate(t, &testCase{
		Name: "Should find deployments",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
			Writer:    mock.NewWriter(),
		},
		Args: []string{"test"},
	})

	validate(t, &testCase{
		Name: "Should return writer errors",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
			Writer:    mock.NewErrorWriter().ErrorAfter(1),
		},
		Args:        []string{"test"},
		ExpectError: true,
	})

	validate(t, &testCase{
		Name: "Should not find a deployment in an empty cluster",
		Opt: &options.Options{
			Client:    mock.NewFakeClient(),
			Namespace: "test",
			Writer:    mock.NewWriter(),
		},
		Args:        []string{"test"},
		ExpectError: true,
	})
}

func TestValidArgs(t *testing.T) {
	type testCase struct {
		Name          string
		Opt           *options.Options
		ExpectedSlice []string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice := validArgs(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
		})
	}

	cl := mock.NewFakeClient(mock.Deployment("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))

	validate(t, &testCase{
		Name: "Should return deployments",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
		},
		ExpectedSlice: []string{"test"},
	})
}
