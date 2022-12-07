package resourcequota

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestResourceQuotaStrategy(t *testing.T) {
	if !Strategy.NamespaceScoped() {
		t.Errorf("ResourceQuota should be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("ResourceQuota should not allow create on update")
	}
	resourceQuota := &api.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Status: api.ResourceQuotaStatus{
			Used: api.ResourceList{
				api.ResourceCPU:                    resource.MustParse("1"),
				api.ResourceMemory:                 resource.MustParse("1Gi"),
				api.ResourcePods:                   resource.MustParse("1"),
				api.ResourceServices:               resource.MustParse("1"),
				api.ResourceReplicationControllers: resource.MustParse("1"),
				api.ResourceQuotas:                 resource.MustParse("1"),
			},
			Hard: api.ResourceList{
				api.ResourceCPU:                    resource.MustParse("100"),
				api.ResourceMemory:                 resource.MustParse("4Gi"),
				api.ResourcePods:                   resource.MustParse("10"),
				api.ResourceServices:               resource.MustParse("10"),
				api.ResourceReplicationControllers: resource.MustParse("10"),
				api.ResourceQuotas:                 resource.MustParse("1"),
			},
		},
	}
	Strategy.PrepareForCreate(genericapirequest.NewContext(), resourceQuota)
	if resourceQuota.Status.Used != nil {
		t.Errorf("ResourceQuota does not allow setting status on create")
	}
}
