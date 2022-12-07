// +build !windows

package dns

import (
	"strings"
)

const etcHostsFile = "/etc/hosts"

// A /etc/resolv.conf file managed by kubelet looks like this:
// nameserver DNS_CLUSTER_IP
// search test-dns.svc.cluster.local svc.cluster.local cluster.local q53aahaikqaehcai3ylfqdtc5b.bx.internal.cloudapp.net
// options ndots:5
func getDNSSuffixList() []string {
	fileData := readFile("/etc/resolv.conf")
	lines := strings.Split(fileData, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "search") {
			// omit the starting "search".
			return strings.Split(line, " ")[1:]
		}
	}

	panic("Could not find DNS search list!")
}

func getDNSServerList() []string {
	fileData := readFile("/etc/resolv.conf")
	lines := strings.Split(fileData, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "nameserver") {
			// omit the starting "nameserver".
			return strings.Split(line, " ")[1:]
		}
	}

	panic("Could not find DNS search list!")
}
