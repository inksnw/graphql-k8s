package lib

import (
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/phuslu/log"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var Document *openapi_v2.Document
var ClientSet *kubernetes.Clientset
var DynamicClient *dynamic.DynamicClient

func init() {

	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	DynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	Document, err = ClientSet.Discovery().OpenAPISchema()
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

}
