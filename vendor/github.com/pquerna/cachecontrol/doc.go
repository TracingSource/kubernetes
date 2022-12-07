// Package cachecontrol implements the logic for HTTP Caching
//
// Deciding if an HTTP Response can be cached is often harder
// and more bug prone than an actual cache storage backend.
// cachecontrol provides a simple interface to determine if
// request and response pairs are cachable as defined under
// RFC 7234 http://tools.ietf.org/html/rfc7234
package cachecontrol
