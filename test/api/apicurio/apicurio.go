package apicurio

import (
	"context"
	"fmt"
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	dynclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetPods(client framework.FrameworkClient) (*v1.PodList, error) {
	opts := &dynclient.ListOptions{}
	err := opts.SetLabelSelector("template=apicurio-studio")
	if err != nil {
		return nil, err
	}
	podList := &v1.PodList{}
	err = client.List(context.TODO(), opts, podList)
	if err != nil {
		return nil, err
	}

	return podList, nil
}

func WaitForDeploymentReadiness(instance *v1alpha1.Apicurio, key types.NamespacedName, client framework.FrameworkClient) func() (runtime.Object, error) {
	return func() (runtime.Object, error) {
		err := client.Get(context.TODO(), key, instance)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil, err
			}

			return nil, nil
		}

		if instance.Status.Type != v1alpha1.ApicurioReady {
			return nil, nil
		}

		podList, err := GetPods(client)
		if err != nil {
			return nil, err
		}

		expectedTotalPods := 4
		if instance.Spec.Type == v1alpha1.ApicurioExternalAuthType {
			expectedTotalPods = 3
		}


		totalPods := len(podList.Items)
		if totalPods < expectedTotalPods {
			return nil, fmt.Errorf("not enough labeled pods: %d", totalPods)
		}

		return instance, nil
	}
}
