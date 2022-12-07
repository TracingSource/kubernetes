package dns

import (
	"strings"
)

const etcHostsFile = "C:/Windows/System32/drivers/etc/hosts"

func getDNSSuffixList() []string {
	output := runCommand("powershell", "-Command", "(Get-DnsClient)[0].SuffixSearchList")
	if len(output) > 0 {
		return strings.Split(output, "\r\n")
	}

	panic("Could not find DNS search list!")
}

func getDNSServerList() []string {
	output := runCommand("powershell", "-Command", "(Get-DnsClientServerAddress).ServerAddresses")
	if len(output) > 0 {
		return strings.Split(output, "\r\n")
	}

	panic("Could not find DNS Server list!")
}
