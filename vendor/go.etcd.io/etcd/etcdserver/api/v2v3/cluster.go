package v2v3

import (
	"go.etcd.io/etcd/etcdserver/api/membership"
	"go.etcd.io/etcd/pkg/types"

	"github.com/coreos/go-semver/semver"
)

func (s *v2v3Server) ID() types.ID {
	// TODO: use an actual member ID
	return types.ID(0xe7cd2f00d)
}
func (s *v2v3Server) ClientURLs() []string                  { panic("STUB") }
func (s *v2v3Server) Members() []*membership.Member         { panic("STUB") }
func (s *v2v3Server) Member(id types.ID) *membership.Member { panic("STUB") }
func (s *v2v3Server) Version() *semver.Version              { panic("STUB") }
