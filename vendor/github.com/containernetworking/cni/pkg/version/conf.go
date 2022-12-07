package version

import (
	"encoding/json"
	"fmt"
)

// ConfigDecoder can decode the CNI version available in network config data
type ConfigDecoder struct{}

func (*ConfigDecoder) Decode(jsonBytes []byte) (string, error) {
	var conf struct {
		CNIVersion string `json:"cniVersion"`
	}
	err := json.Unmarshal(jsonBytes, &conf)
	if err != nil {
		return "", fmt.Errorf("decoding version from network config: %s", err)
	}
	if conf.CNIVersion == "" {
		return "0.1.0", nil
	}
	return conf.CNIVersion, nil
}
