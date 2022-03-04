package mock

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewFakeResource returns an `APIResourceList` to be used with `FakeClient.WithResources`.
func NewFakeResource(groupVersion, name, singularName, kind, group string) *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: groupVersion,
		APIResources: []metav1.APIResource{
			{Name: name, SingularName: singularName, Kind: kind, Namespaced: true, Group: group},
		},
	}
}

// InvalidGroupResource returns an `APIResourceList` with an invalid group.
func InvalidGroupResource() *metav1.APIResourceList {
	return NewFakeResource("a/b/c", "CronJob", "CronJob", "CronJob", "batch/v1")
}

// Jobv1Resource returns a v1 Job resource list.
func Jobv1Resource() *metav1.APIResourceList {
	return NewFakeResource("v1", "Jobs", "Job", "Job", "v1")
}

// CronJobv1Resource returns a v1 CronJob resource list.
func CronJobv1Resource() *metav1.APIResourceList {
	return NewFakeResource("batch/v1", "CronJob", "CronJob", "CronJob", "batch/v1")
}

// CronJobv1beta1Resource returns a v1beta1 CronJob resource list.
func CronJobv1beta1Resource() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "batch/v1beta1",
		APIResources: []metav1.APIResource{
			{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1beta1"},
		},
	}
}

// UnsupportedGroupResource returns a resource list with a group version `batch/unsupported`.
func UnsupportedGroupResource() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "batch/unsupported",
		APIResources: []metav1.APIResource{
			{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/unsupported"},
		},
	}
}
