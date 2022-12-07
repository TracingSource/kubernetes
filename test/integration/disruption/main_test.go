package disruption

import (
	"k8s.io/kubernetes/test/integration/framework"
	"testing"
)

func TestMain(m *testing.M) {
	framework.EtcdMain(m.Run)
}
