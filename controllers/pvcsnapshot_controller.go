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

package controllers

import (
	"context"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	frqv1alpha1 "pvcsnapshotoperator/api/v1alpha1"
)

// PVCSnapshotReconciler reconciles a PVCSnapshot object
type PVCSnapshotReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type PersistentVolumeClaimReconciler struct {
	client.Client
	APIReader client.Reader
	Log       logr.Logger
	Scheme    *runtime.Scheme
}

// +kubebuilder:rbac:groups=frq.akfarooqnaveen.com,resources=pvcsnapshots,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=frq.akfarooqnaveen.com,resources=pvcsnapshots/status,verbs=get;update;patch

func (r *PVCSnapshotReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// 	ctx := context.Background()
	lg := r.Log.WithValues("pvcsnapshot", req.NamespacedName)

	// your logic here
	lg.Info("Reconcile starting")

	// 	var pvc frqv1alpha1.PVCSnapshot

	return ctrl.Result{}, nil
}

func (r *PVCSnapshotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frqv1alpha1.PVCSnapshot{}).
		Complete(r)
}

func (r *PersistentVolumeClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=,resources=persistentvolumeclaim,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=persistentvolumeclaim/status,verbs=get;list;watch;create;update;patch;delete
func (r *PersistentVolumeClaimReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	var VolumeTagRead string = "kubernetes-dynamic-"
	ctx := context.Background()
	lg := r.Log.WithValues("persistentvolumeclaim", req.NamespacedName)

	pvc := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, req.NamespacedName, pvc)
	if err != nil {
		return ctrl.Result{}, err
	}
	lg.Info("pvc created")
	lg.Info(pvc.Name)
	lg.Info(pvc.Namespace)
	lg.Info(pvc.Spec.VolumeName)
	VolumeTagRead += pvc.Spec.VolumeName

	session := session.New(&aws.Config{Region: aws.String("us-east-2")})

	client := dlm.New(session)

	var time string = "09:00"
	times := []*string{&time}
	var cr dlm.CreateRule
	cr.SetInterval(12)
	cr.SetIntervalUnit("HOURS")
	cr.SetLocation("CLOUD")
	cr.SetTimes(times)

	var rr dlm.RetainRule
	rr.SetCount(1)

	schedules := []*dlm.Schedule{}
	var s dlm.Schedule
	s.SetCreateRule(&cr)
	s.SetName("Schedule1")
	s.SetRetainRule(&rr)

	schedules = append(schedules, &s)

	var ResourceType string = "VOLUME"
	var ResourceLocation string = "CLOUD"
	resourceTypes := []*string{&ResourceType}
	resourceLocations := []*string{&ResourceLocation}

	var Key string = "Name"
	var Value string = VolumeTagRead
	var targetTag dlm.Tag = dlm.Tag{
		Key:   &Key,
		Value: &Value,
	}
	targetTags := []*dlm.Tag{}
	targetTags = append(targetTags, &targetTag)

	var p dlm.PolicyDetails
	p.SetPolicyType("EBS_SNAPSHOT_MANAGEMENT")
	p.SetResourceLocations(resourceLocations)
	p.SetResourceTypes(resourceTypes)
	p.SetSchedules(schedules)
	p.SetTargetTags(targetTags)

	var lpi dlm.CreateLifecyclePolicyInput
	lpi.SetDescription(VolumeTagRead)
	lpi.SetExecutionRoleArn("arn:aws:iam::731556103348:role/service-role/AWSDataLifecycleManagerDefaultRole")
	lpi.SetPolicyDetails(&p)
	lpi.SetState("ENABLED")

	reqaws, output := client.CreateLifecyclePolicyRequest(&lpi)

	log.Println("Before send")

	erraws := reqaws.Send()
	if erraws == nil {
		log.Println(output)
	} else {
		log.Fatal(erraws)
	}
	return ctrl.Result{}, nil
}
