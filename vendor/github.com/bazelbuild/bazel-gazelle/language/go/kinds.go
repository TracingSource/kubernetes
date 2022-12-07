package golang

import "github.com/bazelbuild/bazel-gazelle/rule"

var goKinds = map[string]rule.KindInfo{
	"filegroup": {
		NonEmptyAttrs:  map[string]bool{"srcs": true},
		MergeableAttrs: map[string]bool{"srcs": true},
	},
	"go_binary": {
		MatchAny: true,
		NonEmptyAttrs: map[string]bool{
			"deps":  true,
			"embed": true,
			"srcs":  true,
		},
		SubstituteAttrs: map[string]bool{"embed": true},
		MergeableAttrs: map[string]bool{
			"cgo":       true,
			"clinkopts": true,
			"copts":     true,
			"embed":     true,
			"srcs":      true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
	"go_library": {
		MatchAttrs: []string{"importpath"},
		NonEmptyAttrs: map[string]bool{
			"deps":  true,
			"embed": true,
			"srcs":  true,
		},
		SubstituteAttrs: map[string]bool{
			"embed": true,
		},
		MergeableAttrs: map[string]bool{
			"cgo":        true,
			"clinkopts":  true,
			"copts":      true,
			"embed":      true,
			"importmap":  true,
			"importpath": true,
			"srcs":       true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
	"go_proto_library": {
		MatchAttrs: []string{"importpath"},
		NonEmptyAttrs: map[string]bool{
			"deps":  true,
			"embed": true,
			"proto": true,
			"srcs":  true,
		},
		SubstituteAttrs: map[string]bool{"proto": true},
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"importpath": true,
			"importmap":  true,
			"cgo":        true,
			"clinkopts":  true,
			"copts":      true,
			"embed":      true,
			"proto":      true,
			"compilers":  true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
	"go_repository": {
		MatchAttrs: []string{"importpath"},
		NonEmptyAttrs: map[string]bool{
			"importpath": true,
		},
		MergeableAttrs: map[string]bool{
			"commit":       true,
			"importpath":   true,
			"remote":       true,
			"replace":      true,
			"sha256":       true,
			"strip_prefix": true,
			"sum":          true,
			"tag":          true,
			"type":         true,
			"urls":         true,
			"vcs":          true,
			"version":      true,
		},
	},
	"go_test": {
		NonEmptyAttrs: map[string]bool{
			"deps":  true,
			"embed": true,
			"srcs":  true,
		},
		MergeableAttrs: map[string]bool{
			"cgo":       true,
			"clinkopts": true,
			"copts":     true,
			"embed":     true,
			"srcs":      true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
}

var goLoads = []rule.LoadInfo{
	{
		Name: "@io_bazel_rules_go//go:def.bzl",
		Symbols: []string{
			"cgo_library",
			"go_binary",
			"go_library",
			"go_prefix",
			"go_repository",
			"go_test",
		},
	}, {
		Name: "@io_bazel_rules_go//proto:def.bzl",
		Symbols: []string{
			"go_grpc_library",
			"go_proto_library",
		},
	}, {
		Name: "@bazel_gazelle//:deps.bzl",
		Symbols: []string{
			"go_repository",
		},
		After: []string{
			"go_rules_dependencies",
			"go_register_toolchains",
			"gazelle_dependencies",
		},
	},
}

func (_ *goLang) Kinds() map[string]rule.KindInfo { return goKinds }
func (_ *goLang) Loads() []rule.LoadInfo          { return goLoads }
