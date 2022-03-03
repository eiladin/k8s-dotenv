package cronjob

import (
	"bytes"
	"testing"

	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	tests "github.com/eiladin/k8s-dotenv/pkg/testing"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewCmd(t *testing.T) {
	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	cl := fake.NewSimpleClientset(v1mock)
	cl.Fake.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "batch/v1",
			APIResources: []metav1.APIResource{
				{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1"},
			},
		},
	}

	got := NewCmd(&options.Options{Client: client.NewClient(cl), Namespace: "test"})
	assert.NotNil(t, got)

	cronjobs, _ := got.ValidArgsFunction(got, []string{}, "")
	assert.Equal(t, []string{"my-cronjob"}, cronjobs)

	actualError := got.RunE(got, []string{})
	assert.Equal(t, ErrResourceNameRequired, actualError)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name string

		Opt              *options.Options
		Args             []string
		ResultChecker    func() string
		ExpectedResult   string
		ExpectedErrorStr string
		ExpectedError    error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := run(tc.Opt, tc.Args)

			if tc.ExpectedErrorStr != "" {
				assert.Equal(t, tc.ExpectedErrorStr, actualError.Error())
			}

			if tc.ExpectedError != nil {
				assert.EqualError(t, actualError, tc.ExpectedError.Error())
			}

			if tc.ResultChecker != nil {
				assert.Equal(t, tc.ExpectedResult, tc.ResultChecker())
			}
		})
	}

	v1mock := mock.CronJobv1("my-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)
	fakeResources := map[string]*metav1.APIResourceList{
		"invalid": {
			GroupVersion: "a/b/c",
			APIResources: []metav1.APIResource{
				{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1"},
			},
		},
		"unsupported": {
			GroupVersion: "batch/unsupported",
			APIResources: []metav1.APIResource{
				{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/unsupported"},
			},
		},
		"v1": {
			GroupVersion: "batch/v1",
			APIResources: []metav1.APIResource{
				{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1"},
			},
		},
		"v1beta1": {
			GroupVersion: "batch/v1beta1",
			APIResources: []metav1.APIResource{
				{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1beta1"},
			},
		},
	}

	validate(t, &testCase{
		Name:          "Should error with no args",
		ExpectedError: ErrResourceNameRequired,
	})

	cl := fake.NewSimpleClientset()
	cl.Fake.Resources = []*metav1.APIResourceList{fakeResources["invalid"]}

	validate(t, &testCase{
		Name:          "Should return client errors",
		Opt:           &options.Options{Client: client.NewClient(cl)},
		Args:          []string{"test"},
		ExpectedError: newClientError(client.ErrAPIGroup),
	})

	var b bytes.Buffer

	cl = fake.NewSimpleClientset(v1mock)
	cl.Fake.Resources = []*metav1.APIResourceList{fakeResources["v1"]}

	validate(t, &testCase{
		Name: "Should write v1 CronJobs",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    &b,
		},
		Args:           []string{"my-cronjob"},
		ExpectedResult: "export k1=\"v1\"\nexport k2=\"v2\"\n",
		ResultChecker:  b.String,
	})

	cl = fake.NewSimpleClientset(v1beta1mock)
	cl.Fake.Resources = []*metav1.APIResourceList{fakeResources["v1beta1"]}

	b.Reset()

	validate(t, &testCase{
		Name: "Should write v1beta1 CronJobs",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    &b,
		},
		Args:           []string{"my-beta-cronjob"},
		ExpectedResult: "export k1=\"v1\"\nexport k2=\"v2\"\n",
		ResultChecker:  b.String,
	})

	validate(t, &testCase{
		Name: "Should return writer errors",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    tests.NewErrorWriter(&b).ErrorAfter(1),
		},
		Args:          []string{"my-beta-cronjob"},
		ExpectedError: newRunError(environment.NewWriteError(mock.NewError("error"))),
	})

	cl = fake.NewSimpleClientset()
	cl.Fake.Resources = []*metav1.APIResourceList{fakeResources["unsupported"]}

	b.Reset()

	validate(t, &testCase{
		Name: "Should error on unsupported group",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    &b,
		},
		Args:          []string{"test"},
		ExpectedError: ErrUnsupportedGroup,
	})

	cl = fake.NewSimpleClientset()
	cl.Fake.Resources = []*metav1.APIResourceList{fakeResources["v1"]}
	cl.PrependReactor("get", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, assert.AnError
	})

	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
		},
		Args:          []string{"test"},
		ExpectedError: newRunError(v1.NewResourceLoadError(assert.AnError)),
	})
}

func TestValidArgs(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		Group string

		ExpectedSlice []string
	}

	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", nil, nil, nil)
	cl := fake.NewSimpleClientset(v1mock, v1beta1mock)

	validate := func(t *testing.T, tc *testCase) {
		cl.Fake.Resources = []*metav1.APIResourceList{
			{
				GroupVersion: tc.Group,
				APIResources: []metav1.APIResource{
					{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: tc.Group},
				},
			},
		}

		t.Run(tc.Name, func(t *testing.T) {
			actualSlice := validArgs(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
		})
	}

	validate(t, &testCase{
		Name:          "Should find v1 cronjobs",
		Group:         "batch/v1",
		Opt:           &options.Options{Client: client.NewClient(cl), Namespace: "test"},
		ExpectedSlice: []string{"my-cronjob"},
	})

	validate(t, &testCase{
		Name:          "Should find v1beta1 cronjobs",
		Group:         "batch/v1beta1",
		Opt:           &options.Options{Client: client.NewClient(cl), Namespace: "test"},
		ExpectedSlice: []string{"my-beta-cronjob"},
	})

	validate(t, &testCase{
		Name:  "Should not find non-existent groups",
		Group: "batch/not-a-version",
		Opt: &options.Options{Client: client.NewClient(cl),
			Namespace: "test"},
	})
}
