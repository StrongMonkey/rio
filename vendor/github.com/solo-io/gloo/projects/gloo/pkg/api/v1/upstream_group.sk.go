// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"log"
	"sort"

	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewUpstreamGroup(namespace, name string) *UpstreamGroup {
	upstreamgroup := &UpstreamGroup{}
	upstreamgroup.SetMetadata(core.Metadata{
		Name:      name,
		Namespace: namespace,
	})
	return upstreamgroup
}

func (r *UpstreamGroup) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *UpstreamGroup) SetStatus(status core.Status) {
	r.Status = status
}

func (r *UpstreamGroup) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	metaCopy.Generation = 0
	// investigate zeroing out owner refs as well
	return hashutils.HashAll(
		metaCopy,
		r.Destinations,
	)
}

func (r *UpstreamGroup) GroupVersionKind() schema.GroupVersionKind {
	return UpstreamGroupGVK
}

type UpstreamGroupList []*UpstreamGroup

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list UpstreamGroupList) Find(namespace, name string) (*UpstreamGroup, error) {
	for _, upstreamGroup := range list {
		if upstreamGroup.GetMetadata().Name == name {
			if namespace == "" || upstreamGroup.GetMetadata().Namespace == namespace {
				return upstreamGroup, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find upstreamGroup %v.%v", namespace, name)
}

func (list UpstreamGroupList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, upstreamGroup := range list {
		ress = append(ress, upstreamGroup)
	}
	return ress
}

func (list UpstreamGroupList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, upstreamGroup := range list {
		ress = append(ress, upstreamGroup)
	}
	return ress
}

func (list UpstreamGroupList) Names() []string {
	var names []string
	for _, upstreamGroup := range list {
		names = append(names, upstreamGroup.GetMetadata().Name)
	}
	return names
}

func (list UpstreamGroupList) NamespacesDotNames() []string {
	var names []string
	for _, upstreamGroup := range list {
		names = append(names, upstreamGroup.GetMetadata().Namespace+"."+upstreamGroup.GetMetadata().Name)
	}
	return names
}

func (list UpstreamGroupList) Sort() UpstreamGroupList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].GetMetadata().Less(list[j].GetMetadata())
	})
	return list
}

func (list UpstreamGroupList) Clone() UpstreamGroupList {
	var upstreamGroupList UpstreamGroupList
	for _, upstreamGroup := range list {
		upstreamGroupList = append(upstreamGroupList, resources.Clone(upstreamGroup).(*UpstreamGroup))
	}
	return upstreamGroupList
}

func (list UpstreamGroupList) Each(f func(element *UpstreamGroup)) {
	for _, upstreamGroup := range list {
		f(upstreamGroup)
	}
}

func (list UpstreamGroupList) EachResource(f func(element resources.Resource)) {
	for _, upstreamGroup := range list {
		f(upstreamGroup)
	}
}

func (list UpstreamGroupList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *UpstreamGroup) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

// Kubernetes Adapter for UpstreamGroup

func (o *UpstreamGroup) GetObjectKind() schema.ObjectKind {
	t := UpstreamGroupCrd.TypeMeta()
	return &t
}

func (o *UpstreamGroup) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*UpstreamGroup)
}

func (o *UpstreamGroup) DeepCopyInto(out *UpstreamGroup) {
	clone := resources.Clone(o).(*UpstreamGroup)
	*out = *clone
}

var (
	UpstreamGroupCrd = crd.NewCrd(
		"upstreamgroups",
		UpstreamGroupGVK.Group,
		UpstreamGroupGVK.Version,
		UpstreamGroupGVK.Kind,
		"ug",
		false,
		&UpstreamGroup{})
)

func init() {
	if err := crd.AddCrd(UpstreamGroupCrd); err != nil {
		log.Fatalf("could not add crd to global registry")
	}
}

var (
	UpstreamGroupGVK = schema.GroupVersionKind{
		Version: "v1",
		Group:   "gloo.solo.io",
		Kind:    "UpstreamGroup",
	}
)
