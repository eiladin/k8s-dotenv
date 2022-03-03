package resources

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InvalidGroup() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "a/b/c",
		APIResources: []metav1.APIResource{
			{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1"},
		},
	}
}

func Jobv1() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "v1",
		APIResources: []metav1.APIResource{
			{Name: "Jobs", SingularName: "Job", Kind: "Job", Namespaced: false, Group: "v1"},
		},
	}
}

func CronJobv1() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "batch/v1",
		APIResources: []metav1.APIResource{
			{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1"},
		},
	}
}

func CronJobv1beta1() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "batch/v1beta1",
		APIResources: []metav1.APIResource{
			{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/v1beta1"},
		},
	}
}

func UnsupportedGroup() *metav1.APIResourceList {
	return &metav1.APIResourceList{
		GroupVersion: "batch/unsupported",
		APIResources: []metav1.APIResource{
			{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: "batch/unsupported"},
		},
	}
}
