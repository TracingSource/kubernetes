package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubectrlmgrconfigv1alpha1 "k8s.io/kubernetes/pkg/controller/apis/config/v1alpha1"
	serviceconfigv1alpha1 "k8s.io/kubernetes/pkg/controller/service/config/v1alpha1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_CloudControllerManagerConfiguration(obj *CloudControllerManagerConfiguration) {
	zero := metav1.Duration{}
	if obj.NodeStatusUpdateFrequency == zero {
		obj.NodeStatusUpdateFrequency = metav1.Duration{Duration: 5 * time.Minute}
	}

	// These defaults override the recommended defaults from the apimachineryconfigv1alpha1 package that are applied automatically
	// These client-connection defaults are specific to the cloud-controller-manager
	if obj.Generic.ClientConnection.QPS == 0 {
		obj.Generic.ClientConnection.QPS = 20
	}
	if obj.Generic.ClientConnection.Burst == 0 {
		obj.Generic.ClientConnection.Burst = 30
	}

	// Use the default RecommendedDefaultGenericControllerManagerConfiguration options
	kubectrlmgrconfigv1alpha1.RecommendedDefaultGenericControllerManagerConfiguration(&obj.Generic)
	// Use the default RecommendedDefaultServiceControllerConfiguration options
	serviceconfigv1alpha1.RecommendedDefaultServiceControllerConfiguration(&obj.ServiceController)
}
