package watchers

import (
	"context"
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type apicurioWatcher struct {
	Instance *v1alpha1.Apicurio
	client   client.Client
}

func NewAPicurioWatcher(client client.Client, instance *v1alpha1.Apicurio) *apicurioWatcher {
	return &apicurioWatcher{
		Instance: instance,
		client:   client,
	}
}

func (aw *apicurioWatcher) WaitForReadiness(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		readInstance := aw.Instance.DeepCopy()
		key := types.NamespacedName{
			Name: readInstance.Name,
			Namespace: readInstance.Namespace,
		}

		innerErr := aw.client.Get(context.TODO(), key, readInstance)
		if innerErr != nil {
			return true, err
		}

		reasonMacthes := *readInstance.Status.Reason == "ApicurioReady"
		typeMacthes := readInstance.Status.Type == v1alpha1.ApicurioReady
		if typeMacthes || reasonMacthes {
			return true, nil
		}

		return false, nil
	})
}

func (aw *apicurioWatcher) WaitForDeletion(interval time.Duration, timeout time.Duration) error {
	return wait.Poll(interval, timeout, func() (done bool, err error) {
		deletedInstance := aw.Instance.DeepCopy()
		key := types.NamespacedName{
			Name: deletedInstance.Name,
			Namespace: deletedInstance.Namespace,
		}

		innerErr := aw.client.Get(context.TODO(), key, deletedInstance)
		if innerErr != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return true, err
		}

		return false, nil
	})
}
