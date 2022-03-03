package environment

import (
	"bytes"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/configmap"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/secret"
	tests "github.com/eiladin/k8s-dotenv/pkg/testing"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestResultOutput(t *testing.T) {
	type testCase struct {
		Name string

		Result *Result

		Client       *client.Client
		Namespace    string
		ShouldExport bool

		ExpectedString string
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Result.output(tc.Client, tc.Namespace, tc.ShouldExport)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ErrorChecker != nil {
				assert.Equal(t, true, tc.ErrorChecker(actualError))
			}
		})
	}

	objs := []runtime.Object{}
	objs = append(objs, mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}))
	objs = append(objs, mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}))
	cl := fake.NewSimpleClientset(objs...)
	envMap := map[string]string{"env1": "val", "env2": "val2"}

	r1 := FromContainers([]v1.Container{mock.Container(envMap, []string{"test"}, []string{"test"})})
	r2 := FromContainers([]v1.Container{mock.Container(envMap, []string{"test"}, nil)})
	r3 := FromContainers([]v1.Container{mock.Container(envMap, nil, []string{"test"})})
	r4 := FromContainers([]v1.Container{mock.Container(envMap, nil, nil)})
	r5 := FromContainers([]v1.Container{mock.Container(nil, []string{"test"}, nil)})
	r6 := FromContainers([]v1.Container{mock.Container(nil, nil, []string{"test"})})
	r7 := FromContainers([]v1.Container{mock.Container(nil, nil, []string{"test1"})})
	r8 := FromContainers([]v1.Container{mock.Container(nil, []string{"test1"}, nil)})

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	validate(t, &testCase{
		Name:           "Should get env configmaps and secrets",
		Result:         r1,
		Client:         client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: envResult + secResult + cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and configmaps with no secrets",
		Result:         r2,
		Client:         client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: envResult + cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and secrets with no configmaps",
		Result:         r3,
		Client:         client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: envResult + secResult,
	})

	validate(t, &testCase{
		Name:           "Should get env with no secrets or configmaps",
		Result:         r4,
		Client:         client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: envResult,
	})

	validate(t, &testCase{
		Name:           "Should get configmaps with no env or secrets",
		Result:         r5,
		Client:         client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get secrets with no env or configmaps",
		Result:         r6,
		Client:         client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: secResult,
	})

	validate(t, &testCase{
		Name:      "Should error with missing secret",
		Result:    r7,
		Client:    client.NewClient(cl),
		Namespace: "test",
		ErrorChecker: func(err error) bool {
			return assert.ErrorIs(t, err, secret.ErrMissingResource)
		},
	})

	validate(t, &testCase{
		Name:      "Should error with missing configmap",
		Result:    r8,
		Client:    client.NewClient(cl),
		Namespace: "test",
		ErrorChecker: func(err error) bool {
			return assert.ErrorIs(t, err, configmap.ErrMissingResource)
		},
	})
}

func TestResultWrite(t *testing.T) {
	type testCase struct {
		Name string

		Result *Result
		Opt    *options.Options
		Reader func() string

		ExpectedResult string
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.Result.Write(tc.Opt)

			if tc.Reader != nil {
				assert.Equal(t, tc.ExpectedResult, tc.Reader())
			}

			if tc.ErrorChecker != nil {
				assert.True(t, tc.ErrorChecker(actualError))
			}
		})
	}

	objs := []runtime.Object{}
	objs = append(objs, mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}))
	objs = append(objs, mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}))
	cl := fake.NewSimpleClientset(objs...)
	envMap := map[string]string{"env1": "val", "env2": "val2"}

	r1 := FromContainers([]v1.Container{mock.Container(envMap, []string{"test"}, []string{"test"})})
	r2 := FromContainers([]v1.Container{mock.Container(envMap, []string{"test"}, []string{"test1"})})

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	var b bytes.Buffer

	validate(t, &testCase{
		Name:           "Should work",
		Result:         r1,
		Opt:            &options.Options{Client: client.NewClient(cl), Namespace: "test", NoExport: true, Writer: &b},
		Reader:         b.String,
		ExpectedResult: envResult + secResult + cmResult,
	})

	var b2 bytes.Buffer

	validate(t, &testCase{
		Name:   "Should Error with missing secret",
		Result: r2,
		Opt:    &options.Options{Client: client.NewClient(cl), Namespace: "test", NoExport: true, Writer: &b2},
		Reader: b2.String,
		ErrorChecker: func(err error) bool {
			return assert.ErrorIs(t, err, secret.ErrMissingResource)
		},
	})

	defer os.Remove("./test.out")
	validate(t, &testCase{
		Name:   "Should Error with missing writer",
		Result: r1,
		Opt:    &options.Options{Client: client.NewClient(cl), Namespace: "test"},
		ErrorChecker: func(err error) bool {
			return assert.ErrorIs(t, err, err, options.ErrNoFilename)
		},
	})

	validate(t, &testCase{
		Name:   "Should return writer errors",
		Result: r1,
		Opt: &options.Options{
			Client:    client.NewClient(cl),
			Namespace: "test",
			Writer:    tests.NewErrorWriter(&b2).ErrorAfter(1),
		},
		ErrorChecker: func(err error) bool {
			return assert.Equal(t, NewWriteError(mock.NewError("error")), err)
		},
	})
}
