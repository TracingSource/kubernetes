package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apimachineryversion "k8s.io/apimachinery/pkg/version"
)

func TestVersionParsing(t *testing.T) {
	var tests = []struct {
		desc            string
		versionString   string
		expectedVersion string
	}{
		{
			"v1.15.0-alpha-1.12345",
			"v1.15.0-alpha-1.12345",
			"1.15.0",
		},
		{
			"Parse out defaulted string",
			"v0.0.0-master",
			"0.0.0",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			version := apimachineryversion.Info{
				GitVersion: test.versionString,
			}
			parsedV := parseVersion(version)
			assert.Equalf(t, test.expectedVersion, parsedV.String(), "Got %v, wanted %v", parsedV.String(), test.expectedVersion)
		})
	}
}
