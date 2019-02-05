package mock

import (
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewKeycloak() *v1.Pod {
	return &v1.Pod{
		ObjectMeta: v12.ObjectMeta{
			Name: "apicurio-test-auth",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "sso",
					Image: "docker-registry.default.svc:5000/openshift/redhat-sso72-openshift",
					ImagePullPolicy: v1.PullAlways,
					//Resources: v1.ResourceRequirements{
					//	Limits: v1.ResourceList{
					//		"cpu": "500m",
					//	},
					//},
					Env: []v1.EnvVar{
						{
							Name:  "JGROUPS_PING_PROTOCOL",
							Value: "openshift.DNS_PING",
						},
						{
							Name:  "OPENSHIFT_DNS_PING_SERVICE_NAME",
							Value: "sso-ping",
						},
						{
							Name: "OPENSHIFT_DNS_PING_SERVICE_PORT",
							Value: "8888",
						},
						{
							Name: "X509_CA_BUNDLE",
							Value: "/var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt",
						},
						{
							Name: "JGROUPS_CLUSTER_PASSWORD",
							Value: "kEJy7VUrrB32u02rqIF60OhCFQqNClX8",
						},
						{
							Name: "SSO_ADMIN_USERNAME",
							Value: "admin",
						},
						{
							Name: "SSO_ADMIN_PASSWORD",
							Value: "admin",
						},
						{
							Name: "SSO_REALM",
							Value: "apicurio",
						},
					},
				},
			},
		},
	}
}

func SetApicurioAuth(pod *v1.Pod) {
	pod.Spec.Containers[0].Image = "apicurio/apicurio-studio-auth:latest-release"
}
