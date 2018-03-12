/*
Copyright 2017 The Kubernetes Authors.

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

package v1

import (
	scheme "operator/pkg/kubernetes-actors/clientset/versioned/scheme"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "operator/pkg/api/quobyte.com/v1"
)

// QuobyteServicesGetter has a method to return a QuobyteServiceInterface.
// A group's client should implement this interface.
type QuobyteServicesGetter interface {
	QuobyteServices(namespace string) QuobyteServiceInterface
}

// QuobyteServiceInterface has methods to work with QuobyteService resources.
type QuobyteServiceInterface interface {
	Create(*v1.QuobyteService) (*v1.QuobyteService, error)
	Update(*v1.QuobyteService) (*v1.QuobyteService, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.QuobyteService, error)
	List(opts meta_v1.ListOptions) (*v1.QuobyteServiceList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.QuobyteService, err error)
	QuobyteServiceExpansion
}

// quobyteServices implements QuobyteServiceInterface
type quobyteServices struct {
	client rest.Interface
	ns     string
}

// newQuobyteServices returns a QuobyteServices
func newQuobyteServices(c *QuobyteV1Client, namespace string) *quobyteServices {
	return &quobyteServices{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the quobyteService, and returns the corresponding quobyteService object, and an error if there is any.
func (c *quobyteServices) Get(name string, options meta_v1.GetOptions) (result *v1.QuobyteService, err error) {
	result = &v1.QuobyteService{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("quobyteservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of QuobyteServices that match those selectors.
func (c *quobyteServices) List(opts meta_v1.ListOptions) (result *v1.QuobyteServiceList, err error) {
	result = &v1.QuobyteServiceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("quobyteservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested quobyteServices.
func (c *quobyteServices) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("quobyteservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a quobyteService and creates it.  Returns the server's representation of the quobyteService, and an error, if there is any.
func (c *quobyteServices) Create(quobyteService *v1.QuobyteService) (result *v1.QuobyteService, err error) {
	result = &v1.QuobyteService{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("quobyteservices").
		Body(quobyteService).
		Do().
		Into(result)
	return
}

// Update takes the representation of a quobyteService and updates it. Returns the server's representation of the quobyteService, and an error, if there is any.
func (c *quobyteServices) Update(quobyteService *v1.QuobyteService) (result *v1.QuobyteService, err error) {
	result = &v1.QuobyteService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("quobyteservices").
		Name(quobyteService.Name).
		Body(quobyteService).
		Do().
		Into(result)
	return
}

// Delete takes name of the quobyteService and deletes it. Returns an error if one occurs.
func (c *quobyteServices) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("quobyteservices").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *quobyteServices) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("quobyteservices").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched quobyteService.
func (c *quobyteServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.QuobyteService, err error) {
	result = &v1.QuobyteService{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("quobyteservices").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
