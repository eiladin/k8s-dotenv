package environment

import (
	"bytes"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestResultOutput(t *testing.T) {
	type testCase struct {
		Name string

		Result *Result

		Opt *options.Options

		ExpectedString string
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Result.Output(tc.Opt)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ErrorChecker != nil {
				assert.Equal(t, true, tc.ErrorChecker(actualError))
			}
		})
	}

	objs := []runtime.Object{}
	objs = append(objs, mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}))
	objs = append(objs, mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}))
	client := fake.NewSimpleClientset(objs...)

	r1 := FromContainers([]v1.Container{mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test"})})
	r2 := FromContainers([]v1.Container{mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, nil)})
	r3 := FromContainers([]v1.Container{mock.Container(map[string]string{"env1": "val", "env2": "val2"}, nil, []string{"test"})})
	r4 := FromContainers([]v1.Container{mock.Container(map[string]string{"env1": "val", "env2": "val2"}, nil, nil)})
	r5 := FromContainers([]v1.Container{mock.Container(nil, []string{"test"}, nil)})
	r6 := FromContainers([]v1.Container{mock.Container(nil, nil, []string{"test"})})
	r7 := FromContainers([]v1.Container{mock.Container(nil, nil, []string{"test1"})})
	r8 := FromContainers([]v1.Container{mock.Container(nil, []string{"test1"}, nil)})

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	validate(t, &testCase{Name: "Should get env configmaps and secrets", Result: r1, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ExpectedString: envResult + secResult + cmResult})
	validate(t, &testCase{Name: "Should get env and configmaps with no secrets", Result: r2, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ExpectedString: envResult + cmResult})
	validate(t, &testCase{Name: "Should get env and secrets with no configmaps", Result: r3, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ExpectedString: envResult + secResult})
	validate(t, &testCase{Name: "Should get env with no secrets or configmaps", Result: r4, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ExpectedString: envResult})
	validate(t, &testCase{Name: "Should get configmaps with no env or secrets", Result: r5, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ExpectedString: cmResult})
	validate(t, &testCase{Name: "Should get secrets with no env or configmaps", Result: r6, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ExpectedString: secResult})
	validate(t, &testCase{Name: "Should error with missing secret", Result: r7, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ErrorChecker: errors.IsNotFound})
	validate(t, &testCase{Name: "Should error with missing configmap", Result: r8, Opt: &options.Options{Client: client, Namespace: "test", NoExport: true}, ErrorChecker: errors.IsNotFound})
}

func TestResultWrite(t *testing.T) {
	type testCase struct {
		Name string

		Result *Result

		Opt *options.Options

		Reader func() string

		ExpectedResult string

		ErrorChecker func(err error) bool
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
	client := fake.NewSimpleClientset(objs...)

	r1 := FromContainers([]v1.Container{mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test"})})
	r2 := FromContainers([]v1.Container{mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test1"})})

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	var b bytes.Buffer
	validate(t, &testCase{
		Name:   "Should work",
		Result: r1,
		Opt:    &options.Options{Client: client, Namespace: "test", NoExport: true, Writer: &b},
		Reader: func() string {
			return b.String()
		},
		ExpectedResult: envResult + secResult + cmResult,
	})

	var b2 bytes.Buffer
	validate(t, &testCase{
		Name:   "Should Error with missing secret",
		Result: r2,
		Opt:    &options.Options{Client: client, Namespace: "test", NoExport: true, Writer: &b2},
		Reader: func() string {
			return b2.String()
		},
		ErrorChecker: errors.IsNotFound,
	})

	defer os.Remove("./test.out")
	validate(t, &testCase{
		Name:   "Should Error with missing writer",
		Result: r1,
		Opt:    &options.Options{Client: client, Namespace: "test"},
		ErrorChecker: func(err error) bool {
			return err.Error() == options.ErrNoFilename.Error()
		},
	})
}
