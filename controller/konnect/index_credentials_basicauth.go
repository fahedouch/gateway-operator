package konnect

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	configurationv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
)

const (
	// IndexFieldKongCredentialBasicAuthReferencesKongConsumer is the index name for KongCredentialBasicAuth -> Consumer.
	IndexFieldKongCredentialBasicAuthReferencesKongConsumer = "kongCredentialsBasicAuthConsumerRef"
	// IndexFieldKongCredentialBasicAuthReferencesSecret is the index name for KongCredentialBasicAuth -> Secret.
	IndexFieldKongCredentialBasicAuthReferencesSecret = "kongCredentialsBasicAuthSecretRef"
)

// IndexOptionsForCredentialsBasicAuth returns required Index options for KongCredentialBasicAuth.
func IndexOptionsForCredentialsBasicAuth() []ReconciliationIndexOption {
	return []ReconciliationIndexOption{
		{
			IndexObject:  &configurationv1alpha1.KongCredentialBasicAuth{},
			IndexField:   IndexFieldKongCredentialBasicAuthReferencesKongConsumer,
			ExtractValue: kongCredentialBasicAuthReferencesConsumer,
		},
		{
			IndexObject:  &configurationv1alpha1.KongCredentialBasicAuth{},
			IndexField:   IndexFieldKongCredentialBasicAuthReferencesSecret,
			ExtractValue: kongCredentialReferencesSecret[configurationv1alpha1.KongCredentialBasicAuth],
		},
	}
}

// kongCredentialBasicAuthReferencesConsumer returns the name of referenced Consumer.
func kongCredentialBasicAuthReferencesConsumer(obj client.Object) []string {
	cred, ok := obj.(*configurationv1alpha1.KongCredentialBasicAuth)
	if !ok {
		return nil
	}
	return []string{cred.Spec.ConsumerRef.Name}
}
