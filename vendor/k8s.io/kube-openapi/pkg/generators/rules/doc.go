// Package rules contains API rules that are enforced in OpenAPI spec generation
// as part of the machinery. Files under this package implement APIRule interface
// which evaluates Go type and produces list of API rule violations.
//
// Implementations of APIRule should be added to API linter under openAPIGen code-
// generator to get integrated in the generation process.
package rules
