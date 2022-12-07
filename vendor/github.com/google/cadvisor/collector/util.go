
package collector

import "github.com/google/cadvisor/container"

func (endpointConfig *EndpointConfig) configure(containerHandler container.ContainerHandler) {
	//If the exact URL was not specified, generate it based on the ip address of the container.
	endpoint := endpointConfig
	if endpoint.URL == "" {
		ipAddress := containerHandler.GetContainerIPAddress()
		endpointConfig.URL = endpoint.URLConfig.Protocol + "://" + ipAddress + ":" + endpoint.URLConfig.Port.String() + endpoint.URLConfig.Path
	}
}
