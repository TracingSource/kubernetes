package framework

import (
	"bufio"
	"fmt"
	"strings"

	// TODO: Remove the following imports (ref: https://github.com/kubernetes/kubernetes/issues/81245)
	e2essh "k8s.io/kubernetes/test/e2e/framework/ssh"
)

// KubemarkResourceUsage is a struct for tracking the resource usage of kubemark.
type KubemarkResourceUsage struct {
	Name                    string
	MemoryWorkingSetInBytes uint64
	CPUUsageInCores         float64
}

func getMasterUsageByPrefix(prefix string) (string, error) {
	sshResult, err := e2essh.SSH(fmt.Sprintf("ps ax -o %%cpu,rss,command | tail -n +2 | grep %v | sed 's/\\s+/ /g'", prefix), GetMasterHost()+":22", TestContext.Provider)
	if err != nil {
		return "", err
	}
	return sshResult.Stdout, nil
}

// GetKubemarkMasterComponentsResourceUsage returns the resource usage of kubemark which contains multiple combinations of cpu and memory usage for each pod name.
// TODO: figure out how to move this to kubemark directory (need to factor test SSH out of e2e framework)
func GetKubemarkMasterComponentsResourceUsage() map[string]*KubemarkResourceUsage {
	result := make(map[string]*KubemarkResourceUsage)
	// Get kubernetes component resource usage
	sshResult, err := getMasterUsageByPrefix("kube")
	if err != nil {
		Logf("Error when trying to SSH to master machine. Skipping probe. %v", err)
		return nil
	}
	scanner := bufio.NewScanner(strings.NewReader(sshResult))
	for scanner.Scan() {
		var cpu float64
		var mem uint64
		var name string
		fmt.Sscanf(strings.TrimSpace(scanner.Text()), "%f %d /usr/local/bin/kube-%s", &cpu, &mem, &name)
		if name != "" {
			// Gatherer expects pod_name/container_name format
			fullName := name + "/" + name
			result[fullName] = &KubemarkResourceUsage{Name: fullName, MemoryWorkingSetInBytes: mem * 1024, CPUUsageInCores: cpu / 100}
		}
	}
	// Get etcd resource usage
	sshResult, err = getMasterUsageByPrefix("bin/etcd")
	if err != nil {
		Logf("Error when trying to SSH to master machine. Skipping probe")
		return nil
	}
	scanner = bufio.NewScanner(strings.NewReader(sshResult))
	for scanner.Scan() {
		var cpu float64
		var mem uint64
		var etcdKind string
		fmt.Sscanf(strings.TrimSpace(scanner.Text()), "%f %d /bin/sh -c /usr/local/bin/etcd", &cpu, &mem)
		dataDirStart := strings.Index(scanner.Text(), "--data-dir")
		if dataDirStart < 0 {
			continue
		}
		fmt.Sscanf(scanner.Text()[dataDirStart:], "--data-dir=/var/%s", &etcdKind)
		if etcdKind != "" {
			// Gatherer expects pod_name/container_name format
			fullName := "etcd/" + etcdKind
			result[fullName] = &KubemarkResourceUsage{Name: fullName, MemoryWorkingSetInBytes: mem * 1024, CPUUsageInCores: cpu / 100}
		}
	}
	return result
}
