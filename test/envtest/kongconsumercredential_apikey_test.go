package envtest

import (
	"context"
	"testing"

	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"
	sdkkonnectops "github.com/Kong/sdk-konnect-go/models/operations"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/gateway-operator/controller/konnect"
	"github.com/kong/gateway-operator/controller/konnect/ops"
	"github.com/kong/gateway-operator/modules/manager"
	"github.com/kong/gateway-operator/modules/manager/scheme"
	"github.com/kong/gateway-operator/test/helpers/deploy"

	configurationv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	configurationv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	"github.com/kong/kubernetes-configuration/api/konnect/v1alpha1"
)

func TestKongConsumerCredential_APIKey(t *testing.T) {
	t.Parallel()
	ctx, cancel := Context(t, context.Background())
	defer cancel()

	// Setup up the envtest environment.
	cfg, ns := Setup(t, ctx, scheme.Get())

	mgr, logs := NewManager(t, ctx, cfg, scheme.Get())

	clientNamespaced := client.NewNamespacedClient(mgr.GetClient(), ns.Name)

	apiAuth := deploy.KonnectAPIAuthConfigurationWithProgrammed(t, ctx, clientNamespaced)
	cp := deploy.KonnectGatewayControlPlaneWithID(t, ctx, clientNamespaced, apiAuth)

	consumerID := uuid.NewString()
	consumer := deploy.KongConsumerWithProgrammed(t, ctx, clientNamespaced, &configurationv1.KongConsumer{
		Username: "username1",
		Spec: configurationv1.KongConsumerSpec{
			ControlPlaneRef: &configurationv1alpha1.ControlPlaneRef{
				Type: configurationv1alpha1.ControlPlaneRefKonnectNamespacedRef,
				KonnectNamespacedRef: &configurationv1alpha1.KonnectNamespacedRef{
					Name: cp.Name,
				},
			},
		},
	})
	consumer.Status.Konnect = &v1alpha1.KonnectEntityStatusWithControlPlaneRef{
		ControlPlaneID: cp.GetKonnectStatus().GetKonnectID(),
		KonnectEntityStatus: v1alpha1.KonnectEntityStatus{
			ID:        consumerID,
			ServerURL: cp.GetKonnectStatus().GetServerURL(),
			OrgID:     cp.GetKonnectStatus().GetOrgID(),
		},
	}
	require.NoError(t, clientNamespaced.Status().Update(ctx, consumer))

	kongCredentialAPIKey := deploy.KongCredentialAPIKey(t, ctx, clientNamespaced, consumer.Name)
	keyID := uuid.NewString()
	tags := []string{
		"k8s-generation:1",
		"k8s-group:configuration.konghq.com",
		"k8s-kind:KongCredentialAPIKey",
		"k8s-name:" + kongCredentialAPIKey.Name,
		"k8s-namespace:" + ns.Name,
		"k8s-uid:" + string(kongCredentialAPIKey.GetUID()),
		"k8s-version:v1alpha1",
	}

	factory := ops.NewMockSDKFactory(t)
	factory.SDK.KongCredentialsAPIKeySDK.EXPECT().
		CreateKeyAuthWithConsumer(
			mock.Anything,
			sdkkonnectops.CreateKeyAuthWithConsumerRequest{
				ControlPlaneID:              cp.GetKonnectStatus().GetKonnectID(),
				ConsumerIDForNestedEntities: consumerID,
				KeyAuthWithoutParents: sdkkonnectcomp.KeyAuthWithoutParents{
					Key:  lo.ToPtr("key"),
					Tags: tags,
				},
			},
		).
		Return(
			&sdkkonnectops.CreateKeyAuthWithConsumerResponse{
				KeyAuth: &sdkkonnectcomp.KeyAuth{
					ID: lo.ToPtr(keyID),
				},
			},
			nil,
		)
	factory.SDK.KongCredentialsAPIKeySDK.EXPECT().
		UpsertKeyAuthWithConsumer(mock.Anything, mock.Anything, mock.Anything).Maybe().
		Return(
			&sdkkonnectops.UpsertKeyAuthWithConsumerResponse{
				KeyAuth: &sdkkonnectcomp.KeyAuth{
					ID: lo.ToPtr(keyID),
				},
			},
			nil,
		)

	require.NoError(t, manager.SetupCacheIndicesForKonnectTypes(ctx, mgr, false))
	reconcilers := []Reconciler{
		konnect.NewKonnectEntityReconciler(factory, false, mgr.GetClient(),
			konnect.WithKonnectEntitySyncPeriod[configurationv1alpha1.KongCredentialAPIKey](konnectInfiniteSyncTime),
		),
	}

	StartReconcilers(ctx, t, mgr, logs, reconcilers...)

	assert.EventuallyWithT(t,
		assertCollectObjectExistsAndHasKonnectID(t, ctx, clientNamespaced, kongCredentialAPIKey, keyID),
		waitTime, tickTime,
		"KongCredentialAPIKey wasn't created",
	)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.True(c, factory.SDK.KongCredentialsAPIKeySDK.AssertExpectations(t))
	}, waitTime, tickTime)

	factory.SDK.KongCredentialsAPIKeySDK.EXPECT().
		DeleteKeyAuthWithConsumer(
			mock.Anything,
			sdkkonnectops.DeleteKeyAuthWithConsumerRequest{
				ControlPlaneID:              cp.GetKonnectStatus().GetKonnectID(),
				ConsumerIDForNestedEntities: consumerID,
				KeyAuthID:                   keyID,
			},
		).
		Return(
			&sdkkonnectops.DeleteKeyAuthWithConsumerResponse{
				StatusCode: 200,
			},
			nil,
		)

	require.NoError(t, clientNamespaced.Delete(ctx, kongCredentialAPIKey))

	assert.EventuallyWithT(t,
		func(c *assert.CollectT) {
			assert.True(c, k8serrors.IsNotFound(
				clientNamespaced.Get(ctx, client.ObjectKeyFromObject(kongCredentialAPIKey), kongCredentialAPIKey),
			))
		}, waitTime, tickTime,
		"KongCredentialAPIKey wasn't deleted but it should have been",
	)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.True(c, factory.SDK.KongCredentialsAPIKeySDK.AssertExpectations(t))
	}, waitTime, tickTime)
}
