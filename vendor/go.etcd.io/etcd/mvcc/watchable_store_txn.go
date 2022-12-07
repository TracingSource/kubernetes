package mvcc

import (
	"go.etcd.io/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/pkg/traceutil"
)

func (tw *watchableStoreTxnWrite) End() {
	changes := tw.Changes()
	if len(changes) == 0 {
		tw.TxnWrite.End()
		return
	}

	rev := tw.Rev() + 1
	evs := make([]mvccpb.Event, len(changes))
	for i, change := range changes {
		evs[i].Kv = &changes[i]
		if change.CreateRevision == 0 {
			evs[i].Type = mvccpb.DELETE
			evs[i].Kv.ModRevision = rev
		} else {
			evs[i].Type = mvccpb.PUT
		}
	}

	// end write txn under watchable store lock so the updates are visible
	// when asynchronous event posting checks the current store revision
	tw.s.mu.Lock()
	tw.s.notify(rev, evs)
	tw.TxnWrite.End()
	tw.s.mu.Unlock()
}

type watchableStoreTxnWrite struct {
	TxnWrite
	s *watchableStore
}

func (s *watchableStore) Write(trace *traceutil.Trace) TxnWrite {
	return &watchableStoreTxnWrite{s.store.Write(trace), s}
}
