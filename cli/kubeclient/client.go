package kubeclient

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(kubeconfigPath string) (namespace string, client *kubernetes.Clientset, err error) {
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{})
	namespace, _, err = config.Namespace()
	if err != nil {
		return
	}

	var config0 *restclient.Config
	config0, err = config.ClientConfig()
	if err != nil {
		return
	}

	client, err = kubernetes.NewForConfig(config0)
	return
}
