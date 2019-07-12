package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientSet *kubernetes.Clientset

func init() {
	var err error
	clientSet, err = inClusterClient()
	if err != nil {
		panic(err)
	}
}

func inClusterClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	return kubernetes.NewForConfig(config)
}
