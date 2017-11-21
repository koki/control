package koki

import (
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	ControllerName         = "koki-controller"
	SidecarName            = "koki-sidecar"
	ServiceAccountName     = "koki"
	DefaultSidecarPort     = 30547
	DefaultSidecarImage    = "koki/sidecar:latest"
	DefaultControllerImage = "koki/controller:latest"
)

func ConfigMapSelector() labels.Selector {
	sel := labels.NewSelector()
	req, err := labels.NewRequirement("koki", selection.In, []string{"application"})
	if err != nil {
		glog.Fatal(err)
	}

	sel.Add(*req)

	return sel
}

func ConfigMapListOptions() metav1.ListOptions {
	return metav1.ListOptions{LabelSelector: ConfigMapSelector().String()}
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}
