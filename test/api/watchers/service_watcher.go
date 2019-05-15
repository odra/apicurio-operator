package watchers

import (
	"context"
	"errors"
	v12 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type serviceWatcher struct {
	Namespace string
	RawLabelSelector string
	client   client.Client
}

func NewServiceWatcher(client client.Client, labelSelector string, namespace string) *serviceWatcher {
	return &serviceWatcher{
		Namespace: namespace,
		RawLabelSelector: labelSelector,
		client: client,
	}
}

func (sw *serviceWatcher) WaitForReadiness(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		svcList := &v12.ServiceList{}
		listOpts := &client.ListOptions{
			Namespace: sw.Namespace,
		}

		innerError = listOpts.SetLabelSelector(sw.RawLabelSelector)
		if innerError != nil {
			return true, err
		}

		innerError = sw.client.List(context.TODO(), listOpts, svcList)
		if innerError != nil {
			return true, innerError
		}

		if len(svcList.Items) == 0 {
			return false, errors.New("failed to query services")
		}

		endpoints :=  v12.Endpoints{}
		key := types.NamespacedName{
			Namespace: sw.Namespace,
		}

		for _, svc := range svcList.Items {
			key.Name = svc.Name
			innerError = sw.client.Get(context.TODO(), key, &endpoints)
			if innerError != nil {
				if errors2.IsNotFound(innerError) {
					return false, nil
				}
				return true, innerError
			}
		}

		return true, nil
	})
}

func (sw *serviceWatcher) WaitForDeletion(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		svcList := &v12.ServiceList{}
		listOpts := &client.ListOptions{
			Namespace: sw.Namespace,
		}

		innerError = listOpts.SetLabelSelector(sw.RawLabelSelector)
		if innerError != nil {
			return true, err
		}

		innerError = sw.client.List(context.TODO(), listOpts, svcList)
		if innerError != nil {
			return true, innerError
		}

		if len(svcList.Items) == 0 {
			return true, nil
		}

		return false, nil
	})
}