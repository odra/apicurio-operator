package watchers

import (
	"context"
	"errors"
	v13 "github.com/openshift/api/image/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type imageStreamWatcher struct {
	Namespace string
	client   client.Client
}

func NewImageStreamWatcherWatcher(client client.Client, namespace string) *imageStreamWatcher {
	return &imageStreamWatcher{
		Namespace: namespace,
		client: client,
	}
}

func (isw *imageStreamWatcher) WaitForReadiness(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		isList := &v13.ImageStreamList{}
		listOpts := &client.ListOptions{
			Namespace: isw.Namespace,
		}

		innerError = isw.client.List(context.TODO(), listOpts, isList)
		if innerError != nil {
			return true, innerError
		}

		if len(isList.Items) == 0 {
			return false, errors.New("failed to query routes")
		}

		return true, nil
	})
}

func (isw *imageStreamWatcher) WaitForDeletion(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		var innerError error

		isList := &v13.ImageStreamList{}
		listOpts := &client.ListOptions{
			Namespace: isw.Namespace,
		}

		innerError = isw.client.List(context.TODO(), listOpts, isList)
		if innerError != nil {
			return true, innerError
		}

		if len(isList.Items) == 0 {
			return true, nil
		}

		return false, nil
	})
}