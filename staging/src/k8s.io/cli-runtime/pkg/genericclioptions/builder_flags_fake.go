package genericclioptions

import (
	"k8s.io/cli-runtime/pkg/resource"
)

// NewSimpleResourceFinder builds a super simple ResourceFinder that just iterates over the objects you provided
func NewSimpleFakeResourceFinder(infos ...*resource.Info) ResourceFinder {
	return &fakeResourceFinder{
		Infos: infos,
	}
}

type fakeResourceFinder struct {
	Infos []*resource.Info
}

// Do implements the interface
func (f *fakeResourceFinder) Do() resource.Visitor {
	return &fakeResourceResult{
		Infos: f.Infos,
	}
}

type fakeResourceResult struct {
	Infos []*resource.Info
}

// Visit just iterates over info
func (r *fakeResourceResult) Visit(fn resource.VisitorFunc) error {
	for _, info := range r.Infos {
		err := fn(info, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
