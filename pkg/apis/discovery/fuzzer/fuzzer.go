package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/discovery"
)

// Funcs returns the fuzzer functions for the discovery api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *discovery.EndpointSlice, c fuzz.Continue) {
			c.FuzzNoCustom(obj) // fuzz self without calling this function again

			addressTypes := []discovery.AddressType{discovery.AddressTypeIPv4, discovery.AddressTypeIPv6, discovery.AddressTypeFQDN}
			obj.AddressType = addressTypes[c.Rand.Intn(len(addressTypes))]

			for i, endpointPort := range obj.Ports {
				if endpointPort.Name == nil {
					emptyStr := ""
					obj.Ports[i].Name = &emptyStr
				}

				if endpointPort.Protocol == nil {
					protos := []api.Protocol{api.ProtocolTCP, api.ProtocolUDP, api.ProtocolSCTP}
					obj.Ports[i].Protocol = &protos[c.Rand.Intn(len(protos))]
				}
			}
		},
	}
}
