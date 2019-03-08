package main

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func main() {
	var clientSet kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		clientSet = GetClientOutOfCluster()
	} else {
		clientSet = GetClient()
	}

	factory := informers.NewSharedInformerFactory(clientSet, 1000000000)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// The pod won't have an IP at the time of the initial add event
			//  When an IP is assigned to the new pod it will trigger an update event
			mObj := obj.(*v1.Pod)
			log.Printf("New Pod Added: %s %s", mObj.Name, mObj.Status.PodIP)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			mObj := new.(*v1.Pod)
			oldObj := old.(*v1.Pod)
			// Check if the update event assigned an IP to the pod.  If so, log it
			if mObj.Status.PodIP != oldObj.Status.PodIP {
				log.Printf("Pod Updated: %s %s", mObj.Name, mObj.Status.PodIP)
			}
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(*v1.Pod)
			log.Printf("Pod Deleted: %s %s", mObj.Name, mObj.Status.PodIP)
		},
	})

	informer.Run(stopper)

}

// GetClient returns a k8s clientset to the request from inside of cluster
func GetClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	return clientset
}
