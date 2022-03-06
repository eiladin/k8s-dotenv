package cronjob

import (
	"bytes"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCmd(t *testing.T) {
	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	cl := mock.NewFakeClient(v1mock).WithResources(mock.CronJobv1Resource())

	got := NewCmd(&options.Options{Client: cl, Namespace: "test"})
	assert.NotNil(t, got)

	cronjobs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"my-cronjob"}, cronjobs)

	actualError := got.RunE(got, []string{})
	assert.Equal(t, ErrResourceNameRequired, actualError)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name        string
		Opt         *options.Options
		Args        []string
		ExpectError bool
		Comparison  assert.Comparison
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := run(tc.Opt, tc.Args)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			if tc.Comparison != nil {
				assert.Condition(t, tc.Comparison)
			}
		})
	}

	v1mock := mock.CronJobv1("my-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)

	validate(t, &testCase{
		Name:        "Should error with no args",
		ExpectError: true,
	})

	cl := mock.NewFakeClient().WithResources(mock.InvalidGroupResource())

	validate(t, &testCase{
		Name:        "Should return client errors",
		Opt:         &options.Options{Client: cl},
		Args:        []string{"test"},
		ExpectError: true,
	})

	var b bytes.Buffer

	cl = mock.NewFakeClient(v1mock).WithResources(mock.CronJobv1Resource())

	validate(t, &testCase{
		Name: "Should write v1 CronJobs",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
			Writer:    &b,
		},
		Args: []string{"my-cronjob"},
		Comparison: func() (success bool) {
			return assert.Equal(t, "export k1=\"v1\"\nexport k2=\"v2\"\n", b.String())
		},
	})

	cl = mock.NewFakeClient(v1beta1mock).WithResources(mock.CronJobv1beta1Resource())

	b.Reset()

	validate(t, &testCase{
		Name: "Should write v1beta1 CronJobs",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
			Writer:    &b,
		},
		Args: []string{"my-beta-cronjob"},
		Comparison: func() (success bool) {
			return assert.Equal(t, "export k1=\"v1\"\nexport k2=\"v2\"\n", b.String())
		},
	})

	validate(t, &testCase{
		Name: "Should return writer errors",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
			Writer:    mock.NewErrorWriter().ErrorAfter(1),
		},
		Args:        []string{"my-beta-cronjob"},
		ExpectError: true,
	})

	cl = mock.NewFakeClient().WithResources(mock.UnsupportedGroupResource())

	b.Reset()

	validate(t, &testCase{
		Name: "Should error on unsupported group",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
			Writer:    &b,
		},
		Args:        []string{"test"},
		ExpectError: true,
	})

	cl = mock.NewFakeClient().
		WithResources(mock.CronJobv1Resource()).
		PrependReactor("get", "cronjobs", true, nil, assert.AnError)

	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &options.Options{
			Client:    cl,
			Namespace: "test",
		},
		Args:        []string{"test"},
		ExpectError: true,
	})
}

func TestValidArgs(t *testing.T) {
	type testCase struct {
		Name          string
		Opt           *options.Options
		APIResource   *metav1.APIResourceList
		Group         string
		ExpectedSlice []string
	}

	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", nil, nil, nil)
	cl := mock.NewFakeClient(v1mock, v1beta1mock)

	validate := func(t *testing.T, tc *testCase) {
		cl.Fake.Resources = []*metav1.APIResourceList{tc.APIResource}

		t.Run(tc.Name, func(t *testing.T) {
			actualSlice := validArgs(tc.Opt)
			assert.Equal(t, tc.ExpectedSlice, actualSlice)
		})
	}

	validate(t, &testCase{
		Name:          "Should find v1 cronjobs",
		Group:         "batch/v1",
		Opt:           &options.Options{Client: cl, Namespace: "test"},
		APIResource:   mock.CronJobv1Resource(),
		ExpectedSlice: []string{"my-cronjob"},
	})

	validate(t, &testCase{
		Name:          "Should find v1beta1 cronjobs",
		Group:         "batch/v1beta1",
		Opt:           &options.Options{Client: cl, Namespace: "test"},
		APIResource:   mock.CronJobv1beta1Resource(),
		ExpectedSlice: []string{"my-beta-cronjob"},
	})

	validate(t, &testCase{
		Name:        "Should not find non-existent groups",
		Group:       "batch/not-a-version",
		APIResource: mock.InvalidGroupResource(),
		Opt:         &options.Options{Client: cl, Namespace: "test"},
	})
}
