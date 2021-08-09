package k8s

import (
	"github.com/blackstorm/ingress-go/pkg/common"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientWithFeedback(kubeconfigPath string) (*kubernetes.Clientset, error) {
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

func GetClientWithInClusterConfig() (*kubernetes.Clientset, error) {
	return GetClientWithFeedback(common.EMPTY_STRING)
}
