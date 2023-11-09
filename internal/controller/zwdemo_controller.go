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

package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	v1 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v1"

	zwv1 "operator4/api/v1"
)

// ZwdemoReconciler reconciles a Zwdemo object
type ZwdemoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=zw.zw.com,resources=zwdemoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=zw.zw.com,resources=zwdemoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=zw.zw.com,resources=zwdemoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Zwdemo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *ZwdemoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// TODO(user): your logic here
	//用于在日志记录器中添加键值对的上下文信息。在这个例子中，"myapp" 是键，req.NamespacedName 是值。req.NamespacedName的值是一个 "命名空间/资源对象名称"的组合
	logger := r.Log.WithValues("myapp", req.NamespacedName)
	fmt.Println("==========enter================")
	logger.Info("Reconciling myapp")
	//判断MyApp对象是否存在
	instance := &zwv1.Zwdemo{}
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logger.Info("app kind: " + instance.Kind + ", app name: " + instance.Name)
	if instance.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	//
	deployment := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deployment); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		// 1. 不存在，则创建
		// 1-1. 创建 Deployment
		deployment = NewDeployment(instance)
		if err := r.Client.Create(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}

		// 1-2. 创建 Service
		svc := NewService(instance)
		if err := r.Client.Create(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func NewDeployment(app *zwv1.Zwdemo) *appsv1.Deployment {
	labels := map[string]string{"app": app.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    app.Kind,
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            app.Name,
							Image:           app.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: app.Spec.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
}

func NewService(app *zwv1.Zwdemo) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    app.Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       app.Spec.ServicePort,
					TargetPort: intstr.FromInt(int(app.Spec.ContainerPort)),
				},
			},
			Selector: map[string]string{
				"app": app.Name,
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ZwdemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zwv1.Zwdemo{}).
		Complete(r)
}
