package k8s

import (
	"context"
	"errors"
	"github.com/scientificideas/storm/chaos"
	"github.com/scientificideas/storm/container"
	"github.com/scientificideas/storm/container/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"math/rand"
)

type K8S struct {
	cli       *kubernetes.Clientset
	chaos     chaos.Chaos
	filter    string
	namespace string
}

func NewK8sClient(chaosType, filter, namespace, kubeconfig, k8sContext string) (*K8S, error) {
	// get config from kubeconfig
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}, CurrentContext: k8sContext}).ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &K8S{cli: clientset, chaos: chaos.NewChaos(chaosType), filter: filter, namespace: namespace}, nil
}

func (k *K8S) Type() string {
	return "k8s"
}

func (k *K8S) Chaos() chaos.Chaos {
	return k.chaos
}

func (k *K8S) GetContainers(ctx context.Context, all bool) ([]container.Container, error) {
	pods, err := k.cli.CoreV1().Pods(k.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	rand.Shuffle(len(pods.Items), func(i, j int) {
		pods.Items[i], pods.Items[j] = pods.Items[j], pods.Items[i]
	})

	var result []container.Container
	for _, c := range pods.Items {
		result = append(result, k8s.NewPod(c))
	}
	return result, nil
}

func (k *K8S) StopContainer(ctx context.Context, name string) error {
	return errors.New("stop doesn't not implemented for k8s")
}
func (k *K8S) RmContainer(ctx context.Context, name string) error {
	var stopTime int64 = 3
	return k.cli.CoreV1().Pods(k.namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		GracePeriodSeconds: &stopTime,
	})
}
func (k *K8S) StartContainer(ctx context.Context, name string) error {
	return nil
}
