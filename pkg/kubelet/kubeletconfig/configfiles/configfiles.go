package configfiles

import (
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	kubeletscheme "k8s.io/kubernetes/pkg/kubelet/apis/config/scheme"
	utilcodec "k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/codec"
	utilfs "k8s.io/kubernetes/pkg/util/filesystem"
)

// 	@implementBy: fsLoader{}
//
// Loader loads configuration from a storage layer
type Loader interface {
	// Load loads and returns the KubeletConfiguration from the storage layer,
	// or an error if a configuration could not be loaded
	Load() (*kubeletconfig.KubeletConfiguration, error)
}

// fsLoader loads configuration from `configDir`
type fsLoader struct {
	// fs is the filesystem where the config files exist; can be mocked for testing
	fs utilfs.Filesystem
	// kubeletCodecs is the scheme used to decode config files
	kubeletCodecs *serializer.CodecFactory
	// kubeletFile is an absolute path to the file containing a serialized KubeletConfiguration
	kubeletFile string
}

// NewFsLoader returns a Loader that loads a KubeletConfiguration from the `kubeletFile`
func NewFsLoader(fs utilfs.Filesystem, kubeletFile string) (Loader, error) {
	_, kubeletCodecs, err := kubeletscheme.NewSchemeAndCodecs(serializer.EnableStrict)
	if err != nil {
		return nil, err
	}

	return &fsLoader{
		fs:            fs,
		kubeletCodecs: kubeletCodecs,
		kubeletFile:   kubeletFile,
	}, nil
}

// caller:
// 	1. cmd/kubelet/app/server.go -> loadConfigFile()
func (loader *fsLoader) Load() (*kubeletconfig.KubeletConfiguration, error) {
	data, err := loader.fs.ReadFile(loader.kubeletFile)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read kubelet config file %q, error: %v", 
			loader.kubeletFile, err,
		)
	}

	// no configuration is an error, some parameters are required
	if len(data) == 0 {
		return nil, fmt.Errorf("kubelet config file %q was empty", loader.kubeletFile)
	}

	kc, err := utilcodec.DecodeKubeletConfiguration(loader.kubeletCodecs, data)
	if err != nil {
		return nil, err
	}

	// make all paths absolute
	resolveRelativePaths(
		kubeletconfig.KubeletConfigurationPathRefs(kc), 
		filepath.Dir(loader.kubeletFile),
	)
	return kc, nil
}

// resolveRelativePaths makes relative paths absolute by resolving them against `root`
func resolveRelativePaths(paths []*string, root string) {
	for _, path := range paths {
		// leave empty paths alone, "no path" is a valid input
		// do not attempt to resolve paths that are already absolute
		if len(*path) > 0 && !filepath.IsAbs(*path) {
			*path = filepath.Join(root, *path)
		}
	}
}
