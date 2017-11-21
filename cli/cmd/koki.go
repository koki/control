package cmd

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/koki/control/cli/resources"
	"github.com/koki/control/pkg/koki"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeleteAppForPodIfExists(namespace string, client *kubernetes.Clientset, podName string) error {
	err := DeleteConfigMapForPodIfExists(namespace, client, podName)
	if err != nil {
		return err
	}

	DeletePodIfExists(namespace, client, podName)
	if err != nil {
		return err
	}

	return nil
}

func DeleteConfigMapForPodIfExists(namespace string, client *kubernetes.Clientset, podName string) error {
	err := client.CoreV1().ConfigMaps(namespace).Delete(podName, nil)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		glog.Error("Failed deleting configmap")
		return err
	}

	return nil
}

func DeletePodIfExists(namespace string, client *kubernetes.Clientset, podName string) error {
	err := client.CoreV1().Pods(namespace).Delete(podName, nil)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		glog.Error("Failed deleting pod")
		return err
	}

	return nil
}

func DeleteControllerIfExists(namespace string, client *kubernetes.Clientset) error {
	propagation := metav1.DeletePropagationForeground
	err := client.AppsV1beta1().Deployments(namespace).Delete(koki.ControllerName, &metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		glog.Error("Failed deleting controller")
		return err
	}

	return nil
}

func CreateAppForPod(e *Env, pod *v1.Pod) error {
	d, err := CreateControllerIfNeeded(e)
	if err != nil {
		glog.Errorf("Failed creating controller:\n%s", stringYaml(d))
		return err
	}

	cm, err := e.Client.CoreV1().ConfigMaps(e.Namespace).Create(
		resources.ConfigMap(pod.Name))
	if err != nil {
		glog.Errorf("Failed creating configmap:\n%s", stringYaml(cm))
		return err
	}

	p := resources.InsertSidecar(pod, e.SidecarImage, e.SidecarPort)
	_, err = e.Client.CoreV1().Pods(e.Namespace).Create(p)
	if err != nil {
		glog.Errorf("Failed creating pod:\n%s", stringYaml(p))

		// Clean up the configmap we created.
		DeleteConfigMapForPodIfExists(e.Namespace, e.Client, pod.Name)
		return err
	}

	return nil
}

func CreateControllerIfNeeded(e *Env) (*appsv1beta1.Deployment, error) {
	d, err := e.Client.AppsV1beta1().Deployments(e.Namespace).Get(
		koki.ControllerName, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			glog.Error("Failed checking if controller already exists")
			return nil, err
		}
	} else {
		return d, nil
	}

	d = resources.ControllerDeployment(e.Namespace, e.ControllerImage, e.SidecarPort)
	dd, err := e.Client.AppsV1beta1().Deployments(e.Namespace).Create(d)
	if err != nil {
		return d, err
	}

	return dd, nil
}

func PurgeAppsAndController(namespace string, client *kubernetes.Clientset) error {
	err := DeleteControllerIfExists(namespace, client)
	if err != nil {
		return err
	}

	// Delete all the koki configmaps and their pods.
	cmaps, err := client.CoreV1().ConfigMaps(namespace).List(koki.ConfigMapListOptions())
	for _, cm := range cmaps.Items {
		if podName, ok := cm.Data["pod.name"]; ok {
			err := client.CoreV1().ConfigMaps(namespace).Delete(cm.Name, nil)
			if err != nil {
				glog.Error("Failed to delete configmap that was just fetched.")

				return err
			}

			err = DeletePodIfExists(namespace, client, podName)
			if err != nil {
				if errors.IsNotFound(err) {
					continue
				}

				return err
			}
		} else {
			glog.Error("Found configmap without pod.name")
		}
	}

	return nil
}

func printYaml(obj interface{}) {
	b, err := yaml.Marshal(obj)

	if err != nil {
		glog.Error(err)
		return
	}

	fmt.Println(string(b))
}

func stringYaml(obj interface{}) string {
	b, err := yaml.Marshal(obj)

	if err != nil {
		glog.Error(err)
		return "\"yaml error\""
	}

	return string(b)
}
