package fuzzer

import (
	"fmt"
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/storage"
)

// Funcs returns the fuzzer functions for the storage api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *storage.StorageClass, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again
			reclamationPolicies := []api.PersistentVolumeReclaimPolicy{api.PersistentVolumeReclaimDelete, api.PersistentVolumeReclaimRetain}
			obj.ReclaimPolicy = &reclamationPolicies[c.Rand.Intn(len(reclamationPolicies))]
			bindingModes := []storage.VolumeBindingMode{storage.VolumeBindingImmediate, storage.VolumeBindingWaitForFirstConsumer}
			obj.VolumeBindingMode = &bindingModes[c.Rand.Intn(len(bindingModes))]
		},
		func(obj *storage.CSIDriver, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again

			// Custom fuzzing for volume modes.
			switch c.Rand.Intn(7) {
			case 0:
				obj.Spec.VolumeLifecycleModes = nil
			case 1:
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{}
			case 2:
				// Invalid mode.
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{
					storage.VolumeLifecycleMode(fmt.Sprintf("%d", c.Rand.Int31())),
				}
			case 3:
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{
					storage.VolumeLifecyclePersistent,
				}
			case 4:
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{
					storage.VolumeLifecycleEphemeral,
				}
			case 5:
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{
					storage.VolumeLifecyclePersistent,
					storage.VolumeLifecycleEphemeral,
				}
			case 6:
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{
					storage.VolumeLifecycleEphemeral,
					storage.VolumeLifecyclePersistent,
				}
			}

			// match defaulting
			if obj.Spec.AttachRequired == nil {
				obj.Spec.AttachRequired = new(bool)
				*(obj.Spec.AttachRequired) = true
			}
			if obj.Spec.PodInfoOnMount == nil {
				obj.Spec.PodInfoOnMount = new(bool)
				*(obj.Spec.PodInfoOnMount) = false
			}
			if len(obj.Spec.VolumeLifecycleModes) == 0 {
				obj.Spec.VolumeLifecycleModes = []storage.VolumeLifecycleMode{
					storage.VolumeLifecyclePersistent,
				}
			}
		},
	}
}
