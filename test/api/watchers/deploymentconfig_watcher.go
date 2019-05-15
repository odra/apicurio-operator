package watchers

import (
	"context"
	"github.com/openshift/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type deploymentConfigWatcher struct {
	Namespace string
	RawLabelSelector string
	client   client.Client
}

func NewDeploymentConfigWatcher(client client.Client, labelSelector string, namespace string) *deploymentConfigWatcher {
	return &deploymentConfigWatcher{
		Namespace: namespace,
		RawLabelSelector: labelSelector,
		client: client,
	}
}

func (dcw *deploymentConfigWatcher) WaitForReadiness(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		dcList := &v1.DeploymentConfigList{}
		listOpts := &client.ListOptions{
			Namespace: dcw.Namespace,
		}

		innerError = listOpts.SetLabelSelector(dcw.RawLabelSelector)
		if innerError != nil {
			return true, err
		}

		innerError = dcw.client.List(context.TODO(), listOpts, dcList)
		if innerError != nil {
			return true, innerError
		}

		if len(dcList.Items) == 0 {
			return false, nil
		}

		for _, dc := range dcList.Items {
			if dc.Status.Replicas != dc.Status.ReadyReplicas {
				return false, nil
			}

			latestCondition := dc.Status.Conditions[0]
			if latestCondition.Type != v1.DeploymentAvailable && latestCondition.Status != v12.ConditionTrue {
				return false, nil
			}
		}

		return true, nil
	})
}

func (dcw *deploymentConfigWatcher) WaitForDeletion(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		dcList := &v1.DeploymentConfigList{}
		listOpts := &client.ListOptions{
			Namespace: dcw.Namespace,
		}

		innerError = listOpts.SetLabelSelector(dcw.RawLabelSelector)
		if innerError != nil {
			return true, err
		}

		innerError = dcw.client.List(context.TODO(), listOpts, dcList)
		if innerError != nil {
			return true, innerError
		}

		if len(dcList.Items) == 0 {
			return true, nil
		}

		return false, nil
	})
}