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

// QuobyteClientsGetter has a method to return a QuobyteClientInterface.
// A group's client should implement this interface.
type QuobyteClientsGetter interface {
	QuobyteClients(namespace string) QuobyteClientInterface
}

// QuobyteClientInterface has methods to work with QuobyteClient resources.
type QuobyteClientInterface interface {
	Create(*v1.QuobyteClient) (*v1.QuobyteClient, error)
	Update(*v1.QuobyteClient) (*v1.QuobyteClient, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.QuobyteClient, error)
	List(opts meta_v1.ListOptions) (*v1.QuobyteClientList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.QuobyteClient, err error)
	QuobyteClientExpansion
}

// quobyteClients implements QuobyteClientInterface
type quobyteClients struct {
	client rest.Interface
	ns     string
}

// newQuobyteClients returns a QuobyteClients
func newQuobyteClients(c *QuobyteV1Client, namespace string) *quobyteClients {
	return &quobyteClients{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the quobyteClient, and returns the corresponding quobyteClient object, and an error if there is any.
func (c *quobyteClients) Get(name string, options meta_v1.GetOptions) (result *v1.QuobyteClient, err error) {
	result = &v1.QuobyteClient{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("quobyteclients").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of QuobyteClients that match those selectors.
func (c *quobyteClients) List(opts meta_v1.ListOptions) (result *v1.QuobyteClientList, err error) {
	result = &v1.QuobyteClientList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("quobyteclients").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested quobyteClients.
func (c *quobyteClients) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("quobyteclients").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a quobyteClient and creates it.  Returns the server's representation of the quobyteClient, and an error, if there is any.
func (c *quobyteClients) Create(quobyteClient *v1.QuobyteClient) (result *v1.QuobyteClient, err error) {
	result = &v1.QuobyteClient{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("quobyteclients").
		Body(quobyteClient).
		Do().
		Into(result)
	return
}

// Update takes the representation of a quobyteClient and updates it. Returns the server's representation of the quobyteClient, and an error, if there is any.
func (c *quobyteClients) Update(quobyteClient *v1.QuobyteClient) (result *v1.QuobyteClient, err error) {
	result = &v1.QuobyteClient{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("quobyteclients").
		Name(quobyteClient.Name).
		Body(quobyteClient).
		Do().
		Into(result)
	return
}

// Delete takes name of the quobyteClient and deletes it. Returns an error if one occurs.
func (c *quobyteClients) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("quobyteclients").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *quobyteClients) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("quobyteclients").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched quobyteClient.
func (c *quobyteClients) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.QuobyteClient, err error) {
	result = &v1.QuobyteClient{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("quobyteclients").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
