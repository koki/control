package cmd

import (
	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"

	shortutil "github.com/koki/short/util"
)

func CreateKubeObj(kubeObj interface{}, env *Env) error {
	var err error
	switch kubeObj := kubeObj.(type) {
	case *v1.Pod:
		_, err = env.Client.CoreV1().Pods(env.Namespace).Create(kubeObj)
	case *v1.ReplicationController:
		_, err = env.Client.CoreV1().ReplicationControllers(env.Namespace).Create(kubeObj)
	case *exts.ReplicaSet:
		_, err = env.Client.ExtensionsV1beta1().ReplicaSets(env.Namespace).Create(kubeObj)
	case *exts.Deployment:
		_, err = env.Client.ExtensionsV1beta1().Deployments(env.Namespace).Create(kubeObj)
	case *v1.PersistentVolume:
		_, err = env.Client.CoreV1().PersistentVolumes().Create(kubeObj)
	case *v1.Service:
		_, err = env.Client.CoreV1().Services(env.Namespace).Create(kubeObj)
	default:
		err = shortutil.TypeErrorf(kubeObj, "unsupported k8s type")
	}

	return err
}
