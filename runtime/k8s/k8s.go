package k8s

import (
	"context"
	"flag"
	"github.com/docker/docker/api/types"
	"github.com/scientificideas/storm/chaos"
	"github.com/scientificideas/storm/container"
	"github.com/scientificideas/storm/container/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	"path/filepath"
	"time"
)

type K8S struct {
	cli    *kubernetes.Clientset
	chaos  chaos.Chaos
	filter string
}

type loopType int

const (
	stop = iota
	start
	stopAndStartImmediately
)

func NewK8sClient(chaosType, filter string) (*K8S, error) {
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
	return &K8S{cli: clientset, chaos: chaos.NewChaos(chaosType), filter: filter}, nil
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
	stopTime := 1 * time.Millisecond
	return d.cli.ContainerStop(ctx, name, &stopTime)
}
func (k *K8S) RmContainer(ctx context.Context, name string) error {
	return d.cli.ContainerKill(ctx, name, "SIGKILL")
}
func (k *K8S) StartContainer(ctx context.Context, name string) error {
	return d.cli.ContainerStart(ctx, name, types.ContainerStartOptions{})
}
