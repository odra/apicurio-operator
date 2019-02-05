package apicurio

import (
	"github.com/gobuffalo/packr"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/template"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReconcileApiCurio reconciles a Apicurio object
type ReconcileApiCurio struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	config *rest.Config
	scheme *runtime.Scheme
	tmpl   *template.Tmpl
	box    packr.Box
}

var routeParams = map[string]string{
	"UI_ROUTE":   "apicurio-studio",
	"WS_ROUTE":   "apicurio-studio-ws",
	"API_ROUTE":  "apicurio-studio-api",
	"AUTH_ROUTE": "apicurio-studio-auth",
}
