package cmd

import (
	"flag"
	"os"

	"github.com/koki/control/cli/kubeclient"
	"github.com/koki/control/pkg/koki"
	"k8s.io/client-go/kubernetes"
)

type Env struct {
	Namespace       string
	SidecarPort     int32
	ControllerImage string
	SidecarImage    string
	Client          *kubernetes.Clientset
}

func SidecarPortFlag() *int {
	return flag.Int("sidecar-port", koki.DefaultSidecarPort, "(optional) port for the koki sidecar")
}

func ControllerImageFlag() *string {
	return flag.String("controller-image", "", "(optional) docker image for the koki controller")
}

func SidecarImageFlag() *string {
	return flag.String("sidecar-image", "", "(optional) docker image for the koki sidecar")
}

func EnvFromFlags() (*Env, error) {
	kubeconfig := KubeconfigFlag()
	sidecarPort := SidecarPortFlag()
	controllerImage := ControllerImageFlag()
	sidecarImage := SidecarImageFlag()
	flag.Parse()

	namespace, client, err := kubeclient.GetClient(*kubeconfig)
	if err != nil {
		return nil, err
	}

	env := Env{
		Namespace:       namespace,
		SidecarPort:     int32(*sidecarPort),
		ControllerImage: *controllerImage,
		SidecarImage:    *sidecarImage,
		Client:          client,
	}

	if len(env.ControllerImage) == 0 {
		env.ControllerImage = os.Getenv("KOKI_CONTROLLER_IMAGE")
	}

	if len(env.ControllerImage) == 0 {
		env.ControllerImage = koki.DefaultControllerImage
	}

	if len(env.SidecarImage) == 0 {
		env.SidecarImage = os.Getenv("KOKI_SIDECAR_IMAGE")
	}

	if len(env.SidecarImage) == 0 {
		env.SidecarImage = koki.DefaultSidecarImage
	}

	return &env, nil
}
