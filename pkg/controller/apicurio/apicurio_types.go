package apicurio

import (
	"context"
	"github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type watcher struct {
	client           clientSpec       `json:"-"`
	ResourceCheckers []*statusChecker `json:"resource_checkers"`
}

type statusChecker struct {
	Name    string `json:"name"`
	IsReady bool `json:"is_ready"`
	Info *statusResource `json:"info"`
	checker statusCheckerSpec `json:"-"`
}

type statusResource struct {
	Memory map[string]string `json:"memory"`
	CPU map[string]string `json:"cpu"`
	JVM map[string]string `json:"jvm"`
}

type clientSpec interface {
	Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error
}

type statusCheckerSpec interface {
	IsReady() bool
	Reload(dc *v1.DeploymentConfig)
}