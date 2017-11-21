package kubeclient

import (
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(kubeconfigPath string) (namespace string, config *rest.Config, client *kubernetes.Clientset, err error) {
	config0 := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{})
	namespace, _, err = config0.Namespace()
	if err != nil {
		return
	}

	config, err = config0.ClientConfig()
	if err != nil {
		return
	}

	client, err = kubernetes.NewForConfig(config)
	return
}
