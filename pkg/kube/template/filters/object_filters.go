package filters

import (
	"fmt"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/kubernetes"
	"k8s.io/apimachinery/pkg/runtime"
)

func FilterDcs(obj *runtime.Object) error {
	uo, err := kubernetes.UnstructuredFromRuntimeObject(*obj)
	if err != nil {
		return err
	}

	kind := uo.GetKind()
	if uo.GetKind() != "DeploymentConfig" {
		return fmt.Errorf("skipping obejct kind: %s", kind)
	}

	return nil
}

func FilterSkipAuthDc(obj *runtime.Object) error {
	uo, err := kubernetes.UnstructuredFromRuntimeObject(*obj)
	if err != nil {
		return err
	}

	kind := uo.GetKind()
	if uo.GetKind() != "DeploymentConfig" {
		return fmt.Errorf("skipping obejct kind: %s", kind)
	}

	labels := uo.GetLabels()
	if val, ok := labels["app"]; ok && val == "apicurio-studio-auth" {
		return fmt.Errorf("skipping auth object: %s:%s", uo.GetKind(), uo.GetName())
	}

	return nil
}

func FilterSkipAuthObjs(obj *runtime.Object) error {
	uo, err := kubernetes.UnstructuredFromRuntimeObject(*obj)
	if err != nil {
		return err
	}

	labels := uo.GetLabels()
	if val, ok := labels["app"]; ok && val == "apicurio-studio-auth" {
		return fmt.Errorf("skipping auth object: %s:%s", uo.GetKind(), uo.GetName())
	}

	return nil
}
