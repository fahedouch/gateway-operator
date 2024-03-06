// This file is generated by /hack/generators/kic/webhook-config-generator. DO NOT EDIT.

package resources

import (
	"fmt"

	"github.com/Masterminds/semver"
	webhook "github.com/kong/gateway-operator/pkg/utils/kubernetes/resources/validatingwebhookconfig"
	admregv1 "k8s.io/api/admissionregistration/v1"
	pkgapisadmregv1 "k8s.io/kubernetes/pkg/apis/admissionregistration/v1"
)

// GenerateValidatingWebhookConfigurationForControlPlane generates a ValidatingWebhookConfiguration for a control plane
// based on the control plane version. It also overrides all webhooks' client configurations with the provided service
// details.
func GenerateValidatingWebhookConfigurationForControlPlane(webhookName string, cpVersion *semver.Version, clientConfig admregv1.WebhookClientConfig) (*admregv1.ValidatingWebhookConfiguration, error) {
	if webhookName == "" {
		return nil, fmt.Errorf("webhook name is required")
	}
	if cpVersion == nil {
		return nil, fmt.Errorf("control plane version is required")
	}

	var (
		constraint *semver.Constraints
		err        error
	)
	
	constraint, err = semver.NewConstraint(">=3.1")
	if err != nil {
		return nil, err
	}
	if constraint.Check(cpVersion) {
		cfg := webhook.GenerateValidatingWebhookConfigurationForKIC_ge3_1(webhookName, clientConfig)
		pkgapisadmregv1.SetObjectDefaults_ValidatingWebhookConfiguration(cfg)
		LabelObjectAsControlPlaneManaged(cfg)
		return cfg, nil
	}	
	

	return nil, ErrControlPlaneVersionNotSupported 
}
