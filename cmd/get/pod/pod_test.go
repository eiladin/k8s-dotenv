package pod

import (
	"bytes"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/errors/cmd"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewCmd(t *testing.T) {
	client := fake.NewSimpleClientset(mocks.Pod("test", "test", nil, nil, nil))

	got := NewCmd(&options.Options{Client: client, Namespace: "test"})
	assert.NotNil(t, got)

	objs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"test"}, objs)

	actualError := got.RunE(got, []string{})
	assert.Equal(t, cmd.ErrResourceNameRequired, actualError)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name string

		Opt  *options.Options
		Args []string

		ExpectedError error
		ErrorChecker  func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := run(tc.Opt, tc.Args)

			if tc.ErrorChecker != nil {
				assert.True(t, tc.ErrorChecker(actualError))
			} else {
				assert.Equal(t, tc.ExpectedError, actualError)
			}
		})
	}

	validate(t, &testCase{
		Name:          "Should error with no args",
		ExpectedError: cmd.ErrResourceNameRequired,
	})

	var b bytes.Buffer
	client := fake.NewSimpleClientset(mocks.Pod("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))
	validate(t, &testCase{
		Name: "Should find jobs",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
			Writer:    &b,
		},
		Args: []string{"test"},
	})

	b.Reset()
	client = fake.NewSimpleClientset()
	validate(t, &testCase{
		Name: "Should not find a pod in an empty cluster",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
			Writer:    &b,
		},
		Args:         []string{"test"},
		ErrorChecker: errors.IsNotFound,
	})

}

func TestValidArgs(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		ExpectedSlice []string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice := validArgs(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
		})
	}

	validate(t, &testCase{
		Name: "Should return pods",
		Opt: &options.Options{
			Client:    fake.NewSimpleClientset(),
			Namespace: "test",
			Name:      "test",
		},
		ExpectedSlice: []string{},
	})
}
