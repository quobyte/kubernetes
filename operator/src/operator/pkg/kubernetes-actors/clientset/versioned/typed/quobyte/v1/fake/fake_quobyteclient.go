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

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	quobyte_com_v1 "operator/pkg/api/quobyte.com/v1"
)

// FakeQuobyteClients implements QuobyteClientInterface
type FakeQuobyteClients struct {
	Fake *FakeQuobyteV1
	ns   string
}

var quobyteclientsResource = schema.GroupVersionResource{Group: "quobyte.com", Version: "v1", Resource: "quobyteclients"}

var quobyteclientsKind = schema.GroupVersionKind{Group: "quobyte.com", Version: "v1", Kind: "QuobyteClient"}

// Get takes name of the quobyteClient, and returns the corresponding quobyteClient object, and an error if there is any.
func (c *FakeQuobyteClients) Get(name string, options v1.GetOptions) (result *quobyte_com_v1.QuobyteClient, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(quobyteclientsResource, c.ns, name), &quobyte_com_v1.QuobyteClient{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteClient), err
}

// List takes label and field selectors, and returns the list of QuobyteClients that match those selectors.
func (c *FakeQuobyteClients) List(opts v1.ListOptions) (result *quobyte_com_v1.QuobyteClientList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(quobyteclientsResource, quobyteclientsKind, c.ns, opts), &quobyte_com_v1.QuobyteClientList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &quobyte_com_v1.QuobyteClientList{}
	for _, item := range obj.(*quobyte_com_v1.QuobyteClientList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested quobyteClients.
func (c *FakeQuobyteClients) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(quobyteclientsResource, c.ns, opts))

}

// Create takes the representation of a quobyteClient and creates it.  Returns the server's representation of the quobyteClient, and an error, if there is any.
func (c *FakeQuobyteClients) Create(quobyteClient *quobyte_com_v1.QuobyteClient) (result *quobyte_com_v1.QuobyteClient, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(quobyteclientsResource, c.ns, quobyteClient), &quobyte_com_v1.QuobyteClient{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteClient), err
}

// Update takes the representation of a quobyteClient and updates it. Returns the server's representation of the quobyteClient, and an error, if there is any.
func (c *FakeQuobyteClients) Update(quobyteClient *quobyte_com_v1.QuobyteClient) (result *quobyte_com_v1.QuobyteClient, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(quobyteclientsResource, c.ns, quobyteClient), &quobyte_com_v1.QuobyteClient{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteClient), err
}

// Delete takes name of the quobyteClient and deletes it. Returns an error if one occurs.
func (c *FakeQuobyteClients) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(quobyteclientsResource, c.ns, name), &quobyte_com_v1.QuobyteClient{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeQuobyteClients) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(quobyteclientsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &quobyte_com_v1.QuobyteClientList{})
	return err
}

// Patch applies the patch and returns the patched quobyteClient.
func (c *FakeQuobyteClients) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *quobyte_com_v1.QuobyteClient, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(quobyteclientsResource, c.ns, name, data, subresources...), &quobyte_com_v1.QuobyteClient{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteClient), err
}
