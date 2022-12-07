package deviceplugin

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/test/e2e/framework/testfiles"
)

const (
	// sampleResourceName is the name of the example resource which is used in the e2e test
	sampleResourceName = "example.com/resource"
	// sampleDevicePluginDSYAML is the path of the daemonset template of the sample device plugin. // TODO: Parametrize it by making it a feature in TestFramework.
	sampleDevicePluginDSYAML = "test/e2e/testing-manifests/sample-device-plugin.yaml"
	// sampleDevicePluginName is the name of the device plugin pod
	sampleDevicePluginName = "sample-device-plugin"
)

var (
	appsScheme = runtime.NewScheme()
	appsCodecs = serializer.NewCodecFactory(appsScheme)
)

// NumberOfSampleResources returns the number of resources advertised by a node
func NumberOfSampleResources(node *v1.Node) int64 {
	val, ok := node.Status.Capacity[sampleResourceName]

	if !ok {
		return 0
	}

	return val.Value()
}

// GetSampleDevicePluginPod returns the Device Plugin pod for sample resources in e2e tests
func GetSampleDevicePluginPod() *v1.Pod {
	ds := ReadDaemonSetV1OrDie(testfiles.ReadOrDie(sampleDevicePluginDSYAML))
	p := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sampleDevicePluginName,
			Namespace: metav1.NamespaceSystem,
		},

		Spec: ds.Spec.Template.Spec,
	}

	return p
}

// ReadDaemonSetV1OrDie reads daemonset object from bytes. Panics on error.
func ReadDaemonSetV1OrDie(objBytes []byte) *appsv1.DaemonSet {
	appsv1.AddToScheme(appsScheme)
	requiredObj, err := runtime.Decode(appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion), objBytes)
	if err != nil {
		panic(err)
	}
	return requiredObj.(*appsv1.DaemonSet)
}
