package componentstatus

import (
	"crypto/tls"
	"sync"
	"time"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/kubernetes/pkg/probe"
	httpprober "k8s.io/kubernetes/pkg/probe/http"
)

const (
	probeTimeOut = 20 * time.Second
)

type ValidatorFn func([]byte) error

type Server struct {
	Addr        string
	Port        int
	Path        string
	EnableHTTPS bool
	TLSConfig   *tls.Config
	Validate    ValidatorFn
	Prober      httpprober.Prober
	Once        sync.Once
}

type ServerStatus struct {
	// +optional
	Component string `json:"component,omitempty"`
	// +optional
	Health string `json:"health,omitempty"`
	// +optional
	HealthCode probe.Result `json:"healthCode,omitempty"`
	// +optional
	Msg string `json:"msg,omitempty"`
	// +optional
	Err string `json:"err,omitempty"`
}

func (server *Server) DoServerCheck() (probe.Result, string, error) {
	// setup the prober
	server.Once.Do(func() {
		if server.Prober != nil {
			return
		}
		const followNonLocalRedirects = true
		server.Prober = httpprober.NewWithTLSConfig(server.TLSConfig, followNonLocalRedirects)
	})

	scheme := "http"
	if server.EnableHTTPS {
		scheme = "https"
	}
	url := utilnet.FormatURL(scheme, server.Addr, server.Port, server.Path)

	result, data, err := server.Prober.Probe(url, nil, probeTimeOut)

	if err != nil {
		return probe.Unknown, "", err
	}
	if result == probe.Failure {
		return probe.Failure, string(data), err
	}
	if server.Validate != nil {
		if err := server.Validate([]byte(data)); err != nil {
			return probe.Failure, string(data), err
		}
	}
	return result, string(data), nil
}
