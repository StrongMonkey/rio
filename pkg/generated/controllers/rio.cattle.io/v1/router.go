/*
Copyright 2019 Rancher Labs.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"

	v1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	clientset "github.com/rancher/rio/pkg/generated/clientset/versioned/typed/rio.cattle.io/v1"
	informers "github.com/rancher/rio/pkg/generated/informers/externalversions/rio.cattle.io/v1"
	listers "github.com/rancher/rio/pkg/generated/listers/rio.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type RouterHandler func(string, *v1.Router) (*v1.Router, error)

type RouterController interface {
	RouterClient

	OnChange(ctx context.Context, name string, sync RouterHandler)
	OnRemove(ctx context.Context, name string, sync RouterHandler)
	Enqueue(namespace, name string)

	Cache() RouterCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type RouterClient interface {
	Create(*v1.Router) (*v1.Router, error)
	Update(*v1.Router) (*v1.Router, error)
	UpdateStatus(*v1.Router) (*v1.Router, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Router, error)
	List(namespace string, opts metav1.ListOptions) (*v1.RouterList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Router, err error)
}

type RouterCache interface {
	Get(namespace, name string) (*v1.Router, error)
	List(namespace string, selector labels.Selector) ([]*v1.Router, error)

	AddIndexer(indexName string, indexer RouterIndexer)
	GetByIndex(indexName, key string) ([]*v1.Router, error)
}

type RouterIndexer func(obj *v1.Router) ([]string, error)

type routerController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.RoutersGetter
	informer          informers.RouterInformer
	gvk               schema.GroupVersionKind
}

func NewRouterController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.RoutersGetter, informer informers.RouterInformer) RouterController {
	return &routerController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromRouterHandlerToHandler(sync RouterHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Router
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Router))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *routerController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Router))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateRouterOnChange(updater generic.Updater, handler RouterHandler) RouterHandler {
	return func(key string, obj *v1.Router) (*v1.Router, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.Router)
			}
		}

		return copyObj, err
	}
}

func (c *routerController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *routerController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *routerController) OnChange(ctx context.Context, name string, sync RouterHandler) {
	c.AddGenericHandler(ctx, name, FromRouterHandlerToHandler(sync))
}

func (c *routerController) OnRemove(ctx context.Context, name string, sync RouterHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromRouterHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *routerController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *routerController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *routerController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *routerController) Cache() RouterCache {
	return &routerCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *routerController) Create(obj *v1.Router) (*v1.Router, error) {
	return c.clientGetter.Routers(obj.Namespace).Create(obj)
}

func (c *routerController) Update(obj *v1.Router) (*v1.Router, error) {
	return c.clientGetter.Routers(obj.Namespace).Update(obj)
}

func (c *routerController) UpdateStatus(obj *v1.Router) (*v1.Router, error) {
	return c.clientGetter.Routers(obj.Namespace).UpdateStatus(obj)
}

func (c *routerController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Routers(namespace).Delete(name, options)
}

func (c *routerController) Get(namespace, name string, options metav1.GetOptions) (*v1.Router, error) {
	return c.clientGetter.Routers(namespace).Get(name, options)
}

func (c *routerController) List(namespace string, opts metav1.ListOptions) (*v1.RouterList, error) {
	return c.clientGetter.Routers(namespace).List(opts)
}

func (c *routerController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Routers(namespace).Watch(opts)
}

func (c *routerController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Router, err error) {
	return c.clientGetter.Routers(namespace).Patch(name, pt, data, subresources...)
}

type routerCache struct {
	lister  listers.RouterLister
	indexer cache.Indexer
}

func (c *routerCache) Get(namespace, name string) (*v1.Router, error) {
	return c.lister.Routers(namespace).Get(name)
}

func (c *routerCache) List(namespace string, selector labels.Selector) ([]*v1.Router, error) {
	return c.lister.Routers(namespace).List(selector)
}

func (c *routerCache) AddIndexer(indexName string, indexer RouterIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Router))
		},
	}))
}

func (c *routerCache) GetByIndex(indexName, key string) (result []*v1.Router, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.Router))
	}
	return result, nil
}
