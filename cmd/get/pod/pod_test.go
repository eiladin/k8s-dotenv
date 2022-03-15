package pod

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	kubeClient := mock.NewFakeClient(mock.Pod("test", "test", nil, nil, nil))

	got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
	assert.NotNil(t, got)

	objs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"test"}, objs)

	actualError := got.RunE(got, []string{})
	assert.Equal(t, ErrResourceNameRequired, actualError)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name           string
		Opt            *clioptions.CLIOptions
		Args           []string
		ExpectError    bool
		ExpectedResult string
		ResultChecker  func() string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := run(tc.Opt, tc.Args)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			if tc.ResultChecker != nil {
				assert.Equal(t, tc.ExpectedResult, tc.ResultChecker())
			}
		})
	}

	validate(t, &testCase{
		Name:        "Should error with no args",
		ExpectError: true,
	})

	kubeClient := mock.NewFakeClient(mock.Pod("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))
	writer := mock.NewWriter()

	validate(t, &testCase{
		Name: "Should find jobs",
		Opt: &clioptions.CLIOptions{
			KubeClient:   kubeClient,
			Namespace:    "test",
			ResourceName: "test",
			Writer:       writer,
		},
		Args:           []string{"test"},
		ExpectedResult: "export k=\"v\"\nexport k2=\"v2\"\n",
		ResultChecker:  writer.String,
	})

	validate(t, &testCase{
		Name: "Should return writer errors",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
			Writer:     mock.NewErrorWriter().ErrorAfter(1),
		},
		Args:        []string{"test"},
		ExpectError: true,
	})

	validate(t, &testCase{
		Name: "Should not find a pod in an empty cluster",
		Opt: &clioptions.CLIOptions{
			KubeClient: mock.NewFakeClient(),
			Namespace:  "test",
			Writer:     mock.NewWriter(),
		},
		Args:        []string{"test"},
		ExpectError: true,
	})
}

func TestValidArgs(t *testing.T) {
	type testCase struct {
		Name          string
		Opt           *clioptions.CLIOptions
		ExpectedSlice []string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice := validArgs(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
		})
	}

	kubeClient := mock.NewFakeClient(mock.Pod("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))

	validate(t, &testCase{
		Name: "Should return pods",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
		},
		ExpectedSlice: []string{"test"},
	})
}
