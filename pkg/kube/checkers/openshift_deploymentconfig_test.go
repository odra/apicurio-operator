package checkers

import (
	"github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestNewDeploymentConfigCheck(t *testing.T) {
	cases := []struct {
		Name     string
		DC       *v1.DeploymentConfig
		Result   bool
		Validate func(t *testing.T, currentResult bool, expectedResult bool)
	}{
		{
			Name: "Should validate to ready state",
			DC: &v1.DeploymentConfig{
				Status: v1.DeploymentConfigStatus{
					ReadyReplicas:     1,
					Replicas:          1,
					AvailableReplicas: 1,
					Conditions: []v1.DeploymentCondition{
						{
							Type:   v1.DeploymentAvailable,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			Result: true,
			Validate: func(t *testing.T, currentResult bool, expectedResult bool) {
				if currentResult != expectedResult {
					t.Fatalf("Expected: %t but got %t", expectedResult, currentResult)
				}
			},
		},
		{
			Name: "Should not validate to ready state",
			DC: &v1.DeploymentConfig{
				Status: v1.DeploymentConfigStatus{
					ReadyReplicas: 1,
					Replicas:      2,
				},
			},
			Result: false,
			Validate: func(t *testing.T, currentResult bool, expectedResult bool) {
				if currentResult != expectedResult {
					t.Fatalf("Expected: %t but got %t", expectedResult, currentResult)
				}
			},
		},
		{
			Name:   "Should not validate empty dc",
			DC:     &v1.DeploymentConfig{},
			Result: false,
			Validate: func(t *testing.T, currentResult bool, expectedResult bool) {
				if currentResult != expectedResult {
					t.Fatalf("Expected: %t but got %t", expectedResult, currentResult)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			dcc := NewDeploymentConfigCheck(tc.DC)
			tc.Validate(t, dcc.IsReady(), tc.Result)
		})
	}
}
