package cmd

import (
	"flag"
	"os"
	"path/filepath"
)

func KubeconfigFlag() *string {
	if home := homeDir(); home != "" {
		return flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}

	return flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
