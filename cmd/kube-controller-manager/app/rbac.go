package app

import (
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/controller/clusterroleaggregation"
)

func startClusterRoleAggregrationController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"}] {
		return nil, false, nil
	}
	go clusterroleaggregation.NewClusterRoleAggregation(
		ctx.InformerFactory.Rbac().V1().ClusterRoles(),
		ctx.ClientBuilder.ClientOrDie("clusterrole-aggregation-controller").RbacV1(),
	).Run(5, ctx.Stop)
	return nil, true, nil
}
