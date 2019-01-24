package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	maxRetries int = 5
)

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// queueing, and handling of resource changes
type Controller struct {
	logger    *zerolog.Logger
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	server    *http.Server
	handler   Handler
}

// Run is the main path of execution for the controller loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// don't let panics crash the process
	defer utilruntime.HandleCrash()
	// make sure the work queue is quit which will trigger workers to end
	defer c.queue.ShutDown()

	c.logger.Info().Msg("Controller.Run: initiating")

	// run the informer to start listing and watching resources
	go c.informer.Run(stopCh)

	// run http router
	go c.handler.Run(stopCh)

	// wait for the caches to synchronize before starting the worker
	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}

	c.logger.Info().Msg("Controller.Run: cache sync complete")

	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced allows us to satisfy the Controller interface
// by wiring up the informer's HasSynced method to it
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// runWorker executes the loop to process new items added to the queue
func (c *Controller) runWorker() {
	c.logger.Info().Msg("Controller.runWorker: starting")

	// processNextWorkItem will automatically wait until there's work available
	for c.processNextItem() {
		// continue looping
		c.logger.Info().Msg("Controller.runWorker: processing next item")
	}

	c.logger.Info().Msg("Controller.runWorker: completed")
}

// processNextItem retrieves each queued item and takes the
// necessary handler action based off of if the item was
// created or deleted
func (c *Controller) processNextItem() bool {
	c.logger.Info().Msg("Controller.processNextItem: start")

	// pull the next work item from queue.  It should be a key we use to lookup
	// something in a cache
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// you always have to indicate to the queue that you've completed a piece of
	// work
	defer c.queue.Done(key)

	// do your work on the key.
	err := c.processItem(key.(string))
	if err == nil {
		// No error, tell the queue to stop tracking history
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < maxRetries {
		c.logger.Error().Err(err).Msgf("Error processing %s (will retry): %v", key, err)
		// requeue the item to work on later
		c.queue.AddRateLimited(key)
	} else {
		// err != nil and too many retries
		c.logger.Error().Err(err).Msgf("Error processing %s (giving up): %v", key, err)
		c.queue.Forget(key)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(key string) error {
	c.logger.Info().Msgf("Processing change %s", key)

	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
	}

	if !exists {
		c.handler.ObjectDeleted(obj)
		return nil
	}

	c.handler.ObjectUpdated(obj)

	return nil
}
