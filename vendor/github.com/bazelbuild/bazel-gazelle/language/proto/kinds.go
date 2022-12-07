package proto

import "github.com/bazelbuild/bazel-gazelle/rule"

var protoKinds = map[string]rule.KindInfo{
	"proto_library": {
		NonEmptyAttrs: map[string]bool{"srcs": true},
		MergeableAttrs: map[string]bool{
			"srcs": true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	},
}

func (_ *protoLang) Kinds() map[string]rule.KindInfo { return protoKinds }
func (_ *protoLang) Loads() []rule.LoadInfo          { return nil }
