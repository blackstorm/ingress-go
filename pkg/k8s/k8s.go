package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	var conf *rest.Config
	var err error

	if len(kubeconfigPath) > 0 {
		conf, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		conf, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(conf)
}
