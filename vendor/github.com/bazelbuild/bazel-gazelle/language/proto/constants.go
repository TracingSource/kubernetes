package proto

const (
	// PackageInfoKey is the name of a private attribute set on generated
	// proto_library rules. This attribute contains a Package record which
	// describes the library and its sources.
	PackageKey = "_package"

	// wellKnownTypesGoPrefix is the import path for the Go repository containing
	// pre-generated code for the Well Known Types.
	wellKnownTypesGoPrefix = "github.com/golang/protobuf"
)
