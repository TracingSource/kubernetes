// Package configz serves ComponentConfig objects from running components.
//
// Each component that wants to serve its ComponentConfig creates a Config
// object, and the program should call InstallHandler once. e.g.,
//  func main() {
//  	boatConfig := getBoatConfig()
//  	planeConfig := getPlaneConfig()
//
//  	bcz, err := configz.New("boat")
//  	if err != nil {
//  		panic(err)
//  	}
//  	bcz.Set(boatConfig)
//
//  	pcz, err := configz.New("plane")
//  	if err != nil {
//  		panic(err)
//  	}
//  	pcz.Set(planeConfig)
//
//  	configz.InstallHandler(http.DefaultServeMux)
//  	http.ListenAndServe(":8080", http.DefaultServeMux)
//  }
package configz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	configsGuard sync.RWMutex
	configs      = map[string]*Config{}
)

// Config is a handle to a ComponentConfig object. Don't create these directly;
// use New() instead.
type Config struct {
	val interface{}
}

// caller: 
// 	1. cmd/controller-manager/app/serve.go -> NewBaseHandler()
//
// InstallHandler adds an HTTP handler on the given mux for the "/configz"
// endpoint which serves all registered ComponentConfigs in JSON format.
func InstallHandler(m mux) {
	m.Handle("/configz", http.HandlerFunc(handle))
}

type mux interface {
	Handle(string, http.Handler)
}

// caller: 
// 	1. cmd/kubelet/app/server.go -> initConfigz()
//
// New creates a Config object with the given name.
// Each Config is registered with this package's "/configz" handler.
func New(name string) (*Config, error) {
	configsGuard.Lock()
	defer configsGuard.Unlock()
	if _, found := configs[name]; found {
		return nil, fmt.Errorf("register config %q twice", name)
	}
	newConfig := Config{}
	configs[name] = &newConfig
	return &newConfig, nil
}

// Delete removes the named ComponentConfig from this package's "/configz"
// handler.
func Delete(name string) {
	configsGuard.Lock()
	defer configsGuard.Unlock()
	delete(configs, name)
}

// Set sets the ComponentConfig for this Config.
func (v *Config) Set(val interface{}) {
	configsGuard.Lock()
	defer configsGuard.Unlock()
	v.val = val
}

// MarshalJSON marshals the ComponentConfig as JSON data.
func (v *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.val)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if err := write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func write(w http.ResponseWriter) error {
	var b []byte
	var err error
	func() {
		configsGuard.RLock()
		defer configsGuard.RUnlock()
		b, err = json.Marshal(configs)
	}()
	if err != nil {
		return fmt.Errorf("error marshaling json: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, err = w.Write(b)
	return err
}
