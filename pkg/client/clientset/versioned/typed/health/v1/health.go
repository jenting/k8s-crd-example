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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/jenting/k8s-crd-example/pkg/apis/health/v1"
	scheme "github.com/jenting/k8s-crd-example/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HealthsGetter has a method to return a HealthInterface.
// A group's client should implement this interface.
type HealthsGetter interface {
	Healths(namespace string) HealthInterface
}

// HealthInterface has methods to work with Health resources.
type HealthInterface interface {
	Create(*v1.Health) (*v1.Health, error)
	Update(*v1.Health) (*v1.Health, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Health, error)
	List(opts metav1.ListOptions) (*v1.HealthList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Health, err error)
	HealthExpansion
}

// healths implements HealthInterface
type healths struct {
	client rest.Interface
	ns     string
}

// newHealths returns a Healths
func newHealths(c *jentingV1Client, namespace string) *healths {
	return &healths{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the health, and returns the corresponding health object, and an error if there is any.
func (c *healths) Get(name string, options metav1.GetOptions) (result *v1.Health, err error) {
	result = &v1.Health{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("healths").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Healths that match those selectors.
func (c *healths) List(opts metav1.ListOptions) (result *v1.HealthList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HealthList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("healths").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested healths.
func (c *healths) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("healths").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a health and creates it.  Returns the server's representation of the health, and an error, if there is any.
func (c *healths) Create(health *v1.Health) (result *v1.Health, err error) {
	result = &v1.Health{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("healths").
		Body(health).
		Do().
		Into(result)
	return
}

// Update takes the representation of a health and updates it. Returns the server's representation of the health, and an error, if there is any.
func (c *healths) Update(health *v1.Health) (result *v1.Health, err error) {
	result = &v1.Health{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("healths").
		Name(health.Name).
		Body(health).
		Do().
		Into(result)
	return
}

// Delete takes name of the health and deletes it. Returns an error if one occurs.
func (c *healths) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("healths").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *healths) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("healths").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched health.
func (c *healths) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Health, err error) {
	result = &v1.Health{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("healths").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
