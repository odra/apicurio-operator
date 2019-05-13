package checkers

import (
	"github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type deploymentConfigCheck struct {
	ref *v1.DeploymentConfig
}

func NewDeploymentConfigCheck(dc *v1.DeploymentConfig) *deploymentConfigCheck {
	return &deploymentConfigCheck{
		ref: dc,
	}
}

func (dcc *deploymentConfigCheck) Reload(object runtime.Object)  error {
	dcc.ref = object.(*v1.DeploymentConfig)

	return nil
}

func (dcc *deploymentConfigCheck) IsReady() bool {
	ref := dcc.ref
	if ref == nil {
		return false
	}

	conditions := ref.Status.Conditions
	if len(conditions) == 0 {
		return false
	}

	if !dcc.isConditionReady() {
		return false
	}

	if !dcc.isReplicaReady() {
		return false
	}

	if dcc.hasDeletionTS() {
		return false
	}

	return true
}

func (dcc *deploymentConfigCheck) hasDeletionTS() bool {
	return dcc.ref.GetDeletionTimestamp() != nil
}

func (dcc *deploymentConfigCheck) isConditionReady() bool {
	c := &dcc.ref.Status.Conditions[0]
	return c.Type == v1.DeploymentAvailable && c.Status == corev1.ConditionTrue
}

func (dcc *deploymentConfigCheck) isReplicaReady() bool {
	status := dcc.ref.Status
	replicas := status.Replicas
	availableReplicas := status.AvailableReplicas
	readyReplicas := status.ReadyReplicas

	return replicas == availableReplicas && replicas == readyReplicas
}
