package builder

import (
	"sort"

	"github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

type parameters []spec.Parameter

func (s parameters) Len() int      { return len(s) }
func (s parameters) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// byNameIn used in sorting parameters by Name and In fields.
type byNameIn struct {
	parameters
}

func (s byNameIn) Less(i, j int) bool {
	return s.parameters[i].Name < s.parameters[j].Name || (s.parameters[i].Name == s.parameters[j].Name && s.parameters[i].In < s.parameters[j].In)
}

// SortParameters sorts parameters by Name and In fields.
func sortParameters(p []spec.Parameter) {
	sort.Sort(byNameIn{p})
}

func groupRoutesByPath(routes []restful.Route) map[string][]restful.Route {
	pathToRoutes := make(map[string][]restful.Route)
	for _, r := range routes {
		pathToRoutes[r.Path] = append(pathToRoutes[r.Path], r)
	}
	return pathToRoutes
}

func mapKeyFromParam(param *restful.Parameter) interface{} {
	return struct {
		Name string
		Kind int
	}{
		Name: param.Data().Name,
		Kind: param.Data().Kind,
	}
}
