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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1alpha1 "pkg.yezhisheng.me/oyashirosama/api/v1alpha1"
	vcbatch "pkg.yezhisheng.me/volcano/pkg/apis/batch/v1alpha1"
)

// LeaseJobReconciler reconciles a LeaseJob object
type LeaseJobReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	jobOwnerKey             = ".metadata.controller"
	apiGVStr                = batchv1alpha1.GroupVersion.String()
	leaseTerm               = time.Minute
	scheduledTimeAnnotation = "batch.pkg.yezhisheng.me/scheduled-at"
)

// +kubebuilder:rbac:groups=batch.pkg.yezhisheng.me,resources=leasejobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.pkg.yezhisheng.me,resources=leasejobs/status,verbs=get;update;patch

func (r *LeaseJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("leasejob", req.NamespacedName)

	// your logic here
	var leaseJob batchv1alpha1.LeaseJob
	if err := r.Get(ctx, req.NamespacedName, &leaseJob); err != nil {
		log.Error(err, "unable to fetch LeaseJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// list vc jobs match jobOwnerKey
	var childJobs vcbatch.JobList
	if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child vcJobs")
		return ctrl.Result{}, err
	}

	var activeJobs []*vcbatch.Job
	var successfulJobs []*vcbatch.Job
	var failedJobs []*vcbatch.Job
	var mostRecentTime *time.Time

	for i, job := range childJobs.Items {
		if isFinished, phase := isVcJobFinished(&job); isFinished {
			if phase == vcbatch.Failed || phase == vcbatch.Terminated {
				failedJobs = append(failedJobs, &childJobs.Items[i])
			} else {
				successfulJobs = append(successfulJobs, &childJobs.Items[i])
			}
		} else {
			activeJobs = append(activeJobs, &childJobs.Items[i])
		}

		// update mostRecentTime
		scheduledTimeForJob, err := getScheduledTimeForJob(&job)
		if err != nil {
			log.Error(err, "unable to parse schedule time for child job", "job", &job)
			continue
		}
		if scheduledTimeForJob != nil {
			if mostRecentTime == nil {
				mostRecentTime = scheduledTimeForJob
			} else if mostRecentTime.Before(*scheduledTimeForJob) {
				mostRecentTime = scheduledTimeForJob
			}
		}
	}

	// store mostRecentTime in leaseJob
	if mostRecentTime != nil {
		leaseJob.Status.LastScheduleTime = &metav1.Time{Time: *mostRecentTime}
	} else {
		leaseJob.Status.LastScheduleTime = nil
	}

	log.V(1).Info("job count", "active jobs", len(activeJobs), "successful jobs", len(successfulJobs), "failed jobs", len(failedJobs))

	// update leaseJob
	if err := r.Update(ctx, &leaseJob); err != nil {
		log.Error(err, "unable to update leaseJob")
		return ctrl.Result{}, err
	}

	if len(activeJobs) == 1 {
		// check if lease expiry
		for _, job := range activeJobs {
			if isLeaseExpiry(job) {
				newVcJob := r.constructVcJobForLeaseJob(&leaseJob, leaseJob.Spec.ResumeArgs)
				if err := r.Create(ctx, newVcJob); err != nil && !errors.IsAlreadyExists(err) {
					log.Error(err, "unable to create next VC job")
					return ctrl.Result{}, err
				}
				if err := r.Delete(ctx, job, client.GracePeriodSeconds(30)); client.IgnoreNotFound(err) != nil {
					log.Error(err, "unable to delete current lease expiring job")
					return ctrl.Result{}, err
				}
			}
		}
	}

	if len(activeJobs) == 0 {
		newVcJob := r.constructVcJobForLeaseJob(&leaseJob, leaseJob.Spec.ColdStartArgs)
		if err := r.Create(ctx, newVcJob); err != nil && !errors.IsAlreadyExists(err) {
			log.Error(err, "unable to create next VC job")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *LeaseJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1alpha1.LeaseJob{}).
		Complete(r)
}

func isVcJobFinished(job *vcbatch.Job) (bool, vcbatch.JobPhase) {
	if job.Status.State.Phase == vcbatch.Completed ||
		job.Status.State.Phase == vcbatch.Failed ||
		job.Status.State.Phase == vcbatch.Terminated {
		return true, job.Status.State.Phase
	}
	return false, ""
}

func getScheduledTimeForJob(job *vcbatch.Job) (*time.Time, error) {
	return getTimeStamp(job, scheduledTimeAnnotation)
}

func isLeaseExpiry(job *vcbatch.Job) bool {
	startTime, err := getScheduledTimeForJob(job)
	if err != nil {
		return false
	}
	return startTime.Add(leaseTerm).Before(time.Now())
}

func (r *LeaseJobReconciler) constructVcJobForLeaseJob(leaseJob *batchv1alpha1.LeaseJob, startCmd []string) *vcbatch.Job {
	job := &vcbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        "",
			Namespace:   leaseJob.Namespace,
		},
		Spec: vcbatch.JobSpec{
			SchedulerName:           "volcano",
			Tasks:                   make([]vcbatch.TaskSpec, leaseJob.Spec.PodNumber),
			Queue:                   leaseJob.Spec.Queue,
			TTLSecondsAfterFinished: func(a int32) *int32 { return &a }(int32(1)),
		},
	}
	for i := range job.Spec.Tasks {
		job.Spec.Tasks[i] = vcbatch.TaskSpec{
			Name:     "",
			Replicas: 1,
			Template: *leaseJob.Spec.Template.DeepCopy(),
		}
		job.Spec.Tasks[i].Template.Spec.Containers[0].Command = startCmd
	}

	for k, v := range leaseJob.Spec.Template.Annotations {
		job.Annotations[k] = v
	}
	for k, v := range leaseJob.Spec.Template.Labels {
		job.Labels[k] = v
	}

	setTimeStamp(job, scheduledTimeAnnotation, time.Now())
	_ = ctrl.SetControllerReference(leaseJob, job, r.Scheme)

	return job
}

func setTimeStamp(job *vcbatch.Job, anno string, t time.Time) {
	job.Annotations[anno] = t.Format(time.RFC3339)
}

func getTimeStamp(job *vcbatch.Job, anno string) (*time.Time, error) {
	timeRaw := job.Annotations[anno]
	if len(timeRaw) == 0 {
		return nil, nil
	}

	timeParsed, err := time.Parse(time.RFC3339, timeRaw)
	if err != nil {
		return nil, err
	}
	return &timeParsed, nil
}
