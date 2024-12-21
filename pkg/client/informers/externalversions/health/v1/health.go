/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	healthv1 "github.com/jenting/k8s-crd/pkg/apis/health/v1"
	versioned "github.com/jenting/k8s-crd/pkg/client/clientset/versioned"
	internalinterfaces "github.com/jenting/k8s-crd/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/jenting/k8s-crd/pkg/client/listers/health/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// HealthInformer provides access to a shared informer and lister for
// Healths.
type HealthInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.HealthLister
}

type healthInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewHealthInformer constructs a new informer for Health type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewHealthInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredHealthInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredHealthInformer constructs a new informer for Health type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredHealthInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.jentingV1().Healths(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.jentingV1().Healths(namespace).Watch(options)
			},
		},
		&healthv1.Health{},
		resyncPeriod,
		indexers,
	)
}

func (f *healthInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredHealthInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *healthInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&healthv1.Health{}, f.defaultInformer)
}

func (f *healthInformer) Lister() v1.HealthLister {
	return v1.NewHealthLister(f.Informer().GetIndexer())
}
