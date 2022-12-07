package storage

import "github.com/Azure/azure-sdk-for-go/version"

// UserAgent returns the UserAgent string to use when sending http.Requests.
func UserAgent() string {
	return "Azure-SDK-For-Go/" + version.Number + " storage/2019-04-01"
}

// Version returns the semantic version (see http://semver.org) of the client.
func Version() string {
	return version.Number
}
