// +build tools

// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import (
	_ "github.com/bazelbuild/bazel-gazelle/cmd/gazelle"
	_ "github.com/bazelbuild/buildtools/buildozer"
	_ "github.com/cespare/prettybench"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "golang.org/x/lint/golint"
	_ "gotest.tools"
	_ "gotest.tools/gotestsum"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "k8s.io/code-generator/cmd/go-to-protobuf"
	_ "k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo"
	_ "k8s.io/gengo/examples/deepcopy-gen/generators"
	_ "k8s.io/gengo/examples/defaulter-gen/generators"
	_ "k8s.io/gengo/examples/import-boss/generators"
	_ "k8s.io/gengo/examples/set-gen/generators"
	_ "k8s.io/kube-openapi/cmd/openapi-gen"
	_ "k8s.io/repo-infra/cmd/kazel"
)
