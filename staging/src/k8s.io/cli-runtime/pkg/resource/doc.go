// Package resource assists clients in dealing with RESTful objects that match the
// Kubernetes API conventions. The Helper object provides simple CRUD operations
// on resources. The Visitor interface makes it easy to deal with multiple resources
// in bulk for retrieval and operation. The Builder object simplifies converting
// standard command line arguments and parameters into a Visitor that can iterate
// over all of the identified resources, whether on the server or on the local
// filesystem.
package resource // import "k8s.io/cli-runtime/pkg/resource"
