package generators

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

// genExpansion produces a file for a group client, e.g. ExtensionsClient for the extension group.
type genExpansion struct {
	generator.DefaultGen
	groupPackagePath string
	// types in a group
	types []*types.Type
}

// We only want to call GenerateType() once per group.
func (g *genExpansion) Filter(c *generator.Context, t *types.Type) bool {
	return len(g.types) == 0 || t == g.types[0]
}

func (g *genExpansion) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	for _, t := range g.types {
		if _, err := os.Stat(filepath.Join(g.groupPackagePath, strings.ToLower(t.Name.Name+"_expansion.go"))); os.IsNotExist(err) {
			sw.Do(expansionInterfaceTemplate, t)
		}
	}
	return sw.Error()
}

var expansionInterfaceTemplate = `
type $.|public$Expansion interface {}
`
