/*

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

package testjob

import (
	"context"
	"reflect"

	schedulerv1alpha1 "oyashirosama/pkg/apis/scheduler/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")
var cnt_size = 0

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new TestJob Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTestJob{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("testjob-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to TestJob
	err = c.Watch(&source.Kind{Type: &schedulerv1alpha1.TestJob{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by TestJob - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &schedulerv1alpha1.TestJob{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileTestJob{}

// ReconcileTestJob reconciles a TestJob object
type ReconcileTestJob struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a TestJob object and makes changes based on the state read
// and what is in the TestJob.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=scheduler.openi.org,resources=testjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scheduler.openi.org,resources=testjobs/status,verbs=get;update;patch
func (r *ReconcileTestJob) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the TestJob instance
	instance := &schedulerv1alpha1.TestJob{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Deployment object
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-deployment",
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment": instance.Name + "-deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"deployment": instance.Name + "-deployment"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
					NodeSelector: map[string]string{"disktype": "ssd"},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)

		instance.Status.Running = false
		instance.Status.NeedScheduling = true
		err = r.Update(context.TODO(), instance)

		//println("running: ", instance.Status.Running)
		err = r.Create(context.TODO(), deploy)
		return reconcile.Result{}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		found.Spec = deploy.Spec
		log.Info("Updating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}
		//println("!")
		//println(instance.Status.NeedScheduling)
		instance.Status.NeedScheduling = false
		cnt_size++
		//instance.Labels["whoami"] = "yzs"
		for t := range instance.Labels {
			println(t)
		}
		err = r.Update(context.TODO(), instance)
		//err = r.Status().Update(context.Background(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		println(cnt_size)

		println("!")
		if cnt_size > 5 {
			cnt_size = 0
		}
	}
	return reconcile.Result{}, nil
}
