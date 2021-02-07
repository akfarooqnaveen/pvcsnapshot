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
	// 	"strings"

	"github.com/go-logr/logr"
	// 	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	log := r.Log.WithValues("pvcsnapshot", req.NamespacedName)

	// your logic here
	log.Info("Reconcile starting")

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
	ctx := context.Background()
	log := r.Log.WithValues("persistentvolumeclaim", req.NamespacedName)

	pvc := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, req.NamespacedName, pvc)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("pvc created")
	log.Info(pvc.Name)
	log.Info(pvc.Namespace)
	log.Info(pvc.Spec.VolumeName)
	return ctrl.Result{}, nil
}
