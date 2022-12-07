package rule

import (
	"regexp"

	bzl "github.com/bazelbuild/buildtools/build"
)

// Directive is a key-value pair extracted from a top-level comment in
// a build file. Directives have the following format:
//
//     # gazelle:key value
//
// Keys may not contain spaces. Values may be empty and may contain spaces,
// but surrounding space is trimmed.
type Directive struct {
	Key, Value string
}

// TODO(jayconrod): annotation directives will apply to an individual rule.
// They must appear in the block of comments above that rule.

// ParseDirectives scans f for Gazelle directives. The full list of directives
// is returned. Errors are reported for unrecognized directives and directives
// out of place (after the first statement).
func ParseDirectives(f *bzl.File) []Directive {
	return parseDirectives(f.Stmt)
}

// ParseDirectivesFromMacro scans a macro body for Gazelle directives. The
// full list of directives is returned. Errors are reported for unrecognized
// directives and directives out of place (after the first statement).
func ParseDirectivesFromMacro(f *bzl.DefStmt) []Directive {
	return parseDirectives(f.Body)
}

func parseDirectives(stmt []bzl.Expr) []Directive {
	var directives []Directive
	parseComment := func(com bzl.Comment) {
		match := directiveRe.FindStringSubmatch(com.Token)
		if match == nil {
			return
		}
		key, value := match[1], match[2]
		directives = append(directives, Directive{key, value})
	}

	for _, s := range stmt {
		coms := s.Comment()
		for _, com := range coms.Before {
			parseComment(com)
		}
		for _, com := range coms.After {
			parseComment(com)
		}
	}
	return directives
}

var directiveRe = regexp.MustCompile(`^#\s*gazelle:(\w+)\s*(.*?)\s*$`)
