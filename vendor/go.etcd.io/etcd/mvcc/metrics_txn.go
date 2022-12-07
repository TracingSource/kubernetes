package mvcc

import "go.etcd.io/etcd/lease"

type metricsTxnWrite struct {
	TxnWrite
	ranges  uint
	puts    uint
	deletes uint
}

func newMetricsTxnRead(tr TxnRead) TxnRead {
	return &metricsTxnWrite{&txnReadWrite{tr}, 0, 0, 0}
}

func newMetricsTxnWrite(tw TxnWrite) TxnWrite {
	return &metricsTxnWrite{tw, 0, 0, 0}
}

func (tw *metricsTxnWrite) Range(key, end []byte, ro RangeOptions) (*RangeResult, error) {
	tw.ranges++
	return tw.TxnWrite.Range(key, end, ro)
}

func (tw *metricsTxnWrite) DeleteRange(key, end []byte) (n, rev int64) {
	tw.deletes++
	return tw.TxnWrite.DeleteRange(key, end)
}

func (tw *metricsTxnWrite) Put(key, value []byte, lease lease.LeaseID) (rev int64) {
	tw.puts++
	return tw.TxnWrite.Put(key, value, lease)
}

func (tw *metricsTxnWrite) End() {
	defer tw.TxnWrite.End()
	if sum := tw.ranges + tw.puts + tw.deletes; sum > 1 {
		txnCounter.Inc()
		txnCounterDebug.Inc() // TODO: remove in 3.5 release
	}

	ranges := float64(tw.ranges)
	rangeCounter.Add(ranges)
	rangeCounterDebug.Add(ranges) // TODO: remove in 3.5 release

	puts := float64(tw.puts)
	putCounter.Add(puts)
	putCounterDebug.Add(puts) // TODO: remove in 3.5 release

	deletes := float64(tw.deletes)
	deleteCounter.Add(deletes)
	deleteCounterDebug.Add(deletes) // TODO: remove in 3.5 release
}
