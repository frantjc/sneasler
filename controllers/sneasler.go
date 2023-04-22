/*
Copyright 2023.

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
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	kerrors "k8s.io/cri-api/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/frantjc/go-fn"
	"github.com/frantjc/sneasler"
	frantjcv1alpha1 "github.com/frantjc/sneasler/api/v1alpha1"
)

// SneaslerReconciler reconciles a Sneasler object.
type SneaslerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=frantj.cc,resources=sneaslers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=frantj.cc,resources=sneaslers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=frantj.cc,resources=sneaslers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *SneaslerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		_     = sneasler.LoggerFrom(ctx)
		snslr = &frantjcv1alpha1.Sneasler{}
	)

	if err := r.Client.Get(ctx, req.NamespacedName, snslr); err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	var (
		metadata = metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: frantjcv1alpha1.GroupVersion.String(),
					Name:       req.Name,
					Controller: fn.Ptr(true),
					Kind:       "Sneasler",
				},
			},
		}
		selectorLabels = map[string]string{
			"app.kubernetes.io/name": req.Name,
		}
		port       int32 = 8080
		deployment       = &appsv1.Deployment{
			ObjectMeta: metadata,
			Spec: appsv1.DeploymentSpec{
				Replicas: fn.Ptr[int32](1),
				Selector: &metav1.LabelSelector{
					MatchLabels: selectorLabels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: selectorLabels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "sneasler",
								Image: sneasler.GetImageRef(),
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: port,
									},
								},
							},
						},
					},
				},
			},
		}
		service = &corev1.Service{
			ObjectMeta: metadata,
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Port:       port,
						TargetPort: intstr.FromInt(int(port)),
					},
				},
				Selector: selectorLabels,
			},
		}
	)

	return ctrl.Result{}, errors.Join(r.Client.Create(ctx, deployment), r.Client.Create(ctx, service))
}

// SetupWithManager sets up the controller with the Manager.
func (r *SneaslerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frantjcv1alpha1.Sneasler{}).
		Complete(r)
}
