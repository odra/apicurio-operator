package watchers

import (
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/integr8ly/apicurio-operator/test/api/meta"
	"github.com/integr8ly/apicurio-operator/test/mock"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type apicurioWatcher struct {
	Instance *v1alpha1.Apicurio
	client   client.Client
}

func NewAPicurioWatcher(client client.Client) *apicurioWatcher {
	return &apicurioWatcher{
		Instance: mock.NewApicurio(),
		client:   client,
	}
}

func (aw *apicurioWatcher) CompareStatus(status *v1alpha1.ApicurioStatus) bool {
	crStatus := aw.Instance.Status
	reasonMacthes := *crStatus.Reason == *status.Reason
	typeMacthes := crStatus.Type == status.Type

	return typeMacthes && reasonMacthes
}

func (aw *apicurioWatcher) Observe(opts meta.WaitOpts, loader meta.ObjectLoader) error {
	return wait.Poll(opts.RetryInterval, opts.Timeout, func() (done bool, err error) {
		obj, err := loader()
		if err != nil {
			return true, err
		}

		if obj == nil {
			return false, nil
		}

		innerInstance := obj.(*v1alpha1.Apicurio)

		return aw.CompareStatus(&innerInstance.Status), nil
	})
}
