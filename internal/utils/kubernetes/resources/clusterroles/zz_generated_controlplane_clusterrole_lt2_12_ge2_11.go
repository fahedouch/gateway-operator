// This file is generated by /hack/generators/kic-clusterrole-generator. DO NOT EDIT.

package clusterroles

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// -----------------------------------------------------------------------------
// ClusterRole generator
// -----------------------------------------------------------------------------

// GenerateNewClusterRoleForControlPlane_lt2_12_ge2_11 is a helper to generate a ClusterRole
// resource with all the permissions needed by the controlplane deployment.
// It is used for controlplanes that match the semver constraint "<2.12, >=2.11"
func GenerateNewClusterRoleForControlPlane_lt2_12_ge2_11(controlplaneName string) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", controlplaneName),
			Labels: map[string]string{
				"app": controlplaneName,
			},
		},
		Rules: []rbacv1.PolicyRule{

			{
				APIGroups: []string{
					"apiextensions.k8s.io",
				},
				Resources: []string{
					"customresourcedefinitions",
				},
				Verbs: []string{
					"list", "watch",
				},
			},

			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"events",
				},
				Verbs: []string{
					"create", "patch",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
				},
				Verbs: []string{
					"list", "watch",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"secrets",
				},
				Verbs: []string{
					"list", "watch",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"services",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"services/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"ingressclassparameterses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongclusterplugins",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongclusterplugins/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongconsumergroups",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongconsumergroups/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongconsumers",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongconsumers/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongingresses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongingresses/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongplugins",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"kongplugins/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"tcpingresses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"tcpingresses/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"udpingresses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"configuration.konghq.com",
				},
				Resources: []string{
					"udpingresses/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"discovery.k8s.io",
				},
				Resources: []string{
					"endpointslices",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"networking.k8s.io",
				},
				Resources: []string{
					"ingressclasses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"networking.k8s.io",
				},
				Resources: []string{
					"ingresses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"networking.k8s.io",
				},
				Resources: []string{
					"ingresses/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},

			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"configmaps",
				},
				Verbs: []string{
					"get", "list", "watch", "create", "update", "patch", "delete",
				},
			},
			{
				APIGroups: []string{
					"coordination.k8s.io",
				},
				Resources: []string{
					"leases",
				},
				Verbs: []string{
					"get", "list", "watch", "create", "update", "patch", "delete",
				},
			},

			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"namespaces",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"gatewayclasses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"gatewayclasses/status",
				},
				Verbs: []string{
					"get", "update",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"gateways",
				},
				Verbs: []string{
					"get", "list", "update", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"gateways/status",
				},
				Verbs: []string{
					"get", "update",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"grpcroutes",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"grpcroutes/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"httproutes",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"httproutes/status",
				},
				Verbs: []string{
					"get", "update",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"referencegrants",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"referencegrants/status",
				},
				Verbs: []string{
					"get",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"tcproutes",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"tcproutes/status",
				},
				Verbs: []string{
					"get", "update",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"tlsroutes",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"tlsroutes/status",
				},
				Verbs: []string{
					"get", "update",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"udproutes",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"gateway.networking.k8s.io",
				},
				Resources: []string{
					"udproutes/status",
				},
				Verbs: []string{
					"get", "update",
				},
			},

			{
				APIGroups: []string{
					"networking.internal.knative.dev",
				},
				Resources: []string{
					"ingresses",
				},
				Verbs: []string{
					"get", "list", "watch",
				},
			},
			{
				APIGroups: []string{
					"networking.internal.knative.dev",
				},
				Resources: []string{
					"ingresses/status",
				},
				Verbs: []string{
					"get", "patch", "update",
				},
			},
		},
	}
}
