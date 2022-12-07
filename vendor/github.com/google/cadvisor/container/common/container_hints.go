// Unmarshal's a Containers description json file. The json file contains
// an array of ContainerHint structs, each with a container's id and networkInterface
// This allows collecting stats about network interfaces configured outside docker
// and lxc
package common

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var ArgContainerHints = flag.String("container_hints", "/etc/cadvisor/container_hints.json", "location of the container hints file")

type containerHints struct {
	AllHosts []containerHint `json:"all_hosts,omitempty"`
}

type containerHint struct {
	FullName         string            `json:"full_path,omitempty"`
	NetworkInterface *networkInterface `json:"network_interface,omitempty"`
	Mounts           []Mount           `json:"mounts,omitempty"`
}

type Mount struct {
	HostDir      string `json:"host_dir,omitempty"`
	ContainerDir string `json:"container_dir,omitempty"`
}

type networkInterface struct {
	VethHost  string `json:"veth_host,omitempty"`
	VethChild string `json:"veth_child,omitempty"`
}

func GetContainerHintsFromFile(containerHintsFile string) (containerHints, error) {
	dat, err := ioutil.ReadFile(containerHintsFile)
	if os.IsNotExist(err) {
		return containerHints{}, nil
	}
	var cHints containerHints
	if err == nil {
		err = json.Unmarshal(dat, &cHints)
	}

	return cHints, err
}
