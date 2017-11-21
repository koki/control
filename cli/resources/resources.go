package resources

import (
	"fmt"

	"github.com/koki/control/pkg/koki"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SidecarContainer(sidecarImage string, sidecarPort int32) *v1.Container {
	return &v1.Container{
		Image:           sidecarImage,
		ImagePullPolicy: "Always",
		Name:            koki.SidecarName,
		Ports: []v1.ContainerPort{
			v1.ContainerPort{ContainerPort: sidecarPort},
		},
		Args: []string{
			"-port",
			fmt.Sprintf("%d", sidecarPort),
		},
	}
}

func InsertSidecar(pod *v1.Pod, sidecarImage string, sidecarPort int32) *v1.Pod {
	pod.Spec.Containers = append(
		pod.Spec.Containers,
		*SidecarContainer(sidecarImage, sidecarPort))

	return pod
}

func ConfigMap(podName string) *v1.ConfigMap {
	cm := &v1.ConfigMap{
		Data: map[string]string{
			"pod.name": podName,
		},
	}
	cm.Name = podName
	cm.Labels = map[string]string{
		"koki": "application",
	}
	return cm
}

func controllerContainer(namespace string, controllerImage string, sidecarPort int32) *v1.Container {
	return &v1.Container{
		Name:            koki.SidecarName,
		Image:           controllerImage,
		ImagePullPolicy: v1.PullAlways,
		Args: []string{
			"-sidecarPort",
			fmt.Sprintf("%d", sidecarPort),
			"-namespace",
			namespace,
		},
	}
}

func ControllerDeployment(namespace string, controllerImage string, sidecarPort int32) *appsv1beta1.Deployment {
	return &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: koki.ControllerName,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: koki.Int32Ptr(1),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": koki.ControllerName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						*controllerContainer(
							namespace,
							controllerImage,
							sidecarPort,
						),
					},
					ServiceAccountName: koki.ServiceAccountName,
				},
			},
		},
	}
}
