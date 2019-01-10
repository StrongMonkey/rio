package v1

import (
	riocattleiov1 "github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceScaleRecommendation) DeepCopyInto(out *ServiceScaleRecommendation) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceScaleRecommendation.
func (in *ServiceScaleRecommendation) DeepCopy() *ServiceScaleRecommendation {
	if in == nil {
		return nil
	}
	out := new(ServiceScaleRecommendation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceScaleRecommendation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceScaleRecommendationList) DeepCopyInto(out *ServiceScaleRecommendationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceScaleRecommendation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceScaleRecommendationList.
func (in *ServiceScaleRecommendationList) DeepCopy() *ServiceScaleRecommendationList {
	if in == nil {
		return nil
	}
	out := new(ServiceScaleRecommendationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceScaleRecommendationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceScaleRecommendationSpec) DeepCopyInto(out *ServiceScaleRecommendationSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceScaleRecommendationSpec.
func (in *ServiceScaleRecommendationSpec) DeepCopy() *ServiceScaleRecommendationSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceScaleRecommendationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceScaleRecommendationStatus) DeepCopyInto(out *ServiceScaleRecommendationStatus) {
	*out = *in
	if in.DesiredScale != nil {
		in, out := &in.DesiredScale, &out.DesiredScale
		*out = new(int32)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]riocattleiov1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceScaleRecommendationStatus.
func (in *ServiceScaleRecommendationStatus) DeepCopy() *ServiceScaleRecommendationStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceScaleRecommendationStatus)
	in.DeepCopyInto(out)
	return out
}
