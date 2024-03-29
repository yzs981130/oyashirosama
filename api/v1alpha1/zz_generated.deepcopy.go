// +build !ignore_autogenerated

/*
Copyright 2021 yzs.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaseJob) DeepCopyInto(out *LeaseJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaseJob.
func (in *LeaseJob) DeepCopy() *LeaseJob {
	if in == nil {
		return nil
	}
	out := new(LeaseJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LeaseJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaseJobList) DeepCopyInto(out *LeaseJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LeaseJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaseJobList.
func (in *LeaseJobList) DeepCopy() *LeaseJobList {
	if in == nil {
		return nil
	}
	out := new(LeaseJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LeaseJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaseJobPodStatistics) DeepCopyInto(out *LeaseJobPodStatistics) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaseJobPodStatistics.
func (in *LeaseJobPodStatistics) DeepCopy() *LeaseJobPodStatistics {
	if in == nil {
		return nil
	}
	out := new(LeaseJobPodStatistics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaseJobSpec) DeepCopyInto(out *LeaseJobSpec) {
	*out = *in
	if in.ActiveDeadlineSeconds != nil {
		in, out := &in.ActiveDeadlineSeconds, &out.ActiveDeadlineSeconds
		*out = new(int64)
		**out = **in
	}
	if in.ColdStartArgs != nil {
		in, out := &in.ColdStartArgs, &out.ColdStartArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ResumeArgs != nil {
		in, out := &in.ResumeArgs, &out.ResumeArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaseJobSpec.
func (in *LeaseJobSpec) DeepCopy() *LeaseJobSpec {
	if in == nil {
		return nil
	}
	out := new(LeaseJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaseJobStatus) DeepCopyInto(out *LeaseJobStatus) {
	*out = *in
	if in.LastScheduleTime != nil {
		in, out := &in.LastScheduleTime, &out.LastScheduleTime
		*out = (*in).DeepCopy()
	}
	if in.LastLeaseStartTime != nil {
		in, out := &in.LastLeaseStartTime, &out.LastLeaseStartTime
		*out = (*in).DeepCopy()
	}
	out.PodStatistics = in.PodStatistics
	if in.JobGroupCreationTimeStamp != nil {
		in, out := &in.JobGroupCreationTimeStamp, &out.JobGroupCreationTimeStamp
		*out = (*in).DeepCopy()
	}
	if in.CurrentJobScheduledTimeStamp != nil {
		in, out := &in.CurrentJobScheduledTimeStamp, &out.CurrentJobScheduledTimeStamp
		*out = (*in).DeepCopy()
	}
	if in.FormerJobDeletionTimeStamp != nil {
		in, out := &in.FormerJobDeletionTimeStamp, &out.FormerJobDeletionTimeStamp
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaseJobStatus.
func (in *LeaseJobStatus) DeepCopy() *LeaseJobStatus {
	if in == nil {
		return nil
	}
	out := new(LeaseJobStatus)
	in.DeepCopyInto(out)
	return out
}
