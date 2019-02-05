// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DelegatedAuthentication) DeepCopyInto(out *DelegatedAuthentication) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DelegatedAuthentication.
func (in *DelegatedAuthentication) DeepCopy() *DelegatedAuthentication {
	if in == nil {
		return nil
	}
	out := new(DelegatedAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DelegatedAuthorization) DeepCopyInto(out *DelegatedAuthorization) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DelegatedAuthorization.
func (in *DelegatedAuthorization) DeepCopy() *DelegatedAuthorization {
	if in == nil {
		return nil
	}
	out := new(DelegatedAuthorization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenerationHistory) DeepCopyInto(out *GenerationHistory) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenerationHistory.
func (in *GenerationHistory) DeepCopy() *GenerationHistory {
	if in == nil {
		return nil
	}
	out := new(GenerationHistory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericOperatorConfig) DeepCopyInto(out *GenericOperatorConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ServingInfo.DeepCopyInto(&out.ServingInfo)
	out.LeaderElection = in.LeaderElection
	out.Authentication = in.Authentication
	out.Authorization = in.Authorization
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericOperatorConfig.
func (in *GenericOperatorConfig) DeepCopy() *GenericOperatorConfig {
	if in == nil {
		return nil
	}
	out := new(GenericOperatorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GenericOperatorConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingConfig) DeepCopyInto(out *LoggingConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingConfig.
func (in *LoggingConfig) DeepCopy() *LoggingConfig {
	if in == nil {
		return nil
	}
	out := new(LoggingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorCondition) DeepCopyInto(out *OperatorCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorCondition.
func (in *OperatorCondition) DeepCopy() *OperatorCondition {
	if in == nil {
		return nil
	}
	out := new(OperatorCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorSpec) DeepCopyInto(out *OperatorSpec) {
	*out = *in
	out.Logging = in.Logging
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorSpec.
func (in *OperatorSpec) DeepCopy() *OperatorSpec {
	if in == nil {
		return nil
	}
	out := new(OperatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorStatus) DeepCopyInto(out *OperatorStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]OperatorCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CurrentAvailability != nil {
		in, out := &in.CurrentAvailability, &out.CurrentAvailability
		if *in == nil {
			*out = nil
		} else {
			*out = new(VersionAvailablity)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.TargetAvailability != nil {
		in, out := &in.TargetAvailability, &out.TargetAvailability
		if *in == nil {
			*out = nil
		} else {
			*out = new(VersionAvailablity)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorStatus.
func (in *OperatorStatus) DeepCopy() *OperatorStatus {
	if in == nil {
		return nil
	}
	out := new(OperatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionAvailablity) DeepCopyInto(out *VersionAvailablity) {
	*out = *in
	if in.Errors != nil {
		in, out := &in.Errors, &out.Errors
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Generations != nil {
		in, out := &in.Generations, &out.Generations
		*out = make([]GenerationHistory, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionAvailablity.
func (in *VersionAvailablity) DeepCopy() *VersionAvailablity {
	if in == nil {
		return nil
	}
	out := new(VersionAvailablity)
	in.DeepCopyInto(out)
	return out
}
