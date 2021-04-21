package controllers

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	batchv1alpha1 "pkg.yezhisheng.me/oyashirosama/api/v1alpha1"
)

type enqueueFn func(job *batchv1alpha1.LeaseJob)

// RenewalProcessor polls all running leasejobs and enqueues them after their lease expires.
type RenewalProcessor struct {
	pool map[string]*batchv1alpha1.LeaseJob
	// the lock of pool
	lock sync.Mutex

	leaseLength time.Duration

	enqueue enqueueFn
}

func NewRenewalProcessor(length time.Duration, enqueue enqueueFn) *RenewalProcessor {
	return &RenewalProcessor{
		pool:        map[string]*batchv1alpha1.LeaseJob{},
		lock:        sync.Mutex{},
		leaseLength: length,
		enqueue:     enqueue,
	}
}

func (p *RenewalProcessor) Run(stop <-chan struct{}) {
	go wait.Until(p.worker, time.Second, stop)
}

func (p *RenewalProcessor) Add(job *batchv1alpha1.LeaseJob) {
	// TODO: impl
	p.lock.Lock()
	defer p.lock.Unlock()
	jobKey, _ := p.getJobKey(job)
	now := metav1.Time{Time: time.Now()}
	if v, found := p.pool[jobKey]; !found {
		p.pool[jobKey] = job
		job.Status.LastLeaseStartTime = &now
	} else {
		// never know if happens...
		v.Status.LastLeaseStartTime = &now
	}
}

func (p *RenewalProcessor) Delete(job *batchv1alpha1.LeaseJob) {
	// TODO: impl
	p.lock.Lock()
	defer p.lock.Unlock()
	jobKey, _ := p.getJobKey(job)
	if v, found := p.pool[jobKey]; found {
		v.Status.LastLeaseStartTime = nil
		delete(p.pool, jobKey)
	}
}

func (p *RenewalProcessor) worker() {
	for {

		jobs := p.getOvertimeLeaseJobs()

		for _, job := range jobs {
			p.enqueue(job)
		}
	}
}

func (p *RenewalProcessor) getOvertimeLeaseJobs() []*batchv1alpha1.LeaseJob {
	// TODO: impl
	p.lock.Lock()
	defer p.lock.Unlock()

	var ret []*batchv1alpha1.LeaseJob
	expiredTime := time.Now().Add(-p.leaseLength)

	for key, job := range p.pool {
		if job.Status.LastLeaseStartTime.Time.Before(expiredTime) {
			job.Status.Phase = batchv1alpha1.Renewing
			delete(p.pool, key)
			ret = append(ret, job)
		}
	}
	return ret
}

func (p *RenewalProcessor) getJobKey(job *batchv1alpha1.LeaseJob) (string, error) {
	return cache.MetaNamespaceKeyFunc(job)
}

func (p *RenewalProcessor) resetJobTimer(job *batchv1alpha1.LeaseJob) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	jobKey, _ := p.getJobKey(job)
	if _, found := p.pool[jobKey]; found {
		job.Status.LastLeaseStartTime = &metav1.Time{Time: time.Now()}
	} else {
		return fmt.Errorf("%v not found in jobPool", job.Name)
	}
	return nil
}
