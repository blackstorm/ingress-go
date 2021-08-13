package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClient(path *string) (*kubernetes.Clientset, error) {
	var conf *rest.Config
	var err error

	if path == nil {
		conf, err = rest.InClusterConfig()
	} else {
		conf, err = clientcmd.BuildConfigFromFlags("", *path)
	}

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(conf)
}
