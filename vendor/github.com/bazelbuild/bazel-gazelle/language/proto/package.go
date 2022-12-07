package proto

import "path/filepath"

// Package contains metadata for a set of .proto files that have the
// same package name. This translates to a proto_library rule.
type Package struct {
	Name        string
	Files       map[string]FileInfo
	Imports     map[string]bool
	Options     map[string]string
	HasServices bool
}

func newPackage(name string) *Package {
	return &Package{
		Name:    name,
		Files:   map[string]FileInfo{},
		Imports: map[string]bool{},
		Options: map[string]string{},
	}
}

func (p *Package) addFile(info FileInfo) {
	p.Files[info.Name] = info
	for _, imp := range info.Imports {
		p.Imports[imp] = true
	}
	for _, opt := range info.Options {
		p.Options[opt.Key] = opt.Value
	}
	p.HasServices = p.HasServices || info.HasServices
}

func (p *Package) addGenFile(dir, name string) {
	p.Files[name] = FileInfo{
		Name: name,
		Path: filepath.Join(dir, filepath.FromSlash(name)),
	}
}
