// Package schema defines a targeted schema language which allows one to
// represent all the schema information necessary to perform "structured"
// merges and diffs.
//
// Due to the targeted nature of the data model, the schema language can fit in
// just a few hundred lines of go code, making it much more understandable and
// concise than e.g. OpenAPI.
//
// This schema was derived by observing the API objects used by Kubernetes, and
// formalizing a model which allows certain operations ("apply") to be more
// well defined. It is currently missing one feature: one-of ("unions").
package schema
