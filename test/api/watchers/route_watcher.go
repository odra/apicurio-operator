package watchers

import (
	"context"
	"errors"
	"github.com/openshift/api/route/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type routeWatcher struct {
	Namespace string
	RawLabelSelector string
	client   client.Client
}

func NewRouteWatcher(client client.Client, labelSelector string, namespace string) *routeWatcher {
	return &routeWatcher{
		Namespace: namespace,
		RawLabelSelector: labelSelector,
		client: client,
	}
}

func (rw *routeWatcher) WaitForReadiness(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		routeList := &v1.RouteList{}
		listOpts := &client.ListOptions{
			Namespace: rw.Namespace,
		}

		innerError = listOpts.SetLabelSelector(rw.RawLabelSelector)
		if innerError != nil {
			return true, err
		}

		innerError = rw.client.List(context.TODO(), listOpts, routeList)
		if innerError != nil {
			return true, innerError
		}

		if len(routeList.Items) == 0 {
			return false, errors.New("failed to query routes")
		}

		for _, route := range routeList.Items {
			condition := route.Status.Ingress[0].Conditions[0]
			if condition.Type != v1.RouteAdmitted && condition.Status != v12.ConditionTrue {
				return false, nil
			}
		}

		return true, nil
	})
}

func (rw *routeWatcher) WaitForDeletion(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		routeList := &v1.RouteList{}
		listOpts := &client.ListOptions{
			Namespace: rw.Namespace,
		}

		innerError = listOpts.SetLabelSelector(rw.RawLabelSelector)
		if innerError != nil {
			return true, err
		}

		innerError = rw.client.List(context.TODO(), listOpts, routeList)
		if innerError != nil {
			return true, innerError
		}

		if len(routeList.Items) == 0 {
			return true, nil
		}

		return false, nil
	})
}