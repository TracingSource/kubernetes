package tokens

import (
	"fmt"
	"regexp"

	"k8s.io/cluster-bootstrap/token/api"
)

var (
	bootstrapTokenRe = regexp.MustCompile(api.BootstrapTokenPattern)
)

// ParseToken tries and parse a valid token from a string.
// A token ID and token secret are returned in case of success, an error otherwise.
func ParseToken(s string) (tokenID, tokenSecret string, err error) {
	split := bootstrapTokenRe.FindStringSubmatch(s)
	if len(split) != 3 {
		return "", "", fmt.Errorf("token [%q] was not of form [%q]", s, api.BootstrapTokenPattern)
	}
	return split[1], split[2], nil
}
