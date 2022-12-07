package protobuf

import (
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

type ImportTracker struct {
	namer.DefaultImportTracker
}

func NewImportTracker(local types.Name, typesToAdd ...*types.Type) *ImportTracker {
	tracker := namer.NewDefaultImportTracker(local)
	tracker.IsInvalidType = func(t *types.Type) bool { return t.Kind != types.Protobuf }
	tracker.LocalName = func(name types.Name) string { return name.Package }
	tracker.PrintImport = func(path, name string) string { return path }

	tracker.AddTypes(typesToAdd...)
	return &ImportTracker{
		DefaultImportTracker: tracker,
	}
}

// AddNullable ensures that support for the nullable Gogo-protobuf extension is added.
func (tracker *ImportTracker) AddNullable() {
	tracker.AddType(&types.Type{
		Kind: types.Protobuf,
		Name: types.Name{
			Name:    "nullable",
			Package: "gogoproto",
			Path:    "github.com/gogo/protobuf/gogoproto/gogo.proto",
		},
	})
}
