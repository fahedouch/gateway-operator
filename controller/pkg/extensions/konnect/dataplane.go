package konnect

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dputils "github.com/kong/gateway-operator/internal/utils/dataplane"
	"github.com/kong/gateway-operator/pkg/consts"
	k8sutils "github.com/kong/gateway-operator/pkg/utils/kubernetes"
	k8sresources "github.com/kong/gateway-operator/pkg/utils/kubernetes/resources"

	operatorv1beta1 "github.com/kong/kubernetes-configuration/api/gateway-operator/v1beta1"
	konnectv1alpha1 "github.com/kong/kubernetes-configuration/api/konnect/v1alpha1"
)

var (
	// ErrCrossNamespaceReference is returned when a Konnect extension references a different namespace.
	ErrCrossNamespaceReference = errors.New("cross-namespace reference is not currently supported for Konnect extensions")
	// ErrKonnectExtensionNotFound is returned when a Konnect extension is not found.
	ErrKonnectExtensionNotFound = errors.New("konnect extension not found")
	// ErrClusterCertificateNotFound is returned when a cluster certificate secret referenced in the KonnectExtension is not found.
	ErrClusterCertificateNotFound = errors.New("cluster certificate not found")
	// ErrKonnectExtensionNotReady is returned when a Konnect extension is not ready.
	ErrKonnectExtensionNotReady = errors.New("konnect extension is not ready")
)

// ApplyKonnectExtension gets the DataPlane as argument, and in case it references a KonnectExtension, it
// fetches the referenced extension and applies the necessary changes to the DataPlane spec.
func ApplyKonnectExtension(ctx context.Context, cl client.Client, dataplane *operatorv1beta1.DataPlane) error {
	for _, extensionRef := range dataplane.Spec.Extensions {
		if extensionRef.Group != konnectv1alpha1.SchemeGroupVersion.Group || extensionRef.Kind != konnectv1alpha1.KonnectExtensionKind {
			continue
		}
		namespace := dataplane.Namespace
		if extensionRef.Namespace != nil && *extensionRef.Namespace != namespace {
			return errors.Join(ErrCrossNamespaceReference, fmt.Errorf("the cross-namespace reference to the extension %s/%s is not permitted", *extensionRef.Namespace, extensionRef.Name))
		}

		konnectExt := konnectv1alpha1.KonnectExtension{}
		if err := cl.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      extensionRef.Name,
		}, &konnectExt); err != nil {
			if k8serrors.IsNotFound(err) {
				return errors.Join(ErrKonnectExtensionNotFound, fmt.Errorf("the extension %s/%s referenced by the DataPlane is not found", namespace, extensionRef.Name))
			} else {
				return err
			}
		}

		if !k8sutils.HasConditionTrue(konnectv1alpha1.KonnectExtensionReadyConditionType, &konnectExt) {
			return errors.Join(ErrKonnectExtensionNotReady, fmt.Errorf("the extension %s/%s referenced by the DataPlane is not ready", namespace, extensionRef.Name))
		}

		if dataplane.Spec.Deployment.PodTemplateSpec == nil {
			dataplane.Spec.Deployment.PodTemplateSpec = &corev1.PodTemplateSpec{}
		}

		d := k8sresources.Deployment(appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Template: *dataplane.Spec.Deployment.PodTemplateSpec,
			},
		})
		if container := k8sutils.GetPodContainerByName(&d.Spec.Template.Spec, consts.DataPlaneProxyContainerName); container == nil {
			d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, corev1.Container{
				Name: consts.DataPlaneProxyContainerName,
			})
		}

		d.WithVolume(kongInKonnectClusterCertificateVolume())
		d.WithVolumeMount(kongInKonnectClusterCertificateVolumeMount(), consts.DataPlaneProxyContainerName)
		d.WithVolume(kongInKonnectClusterCertVolume(konnectExt.Status.DataPlaneClientAuth.CertificateSecretRef.Name))
		d.WithVolumeMount(kongInKonnectClusterVolumeMount(), consts.DataPlaneProxyContainerName)

		// KonnectID is the only supported type for now, and its presence is guaranteed by a proper CEL rule.
		envSet := dputils.KongInKonnectDefaults(konnectExt.Status)

		dputils.FillDataPlaneProxyContainerEnvs(nil, &d.Spec.Template, envSet)
		dataplane.Spec.Deployment.PodTemplateSpec = &d.Spec.Template
	}
	return nil
}

func kongInKonnectClusterCertVolume(secretName string) corev1.Volume {
	return corev1.Volume{
		Name: consts.KongClusterCertVolume,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  secretName,
				DefaultMode: lo.ToPtr(int32(420)),
			},
		},
	}
}

func kongInKonnectClusterVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      consts.KongClusterCertVolume,
		MountPath: consts.KongClusterCertVolumeMountPath,
	}
}

func kongInKonnectClusterCertificateVolume() corev1.Volume {
	return corev1.Volume{
		Name: consts.ClusterCertificateVolume,
	}
}

func kongInKonnectClusterCertificateVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      consts.ClusterCertificateVolume,
		MountPath: consts.ClusterCertificateVolumeMountPath,
		ReadOnly:  true,
	}
}
