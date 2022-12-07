// +build !windows

package kubelet

func getOSSpecificLabels() (map[string]string, error) {
	return nil, nil
}
