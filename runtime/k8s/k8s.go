package k8s

import (
	"context"
	"errors"
	"flag"
	"github.com/scientificideas/storm/chaos"
	"github.com/scientificideas/storm/container"
	"github.com/scientificideas/storm/container/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	"path/filepath"
)

type K8S struct {
	cli       *kubernetes.Clientset
	chaos     chaos.Chaos
	filter    string
	namespace string
}

func NewK8sClient(chaosType, filter, namespace string) (*K8S, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
	pods, err := k.cli.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
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
	return k.cli.CoreV1().Pods("").Delete(context.TODO(), name, metav1.DeleteOptions{
		GracePeriodSeconds: &stopTime,
	})
}
func (k *K8S) StartContainer(ctx context.Context, name string) error {
	return nil
}
