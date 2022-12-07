package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// ProjectsOps is the manually implemented methods for the Projects service.
type ProjectsOps interface {
	Get(ctx context.Context, projectID string) (*compute.Project, error)
	SetCommonInstanceMetadata(ctx context.Context, projectID string, m *compute.Metadata) error
}

// MockProjectOpsState is stored in the mock.X field.
type MockProjectOpsState struct {
	metadata map[string]*compute.Metadata
}

// Get a project by projectID.
func (m *MockProjects) Get(ctx context.Context, projectID string) (*compute.Project, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if p, ok := m.Objects[*meta.GlobalKey(projectID)]; ok {
		return p.ToGA(), nil
	}
	return nil, &googleapi.Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("MockProjects %v not found", projectID),
	}
}

// Get a project by projectID.
func (g *GCEProjects) Get(ctx context.Context, projectID string) (*compute.Project, error) {
	rk := &RateLimitKey{
		ProjectID: projectID,
		Operation: "Get",
		Version:   meta.Version("ga"),
		Service:   "Projects",
	}
	if err := g.s.RateLimiter.Accept(ctx, rk); err != nil {
		return nil, err
	}
	call := g.s.GA.Projects.Get(projectID)
	call.Context(ctx)
	return call.Do()
}

// SetCommonInstanceMetadata for a given project.
func (m *MockProjects) SetCommonInstanceMetadata(ctx context.Context, projectID string, meta *compute.Metadata) error {
	if m.X == nil {
		m.X = &MockProjectOpsState{metadata: map[string]*compute.Metadata{}}
	}
	state := m.X.(*MockProjectOpsState)
	state.metadata[projectID] = meta
	return nil
}

// SetCommonInstanceMetadata for a given project.
func (g *GCEProjects) SetCommonInstanceMetadata(ctx context.Context, projectID string, m *compute.Metadata) error {
	rk := &RateLimitKey{
		ProjectID: projectID,
		Operation: "SetCommonInstanceMetadata",
		Version:   meta.Version("ga"),
		Service:   "Projects",
	}
	if err := g.s.RateLimiter.Accept(ctx, rk); err != nil {
		return err
	}
	call := g.s.GA.Projects.SetCommonInstanceMetadata(projectID, m)
	call.Context(ctx)

	op, err := call.Do()
	if err != nil {
		return err
	}
	return g.s.WaitForCompletion(ctx, op)
}
