package main

import (
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/language/go"
	"github.com/bazelbuild/bazel-gazelle/language/proto"
)

var languages = []language.Language{
	proto.NewLanguage(),
	golang.NewLanguage(),
}
