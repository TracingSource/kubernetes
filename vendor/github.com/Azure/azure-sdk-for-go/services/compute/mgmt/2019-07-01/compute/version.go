package compute

import "github.com/Azure/azure-sdk-for-go/version"

// UserAgent returns the UserAgent string to use when sending http.Requests.
func UserAgent() string {
	return "Azure-SDK-For-Go/" + version.Number + " compute/2019-07-01"
}

// Version returns the semantic version (see http://semver.org) of the client.
func Version() string {
	return version.Number
}
