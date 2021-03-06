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

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"

	"github.com/rancher/wrangler/pkg/generic"
	v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	clientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	informers "github.com/tektoncd/pipeline/pkg/client/informers/externalversions/pipeline/v1alpha1"
	listers "github.com/tektoncd/pipeline/pkg/client/listers/pipeline/v1alpha1"
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

type TaskRunHandler func(string, *v1alpha1.TaskRun) (*v1alpha1.TaskRun, error)

type TaskRunController interface {
	TaskRunClient

	OnChange(ctx context.Context, name string, sync TaskRunHandler)
	OnRemove(ctx context.Context, name string, sync TaskRunHandler)
	Enqueue(namespace, name string)

	Cache() TaskRunCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type TaskRunClient interface {
	Create(*v1alpha1.TaskRun) (*v1alpha1.TaskRun, error)
	Update(*v1alpha1.TaskRun) (*v1alpha1.TaskRun, error)
	UpdateStatus(*v1alpha1.TaskRun) (*v1alpha1.TaskRun, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.TaskRun, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.TaskRunList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TaskRun, err error)
}

type TaskRunCache interface {
	Get(namespace, name string) (*v1alpha1.TaskRun, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.TaskRun, error)

	AddIndexer(indexName string, indexer TaskRunIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.TaskRun, error)
}

type TaskRunIndexer func(obj *v1alpha1.TaskRun) ([]string, error)

type taskRunController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.TaskRunsGetter
	informer          informers.TaskRunInformer
	gvk               schema.GroupVersionKind
}

func NewTaskRunController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.TaskRunsGetter, informer informers.TaskRunInformer) TaskRunController {
	return &taskRunController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromTaskRunHandlerToHandler(sync TaskRunHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.TaskRun
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.TaskRun))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *taskRunController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.TaskRun))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateTaskRunOnChange(updater generic.Updater, handler TaskRunHandler) TaskRunHandler {
	return func(key string, obj *v1alpha1.TaskRun) (*v1alpha1.TaskRun, error) {
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
				copyObj = newObj.(*v1alpha1.TaskRun)
			}
		}

		return copyObj, err
	}
}

func (c *taskRunController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *taskRunController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *taskRunController) OnChange(ctx context.Context, name string, sync TaskRunHandler) {
	c.AddGenericHandler(ctx, name, FromTaskRunHandlerToHandler(sync))
}

func (c *taskRunController) OnRemove(ctx context.Context, name string, sync TaskRunHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromTaskRunHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *taskRunController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *taskRunController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *taskRunController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *taskRunController) Cache() TaskRunCache {
	return &taskRunCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *taskRunController) Create(obj *v1alpha1.TaskRun) (*v1alpha1.TaskRun, error) {
	return c.clientGetter.TaskRuns(obj.Namespace).Create(obj)
}

func (c *taskRunController) Update(obj *v1alpha1.TaskRun) (*v1alpha1.TaskRun, error) {
	return c.clientGetter.TaskRuns(obj.Namespace).Update(obj)
}

func (c *taskRunController) UpdateStatus(obj *v1alpha1.TaskRun) (*v1alpha1.TaskRun, error) {
	return c.clientGetter.TaskRuns(obj.Namespace).UpdateStatus(obj)
}

func (c *taskRunController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.TaskRuns(namespace).Delete(name, options)
}

func (c *taskRunController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.TaskRun, error) {
	return c.clientGetter.TaskRuns(namespace).Get(name, options)
}

func (c *taskRunController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.TaskRunList, error) {
	return c.clientGetter.TaskRuns(namespace).List(opts)
}

func (c *taskRunController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.TaskRuns(namespace).Watch(opts)
}

func (c *taskRunController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TaskRun, err error) {
	return c.clientGetter.TaskRuns(namespace).Patch(name, pt, data, subresources...)
}

type taskRunCache struct {
	lister  listers.TaskRunLister
	indexer cache.Indexer
}

func (c *taskRunCache) Get(namespace, name string) (*v1alpha1.TaskRun, error) {
	return c.lister.TaskRuns(namespace).Get(name)
}

func (c *taskRunCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.TaskRun, error) {
	return c.lister.TaskRuns(namespace).List(selector)
}

func (c *taskRunCache) AddIndexer(indexName string, indexer TaskRunIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.TaskRun))
		},
	}))
}

func (c *taskRunCache) GetByIndex(indexName, key string) (result []*v1alpha1.TaskRun, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.TaskRun))
	}
	return result, nil
}
