package apicurio

import (
	"fmt"
	integreatlyv1alpha1 "github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	"k8s.io/api/core/v1"
	"strings"
)

func getTemplatePath(cr *integreatlyv1alpha1.Apicurio) string {
	return fmt.Sprintf("%s/template.yml", cr.Spec.Version)
}

func mapResource(params map[string]string, prefix string, resource *integreatlyv1alpha1.ApicurioResource) {
	if resource.Route != "" {
		params[prefixKey(prefix, "ROUTE")] = resource.Route
	}

	if resource.CPU != nil {
		params[prefixKey(prefix, "CPU_REQUEST")] = resource.CPU.Min
		params[prefixKey(prefix, "CPU_LIMIT")] = resource.CPU.Max
	}

	if resource.JVM != nil {
		params[prefixKey(prefix, "JVM_MIN")] = resource.JVM.Min
		params[prefixKey(prefix, "JVM_MAX")] = resource.JVM.Max
	}

	if resource.Memory != nil {
		params[prefixKey(prefix, "MEM_REQUEST")] = resource.Memory.Min
		params[prefixKey(prefix, "MEM_LIMIT")] = resource.Memory.Max
	}
}

func prefixKey(prefix string, key string) string {
	if len(prefix) == 0 {
		return key
	}

	return fmt.Sprintf("%s_%s", prefix, key)
}

func updateRoute(cr *integreatlyv1alpha1.Apicurio, resource *integreatlyv1alpha1.ApicurioResource) {
	if cr.Spec.AppDomain == "" {
		return
	}

	resource.Route = fmt.Sprintf("%s.%s", resource.Route, cr.Spec.AppDomain)
}

func fixAuthUrl(url *string) {
	prefixes := []string{"http://", "https://"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(*url, prefix) {
			*url = strings.Replace(*url, prefix, "", 1)
		}
	}

	if strings.HasSuffix(*url, "/") {
		*url = strings.TrimSuffix(*url, "/")
	}
}

func (r *ReconcileApiCurio) isReady(cr *integreatlyv1alpha1.Apicurio) bool {
	return apicurioWatcher.isReady()
}

func getEnvValue(env []v1.EnvVar, name string) string {
	for _, v := range env {
		if v.Name == name {
			return v.Value
		}
	}

	return ""
}
