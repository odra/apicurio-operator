package mock

import (
	"github.com/integr8ly/apicurio-operator/pkg/apis/integreatly/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewApicurio() *v1alpha1.Apicurio {
	return &v1alpha1.Apicurio{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "integreatly.org/v1alpha1",
			Kind:       "Apicurio",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicurio",
		},
		Spec:   v1alpha1.ApicurioSpec{},
		Status: v1alpha1.ApicurioStatus{},
	}
}

func NewApicurioFull(version string, appDomain string) *v1alpha1.Apicurio {
	instance := NewApicurio()
	instance.Spec = v1alpha1.ApicurioSpec{
		Type:      v1alpha1.ApicurioFullType,
		Version:   version,
		AppDomain: appDomain,
		Auth: &v1alpha1.ApicurioAuthRresource{
			ApicurioResource: v1alpha1.ApicurioResource{
				Route: "apicurio-studio-auth",
				Memory: &v1alpha1.MinMaxCfg{
					Min: "600Mi",
					Max: "1300Mi",
				},
				CPU: &v1alpha1.MinMaxCfg{
					Min: "100m",
					Max: "1",
				},
			},
		},
		Api: &v1alpha1.ApicurioResource{
			Route: "apicurio-studio-api",
			Memory: &v1alpha1.MinMaxCfg{
				Min: "800Mi",
				Max: "4Gi",
			},
			CPU: &v1alpha1.MinMaxCfg{
				Min: "100m",
				Max: "1",
			},
			JVM: &v1alpha1.MinMaxCfg{
				Min: "768m",
				Max: "2048m",
			},
		},
		Studio: &v1alpha1.ApicurioResource{
			Route: "apicurio-studio-ui",
			Memory: &v1alpha1.MinMaxCfg{
				Min: "600Mi",
				Max: "1300Mi",
			},
			CPU: &v1alpha1.MinMaxCfg{
				Min: "100m",
				Max: "1",
			},
		},
		WebSocket: &v1alpha1.ApicurioResource{
			Route: "apicurio-studio-ws",
			Memory: &v1alpha1.MinMaxCfg{
				Min: "600Mi",
				Max: "1300Mi",
			},
			CPU: &v1alpha1.MinMaxCfg{
				Min: "100m",
				Max: "1",
			},
		},
	}

	return instance
}

func SetExternalAuth(instance *v1alpha1.Apicurio, host string, username string, password string, realm string) {
	instance.Spec.Type = v1alpha1.ApicurioExternalAuthType
	instance.Spec.Auth = &v1alpha1.ApicurioAuthRresource{
		Host: host,
		Username: username,
		Password: password,
		Realm: realm,
	}
}
