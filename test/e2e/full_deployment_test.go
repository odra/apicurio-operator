package e2e

import (
	"context"
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	testapi "github.com/integr8ly/apicurio-operator/test/api"
	"github.com/integr8ly/apicurio-operator/test/api/meta"
	"github.com/integr8ly/apicurio-operator/test/api/watchers"
	"github.com/integr8ly/apicurio-operator/test/mock"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/stretchr/testify/assert"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestApicurioStandaloneDeployment(t *testing.T) {
	var err error

	templateLabel := "template=apicurio-studio"
	waitOpts := meta.WaitOpts{
		RetryInterval: meta.DefaultRetryInterval,
		Timeout:       meta.DefaultTimeout,
	}

	ctx := testapi.PrepareContext(t, waitOpts)
	defer ctx.Cleanup()

	err = testapi.RegisterTypes(v1alpha1.SchemeBuilder.AddToScheme, &v1alpha1.ApicurioList{})
	assert.Nil(t, err, "Failed to register crd scheme", err)

	err = testapi.RegisterOpenshiftSchemes()
	assert.Nil(t, err, "Failed to register openshift resources", err)

	ns, err := ctx.GetNamespace()
	assert.Nil(t, err, "Failed to retrieve namespace", err)

	//creation
	cr := mock.NewApicurioFull("0.2.20.Final", "example.com")
	cr.ObjectMeta.Namespace = ns

	err = framework.Global.Client.Create(context.TODO(), cr, &framework.CleanupOptions{})
	assert.Nil(t, err, "Failed to create apicurio resource", err)

	readInstance := &v1alpha1.Apicurio{
		ObjectMeta: v13.ObjectMeta{
			Name: cr.Name,
			Namespace: cr.Namespace,
		},
	}

	//watchers
	apicurioWatcher := watchers.NewAPicurioWatcher(framework.Global.Client.Client, readInstance)
	dccWatcher := watchers.NewDeploymentConfigWatcher(framework.Global.Client.Client, templateLabel, cr.Namespace)
	svcWatcher := watchers.NewServiceWatcher(framework.Global.Client.Client, templateLabel, cr.Namespace)
	routeWatcher := watchers.NewRouteWatcher(framework.Global.Client.Client, templateLabel, cr.Namespace)
	isWatcher := watchers.NewImageStreamWatcherWatcher(framework.Global.Client.Client, cr.Namespace)

	//wait for creation
	err = testapi.WaitForReadiness(apicurioWatcher, waitOpts)
	assert.Nilf(t, err, "Wait for apicurio creation", err)

	err = testapi.WaitForReadiness(dccWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for deployment config readiness", err)

	err = testapi.WaitForReadiness(svcWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for service readiness", err)

	err = testapi.WaitForReadiness(routeWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for route readiness", err)

	err = testapi.WaitForReadiness(isWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for image stream readiness", err)

	//deletion
	err = framework.Global.Client.Delete(context.TODO(), cr)
	assert.Nil(t, err, "Failed to create apicurio resource", err)

	//wait for deletion
	err = testapi.WaitForDeletion(apicurioWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for apicurio deletion", err)

	err = testapi.WaitForDeletion(dccWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for deployment config deletion", err)

	err = testapi.WaitForDeletion(svcWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for service deletion", err)

	err = testapi.WaitForDeletion(routeWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for route deletion", err)

	err = testapi.WaitForDeletion(isWatcher, waitOpts)
	assert.Nil(t, err, "Error while waiting for image stream deletion", err)
}