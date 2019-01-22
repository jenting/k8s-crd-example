package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	healthClientSet "github.com/hsiaoairplane/k8s-crd/pkg/client/clientset/versioned"
	healthInformerV1 "github.com/hsiaoairplane/k8s-crd/pkg/client/informers/externalversions/health/v1"
)

// retrieve the Kubernetes cluster client from outside of the cluster
func getKubernetesClient(logger *zerolog.Logger) (kubernetes.Interface, healthClientSet.Interface) {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logger.Fatal().Err(err).Msgf("getClusterConfig: %v", err)
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal().Err(err).Msgf("getClusterConfig: %v", err)
	}

	myresourceClient, err := healthClientSet.NewForConfig(config)
	if err != nil {
		logger.Fatal().Err(err).Msgf("getClusterConfig: %v", err)
	}

	logger.Info().Msgf("Successfully constructed k8s client")

	return client, myresourceClient
}

// main code path
func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel)

	// get the Kubernetes client for connectivity
	client, myresourceClient := getKubernetesClient(&logger)

	// retrieve our custom resource informer which was generated from
	// the code generator and pass it the custom resource client, specifying
	// we should be looking through all namespaces for listing and watching
	informer := healthInformerV1.NewHealthInformer(
		myresourceClient,
		metaV1.NamespaceAll,
		0,
		cache.Indexers{},
	)

	// create a new queue so that when the informer gets a resource that is either
	// a result of listing or watching, we can add an idenfitying key to the queue
	// so that it can be handled in the handler
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// add event handlers to handle the three types of events for resources:
	//  - adding new resources
	//  - updating existing resources
	//  - deleting resources
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// convert the resource object into a key (in this case
			// we are just doing it in the format of 'namespace/name')
			key, err := cache.MetaNamespaceKeyFunc(obj)
			logger.Info().Msgf("Add health endpoint: %s", key)
			if err == nil {
				// add the key to the queue for the handler to get
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// logger.Info().Msgf("%v", newObj.(*v1.Health).Spec.Action)
			// logger.Info().Msgf("%v", newObj.(*v1.Health).Spec.Switch)

			key, err := cache.MetaNamespaceKeyFunc(newObj)
			logger.Info().Msgf("Update health endpoint: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// DeletionHandlingMetaNamsespaceKeyFunc is a helper function that allows
			// us to check the DeletedFinalStateUnknown existence in the event that
			// a resource was deleted but it is still contained in the index
			//
			// this then in turn calls MetaNamespaceKeyFunc
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			logger.Info().Msgf("Delete health endpoint: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := Controller{
		logger:    &logger,
		clientset: client,
		informer:  informer,
		queue:     queue,
		handler:   &HealthHandler{&logger},
	}

	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// run the controller loop to process items
	go controller.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
