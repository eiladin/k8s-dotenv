package cronjob

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCmd(t *testing.T) {
	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	kubeClient := mock.NewFakeClient(v1mock).WithResources(mock.CronJobv1Resource())

	got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
	assert.NotNil(t, got)

	cronjobs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"my-cronjob"}, cronjobs)

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

	validate := func(t *testing.T, testCase *testCase) {
		t.Run(testCase.Name, func(t *testing.T) {
			actualError := run(testCase.Opt, testCase.Args)

			if testCase.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}

			if testCase.ResultChecker != nil {
				assert.Equal(t, testCase.ExpectedResult, testCase.ResultChecker())
			}
		})
	}

	v1mock := mock.CronJobv1("my-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)

	validate(t, &testCase{
		Name:        "Should error with no args",
		ExpectError: true,
	})

	kubeClient := mock.NewFakeClient().WithResources(mock.InvalidGroupResource())

	validate(t, &testCase{
		Name:        "Should return client errors",
		Opt:         &clioptions.CLIOptions{KubeClient: kubeClient},
		Args:        []string{"test"},
		ExpectError: true,
	})

	writer := mock.NewWriter()
	kubeClient = mock.NewFakeClient(v1mock).WithResources(mock.CronJobv1Resource())

	validate(t, &testCase{
		Name: "Should write v1 CronJobs",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
			Writer:     writer,
		},
		Args:           []string{"my-cronjob"},
		ExpectedResult: "export k1=\"v1\"\nexport k2=\"v2\"\n",
		ResultChecker:  writer.String,
	})

	kubeClient = mock.NewFakeClient(v1beta1mock).WithResources(mock.CronJobv1beta1Resource())
	writer = mock.NewWriter()

	validate(t, &testCase{
		Name: "Should write v1beta1 CronJobs",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
			Writer:     writer,
		},
		Args:           []string{"my-beta-cronjob"},
		ExpectedResult: "export k1=\"v1\"\nexport k2=\"v2\"\n",
		ResultChecker:  writer.String,
	})

	validate(t, &testCase{
		Name: "Should return writer errors",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
			Writer:     mock.NewErrorWriter().ErrorAfter(1),
		},
		Args:        []string{"my-beta-cronjob"},
		ExpectError: true,
	})

	kubeClient = mock.NewFakeClient().WithResources(mock.UnsupportedGroupResource())
	writer = mock.NewWriter()

	validate(t, &testCase{
		Name: "Should error on unsupported group",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
			Writer:     writer,
		},
		Args:        []string{"test"},
		ExpectError: true,
	})

	kubeClient = mock.NewFakeClient().
		WithResources(mock.CronJobv1Resource()).
		PrependReactor("get", "cronjobs", true, nil, assert.AnError)

	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &clioptions.CLIOptions{
			KubeClient: kubeClient,
			Namespace:  "test",
		},
		Args:        []string{"test"},
		ExpectError: true,
	})
}

func TestValidArgs(t *testing.T) {
	type testCase struct {
		Name          string
		Opt           *clioptions.CLIOptions
		APIResource   *metav1.APIResourceList
		Group         string
		ExpectedSlice []string
	}

	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", nil, nil, nil)
	kubeClient := mock.NewFakeClient(v1mock, v1beta1mock)

	validate := func(t *testing.T, tc *testCase) {
		kubeClient.Fake.Resources = []*metav1.APIResourceList{tc.APIResource}

		t.Run(tc.Name, func(t *testing.T) {
			actualSlice := validArgs(tc.Opt)
			assert.Equal(t, tc.ExpectedSlice, actualSlice)
		})
	}

	validate(t, &testCase{
		Name:          "Should find v1 cronjobs",
		Group:         "batch/v1",
		Opt:           &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
		APIResource:   mock.CronJobv1Resource(),
		ExpectedSlice: []string{"my-cronjob"},
	})

	validate(t, &testCase{
		Name:          "Should find v1beta1 cronjobs",
		Group:         "batch/v1beta1",
		Opt:           &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
		APIResource:   mock.CronJobv1beta1Resource(),
		ExpectedSlice: []string{"my-beta-cronjob"},
	})

	validate(t, &testCase{
		Name:        "Should not find non-existent groups",
		Group:       "batch/not-a-version",
		APIResource: mock.InvalidGroupResource(),
		Opt:         &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
	})
}
