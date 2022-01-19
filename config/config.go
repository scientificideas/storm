package config

import (
	"flag"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

type Config struct {
	Filter      string
	Chaos       string
	Targets     string
	Startfast   bool
	RuntimeType string
	Namespace   string
	K8sContext  string
	Kubeconfig  string
}

func GetConfig() *Config {
	filter := flag.String("filter", "", "filter containers by name")
	chaos := flag.String("chaos", "medium", "easy, medium or hard level of chaos")
	targets := flag.String("targets", "", `if you only want to expose certain containers, list them here ("container1,container2,container3")`)
	startfast := flag.Bool("startfast", false, `start stopped containers immediately ("true" or "false")`)
	runtimeType := flag.String("runtime", "k8s", "which orchestrator to interact with")
	namespace := flag.String("kube-namespace", "", "k8s namespace")
	k8sContext := flag.String("kube-context", "", "k8s context")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	if *runtimeType == "k8s" {
		if *startfast == true {
			log.Println(`You can't specify 'startfast' option for k8s environment, because pods can't be stopped, they just deleted and recreated by k8s. 
In k8s, pods always start immediately (no need for 'startfast' option).`)
		}
	}

	return &Config{
		Filter:      *filter,
		Chaos:       *chaos,
		Targets:     *targets,
		Startfast:   *startfast,
		RuntimeType: *runtimeType,
		Namespace:   *namespace,
		K8sContext:  *k8sContext,
		Kubeconfig:  *kubeconfig,
	}
}
