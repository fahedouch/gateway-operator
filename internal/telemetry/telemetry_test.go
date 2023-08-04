package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-telemetry/pkg/forwarders"
	"github.com/kong/kubernetes-telemetry/pkg/serializers"
	"github.com/kong/kubernetes-telemetry/pkg/telemetry"
	"github.com/kong/kubernetes-telemetry/pkg/types"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	testdynclient "k8s.io/client-go/dynamic/fake"
	testk8sclient "k8s.io/client-go/kubernetes/fake"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func prepareScheme(t *testing.T) *runtime.Scheme {
	scheme := runtime.NewScheme()
	require.NoError(t, testk8sclient.AddToScheme(scheme))

	return scheme
}

func TestCreateManager(t *testing.T) {
	var (
		payload = types.ProviderReport{
			"v": "0.6.2",
		}
	)

	objs := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kong",
				Name:      "gateway-operator-abcde",
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-0",
			},
		},
	}
	t.Log("create mock kubernetes clients...")
	k8sclient := testk8sclient.NewSimpleClientset(objs...)
	ctrlClient := fakeclient.NewClientBuilder().Build()
	scheme := prepareScheme(t)
	dyn := testdynclient.NewSimpleDynamicClient(scheme, objs...)

	d, ok := k8sclient.Discovery().(*fakediscovery.FakeDiscovery)
	require.True(t, ok)
	d.FakedServerVersion = &version.Info{
		Major:      "1",
		Minor:      "27",
		GitVersion: "v1.27.2",
		Platform:   "linux/amd64",
	}

	t.Log("create telemetry manager...")
	m, err := createManager(
		types.Signal(SignalPing), k8sclient, ctrlClient, dyn, payload,
		telemetry.OptManagerPeriod(time.Hour),
		telemetry.OptManagerLogger(logr.Discard()),
	)
	require.NoError(t, err)
	t.Log("add consumer to manager...")
	ch := make(chan []byte)
	consumer := telemetry.NewConsumer(
		serializers.NewSemicolonDelimited(),
		forwarders.NewChannelForwarder(ch),
	)
	require.NoError(t, m.AddConsumer(consumer))

	t.Log("trigger a report...")
	require.NoError(t, m.Start())
	ctx := context.Background()
	err = m.TriggerExecute(ctx, "test-signal")
	require.NoError(t, err)

	t.Log("checking received report...")
	select {
	case buf := <-ch:
		checkTelemetryReportContent(t, string(buf),
			"signal=test-signal",
			"v=0.6.2",
			"k8sv=v1.27.2",
			"k8s_nodes_count=1",
			"k8s_pods_count=1",
		)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for report")
	case <-ctx.Done():
		t.Fatalf("context closed with error %v", ctx.Err())
	}
}

func checkTelemetryReportContent(t *testing.T, report string, containValue ...string) {
	for _, v := range containValue {
		require.Containsf(t, report, v, "report should contain %s", v)
	}
}
