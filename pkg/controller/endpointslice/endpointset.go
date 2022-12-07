package endpointslice

import (
	"sort"

	discovery "k8s.io/api/discovery/v1beta1"
	endpointutil "k8s.io/kubernetes/pkg/controller/util/endpoint"
)

// endpointHash is used to uniquely identify endpoints. Only including addresses
// and hostnames as unique identifiers allows us to do more in place updates
// should attributes such as topology, conditions, or targetRef change.
type endpointHash string
type endpointHashObj struct {
	Addresses []string
	Hostname  string
}

func hashEndpoint(endpoint *discovery.Endpoint) endpointHash {
	sort.Strings(endpoint.Addresses)
	hashObj := endpointHashObj{Addresses: endpoint.Addresses}
	if endpoint.Hostname != nil {
		hashObj.Hostname = *endpoint.Hostname
	}

	return endpointHash(endpointutil.DeepHashObjectToString(hashObj))
}

// endpointSet provides simple methods for comparing sets of Endpoints.
type endpointSet map[endpointHash]*discovery.Endpoint

// Insert adds items to the set.
func (s endpointSet) Insert(items ...*discovery.Endpoint) endpointSet {
	for _, item := range items {
		s[hashEndpoint(item)] = item
	}
	return s
}

// Delete removes all items from the set.
func (s endpointSet) Delete(items ...*discovery.Endpoint) endpointSet {
	for _, item := range items {
		delete(s, hashEndpoint(item))
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s endpointSet) Has(item *discovery.Endpoint) bool {
	_, contained := s[hashEndpoint(item)]
	return contained
}

// Returns an endpoint matching the hash if contained in the set.
func (s endpointSet) Get(item *discovery.Endpoint) *discovery.Endpoint {
	return s[hashEndpoint(item)]
}

// UnsortedList returns the slice with contents in random order.
func (s endpointSet) UnsortedList() []*discovery.Endpoint {
	endpoints := make([]*discovery.Endpoint, 0, len(s))
	for _, endpoint := range s {
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

// Returns a single element from the set.
func (s endpointSet) PopAny() (*discovery.Endpoint, bool) {
	for _, endpoint := range s {
		s.Delete(endpoint)
		return endpoint, true
	}
	return nil, false
}

// Len returns the size of the set.
func (s endpointSet) Len() int {
	return len(s)
}
