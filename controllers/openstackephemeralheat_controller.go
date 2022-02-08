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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8s_rand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ospdirectorv1beta1 "github.com/openstack-k8s-operators/osp-director-operator/api/v1beta1"
	common "github.com/openstack-k8s-operators/osp-director-operator/pkg/common"
	openstackephemeralheat "github.com/openstack-k8s-operators/osp-director-operator/pkg/openstackephemeralheat"
)

// OpenStackEphemeralHeatReconciler reconciles a OpenStackEphemeralHeat object
type OpenStackEphemeralHeatReconciler struct {
	client.Client
	Kclient kubernetes.Interface

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// GetClient -
func (r *OpenStackEphemeralHeatReconciler) GetClient() client.Client {
	return r.Client
}

// GetKClient -
func (r *OpenStackEphemeralHeatReconciler) GetKClient() kubernetes.Interface {
	return r.Kclient
}

// GetLogger -
func (r *OpenStackEphemeralHeatReconciler) GetLogger() logr.Logger {
	return r.Log
}

// GetScheme -
func (r *OpenStackEphemeralHeatReconciler) GetScheme() *runtime.Scheme {
	return r.Scheme
}

// +kubebuilder:rbac:groups=osp-director.openstack.org,resources=openstackephemeralheats,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=osp-director.openstack.org,resources=openstackephemeralheats/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=osp-director.openstack.org,resources=openstackephemeralheats/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=create;update;get;list;watch;patch;delete;deletecollection
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=create;update;get;list;watch;patch;deletecollection
// +kubebuilder:rbac:groups=core,resources=services,verbs=create;update;get;list;watch;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=create;update;get;list;watch;patch;deletecollection
// +kubebuilder:rbac:groups=apps,resources=replicasets,verbs=create;update;get;list;watch;patch;deletecollection

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *OpenStackEphemeralHeatReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("openstackephemeralheat", req.NamespacedName)

	// Fetch the controller VM instance
	instance := &ospdirectorv1beta1.OpenStackEphemeralHeat{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			// For additional cleanup logic use finalizers. Return and don't requeue.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	//
	// initialize condition
	//
	cond := instance.Status.Conditions.InitCondition()

	//
	// Used in comparisons below to determine whether a status update is actually needed
	//
	currentStatus := instance.Status.DeepCopy()
	statusChanged := func() bool {
		return !equality.Semantic.DeepEqual(
			r.getNormalizedStatus(&instance.Status),
			r.getNormalizedStatus(currentStatus),
		)
	}

	defer func(cond *ospdirectorv1beta1.Condition) {
		//
		// Update object conditions
		//
		instance.Status.Conditions.UpdateCurrentCondition(
			cond.Type,
			cond.Reason,
			cond.Message,
		)

		if statusChanged() {
			if updateErr := r.Client.Status().Update(context.Background(), instance); updateErr != nil {
				common.LogErrorForObject(r, updateErr, "Update status", instance)
			}
		}

		// log current status message to operator log
		common.LogForObject(r, cond.Message, instance)
	}(cond)

	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(instance, openstackephemeralheat.FinalizerName) {
			controllerutil.AddFinalizer(instance, openstackephemeralheat.FinalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
			r.Log.Info(fmt.Sprintf("Finalizer %s added to CR %s", openstackephemeralheat.FinalizerName, instance.Name))
		}
	} else {
		// 1. check if finalizer is there
		// Reconcile if finalizer got already removed
		if !controllerutil.ContainsFinalizer(instance, openstackephemeralheat.FinalizerName) {
			return ctrl.Result{}, nil
		}

		// 2. Clean up resources used by the operator
		err = r.resourceCleanup(instance)
		if err != nil {
			return ctrl.Result{}, err
		}

		// 3. as last step remove the finalizer on the operator CR to finish delete
		controllerutil.RemoveFinalizer(instance, openstackephemeralheat.FinalizerName)
		err = r.Client.Update(context.TODO(), instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		r.Log.Info(fmt.Sprintf("CR %s deleted", instance.Name))
		return ctrl.Result{}, nil
	}

	// Generate a random password secret
	passwordSecret, res, err := r.generatePasswordSecret(instance, cond)

	if (res != ctrl.Result{}) || err != nil {
		return res, err
	}

	// Generate the config maps for the various services
	err = r.generateServiceConfigMaps(instance, passwordSecret, cond)

	if err != nil {
		return ctrl.Result{}, err
	}

	// MariaDB pod and service
	res, err = r.ensureMariaDB(instance, cond)

	if (res != ctrl.Result{}) || err != nil {
		return res, err
	}

	// RabbitMQ pod and service
	res, err = r.ensureRabbitMQ(instance, cond)

	if (res != ctrl.Result{}) || err != nil {
		return res, err
	}

	// Heat API (this creates the Heat Database and runs DBsync)
	res, err = r.ensureHeat(instance, cond)

	if (res != ctrl.Result{}) || err != nil {
		return res, err
	}

	// If we get here, everything should be ready
	instance.Status.Active = true
	cond.Message = fmt.Sprintf("%s %s is available", instance.Kind, instance.Name)
	cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatReady)
	cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeProvisioned)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpenStackEphemeralHeatReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ospdirectorv1beta1.OpenStackEphemeralHeat{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.ReplicaSet{}).
		Complete(r)
}

func (r *OpenStackEphemeralHeatReconciler) getNormalizedStatus(status *ospdirectorv1beta1.OpenStackEphemeralHeatStatus) *ospdirectorv1beta1.OpenStackEphemeralHeatStatus {
	//
	// set LastHeartbeatTime and LastTransitionTime to a default value as those
	// need to be ignored to compare if conditions changed.
	//
	s := status.DeepCopy()
	for idx := range s.Conditions {
		s.Conditions[idx].LastHeartbeatTime = metav1.Time{}
		s.Conditions[idx].LastTransitionTime = metav1.Time{}
	}

	return s
}

// Generate a password secret for this OpenStackEphemeralHeat if not already present
func (r *OpenStackEphemeralHeatReconciler) generatePasswordSecret(
	instance *ospdirectorv1beta1.OpenStackEphemeralHeat,
	cond *ospdirectorv1beta1.Condition,
) (*corev1.Secret, ctrl.Result, error) {
	cmLabels := common.GetLabels(instance, openstackephemeralheat.AppLabel, map[string]string{})

	// only generate the password secret once
	passwordSecret, _, err := common.GetSecret(r, "ephemeral-heat-"+instance.Name, instance.Namespace)

	if err != nil {
		if k8s_errors.IsNotFound(err) {
			passwordSecret = openstackephemeralheat.PasswordSecret("ephemeral-heat-"+instance.Name, instance.Namespace, cmLabels, k8s_rand.String(10))
			_, op, err := common.CreateOrUpdateSecret(r, instance, passwordSecret)
			if err != nil {
				cond.Message = fmt.Sprintf("Error creating password secret for %s %s", instance.Kind, instance.Name)
				cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.CommonCondReasonSecretError)
				cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
				err = common.WrapErrorForObject(cond.Message, instance, err)

				return nil, ctrl.Result{}, err
			}

			common.LogForObject(
				r,
				fmt.Sprintf("Secret %s successfully reconciled - operation: %s", passwordSecret.Name, string(op)),
				instance,
			)

			cond.Message = "Waiting for password secret to populate"
			cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondWaitOnPassSecret)
			cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeWaiting)

			return nil, ctrl.Result{RequeueAfter: time.Second * 3}, nil
		}

		cond.Message = fmt.Sprintf("Error acquiring password secret for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.CommonCondReasonSecretError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return nil, ctrl.Result{}, err
	}

	return passwordSecret, ctrl.Result{}, nil
}

func (r *OpenStackEphemeralHeatReconciler) generateServiceConfigMaps(
	instance *ospdirectorv1beta1.OpenStackEphemeralHeat,
	passwordSecret *corev1.Secret,
	cond *ospdirectorv1beta1.Condition,
) error {
	cmLabels := common.GetLabels(instance, openstackephemeralheat.AppLabel, map[string]string{})
	envVars := make(map[string]common.EnvSetter)
	templateParameters := make(map[string]interface{})
	templateParameters["MariaDBHost"] = "mariadb-" + instance.Name
	templateParameters["RabbitMQHost"] = "rabbitmq-" + instance.Name
	templateParameters["MariaDBPassword"] = string(passwordSecret.Data["password"])

	// ConfigMaps for all services (MariaDB/Rabbit/Heat)
	cms := []common.Template{
		// ScriptsConfigMap
		{
			Name:               "openstackephemeralheat-" + instance.Name,
			Namespace:          instance.Namespace,
			Type:               common.TemplateTypeScripts,
			InstanceType:       instance.Kind,
			AdditionalTemplate: map[string]string{},
			ConfigOptions:      templateParameters,
			Labels:             cmLabels,
		},
	}

	err := common.EnsureConfigMaps(r, instance, cms, &envVars)

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating service config maps for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.CommonCondReasonConfigMapError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return err
	}

	return nil
}

func (r *OpenStackEphemeralHeatReconciler) ensureMariaDB(
	instance *ospdirectorv1beta1.OpenStackEphemeralHeat,
	cond *ospdirectorv1beta1.Condition,
) (ctrl.Result, error) {
	// MariaDB Pod
	mariadbPod := openstackephemeralheat.MariadbPod(instance)

	op, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, mariadbPod, func() error {
		err := controllerutil.SetControllerReference(instance, mariadbPod, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating MariaDB pod for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondMariaDBError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("MariaDB pod %s successfully reconciled - operation: %s", mariadbPod.Name, string(op)),
			instance,
		)
	}

	if len(mariadbPod.Status.ContainerStatuses) < 1 || !mariadbPod.Status.ContainerStatuses[0].Ready {
		cond.Message = fmt.Sprintf("Waiting on MariaDB to start for: %s", instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondWaitOnMariaDB)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeWaiting)

		return ctrl.Result{RequeueAfter: time.Second * 3}, nil
	}

	// MariaDB Service
	mariadbService := openstackephemeralheat.MariadbService(instance, r.Scheme)

	op, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, mariadbService, func() error {
		err := controllerutil.SetControllerReference(instance, mariadbService, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating MariaDB service for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondMariaDBError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("MariaDB service %s successfully reconciled - operation: %s", mariadbPod.Name, string(op)),
			instance,
		)
	}

	return ctrl.Result{}, nil
}

func (r *OpenStackEphemeralHeatReconciler) ensureRabbitMQ(
	instance *ospdirectorv1beta1.OpenStackEphemeralHeat,
	cond *ospdirectorv1beta1.Condition,
) (ctrl.Result, error) {
	// RabbitMQ Pod
	rabbitmqPod := openstackephemeralheat.RabbitmqPod(instance)

	op, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, rabbitmqPod, func() error {
		err := controllerutil.SetControllerReference(instance, rabbitmqPod, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating RabbitMQ pod for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondRabbitMQError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("RabbitMQ pod %s successfully reconciled - operation: %s", rabbitmqPod.Name, string(op)),
			instance,
		)
	}

	if len(rabbitmqPod.Status.ContainerStatuses) < 1 || !rabbitmqPod.Status.ContainerStatuses[0].Ready {
		cond.Message = fmt.Sprintf("Waiting on RabbitMQ to start for: %s", instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondWaitOnRabbitMQ)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeWaiting)

		return ctrl.Result{RequeueAfter: time.Second * 3}, nil
	}

	// RabbitMQ Service
	rabbitmqService := openstackephemeralheat.RabbitmqService(instance, r.Scheme)

	op, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, rabbitmqService, func() error {
		err := controllerutil.SetControllerReference(instance, rabbitmqService, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating RabbitMQ service for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondRabbitMQError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("RabbitMQ service %s successfully reconciled - operation: %s", rabbitmqPod.Name, string(op)),
			instance,
		)
	}

	return ctrl.Result{}, nil
}

func (r *OpenStackEphemeralHeatReconciler) ensureHeat(
	instance *ospdirectorv1beta1.OpenStackEphemeralHeat,
	cond *ospdirectorv1beta1.Condition,
) (ctrl.Result, error) {
	// Heat Pod
	heatPod := openstackephemeralheat.HeatAPIPod(instance)

	op, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, heatPod, func() error {
		err := controllerutil.SetControllerReference(instance, heatPod, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating Heat pod for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondHeatError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("Heat pod %s successfully reconciled - operation: %s", heatPod.Name, string(op)),
			instance,
		)
	}

	if len(heatPod.Status.ContainerStatuses) < 1 || !heatPod.Status.ContainerStatuses[0].Ready {
		cond.Message = fmt.Sprintf("Waiting on Heat to start for: %s", instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondWaitOnHeat)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeWaiting)

		return ctrl.Result{RequeueAfter: time.Second * 3}, nil
	}

	// Heat Service
	heatService := openstackephemeralheat.HeatAPIService(instance, r.Scheme)

	op, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, heatService, func() error {
		err := controllerutil.SetControllerReference(instance, heatService, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating Heat service for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondHeatError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("Heat service %s successfully reconciled - operation: %s", heatPod.Name, string(op)),
			instance,
		)
	}

	// Heat Engine Replicaset
	heatEngineReplicaset := openstackephemeralheat.HeatEngineReplicaSet(instance)

	op, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, heatEngineReplicaset, func() error {
		err := controllerutil.SetControllerReference(instance, heatEngineReplicaset, r.Scheme)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		cond.Message = fmt.Sprintf("Error creating/updating Heat replicaset for %s %s", instance.Kind, instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondHeatError)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeError)
		err = common.WrapErrorForObject(cond.Message, instance, err)

		return ctrl.Result{}, err
	}

	if op != controllerutil.OperationResultNone {
		common.LogForObject(
			r,
			fmt.Sprintf("Heat Engine replicaset %s successfully reconciled - operation: %s", heatEngineReplicaset.Name, string(op)),
			instance,
		)
	}

	if heatEngineReplicaset.Status.AvailableReplicas < instance.Spec.HeatEngineReplicas {
		cond.Message = fmt.Sprintf("Waiting on Heat Engine Replicas to start for: %s", instance.Name)
		cond.Reason = ospdirectorv1beta1.ConditionReason(ospdirectorv1beta1.EphemeralHeatCondWaitOnHeat)
		cond.Type = ospdirectorv1beta1.ConditionType(ospdirectorv1beta1.CommonCondTypeWaiting)

		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}

	return ctrl.Result{}, nil
}

func (r *OpenStackEphemeralHeatReconciler) resourceCleanup(instance *ospdirectorv1beta1.OpenStackEphemeralHeat) error {

	labelSelector := common.GetLabels(instance, openstackephemeralheat.AppLabel, map[string]string{})

	// delete pods
	if err := common.DeletePodsWithLabel(r, instance, labelSelector); err != nil {
		return err
	}
	// delete secret
	if err := common.DeleteSecretsWithLabel(r, instance, labelSelector); err != nil {
		return err
	}
	// delete service
	if err := common.DeleteServicesWithLabel(r, instance, labelSelector); err != nil {
		return err
	}
	// delete replicaset
	if err := common.DeleteReplicasetsWithLabel(r, instance, labelSelector); err != nil {
		return err
	}

	return nil
}
