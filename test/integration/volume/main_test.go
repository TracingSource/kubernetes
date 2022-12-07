package volume

import (
	"math/rand"
	"testing"
	"time"

	"k8s.io/kubernetes/test/integration/framework"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	framework.EtcdMain(m.Run)
}
