package defaulttolerationseconds

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/apis/core/helper"
	"k8s.io/kubernetes/plugin/pkg/admission/defaulttolerationseconds"
	"k8s.io/kubernetes/test/integration/framework"
)

func TestAdmission(t *testing.T) {
	masterConfig := framework.NewMasterConfig()
	masterConfig.GenericConfig.EnableProfiling = true
	masterConfig.GenericConfig.AdmissionControl = defaulttolerationseconds.NewDefaultTolerationSeconds()
	_, s, closeFn := framework.RunAMaster(masterConfig)
	defer closeFn()

	client := clientset.NewForConfigOrDie(&restclient.Config{Host: s.URL, ContentConfig: restclient.ContentConfig{GroupVersion: &schema.GroupVersion{Group: "", Version: "v1"}}})

	ns := framework.CreateTestingNamespace("default-toleration-seconds", s, t)
	defer framework.DeleteTestingNamespace(ns, s, t)

	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "foo",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "test",
					Image: "an-image",
				},
			},
		},
	}

	updatedPod, err := client.CoreV1().Pods(pod.Namespace).Create(&pod)
	if err != nil {
		t.Fatalf("error creating pod: %v", err)
	}

	var defaultSeconds int64 = 300
	nodeNotReady := v1.Toleration{
		Key:               v1.TaintNodeNotReady,
		Operator:          v1.TolerationOpExists,
		Effect:            v1.TaintEffectNoExecute,
		TolerationSeconds: &defaultSeconds,
	}

	nodeUnreachable := v1.Toleration{
		Key:               v1.TaintNodeUnreachable,
		Operator:          v1.TolerationOpExists,
		Effect:            v1.TaintEffectNoExecute,
		TolerationSeconds: &defaultSeconds,
	}

	found := 0
	tolerations := updatedPod.Spec.Tolerations
	for i := range tolerations {
		if found == 2 {
			break
		}
		if tolerations[i].MatchToleration(&nodeNotReady) {
			if helper.Semantic.DeepEqual(tolerations[i], nodeNotReady) {
				found++
				continue
			}
		}
		if tolerations[i].MatchToleration(&nodeUnreachable) {
			if helper.Semantic.DeepEqual(tolerations[i], nodeUnreachable) {
				found++
				continue
			}
		}
	}

	if found != 2 {
		t.Fatalf("unexpected tolerations: %v\n", updatedPod.Spec.Tolerations)
	}
}
