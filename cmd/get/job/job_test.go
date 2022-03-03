package job

import (
	"bytes"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	tests "github.com/eiladin/k8s-dotenv/pkg/testing"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewCmd(t *testing.T) {
	cl := fake.NewSimpleClientset(mock.Job("test", "test", nil, nil, nil))

	got := NewCmd(&options.Options{Client: client.NewClient(cl), Namespace: "test"})
	assert.NotNil(t, got)

	objs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"test"}, objs)

	actualError := got.RunE(got, []string{})
	assert.Equal(t, ErrResourceNameRequired, actualError)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name string

		Opt  *options.Options
		Args []string

		ExpectError   bool
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := run(tc.Opt, tc.Args)

			checkErrNilFn := assert.Nil
			if tc.ExpectError || tc.ExpectedError != nil {
				checkErrNilFn = assert.NotNil
			}

			checkErrNilFn(t, actualError)

			if tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError, actualError)
			}
		})
	}

	validate(t, &testCase{
		Name:          "Should error with no args",
		ExpectedError: ErrResourceNameRequired,
	})

	var b bytes.Buffer

	cl := fake.NewSimpleClientset(mock.Job("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))

	validate(t, &testCase{
		Name: "Should find jobs",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    &b,
		},
		Args: []string{"test"},
	})

	validate(t, &testCase{
		Name: "Should return writer errors",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    tests.NewErrorWriter(&b).ErrorAfter(1),
		},
		Args:          []string{"test"},
		ExpectedError: newRunError(environment.NewWriteError(mock.NewError("error"))),
	})

	b.Reset()

	cl = fake.NewSimpleClientset()

	validate(t, &testCase{
		Name: "Should not find a job in an empty cluster",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    &b,
		},
		Args:        []string{"test"},
		ExpectError: true,
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

	cl := fake.NewSimpleClientset()

	validate(t, &testCase{
		Name: "Should return jobs",
		Opt: &options.Options{
			Client:       client.NewClient(cl),
			Namespace:    "test",
			ResourceName: "test",
		},
		ExpectedSlice: []string{},
	})
}
