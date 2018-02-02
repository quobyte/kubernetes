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

// FakeQuobyteServices implements QuobyteServiceInterface
type FakeQuobyteServices struct {
	Fake *FakeQuobyteV1
	ns   string
}

var quobyteservicesResource = schema.GroupVersionResource{Group: "quobyte.com", Version: "v1", Resource: "quobyteservices"}

var quobyteservicesKind = schema.GroupVersionKind{Group: "quobyte.com", Version: "v1", Kind: "QuobyteService"}

// Get takes name of the quobyteService, and returns the corresponding quobyteService object, and an error if there is any.
func (c *FakeQuobyteServices) Get(name string, options v1.GetOptions) (result *quobyte_com_v1.QuobyteService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(quobyteservicesResource, c.ns, name), &quobyte_com_v1.QuobyteService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteService), err
}

// List takes label and field selectors, and returns the list of QuobyteServices that match those selectors.
func (c *FakeQuobyteServices) List(opts v1.ListOptions) (result *quobyte_com_v1.QuobyteServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(quobyteservicesResource, quobyteservicesKind, c.ns, opts), &quobyte_com_v1.QuobyteServiceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &quobyte_com_v1.QuobyteServiceList{}
	for _, item := range obj.(*quobyte_com_v1.QuobyteServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested quobyteServices.
func (c *FakeQuobyteServices) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(quobyteservicesResource, c.ns, opts))

}

// Create takes the representation of a quobyteService and creates it.  Returns the server's representation of the quobyteService, and an error, if there is any.
func (c *FakeQuobyteServices) Create(quobyteService *quobyte_com_v1.QuobyteService) (result *quobyte_com_v1.QuobyteService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(quobyteservicesResource, c.ns, quobyteService), &quobyte_com_v1.QuobyteService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteService), err
}

// Update takes the representation of a quobyteService and updates it. Returns the server's representation of the quobyteService, and an error, if there is any.
func (c *FakeQuobyteServices) Update(quobyteService *quobyte_com_v1.QuobyteService) (result *quobyte_com_v1.QuobyteService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(quobyteservicesResource, c.ns, quobyteService), &quobyte_com_v1.QuobyteService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteService), err
}

// Delete takes name of the quobyteService and deletes it. Returns an error if one occurs.
func (c *FakeQuobyteServices) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(quobyteservicesResource, c.ns, name), &quobyte_com_v1.QuobyteService{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeQuobyteServices) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(quobyteservicesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &quobyte_com_v1.QuobyteServiceList{})
	return err
}

// Patch applies the patch and returns the patched quobyteService.
func (c *FakeQuobyteServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *quobyte_com_v1.QuobyteService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(quobyteservicesResource, c.ns, name, data, subresources...), &quobyte_com_v1.QuobyteService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*quobyte_com_v1.QuobyteService), err
}
