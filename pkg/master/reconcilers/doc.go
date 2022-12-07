// Package reconcilers provides objects for managing the list of active masters.
// NOTE: The Lease reconciler is not the intended way for any apiserver other
// than kube-apiserver to accomplish the task of Endpoint registration. This is
// a special case for the time being.
package reconcilers
