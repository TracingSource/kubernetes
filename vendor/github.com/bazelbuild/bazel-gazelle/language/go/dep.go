package golang

import (
	"io/ioutil"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
	toml "github.com/pelletier/go-toml"
)

type depLockFile struct {
	Projects []depProject `toml:"projects"`
}

type depProject struct {
	Name     string `toml:"name"`
	Revision string `toml:"revision"`
	Source   string `toml:"source"`
}

func importReposFromDep(args language.ImportReposArgs) language.ImportReposResult {
	data, err := ioutil.ReadFile(args.Path)
	if err != nil {
		return language.ImportReposResult{Error: err}
	}
	var file depLockFile
	if err := toml.Unmarshal(data, &file); err != nil {
		return language.ImportReposResult{Error: err}
	}

	gen := make([]*rule.Rule, len(file.Projects))
	for i, p := range file.Projects {
		gen[i] = rule.NewRule("go_repository", label.ImportPathToBazelRepoName(p.Name))
		gen[i].SetAttr("importpath", p.Name)
		gen[i].SetAttr("commit", p.Revision)
		if p.Source != "" {
			// TODO(#411): Handle source directives correctly. It may be an import
			// path, or a URL. In the case of an import path, we should resolve it
			// to the correct remote and vcs. In the case of a URL, we should
			// correctly determine what VCS to use (the URL will usually start
			// with "https://", which is used by multiple VCSs).
			gen[i].SetAttr("remote", p.Source)
			gen[i].SetAttr("vcs", "git")
		}
	}
	sort.SliceStable(gen, func(i, j int) bool {
		return gen[i].Name() < gen[j].Name()
	})

	return language.ImportReposResult{Gen: gen}
}
