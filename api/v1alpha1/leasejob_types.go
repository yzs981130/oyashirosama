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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LeaseJobSpec defines the desired state of LeaseJob
type LeaseJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LeaseJob. Edit LeaseJob_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`

	// Specifies the cold start command line arguments.
	ColdStartArgs []string `json:"coldStartArgs"`

	// Specifies the desired number of pods.
	PodNumber int32 `json:"podNumber"`

	// Specifies the resume comand line arguments.
	ResumeArgs []string `json:"resumeArgs"`

	// Specifies
	Queue string `json:"queue,omitempty"`

	// Describes the pod that will be created when executing a job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Template corev1.PodTemplateSpec `json:"template"`
}

// LeaseJobStatus defines the observed state of LeaseJob
type LeaseJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	Phase LeaseJobPhase `json:"phase,omitempty"`

	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// The time of last lease started.
	// +optional
	LastLeaseStartTime *metav1.Time `json:"lastLeaseStartTime,omitempty"`

	// The pods statistics of LeaseJob.
	// +optional
	PodStatistics LeaseJobPodStatistics `json:"podStatistics,omitempty"`

	// total sub lease job it can be divided into
	// +optional
	TotalLeaseJobCnt int32 `json:"total_lease_cnt,omitempty"`

	// +optional
	CurrentLeaseJobCnt int32 `json:"current_lease_job_cnt,omitempty"`

	// +optional
	JobGroupCreationTimeStamp *metav1.Time `json:"job_group_creation_time_stamp,omitempty"`

	// +optional
	CurrentJobScheduledTimeStamp *metav1.Time `json:"current_job_scheduled_time_stamp,omitempty"`

	// +optional
	FormerJobDeletionTimeStamp *metav1.Time `json:"former_job_deletion_time_stamp,omitempty"`

	// +optional
	JobGroupName string `json:"job_group_name,omitempty"`
}

// +kubebuilder:object:root=true

// LeaseJob is the Schema for the leasejobs API
type LeaseJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LeaseJobSpec   `json:"spec,omitempty"`
	Status LeaseJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LeaseJobList contains a list of LeaseJob
type LeaseJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LeaseJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LeaseJob{}, &LeaseJobList{})
}

// The phase of lease job.
type LeaseJobPhase string

const (
	// Pending means the lease job is waiting for scheduling.
	Pending LeaseJobPhase = "Pending"

	// Renewing means the lease job is renewing its lease.
	Renewing LeaseJobPhase = "Renewing"

	// Running means the lease job is running.
	Running LeaseJobPhase = "Running"

	// Succeeded means the lease job finishes successfully.
	Succeeded LeaseJobPhase = "Succeeded"

	// Failed means the lease job finishes with a failure.
	Failed LeaseJobPhase = "Failed"

	// Unknown means the phase of lease job is unknown.
	Unknown LeaseJobPhase = "Unknown"
)

// LeaseJobPodStatistics counts the number of pods in various status.
type LeaseJobPodStatistics struct {
	// The total_executed_time
	// +optional
	TotalExecutedTime time.Duration `json:"total_executed_time,omitempty"`

	// The total_life_time
	// +optional
	TotalLifeTime time.Duration `json:"total_life_time,omitempty"`

	// The lease_expiry_cnt
	// +optional
	LeaseExpiryCnt int32 `json:"lease_expiry_cnt,omitempty"`

	// The executed_time_in_last_lease
	// +optional
	ExecutedTimeInLastLease time.Duration `json:"executed_time_in_last_lease,omitempty"`

	// The checkpoint_cnt
	// +optional
	CheckpointCnt int32 `json:"checkpoint_cnt,omitempty"`

	// The coldstart_cnt
	// +optional
	ColdstartCnt int32 `json:"coldstart_cnt,omitempty"`

	// Utilized
	// +optional
	Utilized float64 `json:"utilized,omitempty"`

	// Deserved
	// +optional
	Deserved float64 `json:"deserved,omitempty"`
}
