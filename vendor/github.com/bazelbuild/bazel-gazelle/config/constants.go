package config

const (
	// RulesGoRepoName is the canonical name of the rules_go repository. It must
	// match the workspace name in WORKSPACE.
	// TODO(jayconrod): move to language/go.
	RulesGoRepoName = "io_bazel_rules_go"

	// GazelleImportsKey is an internal attribute that lists imported packages
	// on generated rules. It is replaced with "deps" during import resolution.
	GazelleImportsKey = "_gazelle_imports"
)
