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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/jenting/k8s-crd-example/pkg/apis/health/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// HealthLister helps list Healths.
type HealthLister interface {
	// List lists all Healths in the indexer.
	List(selector labels.Selector) (ret []*v1.Health, err error)
	// Healths returns an object that can list and get Healths.
	Healths(namespace string) HealthNamespaceLister
	HealthListerExpansion
}

// healthLister implements the HealthLister interface.
type healthLister struct {
	indexer cache.Indexer
}

// NewHealthLister returns a new HealthLister.
func NewHealthLister(indexer cache.Indexer) HealthLister {
	return &healthLister{indexer: indexer}
}

// List lists all Healths in the indexer.
func (s *healthLister) List(selector labels.Selector) (ret []*v1.Health, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Health))
	})
	return ret, err
}

// Healths returns an object that can list and get Healths.
func (s *healthLister) Healths(namespace string) HealthNamespaceLister {
	return healthNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// HealthNamespaceLister helps list and get Healths.
type HealthNamespaceLister interface {
	// List lists all Healths in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Health, err error)
	// Get retrieves the Health from the indexer for a given namespace and name.
	Get(name string) (*v1.Health, error)
	HealthNamespaceListerExpansion
}

// healthNamespaceLister implements the HealthNamespaceLister
// interface.
type healthNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Healths in the indexer for a given namespace.
func (s healthNamespaceLister) List(selector labels.Selector) (ret []*v1.Health, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Health))
	})
	return ret, err
}

// Get retrieves the Health from the indexer for a given namespace and name.
func (s healthNamespaceLister) Get(name string) (*v1.Health, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("health"), name)
	}
	return obj.(*v1.Health), nil
}
